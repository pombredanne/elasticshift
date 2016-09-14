package user

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.com/conspico/esh/core/edge"
)

type createUserRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Team      string `json:"team"`
}

type createUserResponse struct {
	Created bool
	Err     error
}

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var user createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return false, err
	}

	// validate email
	// validate firstname and lastname
	return user, nil
}

func encodeCreateUserResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {

	resp := r.(createUserResponse)
	if resp.Created {
		w.WriteHeader(http.StatusCreated)
		return nil
	}
	return resp.Err
}

func makeCreateUserEdge(s Service) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)
		created, err := s.Create(req.Team, req.Firstname, req.Lastname, req.Email)
		return createUserResponse{Created: created, Err: err}, nil
	}
}
