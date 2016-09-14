package main

import (
	"net/http"

	"context"
	"fmt"

	"github.com/spf13/viper"
	"gitlab.com/conspico/esh/team"
	"gitlab.com/conspico/esh/vcs"

	"gitlab.com/conspico/esh/repository"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

	ctx := context.Background()

	// DB Initialization
	db, err := gorm.Open(conf.GetString("db.dialect"), conf.GetString("db.datasource"))
	if err != nil {
		fmt.Println("Cannot initialize database.", err)
		panic("Cannot connect DB")
	}

	// Ping function checks the database connectivity
	dberr := db.DB().Ping()
	if dberr != nil {
		panic(err)
	}

	db.DB().SetMaxOpenConns(conf.GetInt("db.max_connections"))
	db.DB().SetMaxIdleConns(conf.GetInt("db.idle_connections"))

	db.LogMode(true)
	defer db.Close()

	// TLS

	// register vcs providers
	fmt.Println("Register the VCS providers..")
	vcs.Use(
		vcs.New(conf.GetString("github.key"), conf.GetString("github.secret")),
	)

	// ESH UI pages
	fmt.Println("Map the View directory..")

	var (
		teamRepo = repository.NewTeam(db)
	)

	//team
	var ts team.Service
	ts = team.NewService(teamRepo)

	router := http.NewServeMux()
	//router.PathPrefix("/views/").Handler(http.StripPrefix("/views/", http.FileServer(http.Dir("public"))))

	router.Handle("/teams/v1", accessControl(team.MakeRequestHandler(ctx, ts)))
	router.Handle("/", accessControl(router))
	// Router (includes subdomain)
	fmt.Println("Router Registration")
	//service.RegisterRoutes(appCtx, router)

	// Start the server
	fmt.Println("ESH Server listening on port 5050")
	fmt.Println(http.ListenAndServe(":5050", router))
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
