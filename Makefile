#
# Copyright 2015 The Elasticshift Authors.
#
PROJ=elasticshift
ORG_PATH=github.com/elasticshift
REPO_PATH=$(ORG_PATH)/$(PROJ)
export PATH := $(PWD)/bin:$(PATH)

VERSION=$(shell ./scripts/git-version)

user=$(shell id -u -n)
group=$(shell id -g -n)

export GOBIN=$(PWD)/bin

export PATH=$(GOBIN):$(shell printenv PATH)

LD_FLAGS="-w -X $(REPO_PATH)/version.Version=$(VERSION)"

build: bin/elasticshift

bin/elasticshift: go-version-checker
#	@go install -x
		CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o $(GOBIN)/darwin_386/elasticshift -tags netgo -ldflags '-s -w' ./cmd/elasticshift/elasticshift.go
		CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o $(GOBIN)/linux_386/elasticshift -tags netgo -ldflags '-s -w' ./cmd/elasticshift/elasticshift.go
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(GOBIN)/darwin_amd64/elasticshift -tags netgo -ldflags '-s -w' ./cmd/elasticshift/elasticshift.go
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/linux_amd64/elasticshift -tags netgo -ldflags '-s -w' ./cmd/elasticshift/elasticshift.go
		CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o $(GOBIN)/darwin_386/worker -tags netgo -ldflags '-s -w' ./cmd/worker/worker.go
		CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o $(GOBIN)/linux_386/worker -tags netgo -ldflags '-s -w' ./cmd/worker/worker.go
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(GOBIN)/darwin_amd64/worker -tags netgo -ldflags '-s -w' ./cmd/worker/worker.go
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/linux_amd64/worker -tags netgo -ldflags '-s -w' ./cmd/worker/worker.go

.PHONY: outdated
outdated:
	@go list -m -u all

.PHONY: update
update:
	@go get -u

test:
	@go test -v -i $(shell go list ./... | grep -v '/api')
	@go test -v $(shell go list ./... | grep -v '/api')

testrace:
	@go test -v -i --race $(shell go list ./... | grep -v '/api')
	@go test -v --race $(shell go list ./... | grep -v '/api')

vet:
	@go vet $(shell go list ./... | grep -v '/api')

fmt:
	@go fmt $(shell go list ./... | grep -v '/api')

lint: get-golint
	@for package in $(shell go list ./... | grep -v '/api'); do \
      golint -set_exit_status $$package $$i || exit 1; \
	done

.PHONY: test-cover-html
PACKAGES = $(shell go list ./... | grep -v '/api')

test-cover-html:

	@rm -f bin/coverage.out
	@rm -f bin/coverage-all.out

	echo "mode: count" > bin/coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> bin/coverage-all.out;)
	go tool cover -html=bin/coverage-all.out

.PHONY: docker-image
docker-image:
	@docker build -t elasticshift .

.PHONY: grpc
grpc: pb/gen-pb.go

#pb/gen-pb.go: bin/protoc bin/protoc-gen-go bin/protoc-gen-grpc-gateway bin/protoc-gen-swagger get-proto-descriptor
pb/gen-pb.go: bin/protoc bin/protoc-gen-go get-proto-descriptor
	@./bin/protoc -I/usr/local/include -I. \
	-I${PWD}/bin/tmp/ \
	--go_out=pMgoogle/api/annotations.proto=github.com/google/go-genproto/googleapis/api/annotations,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,plugins=grpc:. \
	api/*.proto
	# @./bin/protoc -I/usr/local/include -I. \
	# -I${PWD}/vendor \
	# --plugin=protoc-gen-grpc-gateway=$(PWD)/bin/protoc-gen-grpc-gateway \
	# --grpc-gateway_out=logtostderr=true:. \
	# api/*.proto

bin/protoc: scripts/get-protoc
	@./scripts/get-protoc bin/protoc

get-golint:
	@go get -u github.com/golang/lint/golint

bin/protoc-gen-go:
	@go install -v github.com/golang/protobuf/protoc-gen-go

bin/protoc-gen-grpc-gateway:
	@go install -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

bin/protoc-gen-swagger:
	@go install -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

PROTOBUF_URL = "https://raw.githubusercontent.com/google/protobuf/master/src/google/protobuf"
GOOGLEAPIS_URL = "https://rawgit.com/googleapis/googleapis/master/google/api"

get-proto-descriptor:
	@mkdir -p ${PWD}/bin/tmp/google/protobuf
	@mkdir -p ${PWD}/bin/tmp/google/api
	@wget -nc -c ${PROTOBUF_URL}/descriptor.proto -P bin/tmp/google/protobuf
	@wget -nc -c ${PROTOBUF_URL}/timestamp.proto -P bin/tmp/google/protobuf
	@wget -nc -c ${GOOGLEAPIS_URL}/annotations.proto -P bin/tmp/google/api
	@wget -nc -c ${GOOGLEAPIS_URL}/http.proto -P bin/tmp/google/api

.PHONY: go-version-checker
go-version-checker:
	@./scripts/go-version-checker

clean:
	@rm -rf bin/

testall: testrace vet fmt lint

all: build vet fmt lint test
FORCE:

.PHONY: test testrace vet fmt lint testall
