.DEFAULT_GOAL := help
.PHONY: run run-help test test-core test-libc test-lint build-libc check cover integration-test-stable integration-test-live install-linters format release clean help

# Static files directory
STATIC_DIR = src/gui/static

# Electron files directory
ELECTRON_DIR = electron

# ./src folder does not have code
# ./src/api folder does not have code
# ./src/util folder does not have code
# ./src/ciper/* are libraries manually vendored by cipher that do not need coverage
# ./src/gui/static* are static assets
# */testdata* folders do not have code
# ./src/consensus/example has no buildable code
PACKAGES = $(shell find ./src -type d -not -path '\./src' \
    							      -not -path '\./src/api' \
    							      -not -path '\./src/util' \
    							      -not -path '\./src/consensus/example' \
    							      -not -path '\./src/gui/static*' \
    							      -not -path '\./src/cipher/*' \
    							      -not -path '*/testdata*' \
    							      -not -path '*/test-fixtures*')

# Compilation output
BUILD_DIR = build
BUILDLIB_DIR = $(BUILD_DIR)/libskycoin
LIB_DIR = lib
LIB_FILES = $(shell find ./lib/cgo -type f -name "*.go")
BIN_DIR = bin
INCLUDE_DIR = include
LIBSRC_DIR = lib/cgo

# Compilation flags
CC = gcc
LIBC_LIBS = -lcriterion
LIBC_FLAGS = -I$(LIBSRC_DIR) -I$(INCLUDE_DIR) -I$(BUILD_DIR)/usr/include -L $(BUILDLIB_DIR) -L$(BUILD_DIR)/usr/lib

# Platform specific checks
OSNAME = $(TRAVIS_OS_NAME)

ifeq ($(shell uname -s),Linux)
  LDLIBS=$(LIBC_LIBS) -lpthread
	LDPATH=$(shell printenv LD_LIBRARY_PATH)
	LDPATHVAR=LD_LIBRARY_PATH
	LDFLAGS=$(LIBC_FLAGS)
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

run:  ## Run the skycoin node. To add arguments, do 'make ARGS="--foo" run'.
	go run cmd/skycoin/skycoin.go --gui-dir="./${STATIC_DIR}" ${ARGS}

run-help: ## Show skycoin node help
	@go run cmd/skycoin/skycoin.go --help

test: ## Run tests for Skycoin
	go test ./cmd/... -timeout=1m
	go test ./src/... -timeout=1m

test-386: ## Run tests for Skycoin with GOARCH=386
	GOARCH=386 go test ./cmd/... -timeout=5m
	GOARCH=386 go test ./src/... -timeout=5m

test-amd64: ## Run tests for Skycoin with GOARCH=amd64
	GOARCH=amd64 go test ./cmd/... -timeout=5m
	GOARCH=amd64 go test ./src/... -timeout=5m

configure-build:
	mkdir -p $(BUILD_DIR)/usr/tmp $(BUILD_DIR)/usr/lib $(BUILD_DIR)/usr/include
	mkdir -p $(BUILDLIB_DIR) $(BIN_DIR) $(INCLUDE_DIR)

build-libc: configure-build ## Build libskycoin C client library
	rm -Rf $(BUILDLIB_DIR)/*
	go build -buildmode=c-shared  -o $(BUILDLIB_DIR)/libskycoin.so $(LIB_FILES)
	go build -buildmode=c-archive -o $(BUILDLIB_DIR)/libskycoin.a  $(LIB_FILES)
	mv $(BUILDLIB_DIR)/libskycoin.h $(INCLUDE_DIR)/

test-libc: build-libc ## Run tests for libskycoin C client library
	cp $(LIB_DIR)/cgo/tests/*.c $(BUILDLIB_DIR)/
	$(CC) -o $(BIN_DIR)/test_libskycoin_shared $(BUILDLIB_DIR)/*.c -lskycoin                    $(LDLIBS) $(LDFLAGS)
	$(CC) -o $(BIN_DIR)/test_libskycoin_static $(BUILDLIB_DIR)/*.c $(BUILDLIB_DIR)/libskycoin.a $(LDLIBS) $(LDFLAGS)
	$(LDPATHVAR)="$(LDPATH):$(BUILD_DIR)/usr/lib:$(BUILDLIB_DIR)" $(BIN_DIR)/test_libskycoin_shared
	$(LDPATHVAR)="$(LDPATH):$(BUILD_DIR)/usr/lib"                 $(BIN_DIR)/test_libskycoin_static

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	gometalinter --disable-all -E vet -E goimports -E varcheck --tests --vendor ./...

check: lint test integration-test-stable ## Run tests and linters

integration-test-stable: ## Run stable integration tests
	./ci-scripts/integration-test-stable.sh -v -w

integration-test-live: ## Run live integration tests
	./ci-scripts/integration-test-live.sh -v -w

integration-test-disable-wallet-api: ## Run disable wallet api integration tests
	./ci-scripts/integration-test-disable-wallet-api.sh -v

cover: ## Runs tests on ./src/ with HTML code coverage
	go test -cover -coverprofile=cover.out -coverpkg=github.com/skycoin/skycoin/... ./src/...
	go tool cover -html=cover.out

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	go get -u github.com/alecthomas/gometalinter
	gometalinter --vendored-linters --install

install-deps-libc: configure-build ## Install locally dependencies for testing libskycoin
	wget -O $(BUILD_DIR)/usr/tmp/criterion-v2.3.2-$(OSNAME)-x86_64.tar.bz2 https://github.com/Snaipe/Criterion/releases/download/v2.3.2/criterion-v2.3.2-$(OSNAME)-x86_64.tar.bz2
	tar -x -C $(BUILD_DIR)/usr/tmp/ -j -f $(BUILD_DIR)/usr/tmp/criterion-v2.3.2-$(OSNAME)-x86_64.tar.bz2
	ls $(BUILD_DIR)/usr/tmp/criterion-v2.3.2/include
	ls -1 $(BUILD_DIR)/usr/tmp/criterion-v2.3.2/lib     | xargs -I NAME mv $(BUILD_DIR)/usr/tmp/criterion-v2.3.2/lib/NAME     $(BUILD_DIR)/usr/lib/NAME
	ls -1 $(BUILD_DIR)/usr/tmp/criterion-v2.3.2/include | xargs -I NAME mv $(BUILD_DIR)/usr/tmp/criterion-v2.3.2/include/NAME $(BUILD_DIR)/usr/include/NAME

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/skycoin ./cmd
	goimports -w -local github.com/skycoin/skycoin ./src
	goimports -w -local github.com/skycoin/skycoin ./lib

release: ## Build electron apps, the builds are located in electron/release folder.
	cd $(ELECTRON_DIR) && ./build.sh
	@echo release files are in the folder of electron/release

clean: ## Clean dist files and delete all builds in electron/release
	rm $(ELECTRON_DIR)/release/*

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
