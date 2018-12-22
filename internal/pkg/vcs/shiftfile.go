/*
Copyright 2018 The Elasticshift Authors.
*/
package vcs

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/elasticshift/elasticshift/internal/shiftserver/identity/oauth2/providers"
	"github.com/elasticshift/elasticshift/pkg/dispatch"
)

var (
	GITHUB_DOT_COM    = "github.com"
	GITLAB_DOT_COM    = "gitlab.com"
	BITBUCKET_DOT_ORG = "bitbucket.org"
)

var (
	// githubComUrl = "https://api.github.com/repos/nshahm/hybrid.test.runner/contents/Shiftfile?ref=master"
	githubComUrl = providers.GithubBaseURL + "/repos/:account/:repo/contents/Shiftfile"
)

func GetSource(provider string) string {

	source := ""
	if provider == providers.GithubProviderName {
		source = GITHUB_DOT_COM
	} else if provider == providers.GitlabProviderName {
		source = GITLAB_DOT_COM
	} else if provider == providers.BitbucketProviderName {
		source = BITBUCKET_DOT_ORG
	}

	// TODO must add for enterprise versions.
	return source
}

func GetShiftFile(source, url, branch string) ([]byte, error) {

	switch source {

	case GITHUB_DOT_COM:
		return getShiftfileFromGithub(url, branch)
	case GITLAB_DOT_COM:
	case BITBUCKET_DOT_ORG:
	}
	return nil, fmt.Errorf("Url not supported")
}

func getShiftfileFromGithub(url, branch string) ([]byte, error) {

	_, account, repo := parseGitUrl(url)

	r := dispatch.NewGetRequestMaker(githubComUrl)

	r.Header("Accept", dispatch.JSON)

	r.PathParams(account, repo)

	r.QueryParam("ref", branch)

	result := struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}{}

	err := r.Scan(&result).Dispatch()
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(result.Content)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func parseGitUrl(uri string) (string, string, string) {

	// parse uri and identify the VCS
	// git@github.com:nshahm/hybrid.test.runner.git
	protoGit := strings.HasPrefix(uri, "git@")
	protoHttps := strings.HasPrefix(uri, "http")

	eIdx := strings.LastIndex(uri, "/")
	var sIdx int
	var source, account, repoName string
	if protoGit {

		sIdx = strings.Index(uri, "@")
		val := uri[sIdx+1 : eIdx]
		valArr := strings.Split(val, ":")
		source = valArr[0]
		account = valArr[1]
		repoName = uri[eIdx+1:]

	} else if protoHttps {

		valArr := strings.Split(uri, "/")
		source = valArr[2]
		account = valArr[3]
		repoName = valArr[4]
	}

	repoName = strings.TrimRight(repoName, ".git")

	return source, account, repoName
}
