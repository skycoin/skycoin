# skycoin build binaries
# reference https://github.com/skycoin/skycoin
FROM golang:1.9-alpine AS build-go

COPY . $GOPATH/src/github.com/skycoin/skycoin

RUN cd $GOPATH/src/github.com/skycoin/skycoin && \
  CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./...


# skycoin gui
FROM node:8.9 AS build-node

COPY . /skycoin

# `unsafe` flag used as work around to prevent infinite loop in Docker
# see https://github.com/nodejs/node-gyp/issues/1236
RUN npm install -g --unsafe @angular/cli && \
    cd /skycoin/src/gui/static && \
    yarn && \
    npm run build


# skycoin image
FROM alpine:3.7

ENV COIN="skycoin" \
    RPC_ADDR="0.0.0.0:6430" \
    DATA_DIR="/data/.$COIN" \
    WALLET_DIR="/wallet" \
    WALLET_NAME="$COIN_cli.wlt"

# create directories and set ownership
RUN mkdir -p $DATA_DIR $WALLET_DIR


# copy binaries
COPY --from=build-go /go/bin/* /usr/bin/

# copy gui
COPY --from=build-node /skycoin/src/gui/static /usr/local/skycoin/src/gui/static

# copy launcher
COPY launcher.sh /usr/local/bin

# volumes
VOLUME $WALLET_DIR
VOLUME $DATA_DIR

EXPOSE 6000 6420 6430

WORKDIR /usr/local/skycoin

CMD ["launcher.sh", "--web-interface-addr=0.0.0.0",  "--rpc-interface-addr=0.0.0.0"]
