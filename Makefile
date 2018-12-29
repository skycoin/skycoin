.DEFAULT_GOAL := help
.PHONY: run run-help test test-core test-libc test-lint build-libc check
.PHONY: integration-test-stable integration-test-stable-disable-csrf
.PHONY: integration-test-live integration-test-live-wallet
.PHONY: integration-test-disable-wallet-api integration-test-disable-seed-api
.PHONY: integration-test-enable-seed-api integration-test-enable-seed-api
.PHONY: integration-test-disable-gui integration-test-disable-gui
.PHONY: integration-test-db-no-unconfirmed integration-test-auth
.PHONY: install-linters format release clean-release clean-coverage
.PHONY: install-deps-ui build-ui help newcoins merge-coverage
.PHONY: generate-mocks update-golden-files

COIN ?= skycoin

# Static files directory
GUI_STATIC_DIR = src/gui/static

# Electron files directory
ELECTRON_DIR = electron

# Compilation output for libskycoin
BUILD_DIR = build
BUILDLIB_DIR = $(BUILD_DIR)/libskycoin
LIB_DIR = lib
LIB_FILES = $(shell find ./lib/cgo -type f -name "*.go")
SRC_FILES = $(shell find ./src -type f -name "*.go")
HEADER_FILES = $(shell find ./include -type f -name "*.h")
BIN_DIR = bin
DOC_DIR = docs
INCLUDE_DIR = include
LIBSRC_DIR = lib/cgo
LIBDOC_DIR = $(DOC_DIR)/libc

# Compilation flags for libskycoin
CC_VERSION = $(shell $(CC) -dumpversion)
STDC_FLAG = $(python -c "if tuple(map(int, '$(CC_VERSION)'.split('.'))) < (6,): print('-std=C99'")
LIBC_LIBS = -lcriterion
LIBC_FLAGS = -I$(LIBSRC_DIR) -I$(INCLUDE_DIR) -I$(BUILD_DIR)/usr/include -L $(BUILDLIB_DIR) -L$(BUILD_DIR)/usr/lib

# Platform specific checks
OSNAME = $(TRAVIS_OS_NAME)

ifeq ($(shell uname -s),Linux)
  LDLIBS=$(LIBC_LIBS) -lpthread
  LDPATH=$(shell printenv LD_LIBRARY_PATH)
  LDPATHVAR=LD_LIBRARY_PATH
  LDFLAGS=$(LIBC_FLAGS) $(STDC_FLAG)
ifndef OSNAME
  OSNAME = linux
endif
else ifeq ($(shell uname -s),Darwin)
ifndef OSNAME
  OSNAME = osx
endif
  LDLIBS = $(LIBC_LIBS)
  LDPATH=$(shell printenv DYLD_LIBRARY_PATH)
  LDPATHVAR=DYLD_LIBRARY_PATH
  LDFLAGS=$(LIBC_FLAGS) -framework CoreFoundation -framework Security
else
  LDLIBS = $(LIBC_LIBS)
  LDPATH=$(shell printenv LD_LIBRARY_PATH)
  LDPATHVAR=LD_LIBRARY_PATH
  LDFLAGS=$(LIBC_FLAGS)
endif

run-client:  ## Run skycoin with desktop client configuration. To add arguments, do 'make ARGS="--foo" run'.
	./run-client.sh ${ARGS}

run-daemon:  ## Run skycoin with server daemon configuration. To add arguments, do 'make ARGS="--foo" run'.
	./run-daemon.sh ${ARGS}

run-help: ## Show skycoin node help
	@go run cmd/$(COIN)/$(COIN).go --help

run-integration-test-live: ## Run the skycoin node configured for live integration tests
	./ci-scripts/run-live-integration-test-node.sh

run-integration-test-live-cover: ## Run the skycoin node configured for live integration tests with coverage
	./ci-scripts/run-live-integration-test-node-cover.sh

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

configure-build:
	mkdir -p $(BUILD_DIR)/usr/tmp $(BUILD_DIR)/usr/lib $(BUILD_DIR)/usr/include
	mkdir -p $(BUILDLIB_DIR) $(BIN_DIR) $(INCLUDE_DIR)

$(BUILDLIB_DIR)/libskycoin.so: $(LIB_FILES) $(SRC_FILES) $(HEADER_FILES)
	rm -Rf $(BUILDLIB_DIR)/libskycoin.so
	go build -buildmode=c-shared  -o $(BUILDLIB_DIR)/libskycoin.so $(LIB_FILES)
	mv $(BUILDLIB_DIR)/libskycoin.h $(INCLUDE_DIR)/

$(BUILDLIB_DIR)/libskycoin.a: $(LIB_FILES) $(SRC_FILES) $(HEADER_FILES)
	rm -Rf $(BUILDLIB_DIR)/libskycoin.a
	go build -buildmode=c-archive -o $(BUILDLIB_DIR)/libskycoin.a  $(LIB_FILES)
	mv $(BUILDLIB_DIR)/libskycoin.h $(INCLUDE_DIR)/

## Build libskycoin C static library
build-libc-static: $(BUILDLIB_DIR)/libskycoin.a

## Build libskycoin C shared library
build-libc-shared: $(BUILDLIB_DIR)/libskycoin.so

## Build libskycoin C client libraries
build-libc: configure-build build-libc-static build-libc-shared

## Build libskycoin C client library and executable C test suites
## with debug symbols. Use this target to debug the source code
## with the help of an IDE
build-libc-dbg: configure-build build-libc-static build-libc-shared
	$(CC) -g -o $(BIN_DIR)/test_libskycoin_shared $(LIB_DIR)/cgo/tests/*.c -lskycoin                    $(LDLIBS) $(LDFLAGS)
	$(CC) -g -o $(BIN_DIR)/test_libskycoin_static $(LIB_DIR)/cgo/tests/*.c $(BUILDLIB_DIR)/libskycoin.a $(LDLIBS) $(LDFLAGS)

test-libc: build-libc ## Run tests for libskycoin C client library
	echo "Compiling with $(CC) $(CC_VERSION) $(STDC_FLAG)"
	$(CC) -o $(BIN_DIR)/test_libskycoin_shared $(LIB_DIR)/cgo/tests/*.c $(LIB_DIR)/cgo/tests/testutils/*.c -lskycoin                    $(LDLIBS) $(LDFLAGS)
	$(CC) -o $(BIN_DIR)/test_libskycoin_static $(LIB_DIR)/cgo/tests/*.c $(LIB_DIR)/cgo/tests/testutils/*.c $(BUILDLIB_DIR)/libskycoin.a $(LDLIBS) $(LDFLAGS)
	$(LDPATHVAR)="$(LDPATH):$(BUILD_DIR)/usr/lib:$(BUILDLIB_DIR)" $(BIN_DIR)/test_libskycoin_shared
	$(LDPATHVAR)="$(LDPATH):$(BUILD_DIR)/usr/lib"                 $(BIN_DIR)/test_libskycoin_static

docs-libc:
	doxygen ./.Doxyfile
	moxygen -o $(LIBDOC_DIR)/API.md $(LIBDOC_DIR)/xml/

docs: docs-libc

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	golangci-lint run -c .golangci.yml ./...
	# lib/cgo needs separate linting rules
	golangci-lint run -c .golangci.libcgo.yml ./lib/cgo/...
	# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

check-newcoin: newcoin ## Check that make newcoin succeeds and no files are changed.
	if [ "$(shell git diff ./ | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after make newcoin' ; exit 2 ; fi

check: lint clean-coverage test integration-test-stable integration-test-stable-disable-csrf \
	integration-test-disable-wallet-api integration-test-disable-seed-api \
	integration-test-enable-seed-api integration-test-disable-gui \
	integration-test-auth integration-test-db-no-unconfirmed check-newcoin ## Run tests and linters

integration-test-stable: ## Run stable integration tests
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -c -n enable-csrf

integration-test-stable-disable-csrf: ## Run stable integration tests with CSRF disabled
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -n disable-csrf

integration-test-live: ## Run live integration tests
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-live.sh -c

integration-test-live-wallet: ## Run live integration tests with wallet
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-live.sh -w

integration-test-live-disable-csrf: ## Run live integration tests against a node with CSRF disabled
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-live.sh

integration-test-live-disable-networking: ## Run live integration tests against a node with networking disabled (requires wallet)
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-live.sh -c -k

integration-test-disable-wallet-api: ## Run disable wallet api integration tests
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-disable-wallet-api.sh

integration-test-enable-seed-api: ## Run enable seed api integration test
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-enable-seed-api.sh

integration-test-disable-gui: ## Run tests with the GUI disabled
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-disable-gui.sh

integration-test-db-no-unconfirmed: ## Run stable tests against the stable database that has no unconfirmed transactions
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -d -n no-unconfirmed

integration-test-auth: ## Run stable tests with HTTP Basic auth enabled
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-auth.sh

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	# For some reason this install method is not recommended, see https://github.com/golangci/golangci-lint#install
	# However, they suggest `curl ... | bash` which we should not do
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

install-deps-libc: configure-build ## Install locally dependencies for testing libskycoin
	git clone --recursive https://github.com/skycoin/Criterion $(BUILD_DIR)/usr/tmp/Criterion
	mkdir $(BUILD_DIR)/usr/tmp/Criterion/build
	cd    $(BUILD_DIR)/usr/tmp/Criterion/build && cmake .. && cmake --build .
	mv    $(BUILD_DIR)/usr/tmp/Criterion/build/libcriterion.* $(BUILD_DIR)/usr/lib/
	cp -R $(BUILD_DIR)/usr/tmp/Criterion/include/* $(BUILD_DIR)/usr/include/

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/skycoin ./cmd
	goimports -w -local github.com/skycoin/skycoin ./src
	goimports -w -local github.com/skycoin/skycoin ./lib

install-deps-ui:  ## Install the UI dependencies
	cd $(GUI_STATIC_DIR) && npm install

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

generate-mocks: ## Regenerate test interface mocks
	go generate ./src/...
	# mockery can't generate the UnspentPooler mock in package visor, patch it
	mv ./src/visor/blockdb/mock_unspent_pooler_test.go ./src/visor/mock_unspent_pooler_test.go
	sed -i "" -e 's/package blockdb/package visor/g' ./src/visor/mock_unspent_pooler_test.go

update-golden-files: ## Run integration tests in update mode
	./ci-scripts/integration-test-stable.sh -u >/dev/null 2>&1 || true
	./ci-scripts/integration-test-stable.sh -c -u >/dev/null 2>&1 || true
	./ci-scripts/integration-test-stable.sh -d -u >/dev/null 2>&1 || true
	./ci-scripts/integration-test-stable.sh -c -d -u >/dev/null 2>&1 || true

merge-coverage: ## Merge coverage files and create HTML coverage output. gocovmerge is required, install with `go get github.com/wadey/gocovmerge`
	@echo "To install gocovmerge do:"
	@echo "go get github.com/wadey/gocovmerge"
	gocovmerge coverage/*.coverage.out > coverage/all-coverage.merged.out
	go tool cover -html coverage/all-coverage.merged.out -o coverage/all-coverage.html
	@echo "Total coverage HTML file generated at coverage/all-coverage.html"
	@echo "Open coverage/all-coverage.html in your browser to view"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
