package user

import (
	"context"

	"gitlab.com/conspico/esh/core/edge"
)

func makeSignupEdge(s Service) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(signupRequest)
		token, err := s.Create(req.Team, req.Firstname, req.Lastname, req.Email, req.Password)
		return signInResponse{Token: token, Err: err}, nil
	}
}

func makeVerifyCodeEdge(s Service) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(verifyCodeRequest)
		valid, err := s.Verify(req.Code)
		return genericResponse{Valid: valid, Err: err}, nil
	}
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
