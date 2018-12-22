/*
Copyright 2017 The Elasticshift Authors.
*/
package dispatch_test

import (
	"testing"

	"github.com/elasticshift/elasticshift/core/dispatch"
)

var (
	url = "http://conspico.elasticshift.com/:team/:user/list"
)

func TestGetPathParams(t *testing.T) {
	dispatch.NewGetRequestMaker(url).PathParams("cp", "ghazni").QueryParam("code", "12345").Dispatch()
}

func TestGetPathParamsInvalid(t *testing.T) {
	err := dispatch.NewGetRequestMaker(url).PathParams("cp", "ghazni", "shahm").QueryParam("code", "12345").Dispatch()
	if err != nil {
		t.Log(err)
	} else {
		t.Fail()
	}
}

func TestInvalidURL(t *testing.T) {

	err := dispatch.NewGetRequestMaker(url).PathParams("cp", "ghazni").QueryParam("code", "12345").Dispatch()
	if err != nil {
		t.Log(err)
	}
}
