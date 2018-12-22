/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
	"gopkg.in/mgo.v2/bson"
)

const (
// INT_ContainerEngine int = iota + 1
// INT_Storage
)

const (

	// Kind
	PROVIDER_ONPREM = iota + 1
	PROVIDER_AZURE
	PROVIDER_DIGITALOCEAN
	PROVIDER_ALIBABACLOUD
)

const (
	// PROVODERS
	KIND_KUBERNETES_TYPE = iota + 1 //refer schema/integration container engine kind
	KIND_DOCKERSWARM
	KIND_DCOS
)

type service struct {
	teamStore        store.Team
	integrationStore store.Integration
	logger           *logrus.Entry
}

// SysconfService ..
type Service interface {
	UploadKubeConfigFile(w http.ResponseWriter, r *http.Request)
}

// NewVCSService ..
func NewService(loggr logger.Loggr, d store.Database, s store.Shift) Service {

	l := loggr.GetLogger("service/integration")
	return &service{
		teamStore:        s.Team,
		integrationStore: s.Integration,
		logger:           l,
	}
}

func (s service) UploadKubeConfigFile(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	name := r.FormValue("name")
	teamID := r.FormValue("team")
	provider := r.FormValue("provider")
	kind := r.FormValue("kind")

	file, _, err := r.FormFile("kubefile")
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ce types.ContainerEngine
	err = s.integrationStore.FindOne(bson.M{"team": teamID, "name": name}, &ce)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		http.Error(w, fmt.Sprintf("Failed to check if the given integration already exist :%v", err), http.StatusBadRequest)
		return
	}

	if ce.ID.Hex() != "" {
		http.Error(w, fmt.Sprintf("The container engine name '%s' already exist for your team", name), http.StatusBadRequest)
		return
	}

	if provider == "" {
		http.Error(w, "Please provide the provider.", http.StatusBadRequest)
	}

	if kind == "" {
		http.Error(w, "Please provide the kind.", http.StatusBadRequest)
	}

	providerVal, _ := strconv.Atoi(provider)

	kindVal, _ := strconv.Atoi(kind)

	i := types.ContainerEngine{}
	i.Name = name
	i.Team = teamID
	i.Kind = kindVal
	i.Provider = providerVal
	i.InternalType = INT_ContainerEngine
	i.KubeFile = buf.Bytes()

	err = s.integrationStore.Save(&i)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to add integration: %v", err), http.StatusInternalServerError)
		return
	}

	http.StatusText(http.StatusOK)
}
