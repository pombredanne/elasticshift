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
	Authorized(subdomain, provider, code string) (VCS, error)
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

func (s service) Authorize(subdomain, provider string, r *http.Request) (AuthorizeResponse, error) {

	p, err := s.vcsProviders.Get(provider)
	if err != nil {
		return AuthorizeResponse{}, err
	}

	url := p.Authorize(subdomain)

	return AuthorizeResponse{Err: nil, URL: url, Request: r}, nil
}

// Authorized ..
// Invoked when authorization finished by oauth app
func (s service) Authorized(subdomain, provider, code string) (VCS, error) {

	p, err := s.vcsProviders.Get(provider)
	if err != nil {
		return VCS{}, err
	}

	teamID, err := s.teamRepository.GetTeamID(subdomain)
	if err != nil {
		return VCS{}, err
	}

	u, err := p.Authorized(code)
	if err != nil {
		fmt.Println(err)
		return VCS{}, err
	}

	// persist user
	fmt.Println(u)
	u.ID, _ = util.NewUUID()
	u.TeamID = teamID
	u.CreatedDt = time.Now()
	u.UpdatedDt = time.Now()
	err = s.vcsRepository.Save(&u)

	return u, err
}
