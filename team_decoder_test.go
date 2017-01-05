// Package esh
// Author Ghazni Nattarshah
// Date: Jan 4, 2017
package esh

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"context"
	"net/http/httptest"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"bytes"
)

func TestTeamDecoder(t *testing.T) {
	suite.Run(t, new(TeamDecoderTestSuite))
}

type decoderTest struct  {
	in string
	out interface{}
	err error
}

type TeamDecoderTestSuite struct {
	suite.Suite
}

func (suite *TeamDecoderTestSuite) TestDecode() {

	testcases := []decoderTest{
		{ "elasticshift", createTeamRequest{Name: "elasticshift"}, nil},
		{ "elasticshift", createTeamRequest{Name: "elasticshift"}, nil},
		{ "", false, errDomainNameIsEmpty},
		{ "e1@$t1csh1ft", false, errDomainNameContainsSymbols},
		{ "esh", false, errDomainNameMinLength},
		{ "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz", false, errDomainNameMaxLength},
	}

	ctx := context.TODO()
	for _, in := range testcases {

		data, err := json.Marshal(Team{Name: in.in})
		assert.Nil(suite.T(), err)

		req := httptest.NewRequest("POST", "http://example.com/foo", bytes.NewBuffer(data))
		res, err := decodeCreateTeamRequest(ctx, req)

		assert.ObjectsAreEqual(in.err, err)
		assert.Equal(suite.T(), in.out, res)
	}

	req := httptest.NewRequest("POST", "http://example.com/foo", bytes.NewBuffer([]byte{0}))
	res, err := decodeCreateTeamRequest(ctx, req)
	assert.NotNil(suite.T(), err)
	assert.ObjectsAreEqual(false, res)

}
