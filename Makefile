#!/usr/bin/make -f

export GO111MODULE = on

###############################################################################
###                                   All                                   ###
###############################################################################

all: lint test-unit install

###############################################################################
###                                Build flags                              ###
###############################################################################

LD_FLAGS =

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'


###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building example binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/example.exe ./example
else
	@echo "building example binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/example ./example
endif
.PHONY: build

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

clean:
	rm -rf $(BUILDDIR)/

.PHONY: go-mod-cache go.sum clean

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

stop-test-db:
	@echo "Stopping Docker container..."
	@docker stop test-db || true && docker rm test-db || true
.PHONY: stop-docker-test

start-test-db: stop-test-db
	@echo "Starting Docker container..."
	@docker run --name test-db --rm -e POSTGRES_USER=milkyway -e POSTGRES_PASSWORD=password -e POSTGRES_DB=milkyway -d -p 6432:5432 postgres
	@sleep 5

.PHONY: start-test-db

coverage:
	@echo "viewing test coverage..."
	@go tool cover --html=coverage.out
.PHONY: coverage

test-unit:
	@echo "Executing unit tests..."
	@go test -mod=readonly -v -coverprofile coverage.txt ./...
.PHONY: test-unit

###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	@echo "--> Running linter"
	@go run $(golangci_lint_cmd) run --timeout=10m

lint-fix:
	@echo "--> Running linter"
	@go run $(golangci_lint_cmd) run --fix --out-format=tab --issues-exit-code=0

.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' -not -path "./venv" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' -not -path "./venv" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' -not -path "./venv" | xargs goimports -w -local github.com/milkyway-labs/flux
.PHONY: format

.PHONY: lint lint-fix format

