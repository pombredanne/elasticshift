package vcs

import (
	"log"

	chttp "gitlab.com/conspico/esh/core/http"
	"golang.org/x/oauth2"
)

const (
	// ProfileURL ...
	ProfileURL = "https://api.github.com/user"

	name = "github"
)

// Github ...
type Github struct {
	ClientID    string
	Secret      string
	CallbackURL string
	Config      *oauth2.Config
}

// GithubUser ..
type githubUser struct {
	RawData     map[string]interface{}
	Type        int
	AccessToken string
	AvatarURL   string
}

// GithubProvider ...
// Creates a new Github provider
func GithubProvider(clientID, secret, callbackURL string) *Github {

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Scopes:       []string{"user,repo"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	return &Github{
		clientID,
		secret,
		callbackURL,
		conf,
	}
}

// Name of the provider
func (g *Github) Name() string {
	return name
}

// Authorize ...
// Provide access to esh app on accessing the github user and repos.
// the elasticshift application to have access to github repo
func (g *Github) Authorize(team string) string {
	g.Config.RedirectURL = g.CallbackURL + "/" + team
	url := g.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return url
}

// Authorized ...
// Finishes the authorize
func (g *Github) Authorized(code string) (VCS, error) {

	tok, err := g.Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	u := VCS{}
	u.AccessToken = tok.AccessToken
	u.Type = GithubType

	us := struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Login   string `json:"login"`
		Picture string `json:"avatar_url"`
	}{}

	err = chttp.NewGetRequestMaker(ProfileURL).QueryParam("access_token", code).Scan(&us).Dispatch()
	if err != nil {
		return u, err
	}

	u.AvatarURL = us.Picture
	u.Name = us.Login

	return u, err
}

func (g *Github) ListRepos(accessToken, repoType string) {

}
