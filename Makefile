.PHONY: default tidy lint test

default: test

tidy:
	@echo "Fetching dependencies..."
	go mod tidy

setup:
	@echo "Getting CI dependencies..."
	@test -x ${GOPATH}/bin/golangci-lint || go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0

lint: setup
	@echo "Linting code..."
	golangci-lint -v run

test: tidy
	go test -v -race -short -cover -coverprofile cover.out ./...
	go tool cover -func cover.out
