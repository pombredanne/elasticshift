/*
Copyright 2017 The Elasticshift Authors.
*/
package providers

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"net/url"

	"time"

	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/pkg/dispatch"

	"golang.org/x/oauth2"

	"strconv"
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
	HookURL     string
	Config      *oauth2.Config
	logger      *logrus.Entry
}

// GitlabProvider ...
// Creates a new Gitlab provider
func GitlabProvider(loggr logger.Loggr, clientID, secret, callbackURL, hookURL string) *Gitlab {

	l := loggr.GetLogger("oauth2/gitlab")

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
		hookURL,
		conf,
		l,
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
	opts := oauth2.SetAuthURLParam("redirect_uri", g.CallbackURL+"?id="+baseURL)
	url := g.Config.AuthCodeURL("state", oauth2.AccessTypeOffline, opts)
	g.logger.Println(url)
	return url
}

// Authorized ...
// Finishes the authorize
func (g *Gitlab) Authorized(id, code string) (types.VCS, error) {

	//tok, err := g.Config.Exchange(oauth2.NoContext, code)
	// Authorize request
	r := dispatch.NewPostRequestMaker(GitlabTokenURL)
	r.SetLogger(g.logger)
	r.SetContentType(dispatch.URLENCODED)

	r.Header("Accept", dispatch.JSON)

	params := make(url.Values)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("redirect_uri", g.Config.RedirectURL)

	r.QueryParam("client_id", g.Config.ClientID)
	r.QueryParam("client_secret", g.Config.ClientSecret)

	r.Body(params)

	var tok Token
	err := r.Scan(&tok).Dispatch()

	u := types.VCS{}
	if err != nil {
		return u, fmt.Errorf("Exchange token failed: ", err)
	}

	u.AccessCode = code
	u.RefreshToken = tok.RefreshToken
	u.AccessToken = tok.AccessToken
	u.TokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	u.Kind = g.Name()

	g.logger.Warn("Token = ", tok)
	// Get user profile
	us := struct {
		ID        int    `json:"id"`
		Name      string `json:"username"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Link      string `json:"web_url"`
	}{}

	r = dispatch.NewGetRequestMaker(GitlabProfileURL)
	r.SetLogger(g.logger)

	r.PathParams()
	r.QueryParam("access_token", tok.AccessToken)

	err = r.Scan(&us).Dispatch()
	if err != nil {
		return u, err
	}

	u.AvatarURL = us.AvatarURL
	u.Name = us.Name
	u.Link = us.Link
	u.ID = strconv.Itoa(us.ID)
	return u, err
}

// RefreshToken ..
func (g *Gitlab) RefreshToken(token string) (*oauth2.Token, error) {

	r := dispatch.NewPostRequestMaker(GitlabTokenURL)
	r.SetLogger(g.logger)

	r.SetBasicAuth(g.Config.ClientID, g.Config.ClientSecret)

	r.Header("Accept", "application/json")
	r.SetContentType(dispatch.URLENCODED)

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
func (g *Gitlab) GetRepos(token, accountName string, ownerType string) ([]types.Repository, error) {

	r := dispatch.NewGetRequestMaker(GitlabGetUserRepoURL)
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

	repos := []types.Repository{}
	for _, rpo := range rp {

		repo := &types.Repository{
			RepoID:        strconv.Itoa(rpo.ID),
			Name:          rpo.Name,
			Link:          rpo.WebURL,
			Description:   rpo.Description,
			DefaultBranch: rpo.DefaultBranch,
		}
		if rpo.Public {
			repo.Private = false
		}

		repos = append(repos, *repo)
	}
	return repos, err
}

func (g *Gitlab) Search(token, vcsName, repoName string) (types.Repository, error) {

	repo := types.Repository{}
	return repo, nil
}

// CreateHook ..
// Create a new hook
func (g *Gitlab) CreateHook(token, owner, repo string) error {

	r := dispatch.NewPostRequestMaker(GithubCreateHookURL)
	r.SetLogger(g.logger)

	r.SetContentType(dispatch.JSON)

	return nil
}
