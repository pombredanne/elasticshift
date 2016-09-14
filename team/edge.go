package team

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gitlab.com/conspico/esh/core/edge"
)

var (
	errBadRoute            = errors.New("bad route")
	errDomainNameIsEmpty   = errors.New("Team name is empty")
	errDomainNameMinLength = errors.New("Team name should be atleast 6 chars")
	errDomainNameMaxLength = errors.New("Team name should not exceed 63 chars")
)

// create team
type createTeamRequest struct {
	Name string
}

type createTeamResponse struct {
	Created bool
	Err     error
	status  int
}

func decodeCreateTeamRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var body struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false, err
	}

	nameLength := len(body.Name)
	if nameLength == 0 {
		return false, errDomainNameIsEmpty
	}

	if nameLength < 6 {
		return false, errDomainNameMinLength
	}

	if nameLength > 63 {
		return false, errDomainNameMaxLength
	}
	return createTeamRequest{Name: body.Name}, nil
}

func encodeCreateTeamResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	resp := response.(createTeamResponse)
	if resp.Created {
		w.WriteHeader(http.StatusCreated)
		return nil
	}
	return resp.Err
}

func makeCreateTeamEdge(s Service) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTeamRequest)
		fmt.Println("domain name", req.Name)
		created, err := s.Create(req.Name)
		return createTeamResponse{Created: created, Err: err}, nil
	}
}
