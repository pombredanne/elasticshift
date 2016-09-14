package service

import (
	"github.com/gorilla/mux"

	"gitlab.com/conspico/esh/model"
	"gitlab.com/conspico/esh/service/api"
	"gitlab.com/conspico/esh/service/handler"
)

// RegisterRoutes ...
func RegisterRoutes(appCtx model.AppContext, router *mux.Router) {

	// oauth
	router.Handle("/auth/{provider}", &handler.ContextAwareHandler{appCtx, handler.RequestHandlerFunc(api.Authorize)})
	router.Handle("/auth/{provider}/callback", &handler.ContextAwareHandler{appCtx, handler.RequestHandlerFunc(api.Authorized)})

	//site or team
	//router.Handle("/sites").Methods("POST")
	//router.Handle("/sites").Methods("PUT")
	//router.Handle("/sites/{id}").Methods("GET")

}
