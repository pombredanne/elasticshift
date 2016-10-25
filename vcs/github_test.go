package vcs_test

import (
	"fmt"
	"testing"

	"gitlab.com/conspico/esh/vcs"
)

func TestAuthorized(t *testing.T) {

	key := "2bb421705ee7f6bb0970"
	secret := "ffd145f08ec0ba1261762f29754ab2a9d12544b7"
	callback := "http://eshift:5050/api/auth/github/callback"

	providers := vcs.NewProviders(
		vcs.GithubProvider(key, secret, callback),
	)

	// p := vcs.GithubProvider(key, secret, callback)

	p, _ := providers.Get("github")

	//oDkyNxdDhY3Fwp3dgdfQ
	//c1 := "124285a434e381f66ee2fca9351747e23055bc48"
	c1 := "dda311e27fb9cfb0b2ca"
	//c2 := "4f03ffd21003502b40a45d3569fc13850ac41f35"
	u, err := p.Authorized(c1)
	if err != nil {
		t.Log("Err = ", err)
	}
	fmt.Println(u)
}
