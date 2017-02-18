rm -f esh
echo "Building..."
go build -o=esh ./cmd/esh/main.go
echo "Running..."
sudo ./esh