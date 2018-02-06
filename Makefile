#
# Copyright 2015 The Elasticshift Authors.
#
PROJ=elasticshift
ORG_PATH=gitlab.com/conspico
REPO_PATH=$(ORG_PATH)/$(PROJ)
export PATH := $(PWD)/bin:$(PATH)

VERSION=$(shell ./scripts/git-version)

user=$(shell id -u -n)
group=$(shell id -g -n)

export GOBIN=$(PWD)/bin

export PATH=$(GOBIN):$(shell printenv PATH)

LD_FLAGS="-w -X $(REPO_PATH)/version.Version=$(VERSION)"

build: bin/elasticshift get-golint

bin/elasticshift: go-version-checker
	@go install -x
		GOOS=linux GOARCH=386 go build -x -o $(GOBIN)/linux_386/elasticshift ./cmd/elasticshift/elasticshift.go
		GOOS=darwin GOARCH=386 go build -x -o $(GOBIN)/darwin_386/elasticshift ./cmd/elasticshift/elasticshift.go
#env GOOS=linux CGO_ENABLED=0 go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)
#GOOS=linux CGO_ENABLED=0 go build -o /bin/elasticshift elasticshift.go
#CGO_ENABLED=0 go build -o /bin/elasticshift elasticshift.go

.PHONY: revendor
revendor:
	@glide up -v
	@glide-vc --use-lock-file --no-tests --only-code --keep '**/*.proto'

test:
	@go test -v -i $(shell go list ./... | grep -v '/vendor\|/api')
	@go test -v $(shell go list ./... | grep -v '/vendor\|/api')

testrace:
	@go test -v -i --race $(shell go list ./... | grep -v '/vendor\|/api')
	@go test -v --race $(shell go list ./... | grep -v '/vendor\|/api')

vet:
	@go vet $(shell go list ./... | grep -v '/vendor\|/api')

fmt:
	@go fmt $(shell go list ./... | grep -v '/vendor\|/api')

lint:
	@for package in $(shell go list ./... | grep -v '/vendor\|/api'); do \
      golint -set_exit_status $$package $$i || exit 1; \
	done

.PHONY: test-cover-html
PACKAGES = $(shell go list ./... | grep -v '/vendor\|/api')

test-cover-html:

	@rm -f bin/coverage.out
	@rm -f bin/coverage-all.out

	#@rm -f bin/coverage.out bin/coverage-all.out
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

pb/gen-pb.go: bin/protoc bin/protoc-gen-go bin/protoc-gen-grpc-gateway bin/protoc-gen-swagger get-proto-descriptor
	@./bin/protoc -I/usr/local/include -I. \
	-I${PWD}/vendor \
	--go_out=pMgoogle/api/annotations.proto=github.com/google/go-genproto/googleapis/api/annotations,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor,plugins=grpc:. \
	api/*.proto
	@./bin/protoc -I/usr/local/include -I. \
	-I${PWD}/vendor \
	--plugin=protoc-gen-grpc-gateway=$(PWD)/bin/protoc-gen-grpc-gateway \
	--grpc-gateway_out=logtostderr=true:. \
	api/*.proto
	@./bin/protoc --go_out=import_path=dex,plugins=grpc:. api/dex/api.proto

bin/protoc: scripts/get-protoc
	@./scripts/get-protoc bin/protoc

get-golint:
	@go get -u github.com/golang/lint/golint

bin/protoc-gen-go:
	@go install -v $(REPO_PATH)/vendor/github.com/golang/protobuf/protoc-gen-go

bin/protoc-gen-grpc-gateway:
	@go install -v $(REPO_PATH)/vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

bin/protoc-gen-swagger:
	@go install -v $(REPO_PATH)/vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

PROTOBUF_URL = "https://raw.githubusercontent.com/google/protobuf/master/src/google/protobuf"
GOOGLEAPIS_URL = "https://rawgit.com/googleapis/googleapis/master/google/api"
DEX_VERSION=master

get-proto-descriptor:
	@mkdir -p ${PWD}/vendor/google/protobuf
	@mkdir -p ${PWD}/vendor/google/api
	@wget -nc -c ${PROTOBUF_URL}/descriptor.proto -P vendor/google/protobuf
	@wget -nc -c ${GOOGLEAPIS_URL}/annotations.proto -P vendor/google/api
	@wget -nc -c ${GOOGLEAPIS_URL}/http.proto -P vendor/google/api
	# @wget -nc -c https://raw.githubusercontent.com/coreos/dex/${DEX_VERSION}/api/api.proto -P api/dex

.PHONY: go-version-checker
go-version-checker:
	@./scripts/go-version-checker

clean:
	@rm -rf bin/

testall: testrace vet fmt lint

all: build vet fmt lint test docker-image
FORCE:

.PHONY: test testrace vet fmt lint testall
