package vcs

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/conspico/esh/core/util"
	"gitlab.com/conspico/esh/team"

	"github.com/spf13/viper"
)

// Service ..
type Service interface {
	Authorize(subdomain, provider string, r *http.Request) (AuthorizeResponse, error)
	Authorized(subdomain, provider, code string, r *http.Request) (AuthorizeResponse, error)
	GetVCS(teamID string) (GetVCSResponse, error)
}

type service struct {
	vcsRepository  Repository
	teamRepository team.Repository
	vcsProviders   *Providers
	config         *viper.Viper
}

// NewService ..
func NewService(v Repository, t team.Repository, conf *viper.Viper) Service {

	providers := NewProviders(
		GithubProvider(conf.GetString("github.key"), conf.GetString("github.secret"), conf.GetString("github.callback")),
	)

	return &service{
		vcsProviders:   providers,
		vcsRepository:  v,
		teamRepository: t,
		config:         conf,
	}
}

func (s service) Authorize(teamID, provider string, r *http.Request) (AuthorizeResponse, error) {

	p, err := s.vcsProviders.Get(provider)
	if err != nil {
		return AuthorizeResponse{}, err
	}

	url := p.Authorize(teamID)

	return AuthorizeResponse{Err: nil, URL: url, Request: r}, nil
}

// Authorized ..
// Invoked when authorization finished by oauth app
func (s service) Authorized(teamID, provider, code string, r *http.Request) (AuthorizeResponse, error) {

	p, err := s.vcsProviders.Get(provider)
	if err != nil {
		return AuthorizeResponse{}, err
	}

	u, err := p.Authorized(code)
	if err != nil {
		fmt.Println(err)
		return AuthorizeResponse{}, err
	}

	// persist user
	u.ID, _ = util.NewUUID()
	u.TeamID = teamID
	u.CreatedDt = time.Now()
	u.UpdatedDt = time.Now()
	err = s.vcsRepository.Save(&u)

	url := r.Referer() + "/api/vcs"
	return AuthorizeResponse{Err: nil, URL: url, Request: r}, err
}

func (s service) GetVCS(teamID string) (GetVCSResponse, error) {

	result, err := s.vcsRepository.GetVCS(teamID)
	return GetVCSResponse{Result: result}, err
}
