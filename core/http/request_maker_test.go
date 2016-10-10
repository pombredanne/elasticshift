package http_test

import (
	"testing"

	chttp "gitlab.com/conspico/esh/core/http"
)

var (
	url = "http://conspico.elasticshift.com/:team/:user/list"
)

func TestGetPathParams(t *testing.T) {
	chttp.NewGetRequestMaker(url).PathParams("cp", "ghazni").QueryParam("code", "12345").Dispatch()
}

func TestGetPathParamsInvalid(t *testing.T) {
	err := chttp.NewGetRequestMaker(url).PathParams("cp", "ghazni", "shahm").QueryParam("code", "12345").Dispatch()
	if err != nil {
		t.Log(err)
	} else {
		t.Fail()
	}
}

func TestInvalidURL(t *testing.T) {

	err := chttp.NewGetRequestMaker(url).PathParams("cp", "ghazni").QueryParam("code", "12345").Dispatch()
	if err != nil {
		t.Log(err)
	}
}
