package vcs

import (
	"fmt"
	"net/http"

	"gitlab.com/conspico/esh/team"

	"github.com/spf13/viper"
)

// Service ..
type Service interface {
	Authorize(subdomain, provider string, r *http.Request) (AuthorizeResponse, error)
	Authorized(subdomain, provider, code string) error
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
		GithubProvider(conf.GetString("github.key"), conf.GetString("github.secret")),
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

	url := p.Authorize()

	return AuthorizeResponse{Err: nil, URL: url, Request: r}, nil
}

// Authorized ..
// Invoked when authorization finished by oauth app
func (s service) Authorized(subdomain, provider, code string) error {

	p, err := s.vcsProviders.Get(provider)
	if err != nil {
		return err
	}

	u, err := p.Authorized(code)

	// persist user
	fmt.Println(u.AccessToken)
	return nil
}
