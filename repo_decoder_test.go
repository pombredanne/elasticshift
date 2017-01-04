// Package esh
// Author Ghazni Nattarshah
// Date: Jan 3, 2017
package esh

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"context"
	"gitlab.com/conspico/esh/core/auth"
)



func TestRepoDecoder(t *testing.T) {
suite.Run(t, new(RepoDecoderTestSuite))
}

type RepoDecoderTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (suite *RepoDecoderTestSuite) TestGetRepoRequest() {

	tok := auth.Token{ Team: "testteam"}
	suite.ctx = context.WithValue(suite.ctx, "token", tok)

	m := make(map[string]string)
	m["id"] = "id"
	suite.ctx = context.WithValue(suite.ctx, "params", m)

	req, err := decodeGetRepoRequest(suite.ctx, nil)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), req)
}