/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"testing"

	"github.com/elasticshift/elasticshift/api"
	"github.com/elasticshift/elasticshift/internal/pkg/shiftfile/parser"
	wtypes "github.com/elasticshift/elasticshift/internal/worker/types"
)

var cachetestfile = `
VERSION "1.0"
NAME "elasticshift/runner"

IMAGE "alphine:latest"

CACHE {
	- ~/.mc
}
`

func TestCache(t *testing.T) {

	f, err := parser.AST([]byte(cachetestfile))
	if err != nil {
		t.Fail()
	}

	proj := &api.GetProjectRes{}
	proj.RepositoryId = "repo_id"
	proj.Branch = "master"
	b := &builder{f: f}
	cfg := wtypes.Config{}
	cfg.TeamID = "team_id"
	cfg.BuildID = "build_id"
	cfg.ShiftDir = "/opt/elasticshift"
	b.config = cfg
	b.project = proj

	//b.saveCache()
	//b.restoreCache()
}
