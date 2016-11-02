package main

import (
	"net/http/pprof"

	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net/http"

	"context"
	"fmt"

	"github.com/spf13/viper"
	"gitlab.com/conspico/esh/team"
	"gitlab.com/conspico/esh/user"
	"gitlab.com/conspico/esh/vcs"

	"github.com/gorilla/mux"
	"gitlab.com/conspico/esh/datastore"

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

	// set the configurations
	db.SingularTable(true)
	db.DB().SetMaxOpenConns(conf.GetInt("db.max_connections"))
	db.DB().SetMaxIdleConns(conf.GetInt("db.idle_connections"))

	db.LogMode(conf.GetBool("db.log"))
	defer db.Close()

	// TLS

	// register vcs providers
	fmt.Println("Register the VCS providers..")

	// Init datastore
	var (
		teamDS = datastore.NewTeam(db)
		userDS = datastore.NewUser(db)
		vcsDS  = datastore.NewVCS(db)
		//repoDS = datastore.NewRepo(db)
	)

	// load keys
	signer, err := loadKey(conf.GetString("key.signer"))
	if err != nil {
		panic(err)
	}

	verifier, err := loadKey(conf.GetString("key.verifier"))
	if err != nil {
		panic(err)
	}

	// Initialize services
	var ts team.Service
	ts = team.NewService(teamDS)

	var us user.Service
	us = user.NewService(userDS, teamDS, conf, signer)

	var vs vcs.Service
	vs = vcs.NewService(vcsDS, teamDS, conf)

	router := mux.NewRouter()

	// Router (includes subdomain)
	team.MakeRequestHandler(ctx, ts, router)
	user.MakeRequestHandler(ctx, us, router, signer, verifier)
	vcs.MakeRequestHandler(ctx, vs, router, signer, verifier)

	// pprof
	router.HandleFunc("/debug/pprof", pprof.Index)
	router.HandleFunc("/debug/symbol", pprof.Symbol)
	router.HandleFunc("/debug/profile", pprof.Profile)
	router.Handle("/debug/heap", pprof.Handler("heap"))
	router.Handle("/debug/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/block", pprof.Handler("block"))

	// ESH UI pages
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))
	router.Handle("/", accessControl(router))

	// Start the server
	fmt.Println("ESH Server listening on port 5050")
	fmt.Println(http.ListenAndServe(":5050", router))
}

func loadKey(path string) (interface{}, error) {

	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyBytes)

	switch keyBlock.Type {
	case "PUBLIC KEY":
		return x509.ParsePKIXPublicKey(keyBlock.Bytes)
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	default:
		return nil, fmt.Errorf("unsupported key type %q", keyBlock.Type)
	}
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
