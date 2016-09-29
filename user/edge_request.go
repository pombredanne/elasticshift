package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// user registration
type signupRequest struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Team     string `json:"team"`
	Domain   string
}

type signInRequest struct {
	Team     string `json:"team"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func decodeSignUpRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var user signupRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return false, err
	}

	subdomain := ctx.Value("team").(string)
	if len(subdomain) >= 6 {
		user.Domain = subdomain
	}

	// validate email
	// validate firstname and lastname
	// validate password
	return user, nil
}

func decodeSignInRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var req signInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return false, err
	}

	team := ctx.Value("team").(string)
	if team == "" {
		req.Team = team
	}

	// team
	if req.Team == "" {
		return false, errNoTeamIDNotExist
	}

	// validate username and password
	if len(req.Email) == 0 || len(req.Password) == 0 {
		return false, errUsernameOrPasswordIsEmpty
	}

	return req, nil
}

// verify code
type verifyCodeRequest struct {
	Code string
}

func decodeVerifyCodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	code := mux.Vars(r)["code"]
	if len(code) == 0 {
		return false, errVerificationCodeIsEmpty
	}
	return verifyCodeRequest{Code: code}, nil
}
