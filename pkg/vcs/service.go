/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"fmt"
	"net/http"
	"strings"

	"bytes"

	"encoding/base64"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/pkg/identity/oauth2/providers"
	"gitlab.com/conspico/elasticshift/pkg/identity/team"
	"gitlab.com/conspico/elasticshift/pkg/secret"
	stypes "gitlab.com/conspico/elasticshift/pkg/store/types"
)

var (

	// VCS
	errNoProviderFound         = "No provider found for %s"
	errGetUpdatedFokenFailed   = "Failed to get updated token %s"
	errGettingRepositories     = "Failed to get repositories for %s"
	errVCSAccountAlreadyLinked = "VCS account already linked"
)

// True or False
const (
	True  = 1
	False = 0
)

// Constants for performing encode decode
const (
	EQUAL        = "="
	DOUBLEEQUALS = "=="
	DOT0         = ".0"
	DOT1         = ".1"
	DOT2         = ".2"
)

// Common constants
const (
	SLASH     = "/"
	SEMICOLON = ";"
)

type service struct {
	store     Store
	teamStore team.Store
	vault     secret.Vault
	logger    logrus.Logger
	providers providers.Providers
}

// Service ..
type Service interface {
	Authorize(w http.ResponseWriter, r *http.Request)
	Authorized(w http.ResponseWriter, r *http.Request)
}

// NewVCSService ..
func NewService(logger logrus.Logger, d stypes.Database, providers providers.Providers, teamStore team.Store, vault secret.Vault) Service {

	return &service{
		store:     NewStore(d),
		teamStore: teamStore,
		vault:     vault,
		logger:    logger,
		providers: providers,
	}
}

func (s service) Authorize(w http.ResponseWriter, r *http.Request) {

	team := mux.Vars(r)["team"]
	exist, err := s.teamStore.CheckExists(team)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch the team: %s, error: %v", team, err), http.StatusBadRequest)
		return
	}

	if !exist {
		http.Error(w, fmt.Sprintf("Team '%s' doesn't exist, please provide the valid name", team), http.StatusBadRequest)
		return
	}

	provider := mux.Vars(r)["provider"]
	p, err := s.providers.Get(provider)

	if err != nil {
		http.Error(w, fmt.Sprintf("Getting provider %s failed: %v", provider, err), http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	buf.WriteString(team)
	buf.WriteString(SEMICOLON)
	buf.WriteString(SLASH)
	buf.WriteString(SLASH)
	buf.WriteString(r.Host)

	url := p.Authorize(s.encode(buf.String()))

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Authorized ..
// Invoked when authorization finished by oauth app
func (s service) Authorized(w http.ResponseWriter, r *http.Request) {

	provider := mux.Vars(r)["provider"]
	p, err := s.providers.Get(provider)
	if err != nil {
		http.Error(w, fmt.Sprintf("Getting provider %s failed: %v", provider, err), http.StatusBadRequest)
	}

	id := r.FormValue("id")
	code := r.FormValue("code")
	u, err := p.Authorized(id, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Finalize the authorization failed : %v", err), http.StatusBadRequest)
	}

	unescID := s.decode(id)
	escID := strings.Split(unescID, SEMICOLON)

	// persist user
	team := escID[0]

	acc, err := s.teamStore.GetVCSByID(team, u.ID)
	if strings.EqualFold(acc.ID, u.ID) {

		sec, err := s.vault.GetByReferenceID(u.ID, secret.RefType_VCS)
		if err != nil {
			http.Error(w, "Failed to update VCS", http.StatusConflict)
		}

		props, err := s.props(u)
		if err != nil {
			http.Error(w, "Failed to update VCS", http.StatusConflict)
		}
		sec.Value = props

		// encrypt and vault the value
		id, err = s.vault.Put(sec)
		if err != nil {
			http.Error(w, "Failed to update VCS", http.StatusConflict)
		}

		// updvcs.UpdatedDt = time.Now()
		acc.OwnerType = u.OwnerType
		acc.TokenExpiry = u.TokenExpiry

		// update the key id
		s.teamStore.UpdateVCS(team, acc)

		http.Error(w, errVCSAccountAlreadyLinked, http.StatusConflict)
	}

	u.Source = r.Host

	var sec types.Secret
	secretID := s.saveSecret(sec, u, team, w)

	if err == nil {

		// TODO sync the repo and setup hook reqeust for the repo
		// go p.CreateHook(u.AccessCode, u.Name, u.OwnerType)
	}

	// u.ID = utils.NewUUID()
	// u.CreatedDt = time.Now()
	// u.UpdatedDt = time.Now()

	u.SecretID = secretID

	err = s.teamStore.SaveVCS(team, &u)
	if err != nil {
		// TODO return http error
		s.logger.Errorln("SAVE VCS: ", err)
	}
	url := escID[1] + "/sysconf/vcs"
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s service) saveSecret(sec types.Secret, u types.VCS, team string, w http.ResponseWriter) string {

	// save token to secret store
	sec.Name = u.Source + "/" + u.Name
	sec.Kind = secret.TYPE_SECRET
	sec.ReferenceKind = secret.RefType_VCS
	sec.ReferenceID = u.ID
	sec.TeamID = team

	props, err := s.props(u)
	if err != nil {
		http.Error(w, "Failed to save VCS info:", http.StatusInternalServerError)
	}
	sec.Value = props

	id, err := s.vault.Put(sec)
	if err != nil {
		http.Error(w, "Failed to save VCS info:", http.StatusInternalServerError)
	}
	return id
}

func (s service) props(u types.VCS) (string, error) {

	p := secret.NewPair()
	p["access_token"] = u.AccessToken
	p["refresh_token"] = u.RefreshToken

	return p.Json()
}

// func (s vcsService) GetVCS(teamID string) (GetVCSResponse, error) {

// 	result, err := s.vcsDS.GetVCS(teamID)
// 	return GetVCSResponse{Result: result}, err
// }

// func (s vcsService) SyncVCS(teamID, userName, providerID string) (bool, error) {

// 	acc, err := s.vcsDS.GetByID(providerID)
// 	if err != nil {
// 		return false, fmt.Errorf("Get by VCS ID failed during sync : %v", err)
// 	}

// 	err = s.sync(acc, userName)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func (s vcsService) sync(acc VCS, userName string) error {

// 	// Get the token
// 	t, err := s.getToken(acc)
// 	if err != nil {
// 		return fmt.Errorf("Get token failed : ", err)
// 	}

// 	// fetch the existing repository
// 	p, err := s.getProvider(acc.Type)
// 	if err != nil {
// 		return fmt.Errorf(errNoProviderFound, err)
// 	}

// 	// repository received from provider
// 	repos, err := p.GetRepos(t, acc.Name, acc.OwnerType)
// 	if err != nil {
// 		return fmt.Errorf("Failed to get repos from provider %s : %v", p.Name(), err)
// 	}

// 	// Fetch the repositories from esh repo store
// 	lrpo, err := s.repoDS.GetReposByVCSID(acc.TeamID, acc.ID)
// 	if err != nil {
// 		return fmt.Errorf("Getting repos by vcs id failed : %v", err)
// 	}

// 	rpo := make(map[string]Repo)
// 	for _, l := range lrpo {
// 		rpo[l.RepoID] = l
// 	}

// 	// combine the result set
// 	for _, rp := range repos {

// 		r, exist := rpo[rp.RepoID]
// 		if exist {

// 			updrepo := Repo{}
// 			updated := false
// 			if r.Name != rp.Name {
// 				updrepo.Name = rp.Name
// 				updated = true
// 			}

// 			if r.Private != rp.Private {
// 				updrepo.Private = rp.Private
// 				updated = true
// 			}

// 			if r.Link != rp.Link {
// 				updrepo.Link = rp.Link
// 				updated = true
// 			}

// 			if r.Description != rp.Description {
// 				updrepo.Description = rp.Description
// 				updated = true
// 			}

// 			if r.Fork != rp.Fork {
// 				updrepo.Fork = rp.Fork
// 				updated = true
// 			}

// 			if r.DefaultBranch != rp.DefaultBranch {
// 				updrepo.DefaultBranch = rp.DefaultBranch
// 				updated = true
// 			}

// 			if r.Language != rp.Language {
// 				updrepo.Language = rp.Language
// 				updated = true
// 			}

// 			if updated {
// 				// perform update
// 				updrepo.UpdatedBy = userName
// 				s.repoDS.Update(r, updrepo)
// 			}
// 		} else {

// 			// perform insert
// 			rp.ID, _ = util.NewUUID()
// 			rp.CreatedDt = time.Now()
// 			rp.UpdatedDt = time.Now()
// 			rp.CreatedBy = userName
// 			rp.TeamID = acc.TeamID
// 			rp.VcsID = acc.ID
// 			s.repoDS.Save(&rp)
// 		}

// 		// removes from the map
// 		if exist {
// 			delete(rpo, r.RepoID)
// 		}
// 	}

// 	var ids []string
// 	// Now iterate thru deleted repositories.
// 	for _, rp := range rpo {
// 		ids = append(ids, rp.ID)
// 	}

// 	err = s.repoDS.DeleteIds(ids)
// 	if err != nil {
// 		return fmt.Errorf("Failed to delete the vcs that does not exist remotly : %v", err)
// 	}

// 	return nil
// }

func (s service) encode(id string) string {

	eid := base64.URLEncoding.EncodeToString([]byte(id))
	if strings.Contains(eid, DOUBLEEQUALS) {
		eid = strings.TrimRight(eid, DOUBLEEQUALS) + DOT2
	} else if strings.Contains(eid, EQUAL) {
		eid = strings.TrimRight(eid, EQUAL) + DOT1
	} else {
		eid = eid + DOT0
	}
	return eid
}

func (s service) decode(id string) string {

	if strings.Contains(id, DOT2) {
		id = strings.TrimRight(id, DOT2) + DOUBLEEQUALS
	} else if strings.Contains(id, DOT1) {
		id = strings.TrimRight(id, DOT1) + EQUAL
	} else {
		id = strings.TrimRight(id, DOT0)
	}
	did, _ := base64.URLEncoding.DecodeString(id)
	return string(did[:])
}
