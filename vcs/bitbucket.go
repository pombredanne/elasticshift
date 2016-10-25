package vcs

import (
	"fmt"
	"log"

	chttp "gitlab.com/conspico/esh/core/http"
	"golang.org/x/oauth2"
	bb "golang.org/x/oauth2/bitbucket"
)

// Bitbucket URL ...
const (
	BitBucketProviderName   = "bitbucket"
	BitbucketBaseURLV2      = "https://api.bitbucket.org/2.0"
	BitbucketProfileURL     = BitbucketBaseURLV2 + "/user"
	BitbucketGetUserRepoURL = BitbucketBaseURLV2 + "/repositories/:username"
)

// Bitbucket ...
type Bitbucket struct {
	CallbackURL string
	Config      *oauth2.Config
}

// BitbucketProvider ...
// Creates a new Github provider
func BitbucketProvider(clientID, secret, callbackURL string) *Bitbucket {

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Scopes:       []string{"repository"},
		Endpoint:     bb.Endpoint,
	}

	return &Bitbucket{
		callbackURL,
		conf,
	}
}

// Name of the provider
func (b *Bitbucket) Name() string {
	return BitBucketProviderName
}

// Authorize ...
// Provide access to esh app on accessing the github user and repos.
// the elasticshift application to have access to github repo
func (b *Bitbucket) Authorize(team string) string {
	b.Config.RedirectURL = b.CallbackURL + "/" + team
	url := b.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return url
}

// Authorized ...
// Finishes the authorize
func (b *Bitbucket) Authorized(code string) (VCS, error) {

	tok, err := b.Config.Exchange(oauth2.NoContext, code)
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
	}
	u.TokenType = tok.TokenType
	u.Type = BitBucketType

	us := struct {
		Name  string `json:"username"`
		Links struct {
			Avatar struct {
				Href string `json:"href"`
			}
		}
	}{}

	err = chttp.NewGetRequestMaker(BitbucketProfileURL).PathParams().QueryParam("access_token", tok.AccessToken).Scan(&us).Dispatch()
	if err != nil {
		return u, err
	}

	u.AvatarURL = us.Links.Avatar.Href
	u.Name = us.Name
	return u, err
}

// RefreshToken ..
func (b *Bitbucket) RefreshToken(token string) (*oauth2.Token, error) {

	r := chttp.NewGetRequestMaker(bb.Endpoint.TokenURL)

	r.Header("Accept", "application/json")
	r.Header("Content-Type", "application/x-www-form-urlencoded")

	r.QueryParam("client_id", b.Config.ClientID)
	r.QueryParam("client_secret", b.Config.ClientSecret)
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
func (b *Bitbucket) GetRepos(token string, ownerType int) ([]Repo, error) {

	var url string
	r := chttp.NewGetRequestMaker(url)

	r.Header("Accept", "application/json")
	r.Header("Content-Type", "application/x-www-form-urlencoded")

	var repos []Repo
	err := r.Scan(&repos).Dispatch()
	return repos, err
}
