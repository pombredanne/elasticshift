CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/worker/worker -a -tags netgo -ldflags '-s -w' ./cmd/worker/worker.go
cd ./bin/worker
tar -zcvf worker-v0.0.1-alpha.tar.gz worker
mc cp worker-v0.0.1-alpha.tar.gz minio/test/sys
cd ..
cd ..
