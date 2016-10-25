package vcs_test

import (
	"encoding/json"
	"testing"
)

func TestBitbucketAuthorized(t *testing.T) {

	/*key := "4EuzbzNEwa7x3SP5yZ"
	secret := "dqHP8bc5qyzwTXuW2dduSSTBJqCa8jzj"
	callback := "http://eshift:5050/api/auth/bitbucket/callback"

	providers := vcs.NewProviders(
		vcs.BitbucketProvider(key, secret, callback),
	)

	// p := vcs.GithubProvider(key, secret, callback)

	p, _ := providers.Get("bitbucket")

	//oDkyNxdDhY3Fwp3dgdfQ
	//c1 := "124285a434e381f66ee2fca9351747e23055bc48"
	//c1 := "bTkhahxSnzseGnG6kY"
	//c2 := "4f03ffd21003502b40a45d3569fc13850ac41f35"*/
	//token := "fKryBJ6xyOAiLqy7K7CEBei4tpasgdcvkyqDkiN-CsGIUt_VdcBO2zthKz-30BTI02wd_swGpUaRUdqC4bQ="

	/*u, err := p.Authorized(c1)
	if err != nil {
		t.Log("Err = ", err)
	}
	fmt.Println(u)
	*/

	us := struct {
		Name  string `json:"username"`
		Links struct {
			Avatar struct {
				Href string
			}
			Repositories struct {
				Href string
			}
		}
	}{}

	result := `
        {"username": "ghazninattarshah", "website": "", "display_name": "Ghazni Nattarshah", "uuid":
"{56c96236-0ae7-47c5-8e2a-1c90efd5b22a}", "links": {"hooks": {"href":
"https://api.bitbucket.org/2.0/users/ghazninattarshah/hooks"}, "self": {"href":
"https://api.bitbucket.org/2.0/users/ghazninattarshah"}, "repositories": {"href":
"https://api.bitbucket.org/2.0/repositories/ghazninattarshah"}, "html": {"href":
"https://bitbucket.org/ghazninattarshah/"}, "followers": {"href":
"https://api.bitbucket.org/2.0/users/ghazninattarshah/followers"}, "avatar": {"href":
"https://bitbucket.org/account/ghazninattarshah/avatar/32/"}, "following": {"href":
"https://api.bitbucket.org/2.0/users/ghazninattarshah/following"}, "snippets": {"href":
"https://api.bitbucket.org/2.0/snippets/ghazninattarshah"}}, "created_on":
"2015-06-20T04:27:38.605697+00:00", "location": null, "type": "user"}
    `

	err := json.Unmarshal([]byte(result), &us)
	if err != nil {
		t.Log(err)
	}

	t.Log("Name = ", us.Name)
	t.Log("Avatar = ", us.Links.Avatar.Href)
	t.Log("repo = ", us.Links.Repositories.Href)

	//err := chttp.NewGetRequestMaker(vcs.BitbucketProfileURL).PathParams().QueryParam("access_token", token).Scan(&us).Dispatch()
	//if err != nil {
	//	t.Log(err)
	//}
}
