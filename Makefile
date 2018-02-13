.DEFAULT_GOAL := help
.PHONY: run run-help test lint check cover integration-test-stable integration-test-live install-linters format release clean help

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
    							      -not -path '*/testdata*')

# Compilation output
BUILD_DIR = dist
BUILDLIB_DIR = $(BUILD_DIR)/skycoinlib
LIB_DIR = lib

run:  ## Run the skycoin node. To add arguments, do 'make ARGS="--foo" run'.
	go run cmd/skycoin/skycoin.go --gui-dir="./${STATIC_DIR}" ${ARGS}

run-help: ## Show skycoin node help
	@go run cmd/skycoin/skycoin.go --help

test: ## Run tests
	go test ./cmd/... -timeout=1m
	go test ./src/... -timeout=1m

build-lib-c: # Build Skycoinlib C
	mkdir -p $(BUILDLIB_DIR)
	rm -Rf $(BUILDLIB_DIR)/*
	go build -buildmode=c-shared  -o $(BUILDLIB_DIR)/libskycoin.so $(LIB_DIR)/cgo/main.go
	go build -buildmode=c-archive -o $(BUILDLIB_DIR)/libskycoin.a  $(LIB_DIR)/cgo/main.go

test-lib-c: build-lib-c
	cp $(LIB_DIR)/cgo/tests/*.c $(BUILDLIB_DIR)/
	rm $(BUILDLIB_DIR)/libskycoin.so	# TODO: Get rid of this step
	gcc -o $(BUILDLIB_DIR)/skycoinlib_test $(BUILDLIB_DIR)/*.c -I$(BUILDLIB_DIR) -lcriterion -lskycoin -L $(BUILDLIB_DIR)
	$(BUILDLIB_DIR)/skycoinlib_test

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	gometalinter --disable-all -E vet -E goimports -E varcheck --tests --vendor ./...

check: lint test ## Run tests and linters

integration-test-stable: ## Run stable integration tests
	./ci-scripts/integration-test-stable.sh

integration-test-live: ## Run live integration tests
	./ci-scripts/integration-test-live.sh

cover: ## Runs tests on ./src/ with HTML code coverage
	@echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	go get -u github.com/alecthomas/gometalinter
	gometalinter --vendored-linters --install

format:  # Formats the code. Must have goimports installed (use make install-linters).
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
