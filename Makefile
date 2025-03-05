MAKEFLAGS = --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c

SERVICE_NAME = guppy
BUILD_DIR = dist
SEED := on

.DEFAULT_GOAL := help
.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-38s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: gen
gen: ## Generate mock files
	go generate ./...

.PHONY: build
build: ## Build the project
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(SERVICE_NAME) cmd/$(SERVICE_NAME)/*.go

.PHONY: run
run: build ## Run the project
	$(BUILD_DIR)/$(SERVICE_NAME)

.PHONY: clean
clean: ## Clean up build artifacts
	rm -rf $(BUILD_FOLDER)

.PHONY: lint-golang
lint-golang: ## Lint Golang source code
	golangci-lint run
	go run golang.org/x/tools/cmd/deadcode@latest -test ./... | tee deadcode.out && [ ! -s deadcode.out ]

.PHONY: lint
lint: lint-golang ## Lint project source

.PHONY: fmt
fmt: ## Format the source code
	go run mvdan.cc/gofumpt@latest -l -w -extra .
	go run golang.org/x/tools/cmd/goimports@latest -l -w .
	go run github.com/daixiang0/gci@latest write \
		--skip-generated \
		--custom-order \
		-s standard \
		-s default \
		-s prefix\(github.com/alkurbatov/guppy\) \
		-s blank \
		-s dot \
		.

.PHONY: test
test: ## Run unit tests
	go test -v -race -shuffle=$(SEED) ./internal/... -coverprofile=coverage.out -covermode atomic
	@grep -v -E "(_mock|.pb).go" coverage.out > coverage.out.tmp
	@mv coverage.out.tmp coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
