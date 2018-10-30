/*
Copyright 2018 The Elasticshift Authors.
*/
package shiftfile

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errNameCantBeEmpty          = errors.New("Shiftfile name cannot be empty")
	errFileContentCannotBeEmpty = errors.New("Shiftfile content is empty")
	errTeamIDCantBeEmpty        = errors.New("Team ID cannot be empty")
)

// Resolver ...
type Resolver interface {
	FetchShiftfile(params graphql.ResolveParams) (interface{}, error)
	AddShiftfile(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	store  store.Shiftfile
	logger *logrus.Entry
	Ctx    context.Context
}

// NewResolver ...
func NewResolver(ctx context.Context, loggr logger.Loggr, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:  s.Shiftfile,
		logger: loggr.GetLogger("graphql/shiftfile"),
		Ctx:    ctx,
	}
	return r, nil
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
