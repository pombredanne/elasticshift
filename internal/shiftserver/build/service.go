/*
Copyright 2018 The Elasticshift Authors.
*/
package build

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/internal/pkg/storage"
	"github.com/elasticshift/elasticshift/internal/pkg/utils"
	"github.com/elasticshift/elasticshift/internal/shiftserver/integration"
	itypes "github.com/elasticshift/elasticshift/internal/shiftserver/integration/types"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
)

var (
	errBuildNotFound = "Build identifier not found"

	retryDuration, _ = time.ParseDuration("1s")
)

type service struct {
	buildStore       store.Build
	teamStore        store.Team
	sysconfStore     store.Sysconf
	integrationStore store.Integration
	defaultsStore    store.Defaults
	logger           *logrus.Entry
	loggr            logger.Loggr
}

// Service ..
type Service interface {
	Viewlog(w http.ResponseWriter, r *http.Request)
}

// NewService ..
func NewService(loggr logger.Loggr, s store.Shift) Service {

	l := loggr.GetLogger("service/build")
	return &service{
		buildStore:       s.Build,
		teamStore:        s.Team,
		sysconfStore:     s.Sysconf,
		integrationStore: s.Integration,
		defaultsStore:    s.Defaults,
		logger:           l,
		loggr:            loggr,
	}
}

func (s service) Viewlog(w http.ResponseWriter, r *http.Request) {

	s.logger.Infoln("Fetching log.. ")

	var follow string
	follow = r.URL.Query().Get("follow")
	if follow != "" {
		_, err := strconv.ParseBool(follow)
		if err != nil {
			http.Error(w, "Query param 'follow' should contain boolean only.", http.StatusBadRequest)
			return
		}
	} else {
		follow = "false"
	}

	buildID := mux.Vars(r)["buildid"]
	if buildID == "" {
		http.Error(w, "URL doesn't container build identifier.", http.StatusBadRequest)
		return
	}
	s.logger.Infoln("BuildID=", buildID)

	subBuildID := mux.Vars(r)["subbuildid"]
	if subBuildID == "" {
		http.Error(w, "URL doesn't container sub-build identifier.", http.StatusBadRequest)
		return
	}

	var b types.Build
	var err error

	b, err = s.buildStore.FetchBuildByID(buildID)
	if err != nil && err.Error() == "not found" {
		http.Error(w, "Build identifier not found.", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Fetching the build failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var sb types.SubBuild
	for _, v := range b.SubBuilds {
		if v.ID == subBuildID {
			sb = v
			break
		}
	}

	// fetch log directly from the container
	if sb.Status == types.BuildStatusWaiting || sb.Status == types.BuildStatusPreparing || sb.Status == types.BuildStatusRunning {

		// Get the default storage based on team defaults
		def, err := s.defaultsStore.FindByReferenceId(b.Team)
		if err != nil {
			http.Error(w, "Failed to get default storage :", http.StatusBadRequest)
			return
		}

		// Get the details of the storeage
		var stor types.Storage
		err = s.integrationStore.FindByID(def.StorageID, &stor)
		if err != nil {
			http.Error(w, "Failed to fetch log", http.StatusInternalServerError)
			return
		}

		for {

			sb, err = s.buildStore.FetchSubBuild(buildID, subBuildID)
			if err != nil || (sb.Metadata != nil && sb.Metadata.PodName != "") {
				break
			}

			time.Sleep(retryDuration)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Fetching the build failed: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Get the details of the integration
		var i types.ContainerEngine
		err = s.integrationStore.FindByID(def.ContainerEngineID, &i)
		if err != nil {
			http.Error(w, "Failed to fetch log", http.StatusInternalServerError)
			return
		}

		// connect to container engine cluster
		ce, err := integration.NewContainerEngine(s.loggr, i, stor)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch log, error when connecting to container engine: %s", err.Error()), http.StatusInternalServerError)
		}

		opts := &itypes.StreamLogOptions{Pod: sb.Metadata.PodName, BuildID: buildID, W: w, Follow: follow}
		if i.Kind == integration.KIND_KUBERNETES_TYPE {
			opts.Pod = sb.Metadata.PodName
		} else if i.Kind == integration.KIND_DOCKERSWARM {
			opts.ContainerID = sb.Metadata.ContainerID
		}

		readCloser, err := ce.StreamLog(opts)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.WriteHeader(http.StatusOK)

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		stream(w, readCloser)
	} else {

		nodeID := mux.Vars(r)["nodeid"]
		if nodeID == "" {
			http.Error(w, "URL doesn't container build identifier.", http.StatusBadRequest)
			return
		}

		sm := &types.StorageMetadata{
			TeamID:       b.Team,
			BuildID:      buildID,
			RepositoryID: b.RepositoryID,
			SubBuildID:   subBuildID,
			Branch:       b.Branch,
			Path:         b.StoragePath,
		}

		// Get the details of the storeage
		var stor types.Storage
		err = s.integrationStore.FindByID(b.StorageID, &stor)
		if err != nil {
			http.Error(w, "Failed to fetch log", http.StatusInternalServerError)
			return
		}

		ss, err := storage.NewWithMetadata(s.logger, &stor, sm)
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot connect to storage: %v", err), http.StatusInternalServerError)
		}

		r, err := ss.GetLog(nodeID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch log from storage: %v", err), http.StatusInternalServerError)
		}

		stream(w, r)
		// Fetch logs from storage
		// w.Write([]byte("Streaming non running builds aren't supported yet. "))
	}

	//if err != nil {
	//	//handle streaming error
	//}
}

func stream(w http.ResponseWriter, r io.ReadCloser) error {

	done := make(chan bool, 1)
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		for {
			select {
			case <-notify:
			case <-done:
			}

			r.Close()
			break
		}
	}()

	// stream the logs
	_, err := io.Copy(utils.StreamWriter(w), r)
	if err != nil {
		done <- true
		return err
	}
	done <- true
	return nil
}
