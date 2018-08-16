echo "Building elasticshift.."
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/linux_386/elasticshift -a -tags netgo -ldflags '-s -w' ./cmd/elasticshift/elasticshift.go

echo "Building elasticshift docker image"
docker build . -t elasticshift

echo "Running docker compose"
docker-compose up --remove-orphans --force-recreate