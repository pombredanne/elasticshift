package vcs

import (
	"fmt"
	"net/url"

	"gitlab.com/conspico/esh/model"
)

// Provider ..
type Provider interface {
	Name() string

	Authorize() string

	Authorized(*url.URL) (model.VCSUser, error)
}

// Providers type
type Providers map[string]Provider

var providers = Providers{}

// Use or initialize the providers
func Use(pvider ...Provider) {
	for _, p := range pvider {
		providers[p.Name()] = p
	}
}

// Get the provider by namee
func Get(name string) (Provider, error) {

	p := providers[name]
	if p == nil {
		return nil, fmt.Errorf("No provider found for %s", name)
	}
	return p, nil
}
