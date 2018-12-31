/*
Copyright 2018 The Elasticshift Authors.
*/
package types

// Identity ..
type Identity struct {
	Issuer      string
	HostAndPort string
	caPath      string
	ID          string
	Secret      string
	RedirectURI string
}
