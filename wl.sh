echo "Building worker.."
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/worker/worker -a -tags netgo -ldflags '-s -w' ./cmd/worker/worker.go

echo "copying to storage"
sudo cp ./bin/worker/worker /lab/elasticshift/sys