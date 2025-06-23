
test:
	go test .
clean:
	go clean
	go mod tidy
	rm -f coverage.out
lint:
	go mod tidy
	gofmt -w -s *.go
	golangci-lint run .

cover:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
