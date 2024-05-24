
test:
	go test .
clean:
	go clean
	go mod tidy

cover:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
