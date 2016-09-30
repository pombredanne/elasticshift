package vcs

import "fmt"

// Provider ..
type Provider interface {
	Name() string

	Authorize() string

	Authorized(code string) (User, error)
}

// Providers type
type Providers struct {
	Providers map[string]Provider
}

// NewProviders ...
func NewProviders(pvider ...Provider) *Providers {

	var prov = make(map[string]Provider)
	for _, p := range pvider {
		prov[p.Name()] = p
	}
	return &Providers{prov}
}

// Get the provider by namee
func (prov Providers) Get(name string) (Provider, error) {

	p := prov.Providers[name]
	if p == nil {
		return nil, fmt.Errorf(errNoProviderFound, name)
	}
	return p, nil
}
