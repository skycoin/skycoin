.DEFAULT_GOAL := help
.PHONY: run-client run-daemon run-help
.PHONY: test test-386 test-amd64
.PHONY: check check-newcoin
.PHONY: run-integration-test-live
.PHONY: run-integration-test-live-disable-csrf
.PHONY: run-integration-test-live-disable-networking
.PHONY: run-integration-test-live-cover
.PHONY: run-integration-test-live-cover-disable-csrf
.PHONY: run-integration-test-live-cover-disable-networking
.PHONY: integration-tests-stable
.PHONY: integration-test-stable
.PHONY: integration-test-stable-disable-csrf
.PHONY: integration-test-stable-disable-wallet-api
.PHONY: integration-test-stable-enable-seed-api
.PHONY: integration-test-stable-disable-gui
.PHONY: integration-test-stable-db-no-unconfirmed
.PHONY: integration-test-stable-auth
.PHONY: integration-test-live integration-test-live-wallet
.PHONY: install-linters format release clean-release clean-coverage
.PHONY: install-deps-ui build-ui build-ui-travis help newcoin merge-coverage
.PHONY: generate update-golden-files
.PHONY: fuzz-base58 fuzz-encoder

COIN ?= skycoin

# Static files directory
GUI_STATIC_DIR = src/gui/static

# Electron files directory
ELECTRON_DIR = electron

# Platform specific checks
OSNAME = $(TRAVIS_OS_NAME)

run-client:  ## Run skycoin with desktop client configuration. To add arguments, do 'make ARGS="--foo" run'.
	./run-client.sh ${ARGS}

run-daemon:  ## Run skycoin with server daemon configuration. To add arguments, do 'make ARGS="--foo" run'.
	./run-daemon.sh ${ARGS}

run-help: ## Show skycoin node help
	@go run cmd/$(COIN)/$(COIN).go --help

run-integration-test-live: ## Run the skycoin node configured for live integration tests
	./ci-scripts/run-live-integration-test-node.sh

run-integration-test-live-disable-csrf: ## Run the skycoin node configured for live integration tests with CSRF disabled
	./ci-scripts/run-live-integration-test-node.sh -disable-csrf

run-integration-test-live-disable-networking: ## Run the skycoin node configured for live integration tests with networking disabled
	./ci-scripts/run-live-integration-test-node.sh -disable-networking

run-integration-test-live-cover: ## Run the skycoin node configured for live integration tests with coverage
	./ci-scripts/run-live-integration-test-node-cover.sh

run-integration-test-live-cover-disable-csrf: ## Run the skycoin node configured for live integration tests with CSRF disabled and with coverage
	./ci-scripts/run-live-integration-test-node-cover.sh -disable-csrf

run-integration-test-live-cover-disable-networking: ## Run the skycoin node configured for live integration tests with networking disabled and with coverage
	./ci-scripts/run-live-integration-test-node-cover.sh -disable-networking

test: ## Run tests for Skycoin
	@mkdir -p coverage/
	COIN=$(COIN) go test -coverpkg="github.com/$(COIN)/$(COIN)/..." -coverprofile=coverage/go-test-cmd.coverage.out -timeout=5m ./cmd/...
	COIN=$(COIN) go test -coverpkg="github.com/$(COIN)/$(COIN)/..." -coverprofile=coverage/go-test-src.coverage.out -timeout=5m ./src/...

test-386: ## Run tests for Skycoin with GOARCH=386
	GOARCH=386 COIN=$(COIN) go test ./cmd/... -timeout=5m
	GOARCH=386 COIN=$(COIN) go test ./src/... -timeout=5m

test-amd64: ## Run tests for Skycoin with GOARCH=amd64
	GOARCH=amd64 COIN=$(COIN) go test ./cmd/... -timeout=5m
	GOARCH=amd64 COIN=$(COIN) go test ./src/... -timeout=5m

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	golangci-lint run -c .golangci.yml ./...
	@# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

check-newcoin: newcoin ## Check that make newcoin succeeds and no templated files are changed.
	@if [ "$(shell git diff ./cmd/skycoin/skycoin.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after make newcoin' ; exit 2 ; fi
	@if [ "$(shell git diff ./cmd/skycoin/skycoin_test.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after make newcoin' ; exit 2 ; fi
	@if [ "$(shell git diff ./src/params/params.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after make newcoin' ; exit 2 ; fi

check: lint clean-coverage test test-386 integration-tests-stable check-newcoin ## Run tests and linters

integration-tests-stable: integration-test-stable \
	integration-test-stable-disable-csrf \
	integration-test-stable-disable-wallet-api \
	integration-test-stable-enable-seed-api \
	integration-test-stable-disable-gui \
	integration-test-stable-auth \
	integration-test-stable-db-no-unconfirmed ## Run all stable integration tests

integration-test-stable: ## Run stable integration tests
	COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -c -x -n enable-csrf-header-check

integration-test-stable-disable-header-check: ## Run stable integration tests with header check disabled
	COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -n disable-header-check

integration-test-stable-disable-csrf: ## Run stable integration tests with CSRF disabled
	COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -n disable-csrf

integration-test-stable-disable-wallet-api: ## Run disable wallet api integration tests
	COIN=$(COIN) ./ci-scripts/integration-test-disable-wallet-api.sh

integration-test-stable-enable-seed-api: ## Run enable seed api integration test
	COIN=$(COIN) ./ci-scripts/integration-test-enable-seed-api.sh

integration-test-stable-disable-gui: ## Run tests with the GUI disabled
	COIN=$(COIN) ./ci-scripts/integration-test-disable-gui.sh

integration-test-stable-db-no-unconfirmed: ## Run stable tests against the stable database that has no unconfirmed transactions
	COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -d -n no-unconfirmed

integration-test-stable-auth: ## Run stable tests with HTTP Basic auth enabled
	COIN=$(COIN) ./ci-scripts/integration-test-auth.sh

integration-test-live: ## Run live integration tests
	COIN=$(COIN) ./ci-scripts/integration-test-live.sh -c

integration-test-live-wallet: ## Run live integration tests with wallet
	COIN=$(COIN) ./ci-scripts/integration-test-live.sh -w

integration-test-live-enable-header-check: ## Run live integration tests against a node with header check enabled
	COIN=$(COIN) ./ci-scripts/integration-test-live.sh

integration-test-live-disable-csrf: ## Run live integration tests against a node with CSRF disabled
	COIN=$(COIN) ./ci-scripts/integration-test-live.sh

integration-test-live-disable-networking: ## Run live integration tests against a node with networking disabled (requires wallet)
	COIN=$(COIN) ./ci-scripts/integration-test-live.sh -c -k

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	# For some reason this install method is not recommended, see https://github.com/golangci/golangci-lint#install
	# However, they suggest `curl ... | bash` which we should not do
	# go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	# Change to use go get -u with version when go is v1.12+
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.18.0

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/SkycoinProject/skycoin ./cmd
	goimports -w -local github.com/SkycoinProject/skycoin ./src

install-deps-ui:  ## Install the UI dependencies
	cd $(GUI_STATIC_DIR) && npm ci

lint-ui:  ## Lint the UI code
	cd $(GUI_STATIC_DIR) && npm run lint

test-ui:  ## Run UI tests
	cd $(GUI_STATIC_DIR) && npm run test

test-ui-e2e:  ## Run UI e2e tests
	./ci-scripts/ui-e2e.sh

build-ui:  ## Builds the UI
	cd $(GUI_STATIC_DIR) && npm run build

build-ui-travis:  ## Builds the UI for travis
	cd $(GUI_STATIC_DIR) && npm run build-travis

release: ## Build electron, standalone and daemon apps. Use osarch=${osarch} to specify the platform. Example: 'make release osarch=darwin/amd64', multiple platform can be supported in this way: 'make release osarch="darwin/amd64 windows/amd64"'. Supported architectures are: darwin/amd64 windows/amd64 windows/386 linux/amd64 linux/arm, the builds are located in electron/release folder.
	cd $(ELECTRON_DIR) && ./build.sh ${osarch}
	@echo release files are in the folder of electron/release

release-standalone: ## Build standalone apps. Use osarch=${osarch} to specify the platform. Example: 'make release-standalone osarch=darwin/amd64' Supported architectures are the same as 'release' command.
	cd $(ELECTRON_DIR) && ./build-standalone-release.sh ${osarch}
	@echo release files are in the folder of electron/release

release-electron: ## Build electron apps. Use osarch=${osarch} to specify the platform. Example: 'make release-electron osarch=darwin/amd64' Supported architectures are the same as 'release' command.
	cd $(ELECTRON_DIR) && ./build-electron-release.sh ${osarch}
	@echo release files are in the folder of electron/release

release-daemon: ## Build daemon apps. Use osarch=${osarch} to specify the platform. Example: 'make release-daemon osarch=darwin/amd64' Supported architectures are the same as 'release' command.
	cd $(ELECTRON_DIR) && ./build-daemon-release.sh ${osarch}
	@echo release files are in the folder of electron/release

release-cli: ## Build CLI apps. Use osarch=${osarch} to specify the platform. Example: 'make release-cli osarch=darwin/amd64' Supported architectures are the same as 'release' command.
	cd $(ELECTRON_DIR) && ./build-cli-release.sh ${osarch}
	@echo release files are in the folder of electron/release

clean-release: ## Remove all electron build artifacts
	rm -rf $(ELECTRON_DIR)/release
	rm -rf $(ELECTRON_DIR)/.gox_output
	rm -rf $(ELECTRON_DIR)/.daemon_output
	rm -rf $(ELECTRON_DIR)/.cli_output
	rm -rf $(ELECTRON_DIR)/.standalone_output
	rm -rf $(ELECTRON_DIR)/.electron_output

clean-coverage: ## Remove coverage output files
	rm -rf ./coverage/

newcoin: ## Rebuild cmd/$COIN/$COIN.go file from the template. Call like "make newcoin COIN=foo".
	go run cmd/newcoin/newcoin.go createcoin --coin $(COIN)

generate: ## Generate test interface mocks and struct encoders
	go generate ./src/...
	# mockery can't generate the UnspentPooler mock in package visor, patch it
	mv ./src/visor/blockdb/mock_unspent_pooler_test.go ./src/visor/mock_unspent_pooler_test.go
	sed -i "" -e 's/package blockdb/package visor/g' ./src/visor/mock_unspent_pooler_test.go
	sed -i "" -e 's/AddressHashes/blockdb.AddressHashes/g' ./src/visor/mock_unspent_pooler_test.go
	goimports -w -local github.com/SkycoinProject/skycoin ./src/visor/mock_unspent_pooler_test.go

install-generators: ## Install tools used by go generate
	go get github.com/vektra/mockery/.../
	go get github.com/SkycoinProject/skyencoder/cmd/skyencoder

update-golden-files: ## Run integration tests in update mode
	./ci-scripts/integration-test-stable.sh -u >/dev/null 2>&1 || true
	./ci-scripts/integration-test-stable.sh -c -x -u >/dev/null 2>&1 || true
	./ci-scripts/integration-test-stable.sh -d -u >/dev/null 2>&1 || true
	./ci-scripts/integration-test-stable.sh -c -x -d -u >/dev/null 2>&1 || true

merge-coverage: ## Merge coverage files and create HTML coverage output. gocovmerge is required, install with `go get github.com/wadey/gocovmerge`
	@echo "To install gocovmerge do:"
	@echo "go get github.com/wadey/gocovmerge"
	gocovmerge coverage/*.coverage.out > coverage/all-coverage.merged.out
	go tool cover -html coverage/all-coverage.merged.out -o coverage/all-coverage.html
	@echo "Total coverage HTML file generated at coverage/all-coverage.html"
	@echo "Open coverage/all-coverage.html in your browser to view"

fuzz-base58: ## Fuzz the base58 package. Requires https://github.com/dvyukov/go-fuzz
	go-fuzz-build github.com/SkycoinProject/skycoin/src/cipher/base58/internal
	go-fuzz -bin=base58fuzz-fuzz.zip -workdir=src/cipher/base58/internal

fuzz-encoder: ## Fuzz the encoder package. Requires https://github.com/dvyukov/go-fuzz
	go-fuzz-build github.com/SkycoinProject/skycoin/src/cipher/encoder/internal
	go-fuzz -bin=encoderfuzz-fuzz.zip -workdir=src/cipher/encoder/internal

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
