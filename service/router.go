package service

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"gitlab.com/conspico/esh/service/api"
)

// RegisterRoutes ...
func RegisterRoutes(router *mux.Router) {

	a := alice.New()
	chain := a

	// oauth
	router.Handle("/auth/{provider}", chain.ThenFunc(api.Authorize))
	router.Handle("/auth/{provider}/callback", chain.ThenFunc(api.Autorized))
}
