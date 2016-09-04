package main

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gitlab.com/conspico/esh/model"
	"gitlab.com/conspico/esh/service"
	"gitlab.com/conspico/esh/vcs"
)

var (
	appCtx *model.AppContext
)

func main() {

	fmt.Println("Starting ESH Server...")
	// CLI args

	// Logger

	// App configuration
	conf := viper.New()
	conf.SetConfigType("yml")
	conf.SetConfigFile("esh.yml")
	conf.ReadInConfig()

	// Unwrap DEK - data encryption key

	// DB Initialization

	// TLS

	// register vcs providers
	fmt.Println("Register the VCS providers..")
	vcs.Use(
		vcs.New(conf.GetString("github.key"), conf.GetString("github.secret")),
	)

	// App
	router := mux.NewRouter()

	// ESH UI pages
	fmt.Println("Map the View directory.")
	router.PathPrefix("/views/").Handler(http.StripPrefix("/views/", http.FileServer(http.Dir("public"))))

	// Router (includes subdomain)
	fmt.Println("Router Registration")
	service.RegisterRoutes(router)

	// API

	// appCtx = &model.AppContext{
	// 	conf,
	// }

	// Start the server
	fmt.Println("ESH Server listening on port 5050")
	fmt.Println(http.ListenAndServe(":5050", router))
}
