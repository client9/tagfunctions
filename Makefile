
test:
	go test .
cover:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
