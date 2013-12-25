default: run

.PHONY: default run

run:
	GOPATH=`pwd` go run main.go
build:
	GOPATH=`pwd` go build main.go
get:
	GOPATH=`pwd` go get ...

.PHONY: run get build
