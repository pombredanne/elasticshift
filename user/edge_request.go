package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// user registration
type signupRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Team      string `json:"team"`
}

func decodeSignupRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var user signupRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return false, err
	}

	// team
	if len(user.Team) == 0 {
		return false, errNoTeamIDNotExist
	}
	// validate email
	// validate firstname and lastname
	// validate password
	return user, nil
}

// verify code
type verifyCodeRequest struct {
	Code string
}

func decodeVerifyCodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	//code := r.FormValue("code")
	code := mux.Vars(r)["code"]
	fmt.Println("code = ", code)
	if len(code) == 0 {
		return false, errVerificationCodeIsEmpty
	}
	return verifyCodeRequest{Code: code}, nil
}