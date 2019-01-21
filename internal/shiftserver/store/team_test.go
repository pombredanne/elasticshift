package store

import (
	"testing"

	"github.com/elasticshift/elasticshift/api/types"
	"gopkg.in/mgo.v2/bson"
)

var (
	teamName = "team001"
	teamID   = "5ace310919177eb5f314d43b"
)

func TestCreateTeam(t *testing.T) {

	shift, session := testInitShiftStore()
	defer session.Close()

	team := types.Team{}
	team.ID = bson.ObjectIdHex(teamID)
	team.Name = teamName

	err := shift.Team.Save(team)
	if err != nil {
		t.Logf("Expected to save team %s, but failed.", teamName)
		t.Fail()
	}
}

func TestCheckTeamExist(t *testing.T) {

	shift, session := testInitShiftStore()
	defer session.Close()

	exist, err := shift.Team.CheckExists(teamName)
	if err != nil || !exist {
		t.Logf("Expected the team %s existance, but can't found.", teamName)
		t.Fail()
	}
}

func TestGetTeam(t *testing.T) {

	shift, session := testInitShiftStore()
	defer session.Close()

	team, err := shift.Team.GetTeam(teamID, teamName)
	if err != nil || team.ID != "" {
		t.Logf("Expected to fetch the team %s, but failed", teamName)
		t.Fail()
	}
}
