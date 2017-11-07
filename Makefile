.DEFAULT_GOAL := help
.PHONY: run run-help test lint check format install-linters release clean help

# Static files directory
STATIC_DIR = src/gui/static

# Electron files directory
ELECTRON_DIR = electron

run:  ## Run the skycoin node. To add arguments, do 'make ARGS="--foo" run'.
	go run cmd/skycoin/skycoin.go --gui-dir="./${STATIC_DIR}" ${ARGS}

run-help: ## Show skycoin node help
	@go run cmd/skycoin/skycoin.go --help

test: ## Run tests
	go test ./cmd/... -timeout=1m
	go test ./src/... -timeout=1m

lint: ## Run linters. requires vendorcheck, gometalinter
	vendorcheck ./...
	gometalinter --enable-gc --warn-unmatched-nolint --disable-all -E goimports -E unparam --tests --vendor ./...

check: lint test ## Run tests and linters

install-linters: ## Install linters
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --vendored-linters
	go get -u github.com/FiloSottile/vendorcheck

format:  # Formats the code. Must have goimports installed (use make install-linters).
	goimports -w ./cmd/...
	goimports -w ./src/...

release: ## Build electron apps, the builds are located in electron/release folder.
	cd $(ELECTRON_DIR) && ./build.sh
	@echo release files are in the folder of electron/release

clean: ## Clean dist files and delete all builds in electron/release
	rm $(ELECTRON_DIR)/release/*

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
