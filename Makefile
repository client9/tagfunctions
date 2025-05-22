
test:
	go test .
clean:
	go clean
	go mod tidy
	rm -f coverage.out

cover:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
