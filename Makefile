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
.PHONY: install-linters format release clean-release clean-coverage dep-github-release
.PHONY: install-deps-ui build-ui build-ui help newcoin merge-coverage
.PHONY: build build-skycoin build-skyhw build-skyhw-static
.PHONY: generate update-golden-files
.PHONY: fuzz-base58 fuzz-encoder
.PHONY: check-lang check-lang-es check-lang-zh

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
ifeq ($(shell go env GOOS),darwin)
	@echo "Skipping test-386 on macOS (32-bit not supported)"
else
	GOARCH=386 COIN=$(COIN) go test ./cmd/... -timeout=5m
	GOARCH=386 COIN=$(COIN) go test ./src/... -timeout=5m
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

lint: ## Run linters. Use make install-linters first.
	go mod vendor -v
	golangci-lint run -c .golangci.yml ./...
	@# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

check-newcoin: newcoin ## Check that make newcoin succeeds and no templated files are changed.
	@if [ "$(shell git diff ./cmd/skycoin/skycoin.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after make newcoin' ; exit 2 ; fi
	@if [ "$(shell git diff ./cmd/skycoin/skycoin_test.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after make newcoin' ; exit 2 ; fi
	@if [ "$(shell git diff ./src/params/params.go | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected after make newcoin' ; exit 2 ; fi

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
	go install github.com/FiloSottile/vendorcheck@latest

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/skycoin ./cmd
	goimports -w -local github.com/skycoin/skycoin ./src

install-deps-ui:  ## Install the UI dependencies
	cd $(GUI_STATIC_DIR) && npm ci

lint-ui:  ## Lint the UI code
	cd $(GUI_STATIC_DIR) && npm run lint

check-lang-es: ## Check the Spanish translation
	cd $(GUI_STATIC_DIR)/src/assets/i18n &&node check.js es

check-lang-zh: ## Check the Chinese translation
	cd $(GUI_STATIC_DIR)/src/assets/i18n &&node check.js zh

check-lang: check-lang-es \
	check-lang-zh

test-ui:  ## Run UI tests
	cd $(GUI_STATIC_DIR) && npm run test

test-ui-e2e:  ## Run UI e2e tests
	./ci-scripts/ui-e2e.sh

build-ui:  ## Builds the UI
	cd $(GUI_STATIC_DIR) && npm run build


snapshot: ## Build snapshot release with goreleaser (all platforms)
	go run github.com/goreleaser/goreleaser/v2@main --snapshot --clean --skip=publish --config .goreleaser-linux.yml

snapshot-linux: ## Build snapshot release for Linux only
	go run github.com/goreleaser/goreleaser/v2@main --snapshot --clean --skip=publish --config .goreleaser-linux.yml

snapshot-darwin: ## Build snapshot release for macOS only
	go run github.com/goreleaser/goreleaser/v2@main --snapshot --clean --skip=publish --config .goreleaser-darwin.yml

snapshot-windows: ## Build snapshot release for Windows only
	go run github.com/goreleaser/goreleaser/v2@main --snapshot --clean --skip=publish --config .goreleaser-windows.yml

github-prepare-release:
	$(eval GITHUB_TAG=$(shell git describe --abbrev=0 --tags | sed 's/-.*//'))
	sed '/^## ${GITHUB_TAG}$$/,/^## .*/!d;//d;/^$$/d' ./CHANGELOG.md > releaseChangelog.md

github-release: github-prepare-release ## Create GitHub release for Linux (triggered by GitHub Actions on tag push)
	go run github.com/goreleaser/goreleaser/v2@main --clean --config .goreleaser-linux.yml --release-notes releaseChangelog.md --skip=publish
	$(eval GITHUB_TAG=$(shell git describe --abbrev=0 --tags))
	gh release create ${GITHUB_TAG} --repo skycoin/skycoin --title ${GITHUB_TAG} --notes-file releaseChangelog.md || true
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-linux-amd64.tar.gz --clobber
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-linux-arm64.tar.gz --clobber
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-linux-386.tar.gz --clobber
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-linux-arm.tar.gz --clobber
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-linux-armhf.tar.gz --clobber
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/checksums.txt --clobber

github-release-darwin-amd64: ## Create GitHub release for macOS Intel (triggered by GitHub Actions)
	go run github.com/goreleaser/goreleaser/v2@main --clean --config .goreleaser-darwin-amd64.yml --skip=publish
	$(eval GITHUB_TAG=$(shell git describe --abbrev=0 --tags))
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-darwin-amd64.tar.gz
	gh release download ${GITHUB_TAG} --repo skycoin/skycoin --pattern 'checksums*'
	cat ./dist/checksums.txt >> checksums.txt
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} --clobber ./checksums.txt

github-release-darwin-arm64: ## Create GitHub release for macOS ARM (triggered by GitHub Actions)
	go run github.com/goreleaser/goreleaser/v2@main --clean --config .goreleaser-darwin-arm64.yml --skip=publish
	$(eval GITHUB_TAG=$(shell git describe --abbrev=0 --tags))
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-darwin-arm64.tar.gz
	gh release download ${GITHUB_TAG} --repo skycoin/skycoin --pattern 'checksums*'
	cat ./dist/checksums.txt >> checksums.txt
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} --clobber ./checksums.txt

github-release-darwin: ## Create GitHub release for macOS (triggered by GitHub Actions)
	go run github.com/goreleaser/goreleaser/v2@main --clean --config .goreleaser-darwin.yml --skip=publish
	$(eval GITHUB_TAG=$(shell git describe --abbrev=0 --tags))
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-darwin-amd64.tar.gz
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-darwin-arm64.tar.gz
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-darwin-amd64.pkg || true
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-darwin-arm64.pkg || true
	gh release download ${GITHUB_TAG} --repo skycoin/skycoin --pattern 'checksums*'
	cat ./dist/checksums.txt >> checksums.txt
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} --clobber ./checksums.txt

github-release-windows: ## Create GitHub release for Windows (triggered by GitHub Actions)
	go run github.com/goreleaser/goreleaser/v2@main --clean --config .goreleaser-windows.yml --skip=publish
	$(eval GITHUB_TAG=$(shell git describe --abbrev=0 --tags))
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-windows-amd64.zip --clobber
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-windows-386.zip --clobber
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./dist/skycoin-${GITHUB_TAG}-windows-arm64.zip --clobber
	gh release download ${GITHUB_TAG} --repo skycoin/skycoin --pattern 'checksums*'
	cat ./dist/checksums.txt >> checksums.txt
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} --clobber ./checksums.txt

win-installer: ## Build the windows .msi (installer) custom version
	@powershell '.\scripts\win_installer\script.ps1 $(CUSTOM_VERSION) amd64'
	@powershell '.\scripts\win_installer\script.ps1 $(CUSTOM_VERSION) 386'

windows-installer-release: ## Upload Windows .msi installers to GitHub release
	$(eval GITHUB_TAG=$(shell git describe --abbrev=0 --tags))
	make win-installer CUSTOM_VERSION=$(GITHUB_TAG)
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./skycoin-installer-${GITHUB_TAG}-windows-amd64.msi --clobber
	gh release upload --repo skycoin/skycoin ${GITHUB_TAG} ./skycoin-installer-${GITHUB_TAG}-windows-386.msi --clobber

dep-github-release:
	rm -rf musl-data
	mkdir -p musl-data
	go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/aarch64-linux-musl-cross.tgz
	tar -xzf aarch64-linux-musl-cross.tgz -C ./musl-data && rm aarch64-linux-musl-cross.tgz
	go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/arm-linux-musleabi-cross.tgz
	tar -xzf arm-linux-musleabi-cross.tgz -C ./musl-data && rm arm-linux-musleabi-cross.tgz
	go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/arm-linux-musleabihf-cross.tgz
	tar -xzf arm-linux-musleabihf-cross.tgz -C ./musl-data && rm arm-linux-musleabihf-cross.tgz
	go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/i686-linux-musl-cross.tgz
	tar -xzf i686-linux-musl-cross.tgz -C ./musl-data && rm i686-linux-musl-cross.tgz
	go run github.com/melbahja/got/cmd/got@latest https://github.com/skycoin/skywire/releases/download/v1.3.29/x86_64-linux-musl-cross.tgz
	tar -xzf x86_64-linux-musl-cross.tgz -C ./musl-data && rm x86_64-linux-musl-cross.tgz
	# Build libusb-1.0 static libraries for each musl target
	./ci-scripts/build-libusb-musl.sh amd64 x86_64-linux-musl ./musl-data/x86_64-linux-musl-cross
	./ci-scripts/build-libusb-musl.sh arm64 aarch64-linux-musl ./musl-data/aarch64-linux-musl-cross
	./ci-scripts/build-libusb-musl.sh arm arm-linux-musleabi ./musl-data/arm-linux-musleabi-cross
	./ci-scripts/build-libusb-musl.sh armhf arm-linux-musleabihf ./musl-data/arm-linux-musleabihf-cross
	./ci-scripts/build-libusb-musl.sh 386 i686-linux-musl ./musl-data/i686-linux-musl-cross

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

clean-release: ## Remove all electron build artifacts and goreleaser dist
	rm -rf $(ELECTRON_DIR)/release
	rm -rf $(ELECTRON_DIR)/.gox_output
	rm -rf $(ELECTRON_DIR)/.daemon_output
	rm -rf $(ELECTRON_DIR)/.cli_output
	rm -rf $(ELECTRON_DIR)/.standalone_output
	rm -rf $(ELECTRON_DIR)/.electron_output
	rm -rf ./dist

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
	go-fuzz-build github.com/skycoin/skycoin/src/cipher/base58/internal
	go-fuzz -bin=base58fuzz-fuzz.zip -workdir=src/cipher/base58/internal

fuzz-encoder: ## Fuzz the encoder package. Requires https://github.com/dvyukov/go-fuzz
	go-fuzz-build github.com/skycoin/skycoin/src/cipher/encoder/internal
	go-fuzz -bin=encoderfuzz-fuzz.zip -workdir=src/cipher/encoder/internal

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
