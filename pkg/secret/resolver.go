/*
Copyright 2018 The Elasticshift Authors.
*/
package secret

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/store"
	"gopkg.in/mgo.v2/bson"
)

var (
	errTeamCannotBeEmpty       = errors.New("Team must be provided")
	errSecretIDCannotBeEmpty   = errors.New("Secret identifier cannot be empty")
	errSecretNameCannotBeEmpty = errors.New("Secret name cannot be empty")
)

const (
	TYPE_SECRET string = "secret"
	TYPE_SSHKEY string = "sshkey"
	TYPE_PGP    string = "pgp"
)

const (
	RefType_SYS  string = "sys"
	RefType_VCS  string = "vcs"
	RefType_TEAM string = "team"
	RefType_USER string = "user"
)

type resolver struct {
	store     store.Secret
	teamStore store.Team
	logger    logrus.Logger
	Ctx       context.Context
}

func (r *resolver) FetchSecret(params graphql.ResolveParams) (interface{}, error) {

	teamID, _ := params.Args["team_id"].(string)
	if teamID == "" {
		return nil, errTeamCannotBeEmpty
	}

	q := bson.M{"team_id": teamID, "internal_type": TYPE_SECRET}

	var result []types.Secret
	err := r.store.FindAll(q, &result)

	var res types.SecretList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}

func (r *resolver) FetchSecretByID(params graphql.ResolveParams) (interface{}, error) {

	team_id, _ := params.Args["team_id"].(string)
	if team_id == "" {
		return nil, errTeamCannotBeEmpty
	}

	id, _ := params.Args["id"].(string)
	if id == "" {
		return nil, errSecretIDCannotBeEmpty
	}

	q := bson.M{"team_id": team_id, "internal_type": TYPE_SECRET}
	q["_id"] = bson.ObjectIdHex(id)

	var result types.Secret
	err := r.store.FindOne(q, &result)

	return &result, err
}

func (r *resolver) FetchSecretByName(params graphql.ResolveParams) (interface{}, error) {

	team_id, _ := params.Args["team_id"].(string)
	if team_id == "" {
		return nil, errTeamCannotBeEmpty
	}

	name, _ := params.Args["name"].(string)
	if name == "" {
		return nil, errSecretNameCannotBeEmpty
	}

	q := bson.M{"team_id": team_id, "internal_type": TYPE_SECRET}
	q["name"] = name

	var result types.Secret
	err := r.store.FindOne(q, &result)

	return &result, err
}

func (r *resolver) AddSecret(params graphql.ResolveParams) (interface{}, error) {

	teamID, _ := params.Args["team_id"].(string)
	name, _ := params.Args["name"].(string)
	kind, _ := params.Args["kind"].(string)
	referenceKind, _ := params.Args["reference_kind"].(string)
	referenceID, _ := params.Args["reference_id"].(string)

	q := bson.M{"team_id": teamID}
	q["name"] = name
	q["kind"] = kind
	q["reference_kind"] = referenceKind
	q["reference_id"] = referenceID

	var sec types.Secret
	err := r.store.FindOne(q, &sec)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to check if the given secret already exist :%v", err)
	}

	if sec.ID.Hex() != "" {
		return nil, fmt.Errorf("The secret name '%s' already exist under the reference kind for your team", name)
	}

	value, _ := params.Args["value"].(string)
	if value == "" {
		return nil, fmt.Errorf("The actual value of secret can't be empty.")
	}

	sec.Name = name
	sec.TeamID = teamID
	sec.Kind = kind
	sec.ReferenceKind = referenceKind
	sec.ReferenceID = referenceID
	sec.Value = value
	sec.InternalType = TYPE_SECRET

	err = r.store.Save(&sec)
	if err != nil {
		return nil, fmt.Errorf("Failed to add integration: %v", err)
	}
	return sec, nil
}
