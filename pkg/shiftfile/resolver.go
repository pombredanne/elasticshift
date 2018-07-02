/*
Copyright 2018 The Elasticshift Authors.
*/
package shiftfile

import (
	"context"
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errNameCantBeEmpty          = errors.New("Shiftfile name cannot be empty")
	errFileContentCannotBeEmpty = errors.New("Shiftfile content is empty")
	errTeamIDCantBeEmpty        = errors.New("Team ID cannot be empty")
)

type resolver struct {
	store  store.Shiftfile
	logger logrus.Logger
	Ctx    context.Context
}

func (r *resolver) FetchShiftfile(params graphql.ResolveParams) (interface{}, error) {

	name := params.Args["name"].(string)
	if name == "" {
		return nil, errNameCantBeEmpty
	}

	q := bson.M{}

	if name != "" {
		q["name"] = name
	}

	var err error
	var result []types.Shiftfile
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	var res types.ShiftfileList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}

func (r *resolver) AddShiftfile(params graphql.ResolveParams) (interface{}, error) {

	name, _ := params.Args["name"].(string)
	if name == "" {
		return nil, errNameCantBeEmpty
	}

	teamID, _ := params.Args["team_id"].(string)
	if teamID == "" {
		return nil, errTeamIDCantBeEmpty
	}

	description, _ := params.Args["description"].(string)
	file, _ := params.Args["file"].(string)
	if file == "" {
		return nil, errFileContentCannotBeEmpty
	}

	// TODO validate the file content

	sf := types.Shiftfile{}
	sf.Name = name
	sf.Description = description
	sf.File = []byte(file)
	sf.TeamID = teamID

	err := r.store.Save(&sf)
	return sf, err
}
