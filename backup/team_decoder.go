// Package esh ...
// Author: Ghazni Nattarshah
// Date: Sep 15, 2016
package esh

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.com/conspico/esh/core/util"
)

// create team
type createTeamRequest struct {
	Name string `json:"name"`
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
