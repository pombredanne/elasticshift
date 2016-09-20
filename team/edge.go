package team

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.com/conspico/esh/core/edge"
	"gitlab.com/conspico/esh/core/util"
)

// create team
type createTeamRequest struct {
	Name string `json:"name"`
}

type createTeamResponse struct {
	Created bool
	Code    string
	Err     error
}

func decodeCreateTeamRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var team createTeamRequest

	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		return false, err
	}

	// team name validation
	nameLength := len(team.Name)
	if nameLength == 0 {
		return false, errDomainNameIsEmpty
	}

	if !util.IsAlphaNumericOnly(team.Name) {
		return false, errDomainNameContainsSymbols
	}

	if nameLength < 6 {
		return false, errDomainNameMinLength
	}

	if nameLength > 63 {
		return false, errDomainNameMaxLength
	}
	return team, nil
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
		created, err := s.Create(req.Name)
		return createTeamResponse{Created: created, Err: err}, nil
	}
}
