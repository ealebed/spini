VERSION=$(shell date -u '+%y.%m.%d-%H.%M')

GO            := go
GO111MODULE   := on
CGO_ENABLED   := 0
GOBUILD       := CGO_ENABLED=$(CGO_ENABLED) GO111MODULE=$(GO111MODULE) $(GO) build
GOTEST        := $(GO) test -gcflags='-l' -p 3
GOLANGCI_LINT := golangci-lint
BIN           := bin/spini

.PHONY: all
all: build

.PHONY: build
build:
	$(GOBUILD) -o $(BIN) ./

.PHONY: install
install:
	$(GO) install ./

.PHONY: image
image:
	docker build -t ealebed/spini:${VERSION} .
	docker push ealebed/spini:${VERSION}

.PHONY: test
test:
	$(GOTEST) ./...

.PHONY: test-race
test-race:
	$(GO) test ./... -race

.PHONY: fmt
fmt:
	$(GO) fmt ./...
	gofmt -s -w .

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run --timeout 4m --config .golangci.yaml

.PHONY: clean
clean:
	rm -f $(BIN)
	rm -rf bin/
	rm -f coverage.out
	rm -f *.test

.PHONY: update
update:
	$(GO) get -u -v ./...
	$(GO) mod verify
	$(GO) mod tidy

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - build the project (default)"
	@echo "  build      - build the project"
	@echo "  install    - install the project"
	@echo "  image      - build and push Docker image"
	@echo "  test       - run tests"
	@echo "  test-race  - run tests with race detector"
	@echo "  fmt        - format code"
	@echo "  lint       - run linter"
	@echo "  clean      - clean build artifacts"
	@echo "  update     - update dependencies"
	@echo "  help       - show this help message"
