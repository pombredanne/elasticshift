package esh_test

import (
	"testing"

	"gitlab.com/conspico/esh"
)

func TestBitbucketAuthorized(t *testing.T) {

	key := "4EuzbzNEwa7x3SP5yZ"
	secret := "dqHP8bc5qyzwTXuW2dduSSTBJqCa8jzj"
	callback := "http://eshift:5050/api/auth/bitbucket/callback"

	providers := esh.NewProviders(
		esh.BitbucketProvider(key, secret, callback),
	)

	// p := vcs.GithubProvider(key, secret, callback)

	p, _ := providers.Get("bitbucket")

	tok, err := p.RefreshToken("Xx9FLeba9quF9cmn2R")

	t.Log("Err = ", err)
	t.Log("Token = ", tok)
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

	//err := chttp.NewGetRequestMaker(vcs.BitbucketProfileURL).PathParams().QueryParam("access_token", token).Scan(&us).Dispatch()
	//if err != nil {
	//	t.Log(err)
	//}
}
