.DEFAULT_GOAL := help
.PHONY: lint check help

install-deps-ui: ## install npm dependences
	npm install
	cd electron && npm install

build: ## builds the GUI to src/gui/dist
	npm run build

build-for-electron: ## compiles a version to be used with Electron
	npm run build-for-local-fs

build-for-local-fs: ## compiles a version to be used from the local file system
	npm run build-for-local-fs

lint: ## runs lint
	npm run lint

unit-test: ## runs unit tests
	npm run test

cipher-test:
	npm run cipher-test

run-docker: ## runs docker container
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
	-disable-networking \
	-disable-csrf

stop-docker: ## stops docker container
	docker stop skycoin-backend

e2e-test: ## runs e2e tests using a node running in Docker
	npm run e2e-docker

e2e-prod-test: ## runs e2e prod tests using a node running in Docker
	npm run e2e-docker-prod

check: run-docker e2e-test e2e-prod-test stop-docker ## runs linter, unit tests, e2e tests

build-electron: ## creates Electron wallets for Mac and Linux.
	./ci-scripts/build-wallets.sh

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
