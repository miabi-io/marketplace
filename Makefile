BINARY     := marketplace
PKG        := github.com/miabi-io/marketplace
VERSION    ?= dev
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -X $(PKG)/internal/buildinfo.Version=$(VERSION) -X $(PKG)/internal/buildinfo.CommitID=$(COMMIT) -X $(PKG)/internal/buildinfo.BuildDate=$(BUILD_DATE)
IMAGE      ?= ghcr.io/miabi-io/marketplace
PORT       ?= 8088

WEB_DIR       := web
EMBED_WEB_DIR := internal/web/dist

.PHONY: run serve generate lint build build-ui build-all dev-ui test vet check tidy docker help

run: ## Serve the API + storefront (MARKETPLACE_PORT to override :8088)
	MARKETPLACE_PORT=$(PORT) go run -ldflags "$(LDFLAGS)" ./cmd/marketplace server

serve: run ## Alias of run

generate: ## Rewrite export.json + registry/index.json (CI runs + diffs this)
	go run ./cmd/marketplace generate

lint: ## Validate every embedded template (catalog drift check)
	go run ./cmd/marketplace lint

build: ## Build the marketplace binary into bin/
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/marketplace

build-ui: ## Build the storefront (Vue) and stage it for embedding
	npm --prefix $(WEB_DIR) ci
	npm --prefix $(WEB_DIR) run build
	# Stage the build output where `go build` embeds it, keeping the committed
	# .gitkeep so `go build` still works on a clean tree.
	rm -rf $(EMBED_WEB_DIR)
	cp -r $(WEB_DIR)/dist $(EMBED_WEB_DIR)
	touch $(EMBED_WEB_DIR)/.gitkeep

build-all: build-ui build ## Build the storefront and the server binary

dev-ui: ## Run the Vite dev server (proxies /v1 -> :8088)
	npm --prefix $(WEB_DIR) run dev

test: ## Run tests
	go test ./...

vet: ## Static analysis
	go vet ./...
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipping"

check: vet test lint ## Full local CI: vet, test, and lint the catalog

tidy: ## Sync go.mod / go.sum
	go mod tidy

docker: ## Build the Docker image
	docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) --build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

help: ## List targets
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
