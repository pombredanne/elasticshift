package sysconf

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/pkg/identity/team"
	stypes "gitlab.com/conspico/elasticshift/pkg/store/types"
)

type service struct {
	teamStore team.Store
	logger    logrus.Logger
}

// SysconfService ..
type Service interface {
	UploadKubeConfigFile(w http.ResponseWriter, r *http.Request)
}

// NewVCSService ..
func NewService(logger logrus.Logger, d stypes.Database, teamStore team.Store) Service {

	return &service{
		teamStore: teamStore,
		logger:    logger,
	}
}

func (s service) UploadKubeConfigFile(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	teamName := r.FormValue("team")
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
	}

	err = s.teamStore.SaveKubeConfig(teamName, buf.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.StatusText(http.StatusOK)
}
