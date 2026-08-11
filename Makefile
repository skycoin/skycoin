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
.PHONY: integration-test-fibercoin
.PHONY: check-newcoin-templates check-release-nocgo
.PHONY: update-dep sync-upstream-develop
.PHONY: install-linters format release clean-release clean-coverage dep-github-release
.PHONY: docs-install docs-serve docs-build
.PHONY: install-deps-ui build-ui build-ui help newcoin merge-coverage
.PHONY: install-deps-skydex-ui build-skydex-ui check-skydex-ui
.PHONY: build build-skycoin build-skyhw build-skyhw-static
.PHONY: test-skyhw test-skyhw-race lint-skyhw check-skyhw
.PHONY: generate update-golden-files
.PHONY: fuzz-base58 fuzz-encoder
.PHONY: check-lang check-lang-es check-lang-zh
.PHONY: install-deps-ui lint-ui test-ui build-ui check-ui check-onpush
.PHONY: test-ui-e2e test-explorer-e2e build-wasm check-wasm

COIN ?= skycoin

# Static files directory
GUI_STATIC_DIR = src/gui/static
# The other two Angular front-ends
EXPLORER_DIR = explorer
SKYCOIN_WEB_DIR = src/skycoin-web
# Holds the wasm cipher the web wallet serves, one copy per toolchain
SKYCOIN_LITE_DIR = src/skycoin-lite

# All Angular front-ends, installed/linted/tested/built as a set so a change to
# one cannot silently break the others.
ANGULAR_UI_DIRS = $(GUI_STATIC_DIR) $(EXPLORER_DIR) $(SKYCOIN_WEB_DIR)

# Every front-end now has a unit suite that runs. src/skycoin-web's ~60 original
# specs still import APIs Angular removed years ago and remain excluded from its
# build (see src/skycoin-web/src/tsconfig.spec.json), but the project is no
# longer without tests: the specs listed there run, so a change-detection
# regression fails CI.
ANGULAR_UI_TEST_DIRS = $(ANGULAR_UI_DIRS)

# Production bundles that are committed to the repository. src/gui/static/dist
# and src/skycoin-web/src/gui/dist are embedded into the Go binary, so a stale
# bundle ships to users.
UI_DIST_DIRS = $(GUI_STATIC_DIR)/dist $(EXPLORER_DIR)/dist $(SKYCOIN_WEB_DIR)/src/gui/dist
# skydex-client trading UI: Vite builds straight into the Go embed dir
# (cmd/skydex-client/commands/static, per its vite.config.js) — no separate copy step.
SKYDEX_UI_DIR = cmd/skydex-client
SKYDEX_UI_EMBED_DIR = cmd/skydex-client/commands/static


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
ifeq ($(shell go env GOOS),darwin)
	@echo "Skipping test-386 on macOS (32-bit not supported)"
else
	GOARCH=386 COIN=$(COIN) go test ./cmd/... -timeout=5m
	GOARCH=386 COIN=$(COIN) go test $$(go list ./src/... | grep -v hardware-wallet) -timeout=5m
endif

test-amd64: ## Run tests for Skycoin with GOARCH=amd64
	GOARCH=amd64 COIN=$(COIN) go test ./cmd/... -timeout=5m
	GOARCH=amd64 COIN=$(COIN) go test ./src/... -timeout=5m

build: build-skycoin build-skyhw ## Build skycoin and skyhw binaries

build-skycoin: ## Build skycoin binary
	go build -o skycoin .

build-skyhw: ## Build skyhw hardware wallet binary with CGO (requires libusb-1.0-dev)
	CGO_ENABLED=1 go build -tags=cgo -o skyhw ./cmd/hardware-wallet/

build-skyhw-static: ## Build statically-linked skyhw binary (requires libusb-1.0-dev)
	CGO_ENABLED=1 go build -tags=cgo -trimpath -ldflags '-linkmode external -extldflags "-static"' -o skyhw ./cmd/hardware-wallet/

test-skyhw: ## Run hardware wallet unit tests
	@mkdir -p coverage/
	go test -v -coverprofile=coverage/skyhw.coverage.out -timeout=2m ./src/hardware-wallet/...

test-skyhw-race: ## Run hardware wallet unit tests with race detector
	go test -v -race -timeout=2m ./src/hardware-wallet/...

lint-skyhw: ## Run linters on hardware wallet code
	go vet ./src/hardware-wallet/...
	go vet ./cmd/hardware-wallet/...

check-skyhw: lint-skyhw test-skyhw build-skyhw ## Run all hardware wallet checks (lint, test, build)

lint: ## Run linters. Use make install-linters first.
	go mod vendor -v
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run -c .golangci.yml ./...
	@# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

check-newcoin: newcoin ## Check that make newcoin succeeds and no templated files are changed.
	@if [ "$(shell git diff ./cmd/skycoin/skycoin.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after make newcoin' ; exit 2 ; fi
	@if [ "$(shell git diff ./cmd/skycoin/skycoin_test.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after make newcoin' ; exit 2 ; fi
	@if [ "$(shell git diff ./src/params/params.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after make newcoin' ; exit 2 ; fi

check-newcoin-templates: ## Verify newcoin templates export and round-trip with --template-dir
	@TMPDIR=$$(mktemp -d) && \
	go run . newcoin templates "$$TMPDIR" && \
	go run . newcoin createcoin --coin skycoin --template-dir "$$TMPDIR" && \
	FAIL=0; \
	if [ "$$(git diff ./cmd/skycoin/skycoin.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after template round-trip in cmd/skycoin/skycoin.go' ; FAIL=1 ; fi; \
	if [ "$$(git diff ./cmd/skycoin/skycoin_test.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after template round-trip in cmd/skycoin/skycoin_test.go' ; FAIL=1 ; fi; \
	if [ "$$(git diff ./src/params/params.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after template round-trip in src/params/params.go' ; FAIL=1 ; fi; \
	rm -rf "$$TMPDIR"; \
	if [ "$$FAIL" = "1" ] ; then exit 2 ; fi; \
	echo "newcoin templates round-trip OK"

check-release-nocgo: ## Verify cmd/release builds without CGO (no hardware wallet)
	CGO_ENABLED=0 go build -o /dev/null ./cmd/release/
	@echo "cmd/release builds OK with CGO_ENABLED=0"

check: lint clean-coverage test test-386 integration-tests-stable check-newcoin ## Run tests and linters

integration-tests-stable: integration-test-stable \
	integration-test-stable-disable-wallet-api \
	integration-test-stable-enable-seed-api \
	integration-test-stable-disable-gui \
	integration-test-stable-auth \
	integration-test-stable-disable-csrf ## Run all stable integration tests

integration-test-stable: ## Run stable integration tests use CSRF, with header check disabled
	COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -c -x -n enable-csrf-header-check

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

integration-test-fibercoin: ## Run fibercoin genesis creation and distribution integration test
	COIN=$(COIN) ./ci-scripts/integration-test-fibercoin.sh

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
	go install golang.org/x/tools/cmd/goimports@latest

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/skycoin ./cmd
	goimports -w -local github.com/skycoin/skycoin ./src

docs-install: ## Install the MkDocs toolchain used to build the documentation site
	pip install -r docs/requirements.txt

docs-serve: ## Stage doc sources and serve the documentation site locally at http://127.0.0.1:8000
	bash scripts/docs-prepare.sh
	mkdocs serve

docs-build: ## Stage doc sources and build the static documentation site into site/
	bash scripts/docs-prepare.sh
	mkdocs build

install-deps-ui:  ## Install the dependencies of every Angular front-end
	@set -e; for d in $(ANGULAR_UI_DIRS); do \
		echo "==> npm ci ($$d)"; \
		(cd $$d && npm ci); \
	done

lint-ui:  ## Lint every Angular front-end
	@set -e; for d in $(ANGULAR_UI_DIRS); do \
		echo "==> npm run lint ($$d)"; \
		(cd $$d && npm run lint); \
	done

check-lang-es: ## Check the Spanish translation
	cd $(GUI_STATIC_DIR)/src/assets/i18n &&node check.js es

check-lang-zh: ## Check the Chinese translation
	cd $(GUI_STATIC_DIR)/src/assets/i18n &&node check.js zh

check-lang: check-lang-es \
	check-lang-zh

test-ui:  ## Run the unit tests of the Angular front-ends that have a working suite
	@set -e; for d in $(ANGULAR_UI_TEST_DIRS); do \
		echo "==> npm run test ($$d)"; \
		(cd $$d && npm run test); \
	done
	@# @angular/build's karma builder writes a scratch bundle to
	@# <project>/dist/test-out/<uuid>, which it does not always clean up and
	@# which cannot be configured elsewhere. In /explorer that lands inside the
	@# committed dist/ that explorer.go embeds with //go:embed dist/*, and an
	@# empty directory there fails the Go build with "contains no embeddable
	@# files". Remove it so `make test-ui && go build .` works in any order.
	@for d in $(ANGULAR_UI_DIRS); do rm -rf "$$d/dist/test-out"; done

test-ui-e2e:  ## Run the desktop GUI e2e tests
	./ci-scripts/ui-e2e.sh

test-explorer-e2e:  ## Run the explorer e2e tests against the pinned blockchain database
	./ci-scripts/explorer-e2e.sh

build-ui:  ## Build the production bundle of every Angular front-end
	@set -e; for d in $(ANGULAR_UI_DIRS); do \
		echo "==> npm run build ($$d)"; \
		(cd $$d && npm run build); \
	done

# TinyGo used to build the wasm cipher. It must be one that records build
# information: upstream TinyGo writes no vcs.* settings and no module version,
# so its output cannot be told apart from any other build and check-wasm rejects
# it. github.com/0magnet/tinygo does record them.
#
#   make build-wasm TINYGO=/path/to/0magnet/tinygo/build/tinygo
TINYGO ?= tinygo

build-wasm:  ## Rebuild the skycoin-lite wasm cipher (needs go and a stamping tinygo)
	@# The repository keeps one copy of each build, in the Go package that embeds
	@# it. cmd/skycoin-web selects between them by build tag and serves the bytes
	@# from memory, so nothing needs copying into a bundle afterwards.
	@#
	@# Each wasm_exec.js is copied from the toolchain that built the wasm beside
	@# it. They are not interchangeable.
	GOOS=js GOARCH=wasm go build -o $(SKYCOIN_LITE_DIR)/wasm-go/skycoin-lite.wasm ./$(SKYCOIN_LITE_DIR)/wasm/
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" $(SKYCOIN_LITE_DIR)/wasm-go/wasm_exec.js
	$(TINYGO) build -o $(SKYCOIN_LITE_DIR)/wasm-tinygo/skycoin-lite.wasm -target wasm ./$(SKYCOIN_LITE_DIR)/wasm/
	cp "$$($(TINYGO) env TINYGOROOT)/targets/wasm_exec.js" $(SKYCOIN_LITE_DIR)/wasm-tinygo/wasm_exec.js
	@node ci-scripts/check-wasm-version.js

check-wasm:  ## Fail if a committed wasm was built from a dirty working tree
	node ci-scripts/check-wasm-version.js

check-onpush:  ## Fail if an OnPush component has an unmarked asynchronous callback
	node ci-scripts/check-onpush-marks.js $(addsuffix /src,$(ANGULAR_UI_DIRS))

check-ui: build-ui  ## Fail if any committed Angular bundle is stale vs a fresh build
	@git diff --quiet -- $(UI_DIST_DIRS) && git diff --quiet --cached -- $(UI_DIST_DIRS) \
		&& test -z "$$(git ls-files --others --exclude-standard -- $(UI_DIST_DIRS))" || { \
		echo "ERROR: a committed Angular production bundle is stale."; \
		echo "Run 'make build-ui' and commit the result."; \
		git --no-pager diff --stat -- $(UI_DIST_DIRS); \
		git ls-files --others --exclude-standard -- $(UI_DIST_DIRS); \
		exit 1; \
	}
	@echo "Committed Angular bundles are up to date."


# skydex-client trading UI (React/Vite, embedded via //go:embed static)
install-deps-skydex-ui:  ## Install the skydex-client UI dependencies
	cd $(SKYDEX_UI_DIR) && npm ci

build-skydex-ui: install-deps-skydex-ui  ## Build the skydex-client UI into the Go embed (commands/static)
	cd $(SKYDEX_UI_DIR) && npm run build

check-skydex-ui: install-deps-skydex-ui  ## Fail if the committed skydex-client UI embed is stale vs a fresh build
	cd $(SKYDEX_UI_DIR) && npm run build
	@git diff --quiet -- $(SKYDEX_UI_EMBED_DIR) || { \
		echo "ERROR: the committed skydex-client UI embed is stale."; \
		echo "Run 'make build-skydex-ui' and commit $(SKYDEX_UI_EMBED_DIR)."; \
		git --no-pager diff --stat -- $(SKYDEX_UI_EMBED_DIR); \
		exit 1; \
	}
	@echo "skydex-client UI embed is up to date."











win-installer: ## Build the windows .msi (installer) custom version
	@powershell '.\scripts\win_installer\script.ps1 $(CUSTOM_VERSION) amd64'
	@powershell '.\scripts\win_installer\script.ps1 $(CUSTOM_VERSION) 386'

windows-installer-release: ## Upload Windows .msi installers to GitHub release
	$(eval GITHUB_TAG=$(shell git describe --abbrev=0 --tags))
	make win-installer CUSTOM_VERSION=$(GITHUB_TAG)
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./skycoin-installer-${GITHUB_TAG}-windows-amd64.msi --clobber
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./skycoin-installer-${GITHUB_TAG}-windows-386.msi --clobber

mac-installer: ## Create unsigned macOS .pkg installers for both architectures
	./scripts/mac_installer/create_installer.sh

mac-installer-release: mac-installer ## Upload macOS .pkg installers to GitHub release
	$(eval GITHUB_TAG=$(shell git describe --abbrev=0 --tags))
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./skycoin-installer-${GITHUB_TAG}-darwin-amd64.pkg --clobber
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./skycoin-installer-${GITHUB_TAG}-darwin-arm64.pkg --clobber

dep-github-release:
	@# Check if musl toolchains are already cached
	@if [ -d "./musl-data/x86_64-linux-musl-cross" ] && \
	    [ -d "./musl-data/aarch64-linux-musl-cross" ] && \
	    [ -d "./musl-data/arm-linux-musleabi-cross" ] && \
	    [ -d "./musl-data/arm-linux-musleabihf-cross" ] && \
	    [ -d "./musl-data/i686-linux-musl-cross" ] && \
	    [ -d "./musl-data/riscv64-linux-musl-cross" ]; then \
		echo "Using cached musl toolchains..."; \
	else \
		echo "Downloading musl toolchains..."; \
		rm -rf musl-data; \
		mkdir -p musl-data; \
		for i in 1 2 3 4 5; do \
			echo "Attempt $$i/5: Downloading aarch64-linux-musl-cross.tgz..."; \
			go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/aarch64-linux-musl-cross.tgz && \
			tar -xzf aarch64-linux-musl-cross.tgz -C ./musl-data && rm aarch64-linux-musl-cross.tgz && break || \
			{ [ $$i -lt 5 ] && { echo "Failed, retrying in 10 seconds..."; sleep 10; } || { echo "All retries failed"; exit 1; }; }; \
		done; \
		for i in 1 2 3 4 5; do \
			echo "Attempt $$i/5: Downloading arm-linux-musleabi-cross.tgz..."; \
			go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/arm-linux-musleabi-cross.tgz && \
			tar -xzf arm-linux-musleabi-cross.tgz -C ./musl-data && rm arm-linux-musleabi-cross.tgz && break || \
			{ [ $$i -lt 5 ] && { echo "Failed, retrying in 10 seconds..."; sleep 10; } || { echo "All retries failed"; exit 1; }; }; \
		done; \
		for i in 1 2 3 4 5; do \
			echo "Attempt $$i/5: Downloading arm-linux-musleabihf-cross.tgz..."; \
			go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/arm-linux-musleabihf-cross.tgz && \
			tar -xzf arm-linux-musleabihf-cross.tgz -C ./musl-data && rm arm-linux-musleabihf-cross.tgz && break || \
			{ [ $$i -lt 5 ] && { echo "Failed, retrying in 10 seconds..."; sleep 10; } || { echo "All retries failed"; exit 1; }; }; \
		done; \
		for i in 1 2 3 4 5; do \
			echo "Attempt $$i/5: Downloading i686-linux-musl-cross.tgz..."; \
			go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/i686-linux-musl-cross.tgz && \
			tar -xzf i686-linux-musl-cross.tgz -C ./musl-data && rm i686-linux-musl-cross.tgz && break || \
			{ [ $$i -lt 5 ] && { echo "Failed, retrying in 10 seconds..."; sleep 10; } || { echo "All retries failed"; exit 1; }; }; \
		done; \
		for i in 1 2 3 4 5; do \
			echo "Attempt $$i/5: Downloading x86_64-linux-musl-cross.tgz..."; \
			go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/x86_64-linux-musl-cross.tgz && \
			tar -xzf x86_64-linux-musl-cross.tgz -C ./musl-data && rm x86_64-linux-musl-cross.tgz && break || \
			{ [ $$i -lt 5 ] && { echo "Failed, retrying in 10 seconds..."; sleep 10; } || { echo "All retries failed"; exit 1; }; }; \
		done; \
		for i in 1 2 3 4 5; do \
			echo "Attempt $$i/5: Downloading riscv64-linux-musl-cross.tgz..."; \
			go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/riscv64-linux-musl-cross.tgz && \
			tar -xzf riscv64-linux-musl-cross.tgz -C ./musl-data && rm riscv64-linux-musl-cross.tgz && break || \
			{ [ $$i -lt 5 ] && { echo "Failed, retrying in 10 seconds..."; sleep 10; } || { echo "All retries failed"; exit 1; }; }; \
		done; \
	fi
	# Build libusb-1.0 static libraries for each musl target
	./ci-scripts/build-libusb-musl.sh amd64 x86_64-linux-musl ./musl-data/x86_64-linux-musl-cross
	./ci-scripts/build-libusb-musl.sh arm64 aarch64-linux-musl ./musl-data/aarch64-linux-musl-cross
	./ci-scripts/build-libusb-musl.sh arm arm-linux-musleabi ./musl-data/arm-linux-musleabi-cross
	./ci-scripts/build-libusb-musl.sh armhf arm-linux-musleabihf ./musl-data/arm-linux-musleabihf-cross
	./ci-scripts/build-libusb-musl.sh 386 i686-linux-musl ./musl-data/i686-linux-musl-cross
	./ci-scripts/build-libusb-musl.sh riscv64 riscv64-linux-musl ./musl-data/riscv64-linux-musl-cross







clean-coverage: ## Remove coverage output files
	rm -rf ./coverage/

newcoin: ## Rebuild cmd/$COIN/$COIN.go file from the template. Call like "make newcoin COIN=foo".
	go run -mod=mod . newcoin createcoin --coin $(COIN)

generate: ## Generate test interface mocks and struct encoders
	go generate ./src/...
	# mockery can't generate the UnspentPooler mock in package visor, patch it
	mv ./src/visor/blockdb/mock_unspent_pooler_test.go ./src/visor/mock_unspent_pooler_test.go
	sed -i "" -e 's/package blockdb/package visor/g' ./src/visor/mock_unspent_pooler_test.go
	sed -i "" -e 's/AddressHashes/blockdb.AddressHashes/g' ./src/visor/mock_unspent_pooler_test.go
	goimports -w -local github.com/skycoin/skycoin ./src/visor/mock_unspent_pooler_test.go

install-generators: ## Install tools used by go generate
	go get github.com/vektra/mockery/.../
	go get github.com/skycoin/skyencoder/cmd/skyencoder

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
	go-fuzz-build github.com/skycoin/skycoin/src/cipher/base58/fuzz
	go-fuzz -bin=base58fuzz-fuzz.zip -workdir=src/cipher/base58/fuzz

fuzz-encoder: ## Fuzz the encoder package. Requires https://github.com/dvyukov/go-fuzz
	go-fuzz-build github.com/skycoin/skycoin/src/cipher/encoder/fuzz
	go-fuzz -bin=encoderfuzz-fuzz.zip -workdir=src/cipher/encoder/fuzz

update-dep: ## Update vendor deps, commit, and push
	go get -v -u ./...
	go mod tidy
	go mod vendor
	git add go.mod go.sum vendor
	git commit -m "update deps"
	git push

## Sync local develop branch with upstream skycoin/skycoin develop.
## Requires: origin = your fork, upstream = skycoin/skycoin
sync-upstream-develop: ## Sync develop branch with upstream develop branch for forks
	@normalize() { \
		echo "$$1" | sed \
			-e 's|git@github.com:|https://github.com/|' \
			-e 's|\.git$$||' \
			-e 's|https://github.com/||' \
			| tr '[:upper:]' '[:lower:]'; \
	}; \
	UPSTREAM_URL=$$(git remote get-url upstream 2>/dev/null); \
	if [ -z "$$UPSTREAM_URL" ]; then \
		echo "[error] no 'upstream' remote found. Add it with:"; \
		echo "  git remote add upstream https://github.com/skycoin/skycoin.git"; \
		exit 1; \
	fi; \
	UPSTREAM_NORM=$$(normalize "$$UPSTREAM_URL"); \
	if [ "$$UPSTREAM_NORM" != "skycoin/skycoin" ]; then \
		echo "[error] upstream remote does not point to skycoin/skycoin."; \
		echo "  Found: $$UPSTREAM_URL"; \
		exit 1; \
	fi; \
	ORIGIN_URL=$$(git remote get-url origin 2>/dev/null); \
	if [ -z "$$ORIGIN_URL" ]; then \
		echo "[error] no 'origin' remote found."; \
		exit 1; \
	fi; \
	ORIGIN_NORM=$$(normalize "$$ORIGIN_URL"); \
	if [ "$$ORIGIN_NORM" = "skycoin/skycoin" ]; then \
		echo "[error] origin points to skycoin/skycoin directly."; \
		echo "  This target must be run from a fork, not a clone of the canonical repo."; \
		exit 1; \
	fi; \
	echo "[ok] origin is a fork ($$ORIGIN_NORM), upstream is skycoin/skycoin — syncing develop..."; \
	git checkout develop && \
	git pull && \
	git fetch upstream && \
	git merge upstream/develop && \
	git push

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
