// Package esh ...
// Author: Ghazni Nattarshah
// Date: Nov 22, 2016
package esh

import (
	"github.com/Sirupsen/logrus"

	"net/url"

	"time"

	chttp "gitlab.com/conspico/esh/core/http"

	"golang.org/x/oauth2"

	"strconv"

	"github.com/palantir/stacktrace"
)

// Gitlab URL ...
const (
	GitlabProviderName = "gitlab"
	GitlabAuthURL      = "https://gitlab.com/oauth/authorize"
	GitlabTokenURL     = "https://gitlab.com/oauth/token"

	GitlabBaseURLV3      = "https://gitlab.com/api/v3"
	GitlabProfileURL     = GitlabBaseURLV3 + "/user"
	GitlabGetUserRepoURL = GitlabBaseURLV3 + "/projects"
)

// Gitlab ...
type Gitlab struct {
	CallbackURL string
	Config      *oauth2.Config
	logger      *logrus.Logger
}

// GitlabProvider ...
// Creates a new Gitlab provider
func GitlabProvider(logger *logrus.Logger, clientID, secret, callbackURL string) *Gitlab {

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  GitlabAuthURL,
			TokenURL: GitlabTokenURL,
		},
	}

	return &Gitlab{
		callbackURL,
		conf,
		logger,
	}
}

// Name of the provider
func (g *Gitlab) Name() string {
	return GitlabProviderName
}

// Authorize ...
// Provide access to esh app on accessing the github user and repos.
// the elasticshift application to have access to github repo
func (g *Gitlab) Authorize(baseURL string) string {
	g.Config.RedirectURL = g.CallbackURL + "?id=" + baseURL
	url := g.Config.AuthCodeURL("state")
	g.logger.Println(url)
	return url
}

// Authorized ...
// Finishes the authorize
func (g *Gitlab) Authorized(code string) (VCS, error) {

	//tok, err := g.Config.Exchange(oauth2.NoContext, code)
	// Authorize request
	r := chttp.NewPostRequestMaker(GitlabTokenURL)
	r.SetLogger(g.logger)
	r.SetContentType(chttp.URLENCODED)

	r.Header("Accept", chttp.JSON)

	params := make(url.Values)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("redirect_uri", g.Config.RedirectURL)

	r.QueryParam("client_id", g.Config.ClientID)
	r.QueryParam("client_secret", g.Config.ClientSecret)

	r.Body(params)

	var tok Token
	err := r.Scan(&tok).Dispatch()

	u := VCS{}
	if err != nil {
		return u, stacktrace.Propagate(err, "Exchange token failed")
	}

	u.AccessCode = code
	u.RefreshToken = tok.RefreshToken
	u.AccessToken = tok.AccessToken
	u.TokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	u.TokenType = tok.TokenType
	u.Type = GitlabType

	g.logger.Warn("Token = ", tok)
	// Get user profile
	us := struct {
		ID        int    `json:"id"`
		Name      string `json:"username"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}{}

	r = chttp.NewGetRequestMaker(GitlabProfileURL)
	r.SetLogger(g.logger)

	r.PathParams()
	r.QueryParam("access_token", tok.AccessToken)

	err = r.Scan(&us).Dispatch()
	if err != nil {
		return u, err
	}

	u.AvatarURL = us.AvatarURL
	u.Name = us.Name
	u.VcsID = strconv.Itoa(us.ID)
	return u, err
}

// RefreshToken ..
func (g *Gitlab) RefreshToken(token string) (*oauth2.Token, error) {

	r := chttp.NewPostRequestMaker(GitlabTokenURL)
	r.SetLogger(g.logger)

	r.SetBasicAuth(g.Config.ClientID, g.Config.ClientSecret)

	r.Header("Accept", "application/json")
	r.SetContentType(chttp.URLENCODED)

	params := make(url.Values)
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", token)
	params.Set("scope", "api")

	r.Body(params)

	var tok Token
	err := r.Scan(&tok).Dispatch()

	if err != nil {
		return nil, err
	}

	g.logger.Infoln("Token Created at ", tok.CreatedAt)

	if tok.ExpiresIn == 0 {
		tok.ExpiresIn = 7200
	}

	otok := &oauth2.Token{
		AccessToken:  tok.AccessToken,
		Expiry:       time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second),
		RefreshToken: tok.RefreshToken,
		TokenType:    tok.TokenType,
	}

	return otok, nil
}

// GetRepos ..
// returns the list of repositories
func (g *Gitlab) GetRepos(token, accountName string, ownerType int) ([]Repo, error) {

	r := chttp.NewGetRequestMaker(GitlabGetUserRepoURL)
	r.SetLogger(g.logger)

	r.Header("Accept", "application/json")
	r.Header("Content-Type", "application/x-www-form-urlencoded")

	//r.PathParams(accountName)

	r.QueryParam("access_token", token)

	rp := []struct {
		ID            int    `json:"id,omitempty"`
		Name          string `json:"name,omitempty"`
		Description   string `json:"description,omitempty"`
		DefaultBranch string `json:"default_branch,omitempty"`
		Public        bool   `json:"public,omitempty"`
		WebURL        string `json:"web_url"`
		AvatarURL     string `json:"avatar_url"`
	}{}
	err := r.Scan(&rp).Dispatch()

	repos := []Repo{}
	for _, rpo := range rp {

		repo := &Repo{
			RepoID:        strconv.Itoa(rpo.ID),
			Name:          rpo.Name,
			Link:          rpo.WebURL,
			Description:   rpo.Description,
			DefaultBranch: rpo.DefaultBranch,
		}
		if rpo.Public {
			repo.Private = False
		}

		repos = append(repos, *repo)
	}
	return repos, err
}
