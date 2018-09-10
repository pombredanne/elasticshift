source .env
name=registry.gitlab.com/conspico/esh:${SHIFT_VERSION}_${ALPINE_VERSION}
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/linux_386/elasticshift -a -tags netgo -ldflags '-s -w' ./cmd/elasticshift/elasticshift.go
docker build -t $name --build-arg ALPINE_VERSION=${ALPINE_VERSION} .
docker push $name