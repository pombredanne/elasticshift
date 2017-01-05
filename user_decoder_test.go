// Package esh
// Author Ghazni Nattarshah
// Date: Jan 4, 2017
package esh

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"bytes"
	"context"
)

func TestUserDecoder(t *testing.T) {
	suite.Run(t, new(UserDecoderTestSuite))
}

type UserDecoderTestSuite struct {
	suite.Suite
}

func (suite *UserDecoderTestSuite) TestSignup() {

	sr  := signupRequest{Fullname:"GN", Email:"g.n@gmai.com", Password:"Pwd", Team:"testteam"}
	sr2 := signupRequest{Fullname:"GN", Email:"g.n@gmai.com", Password:"Pwd"}
	sr3 := signupRequest{Fullname:"GN", Email:"g.n@gmai.com", Password:"Pwd", Team:"elasticshift"}
	testcases := []Testcase{
		{sr, sr, nil},
		{sr2, sr3, nil},
		{sr3, sr3, nil},
	}

	ctx := context.TODO()
	ctx = context.WithValue(ctx, "subdomain", "elasticshift")

	for _, testcase := range testcases {

		data, err := json.Marshal(testcase.In)
		assert.Nil(suite.T(), err)

		req := httptest.NewRequest("POST", "http://example.com/foo", bytes.NewBuffer(data))
		res, err := decodeSignUpRequest(ctx, req)

		assert.ObjectsAreEqual(testcase.Err, err)
		assert.Equal(suite.T(), testcase.Out, res)
	}

	req := httptest.NewRequest("POST", "http://example.com/foo", bytes.NewBuffer([]byte{0}))
	res, err := decodeSignUpRequest(ctx, req)
	assert.NotNil(suite.T(), err)
	assert.ObjectsAreEqual(false, res)
}


func (suite *UserDecoderTestSuite) TestSignin() {

	sr  := signInRequest{Email:"g.n@gmai.com", Password:"Pwd", Team:"testteam"}
	sr2 := signInRequest{Email:"g.n@gmai.com", Password:"Pwd"}
	sr3 := signInRequest{Email:"g.n@gmai.com", Password:"Pwd", Team:"elasticshift"}
	sr4 := signInRequest{Email:"g.n@gmai.com"}
	testcases := []Testcase{
		{sr, sr, nil},
		{sr2, sr3, nil},
		{sr3, sr3, nil},
		{sr4, false, errUsernameOrPasswordIsEmpty},
	}

	ctx := context.TODO()
	ctx = context.WithValue(ctx, "subdomain", "elasticshift")

	for _, testcase := range testcases {

		data, err := json.Marshal(testcase.In)
		assert.Nil(suite.T(), err)

		req := httptest.NewRequest("POST", "http://example.com/foo", bytes.NewBuffer(data))
		res, err := decodeSignInRequest(ctx, req)

		assert.ObjectsAreEqual(testcase.Err, err)
		assert.Equal(suite.T(), testcase.Out, res)
	}

	req := httptest.NewRequest("POST", "http://example.com/foo", bytes.NewBuffer([]byte{0}))
	res, err := decodeSignInRequest(ctx, req)
	assert.NotNil(suite.T(), err)
	assert.ObjectsAreEqual(false, res)
}

func (suite *UserDecoderTestSuite) TestSignOut() {

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	res, err := decodeSignOutRequest(context.TODO(), req)

	assert.Nil(suite.T(), err)
	assert.ObjectsAreEqual(req, res)
}
