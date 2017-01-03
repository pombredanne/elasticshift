go build -o=bin/esh ./cmd/esh/main.go
go test -coverprofile cover.out -v
go tool cover -html=cover.out -o cover.html
