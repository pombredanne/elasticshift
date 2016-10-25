package vcs

import (
	"fmt"
	"log"
	"time"

	chttp "gitlab.com/conspico/esh/core/http"
	"golang.org/x/oauth2"
	gh "golang.org/x/oauth2/github"
)

// Github related properties
const (
	GithubBaseURL        = "https://api.github.com"
	GithubProfileURL     = GithubBaseURL + "/user"
	GithubGetUserRepoURL = GithubBaseURL + "/users/:user/repos"
	GithubGetOrgRepoURL  = GithubBaseURL + "/orgs/:org/repos"
	GithubProviderName   = "github"
)

// Github ...
type Github struct {
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
		Endpoint:     gh.Endpoint,
	}

	return &Github{
		callbackURL,
		conf,
	}
}

// Name of the provider
func (g *Github) Name() string {
	return GithubProviderName
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

	fmt.Println("Extracted token = ", tok)
	u := VCS{}
	u.AccessCode = code
	u.RefreshToken = tok.RefreshToken
	u.AccessToken = tok.AccessToken
	if !tok.Expiry.IsZero() { // zero never expires
		u.TokenExpiry = tok.Expiry
	} else {
		u.TokenExpiry = time.Now()
	}
	u.TokenType = tok.TokenType
	u.Type = GithubType

	us := struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Login   string `json:"login"`
		Picture string `json:"avatar_url"`
	}{}

	err = chttp.NewGetRequestMaker(GithubProfileURL).Header("Accept", "application/json").QueryParam("access_token", tok.AccessToken).Scan(&us).Dispatch()
	if err != nil {
		return u, err
	}

	u.AvatarURL = us.Picture
	u.Name = us.Login
	return u, err
}

// RefreshToken ..
func (g *Github) RefreshToken(token string) (*oauth2.Token, error) {

	r := chttp.NewGetRequestMaker(gh.Endpoint.TokenURL)

	r.Header("Accept", "application/json")
	r.Header("Content-Type", "application/x-www-form-urlencoded")

	r.QueryParam("client_id", g.Config.ClientID)
	r.QueryParam("client_secret", g.Config.ClientSecret)
	r.QueryParam("grant_type", "refresh_token")
	r.QueryParam("refresh_token", token)

	var tok oauth2.Token
	err := r.Scan(&tok).Dispatch()

	if err != nil {
		return nil, err
	}
	return &tok, nil
}

// GetRepos ..
// returns the list of repositories
func (g *Github) GetRepos(token string, ownerType int) ([]Repo, error) {

	var url string
	if OwnerTypeUser == ownerType {
		url = GithubGetUserRepoURL
	} else if OwnerTypeOrg == ownerType {
		url = GithubGetUserRepoURL
	}

	r := chttp.NewGetRequestMaker(url)

	r.Header("Accept", "application/json")
	r.Header("Content-Type", "application/x-www-form-urlencoded")

	var rpo []struct {
		RepoID        string
		Name          string
		Private       string
		Link          string `json:"html_url"`
		Description   string
		Fork          string
		DefaultBranch string
		Language      string
	}
	err := r.Scan(&rpo).Dispatch()
	if err != nil {

	}

	var repos []Repo
	for _, repo := range rpo {

		rp := &Repo{
			RepoID:        repo.RepoID,
			Name:          repo.Name,
			Private:       repo.Private,
			Link:          repo.Link,
			Description:   repo.Description,
			Fork:          repo.Fork,
			DefaultBranch: repo.DefaultBranch,
			Language:      repo.Language,
		}
		repos = append(repos, *rp)
	}
	return repos, err
}
