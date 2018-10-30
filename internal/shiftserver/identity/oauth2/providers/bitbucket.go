package providers

import (
	"github.com/sirupsen/logrus"

	"net/url"

	"time"

	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
	"gitlab.com/conspico/elasticshift/pkg/dispatch"

	"golang.org/x/oauth2"

	bb "golang.org/x/oauth2/bitbucket"
)

// Bitbucket URL ...
const (
	BitbucketProviderName   = "bitbucket"
	BitbucketBaseURLV2      = "https://api.bitbucket.org/2.0"
	BitbucketProfileURL     = BitbucketBaseURLV2 + "/user"
	BitbucketGetUserRepoURL = BitbucketBaseURLV2 + "/repositories/:username"
)

// Bitbucket ...
type Bitbucket struct {
	CallbackURL string
	HookURL     string
	Config      *oauth2.Config
	logger      *logrus.Entry
}

// BitbucketProvider ...
// Creates a new Github provider
func BitbucketProvider(loggr logger.Loggr, clientID, secret, callbackURL, HookURL string) *Bitbucket {

	l := loggr.GetLogger("oauth2/bitbucket")

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Scopes:       []string{"repository"},
		Endpoint:     bb.Endpoint,
	}

	return &Bitbucket{
		callbackURL,
		HookURL,
		conf,
		l,
	}
}

// Name of the provider
func (b *Bitbucket) Name() string {
	return BitbucketProviderName
}

// Authorize ...
// Provide access to esh app on accessing the github user and repos.
// the elasticshift application to have access to github repo
func (b *Bitbucket) Authorize(baseURL string) string {

	opts := oauth2.SetAuthURLParam("redirect_uri", b.CallbackURL+"?id="+baseURL)
	url := b.Config.AuthCodeURL("state", oauth2.AccessTypeOffline, opts)
	return url
}

// Authorized ...
// Finishes the authorize
func (b *Bitbucket) Authorized(id, code string) (types.VCS, error) {

	tok, err := b.Config.Exchange(oauth2.NoContext, code)
	u := types.VCS{}
	// if err != nil {
	// 	return u, stacktrace.Propagate(err, "Exchange token after bitbucket auth failed")
	// }

	u.AccessCode = code
	u.RefreshToken = tok.RefreshToken
	u.AccessToken = tok.AccessToken
	if !tok.Expiry.IsZero() { // zero never expires
		u.TokenExpiry = tok.Expiry
	}
	u.Kind = b.Name()

	us := struct {
		UUID  string
		Name  string `json:"username"`
		Links struct {
			Avatar struct {
				Href string `json:"href"`
			}

			Html struct {
				Href string `json:"href"`
			}
		}
	}{}

	r := dispatch.NewGetRequestMaker(BitbucketProfileURL)
	r.SetLogger(b.logger)

	r.PathParams()
	r.QueryParam("access_token", tok.AccessToken)

	err = r.Scan(&us).Dispatch()
	if err != nil {
		return u, err
	}

	u.AvatarURL = us.Links.Avatar.Href
	u.Name = us.Name
	u.Link = us.Links.Html.Href
	u.ID = us.UUID
	return u, err
}

// RefreshToken ..
func (b *Bitbucket) RefreshToken(token string) (*oauth2.Token, error) {

	r := dispatch.NewPostRequestMaker(b.Config.Endpoint.TokenURL)

	r.SetBasicAuth(b.Config.ClientID, b.Config.ClientSecret)

	r.Header("Accept", "application/json")
	r.SetContentType(dispatch.URLENCODED)

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
func (b *Bitbucket) GetRepos(token, accountName string, ownerType string) ([]types.Repository, error) {

	r := dispatch.NewGetRequestMaker(BitbucketGetUserRepoURL)
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

	repos := []types.Repository{}
	for _, rpo := range rp.Values {

		repo := &types.Repository{
			RepoID:      rpo.UUID,
			Name:        rpo.Name,
			Language:    rpo.Language,
			Link:        rpo.Links.HTML.Href,
			Description: rpo.Description,
		}
		if rpo.Private {
			repo.Private = true
		}

		repos = append(repos, *repo)
	}
	return repos, err
}

func (g *Bitbucket) Search(token, vcsName, repoName string) (types.Repository, error) {

	repo := types.Repository{}
	return repo, nil
}

// CreateHook ..
// Create a new hook
func (b *Bitbucket) CreateHook(token, owner, repo string) error {

	r := dispatch.NewPostRequestMaker(GithubCreateHookURL)
	r.SetLogger(b.logger)

	r.SetContentType(dispatch.JSON)

	return nil
}
