echo "Building worker.."
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/worker/worker -a -tags netgo -ldflags '-s -w' ./cmd/worker/worker.go

echo "Removing existing tar.gz file"
rm -f worker-linux-x86_64-1.0.0.tar.gz

echo "Creating tar file"
tar -zcvf worker-linux-x86_64-1.0.0.tar.gz -C ./bin/worker/ .

echo "copying to minio downloads directory."
cp worker-linux-x86_64-1.0.0.tar.gz ~/.elasticshift/minio/data/downloads

echo "copying worker-startup script"
cp ~/sandbox/projects/cp/scripts/worker.sh ~/.elasticshift/minio/data/downloads

cp ./bin/worker/worker ~/mkshare/sys

