	go test -v -cover $1 -coverprofile=/tmp/coverage.out
	go tool cover -html=/tmp/coverage.out -o /tmp/coverage.html