// Package esh
// Author Ghazni Nattarshah
// Date: Jan 5, 2017
package esh

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"context"
	"gitlab.com/conspico/esh/core/auth"
	"net/url"
	"bytes"
)

func TestVCSDecoder(t *testing.T) {
	suite.Run(t, new(VCSDecoderTestSuite))
}

type VCSDecoderTestSuite struct {
	suite.Suite
}

func (suite *VCSDecoderTestSuite) TestAuthorize() {

	ctx := context.TODO()
	ctx = context.WithValue(ctx, "token", auth.Token{Team:"elasticshift"})

	m := make(map[string]string)
	m["provider"] = "github"
	ctx = context.WithValue(ctx, "params", m)

	req := httptest.NewRequest("POST", "http://example.com/foo", nil)
	res, err := decodeAuthorizeRequest(ctx, req)

	assert.Nil(suite.T(), err)
	assert.ObjectsAreEqual(AuthorizeRequest{Team: "elasticshift", Provider:"github", Request: req}, res)
}

func (suite *VCSDecoderTestSuite) TestAuthorized() {

	ctx := context.TODO()

	m := make(map[string]string)
	m["provider"] = "github"
	ctx = context.WithValue(ctx, "params", m)

	params:= url.Values{}
	params.Add("id", "id")
	params.Add("code", "code")

	req := httptest.NewRequest("POST", "http://example.com/foo", bytes.NewBufferString(params.Encode()))

	res, err := decodeAuthorizedRequest(ctx, req)

	assert.Nil(suite.T(), err)
	assert.ObjectsAreEqual(AuthorizeRequest{ID: "id", Provider:"github", Request: req, Code:"code"}, res)
}

func (suite *VCSDecoderTestSuite) TestGetVCS() {

	ctx := context.TODO()
	ctx = context.WithValue(ctx, "token", auth.Token{Team:"elasticshift"})
	req := httptest.NewRequest("POST", "http://example.com/foo", nil)

	res, err := decodeGetVCSRequest(ctx, req)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "elasticshift", res)
}

func (suite *VCSDecoderTestSuite) TestSyncVCS() {

	ctx := context.TODO()
	ctx = context.WithValue(ctx, "token", auth.Token{Team:"elasticshift", Username: "testuser"})

	m := make(map[string]string)
	m["id"] = "id"
	ctx = context.WithValue(ctx, "params", m)

	req := httptest.NewRequest("POST", "http://example.com/foo", nil)
	res, err := decodeSyncVCSRequest(ctx, req)

	assert.Nil(suite.T(), err)
	assert.ObjectsAreEqual(SyncVCSRequest{Team: "elasticshift", Username:"testuser", ProviderID:"id" }, res)
}