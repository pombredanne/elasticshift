// Package esh
// Author Ghazni Nattarshah
// Date: Jan 4, 2017
package esh

import (
	"context"
	"github.com/iris-contrib/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

func TestTeamEncoder(t *testing.T) {
	suite.Run(t, new(TeamEncoderTestSuite))
}

type TeamEncoderTestSuite struct {
	suite.Suite
}

func (suite *TeamEncoderTestSuite) TestEncode() {

	testcases := []Testcase{
		{In: createTeamResponse{Created: true}, Err: nil},
		{In: createTeamResponse{Err: errors.New("Can't createuser")}, Err: errors.New("Can't createuser")},
	}

	ctx := context.TODO()
	for _, testcase := range testcases {
		err := encodeCreateTeamResponse(ctx, httptest.NewRecorder(), testcase.In)
		assert.ObjectsAreEqual(testcase.Err, err)
	}
}
