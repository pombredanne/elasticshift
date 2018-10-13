/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"os"
	"testing"
)

var (
	testlog = `docker registry secret

CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/linux_386/shiftctl -a -tags netgo -ldflags '-s -w' ./shiftctl.go
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o shiftctl -a -tags netgo -ldflags '-s -w' ./shiftctl.go


kubectl create secret docker-registry gitlabreg --docker-server=registry.gitlab.com --docker-username=shiftapp --docker-password=L3RKXm67q-y9G2KhqjeB

apiVersion: v1
data:
  .dockerconfigjson: eyJhdXRocyI6eyJyZWdpc3RyeS5naXRsYWIuY29tIjp7InVzZXJuYW1lIjoic2hpZnRhcHAiLCJwYXNzd29yZCI6IkwzUktYbTY3cS15OUcyS2hxamVCIiwiYXV0aCI6ImMyaHBablJoY0hBNlRETlNTMWh0TmpkeExYazVSekpMYUhGcVpVST0ifX19
kind: Secret
metadata:
  creationTimestamp: 2018-09-16T15:43:23Z
  name: gitlabreg
  namespace: default
  resourceVersion: "987485"
  selfLink: /api/v1/namespaces/default/secrets/gitlabreg
  uid: 3e5d7e5c-b9c7-11e8-8fb0-000c296a3b4f
type: kubernetes.io/dockerconfigjson
`
)

func TestPrefixWriter(t *testing.T) {

	pw := newPrefixWriter("4.1[!@#$]", os.Stdout)
	pw.Write([]byte(testlog))
}
