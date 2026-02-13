.PHONY: help build test lint coverage clean install fmt vet

BINARY_NAME=gosh
VERSION ?= 0.1.1
LDFLAGS=-ldflags "-X github.com/gosh/pkg/version.Version=$(VERSION)"
GO=go

help:
	@echo "gosh - HTTPie CLI alternative"
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build the gosh binary"
	@echo "  make install        - Install gosh to \$$GOBIN"
	@echo "  make test           - Run all tests"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make coverage       - Generate coverage report"
	@echo "  make coverage-html  - Generate HTML coverage report"
	@echo "  make lint           - Run linters"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make help           - Show this help message"

build:
	$(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/gosh

install:
	$(GO) install $(LDFLAGS) ./cmd/gosh

test:
	$(GO) test -race -timeout=5m ./...

test-verbose:
	$(GO) test -v -race -timeout=5m ./...

coverage:
	$(GO) test -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -func=coverage.out

coverage-html: coverage
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Installing..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run --timeout=5m

fmt:
	$(GO) fmt ./...
	gofmt -s -w .

vet:
	$(GO) vet ./...

clean:
	$(GO) clean
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f $(BINARY_NAME)

bench:
	$(GO) test -bench=. -benchmem ./...

deps:
	$(GO) mod download
	$(GO) mod tidy

docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

docker-run:
	docker run --rm $(BINARY_NAME):$(VERSION) get https://api.example.com/users

.PHONY: all
all: clean lint test coverage build

.DEFAULT_GOAL := help
