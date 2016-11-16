package main

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net/http"

	"context"
	"fmt"

	"github.com/spf13/viper"
	"gitlab.com/conspico/esh"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	fmt.Println("Starting ESH Server...")
	// CLI args

	// Logger

	// App configuration
	vip := viper.New()
	vip.SetConfigType("yml")
	vip.SetConfigFile("esh.yml")
	vip.ReadInConfig()

	config := esh.Config{}
	vip.Unmarshal(&config)

	ctx := esh.AppContext{}
	ctx.Config = config

	// Unwrap DEK - data encryption key
	ctx.Context = context.Background()

	// DB Initialization
	db, err := gorm.Open(config.DB.Dialect, config.DB.Datasource)
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
	db.DB().SetMaxOpenConns(config.DB.MaxConnection)
	db.DB().SetMaxIdleConns(config.DB.IdleConnection)

	db.LogMode(config.DB.Log)
	defer db.Close()

	// TLS

	// register vcs providers
	fmt.Println("Register the VCS providers..")

	// Init datastore
	ctx.TeamDatastore = esh.NewTeamDatastore(db)
	ctx.UserDatastore = esh.NewUserDatastore(db)
	ctx.VCSDatastore = esh.NewVCSDatastore(db)
	ctx.RepoDatastore = esh.NewRepoDatastore(db)

	// load keys
	signer, err := loadKey(config.Key.Signer)
	if err != nil {
		panic(err)
	}
	ctx.Signer = signer

	verifier, err := loadKey(config.Key.Verifier)
	if err != nil {
		panic(err)
	}
	ctx.Verifier = verifier

	// Initialize services
	ctx.TeamService = esh.NewTeamService(ctx)
	ctx.UserService = esh.NewUserService(ctx)
	ctx.VCSService = esh.NewVCSService(ctx)
	ctx.RepoService = esh.NewRepoService(ctx)

	router := mux.NewRouter()
	ctx.Router = router

	// Router (includes subdomain)
	esh.MakeHandlers(ctx)

	// ESH UI pages
	router.PathPrefix("/").Handler(ctx.PublicChain.Then(http.FileServer(http.Dir("./dist/"))))
	router.Handle("/", ctx.PublicChain.Then(router))

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
