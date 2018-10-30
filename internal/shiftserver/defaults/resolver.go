/*
Copyright 2018 The Elasticshift Authors.
*/
package defaults

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty     = errors.New("Container ID cannot be empty")
	errTeamCannotBeEmpty = errors.New("Team must be provided")
)

// Default kinds
const (
	DK_Team int = iota + 1
	DK_User
)

type Resolver interface {
	FetchDefault(params graphql.ResolveParams) (interface{}, error)
	SetDefaults(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	store     store.Defaults
	logger    *logrus.Entry
	Ctx       context.Context
	teamStore store.Team
}

// NewResolver ...
func NewResolver(ctx context.Context, loggr logger.Loggr, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:     s.Defaults,
		logger:    loggr.GetLogger("graphql/default"),
		Ctx:       ctx,
		teamStore: s.Team,
	}
	return r, nil
}

func (r *resolver) FetchDefault(params graphql.ResolveParams) (interface{}, error) {

	q := bson.M{}

	id, _ := params.Args["id"].(string)
	if id != "" {
		q["_id"] = bson.ObjectIdHex(id)
	}

	kind, _ := params.Args["kind"].(int)
	if kind > 0 {
		q["kind"] = kind
	}

	referenceID, _ := params.Args["reference_id"].(string)
	if referenceID != "" {
		q["reference_id"] = referenceID
	}

	var result types.Default
	err := r.store.FindOne(q, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the defaults: %v", err)
	}

	return &result, err
}

func (r *resolver) SetDefaults(params graphql.ResolveParams) (interface{}, error) {

	reference_id, _ := params.Args["reference_id"].(string)
	kind, _ := params.Args["kind"].(int)
	if kind < 1 {
		return nil, fmt.Errorf("Default kind should either be a type of user or team")
	}

	storage_id, _ := params.Args["storage_id"].(string)
	container_engine_id, _ := params.Args["container_engine_id"].(string)
	languages, _ := params.Args["languages"].(string)
	if storage_id == "" && container_engine_id == "" && languages == "" {
		return nil, fmt.Errorf("No default value set, please set storage_id or container_engine_id or language.")
	}

	if kind == DK_Team {

		// check the team existance
		t, err := r.teamStore.GetTeam(reference_id, "")
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch the team with reference_id provided: %v", err)
		}

		if t.ID.Hex() == "" {
			return nil, fmt.Errorf("No team found for given reference_id.")
		}

	} else if kind == DK_User {

		// TODO check the user existance
	}

	// fetch the defaults
	def, err := r.store.FindByReferenceId(reference_id)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to save the defaults: %v", err)
	}

	upd := def.ID.Hex() != ""

	updfields := bson.M{}

	if storage_id != "" {
		def.StorageID = storage_id
		if upd {
			updfields["storage_id"] = storage_id
		}
	}

	if container_engine_id != "" {
		def.ContainerEngineID = container_engine_id
		if upd {
			updfields["container_engine_id"] = container_engine_id
		}
	}

	if languages != "" {

		var props map[string]string
		err = json.Unmarshal([]byte(languages), &props)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse the language json: %v", err)
		}

		def.Languages = props

		if upd {
			updfields["languages"] = props
		}
	}

	if upd {
		err = r.store.UpdateDefaults(reference_id, updfields)
	} else {
		def.Kind = kind
		def.ReferenceID = reference_id
		err = r.store.Save(&def)
	}

	return &def, nil
}
