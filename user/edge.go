package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/conspico/esh/core/edge"
)

// user registration
type createUserRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Team      string `json:"team"`
}

type createUserResponse struct {
	Code string
	Err  error
}

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var user createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return false, err
	}

	// team
	if user.Team == "" {
		return false, errNoTeamIDNotExist
	}
	// validate email
	// validate firstname and lastname
	return user, nil
}

func encodeCreateUserResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {

	resp := r.(createUserResponse)
	if len(resp.Code) > 0 {

		var body struct {
			Code string
		}
		body.Code = resp.Code

		data, err := json.Marshal(body)
		if err != nil {
			return err
		}

		w.Write(data)

		w.WriteHeader(http.StatusCreated)
		return nil
	}
	return resp.Err
}

func makeCreateUserEdge(s Service) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)
		code, err := s.Create(req.Team, req.Firstname, req.Lastname, req.Email)
		return createUserResponse{Code: code, Err: err}, nil
	}
}

// verify and signin
// type verifyAndSignInRequest struct {
// 	Code     string
// 	Password string
// }

type genericResponse struct {
	Valid bool
	Err   error
}

// func decodeVerifyAndSignInRequest(ctx context.Context, r *http.Request) (interface{}, error) {

// 	var verify verifyAndSignInRequest
// 	if err := json.NewDecoder(r.Body).Decode(&verify); err != nil {
// 		return false, err
// 	}

// 	// validate email
// 	// validate firstname and lastname
// 	return verify, nil
// }

// func makeVerifyUserEdge(s Service) edge.Edge {

// 	return func(ctx context.Context, request interface{}) (interface{}, error) {
// 		req := request.(verifyAndSignInRequest)
// 		valid, err := s.Verify(req.Code)
// 		return genericResponse{Valid: valid, Err: err}, nil
// 	}
// }

// verify code
type verifyCodeRequest struct {
	Code string
}

func decodeVerifyCodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	//code := r.FormValue("code")
	code := mux.Vars(r)["code"]
	if len(code) == 0 {
		return false, errVerificationCodeIsEmpty
	}
	return verifyCodeRequest{Code: code}, nil
}

func encodeVerifyCodeRequest(ctx context.Context, w http.ResponseWriter, r interface{}) error {

	resp := r.(genericResponse)
	if resp.Err != nil {
		return resp.Err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func makeVerifyCodeEdge(s Service) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(verifyCodeRequest)
		valid, err := s.Verify(req.Code)
		return genericResponse{Valid: valid, Err: err}, nil
	}
}
