package vcs

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/oauth2"
)

const (
	// ProfileURL ...
	ProfileURL = "https://api.github.com/user"

	name = "github"
)

// Github ...
type Github struct {
	ClientID string
	Secret   string
	Config   *oauth2.Config
}

// GithubProvider ...
// Creates a new Github provider
func GithubProvider(clientID, secret string) *Github {

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
func (g *Github) Authorize() string {

	url := g.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return url
}

// Authorized ...
// Finishes the authorize
func (g *Github) Authorized(code string) (User, error) {

	tok, err := g.Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	user := User{
		AccessToken: tok.AccessToken,
		Provider:    name,
	}

	res, err := http.Get(ProfileURL + "?access_token" + url.QueryEscape(tok.AccessToken))
	if err != nil {
		if res != nil {
			res.Body.Close()
		}
		return user, err
	}
	defer res.Body.Close()

	bits, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return user, err
	}

	err = json.NewDecoder(bytes.NewReader(bits)).Decode(&user.RawData)
	if err != nil {
		return user, err
	}

	err = readUser(bytes.NewReader(bits), &user)
	return user, err
}

// Helper method to convert reader to user
func readUser(reader io.Reader, user *User) error {

	u := struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Bio      string `json:"bio"`
		Name     string `json:"name"`
		Login    string `json:"login"`
		Picture  string `json:"avatar_url"`
		Location string `json:"location"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Name = u.Name
	user.NickName = u.Login
	user.Email = u.Email
	user.Description = u.Bio
	user.AvatarURL = u.Picture
	user.UserID = strconv.Itoa(u.ID)
	user.Location = u.Location

	return err
}
