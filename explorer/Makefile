EXPLORER := explorer

.DEFAULT_GOAL := help

.PHONY: run build-go build-ng all setcap deploy help lint check install-linters install-linters-local format verify

run: ## Run explorer.go
	go run explorer.go

run-api: ## Run explorer.go -api-only
	go run explorer.go -api-only

build-go: ## Build explorer.go
	go build -o $(EXPLORER) explorer.go

build-ng: ## Build the angular frontend
	npm run build

all: build-go build-ng ## Build explorer.go and the angular frontend

setcap: ## Use setcap to allow the binary to bind to port 80
	sudo setcap 'cap_net_bind_service=+ep' $(EXPLORER)

protractor: ## Run e2e tests. Must have explorer.go running (call "make run").
	npm run e2e

test: ## Run tests
	go test . -timeout=1m

verify: ## Run explorer self-verification
	go build -o explorer.verify explorer.go
	@./explorer.verify -verify; \
	status=$$?; \
	rm -f explorer.verify; \
	exit $$status

lint: ## Run linters. Use make install-linters or install-linters-local first.
	# go mod vendor ./*.go
	golangci-lint run -c .golangci.yml ./*.go
	# @# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./*.go

check: lint test verify ## Run tests, linters and self-verification

install-linters: ## Install linters for CI testing
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.52.2

install-linters-local: ## Install linters for local testing
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w explorer.go

lint-ts: ## runs ts lint
	npm run lint

check-ui: ## runs e2e tests connecting to a containerized node
	go run explorer.go &>/dev/null &
	sleep 10

	docker volume create skycoin-data
	docker volume create skycoin-wallet
	chmod 777 $(PWD)/e2e/test-fixtures/blockchain-180.db

	docker run -d --rm \
	-v skycoin-data:/data \
	-v skycoin-wallet:/wallet \
	-v $(PWD)/e2e/:/project-root \
	--name skycoin-backend \
	-p 6000:6000 \
	-p 6420:6420 \
	skycoin/skycoin:develop \
	-web-interface-addr 172.17.0.2 \
	-db-path=project-root/test-fixtures/blockchain-180.db \
	-disable-networking

	npm run e2e-blockchain-180

	docker stop skycoin-backend

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
