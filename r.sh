rm -f esh
echo "Building..."
go build ./cmd/esh.go
echo "Running..."
./esh