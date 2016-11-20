package esh

import (
	"github.com/Sirupsen/logrus"

	"net/url"

	"time"

	chttp "gitlab.com/conspico/esh/core/http"

	"golang.org/x/oauth2"

	"github.com/palantir/stacktrace"
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
	logger      *logrus.Logger
}

// BitbucketProvider ...
// Creates a new Github provider
func BitbucketProvider(logger *logrus.Logger, clientID, secret, callbackURL string) *Bitbucket {

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Scopes:       []string{"repository"},
		Endpoint:     bb.Endpoint,
	}

	return &Bitbucket{
		callbackURL,
		conf,
		logger,
	}
}

// Name of the provider
func (b *Bitbucket) Name() string {
	return BitBucketProviderName
}

// Authorize ...
// Provide access to esh app on accessing the github user and repos.
// the elasticshift application to have access to github repo
func (b *Bitbucket) Authorize(baseURL string) string {
	b.Config.RedirectURL = b.CallbackURL + "/" + baseURL
	url := b.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return url
}

// Authorized ...
// Finishes the authorize
func (b *Bitbucket) Authorized(code string) (VCS, error) {

	tok, err := b.Config.Exchange(oauth2.NoContext, code)
	u := VCS{}
	if err != nil {
		return u, stacktrace.Propagate(err, "Exchange token after bitbucket auth failed")
	}

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

	r := chttp.NewGetRequestMaker(BitbucketProfileURL)
	r.SetLogger(b.logger)

	r.PathParams()
	r.QueryParam("access_token", tok.AccessToken)

	err = r.Scan(&us).Dispatch()
	if err != nil {
		return u, err
	}

	u.AvatarURL = us.Links.Avatar.Href
	u.Name = us.Name
	return u, err
}

// RefreshToken ..
func (b *Bitbucket) RefreshToken(token string) (*oauth2.Token, error) {

	r := chttp.NewPostRequestMaker(b.Config.Endpoint.TokenURL)

	r.SetBasicAuth(b.Config.ClientID, b.Config.ClientSecret)

	r.Header("Accept", "application/json")
	r.SetContentType(chttp.URLENCODED)

	// body := struct {
	// 	GrantType    string `json:"grant_type"`
	// 	RefreshToken string `json:"refresh_token"`
	// }{}

	// body.GrantType = "refresh_token"
	// body.RefreshToken = token

	params := make(url.Values)
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", token)

	r.Body(params)

	var tok Token
	err := r.Scan(&tok).Dispatch()

	if err != nil {
		return nil, err
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
func (b *Bitbucket) GetRepos(token, accountName string, ownerType int) ([]Repo, error) {

	r := chttp.NewGetRequestMaker(BitbucketGetUserRepoURL)
	r.SetLogger(b.logger)

	r.Header("Accept", "application/json")
	r.Header("Content-Type", "application/x-www-form-urlencoded")

	r.PathParams(accountName)

	r.QueryParam("access_token", token)

	rp := struct {
		Values []struct {
			Name     string
			UUID     string
			Language string
			Links    struct {
				HTML struct {
					Href string
				}
				Avatar struct {
					Href string
				}
			}
			Owner struct {
				Type string
			}
			Description string
			Private     bool `json:"is_private"`
		}
	}{}
	err := r.Scan(&rp).Dispatch()

	repos := []Repo{}
	for _, rpo := range rp.Values {

		repo := &Repo{
			RepoID:      rpo.UUID,
			Name:        rpo.Name,
			Language:    rpo.Language,
			Link:        rpo.Links.HTML.Href,
			Description: rpo.Description,
		}
		if rpo.Private {
			repo.Private = True
		}

		repos = append(repos, *repo)
	}
	return repos, err
}
