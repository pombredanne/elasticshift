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

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	ctx := esh.AppContext{}

	// Unwrap DEK - data encryption key
	ctx.Context = context.Background()

	// Logger
	logger := logrus.New()
	ctx.Logger = logger

	logger.Infoln("Starting ESH Server...")
	// CLI args

	// App configuration
	logger.Infoln("Reading properties...")
	vip := viper.New()
	vip.SetConfigType("yml")
	vip.SetConfigFile("esh.yml")
	vip.ReadInConfig()

	config := esh.Config{}
	vip.Unmarshal(&config)

	ctx.Config = config

	// DB Initialization
	logger.Infoln("Opening DB Connection...")
	db, err := gorm.Open(config.DB.Dialect, config.DB.Datasource)
	if err != nil {
		logger.Fatalln("Cannot initialize database.", err)
		os.Exit(-1)
	}

	// Ping function checks the database connectivity
	dberr := db.DB().Ping()
	if dberr != nil {
		logger.Fatalln("Cannot initialize database.", err)
		os.Exit(-1)
	}

	// set the configurations
	db.SingularTable(true)
	db.DB().SetMaxOpenConns(config.DB.MaxConnection)
	db.DB().SetMaxIdleConns(config.DB.IdleConnection)

	db.LogMode(config.DB.Log)
	defer db.Close()

	// TLS

	// Init datastore
	ctx.TeamDatastore = esh.NewTeamDatastore(db)
	ctx.UserDatastore = esh.NewUserDatastore(db)
	ctx.VCSDatastore = esh.NewVCSDatastore(db)
	ctx.RepoDatastore = esh.NewRepoDatastore(db)

	// load keys
	signer, err := loadKey(config.Key.Signer)
	if err != nil {
		logger.Fatalln("Cannot load signer key", err)
		os.Exit(-1)
	}
	ctx.Signer = signer

	verifier, err := loadKey(config.Key.Verifier)
	if err != nil {
		logger.Fatalln("Cannot load verifier key", err)
		os.Exit(-1)
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
	csrf.Secure(config.CSRF.Secure)
	csrfHandler := csrf.Protect([]byte(config.CSRF.Key))
	router.PathPrefix("/").Handler(ctx.PublicChain.Then(http.FileServer(http.Dir("./dist/"))))
	router.Handle("/", ctx.PublicChain.Then(csrfHandler(router)))

	// Start the server
	logger.Info("ESH Server listening on port 5050")
	logger.Info(http.ListenAndServe(":5050", router))
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
