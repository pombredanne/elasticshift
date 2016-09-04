package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/conspico/esh/vcs"
)

// Authorize vcs repository
func Authorize(w http.ResponseWriter, r *http.Request) {

	provider, err := getProvider(r)
	if err != nil {
		fmt.Println(err)
	}

	url := provider.Authorize()

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Autorized ..
// Invoked when authorization finished by oauth app
func Autorized(w http.ResponseWriter, r *http.Request) {

	provider, err := getProvider(r)
	if err != nil {
		fmt.Println(err)
	}

	u, err := provider.Authorized(r.URL)

	// persist user
	fmt.Println(u.AccessToken)
}

// getProvider fetches the provider by name
func getProvider(r *http.Request) (vcs.Provider, error) {

	providerName := mux.Vars(r)["provider"]
	fmt.Println(providerName)
	return vcs.Get(providerName)
}
