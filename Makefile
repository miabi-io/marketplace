BINARY  := marketplace
VERSION ?= dev
IMAGE   ?= ghcr.io/miabi-io/marketplace
PORT    ?= 8088

.PHONY: run serve generate lint build test vet check tidy docker help

run: ## Serve the API + storefront (MARKETPLACE_PORT to override :8088)
	MARKETPLACE_PORT=$(PORT) go run ./cmd/marketplace server

serve: run ## Alias of run

generate: ## Rewrite export.json + registry/index.json (CI runs + diffs this)
	go run ./cmd/marketplace generate

lint: ## Validate every embedded template (catalog drift check)
	go run ./cmd/marketplace lint

build: ## Build the marketplace binary into bin/
	go build -o bin/$(BINARY) ./cmd/marketplace

test: ## Run tests
	go test ./...

vet: ## Static analysis
	go vet ./...
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipping"

check: vet test lint ## Full local CI: vet, test, and lint the catalog

tidy: ## Sync go.mod / go.sum
	go mod tidy

docker: ## Build the Docker image
	docker build -t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

help: ## List targets
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
