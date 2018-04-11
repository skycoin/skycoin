# skycoin build binaries
# reference https://github.com/skycoin/skycoin
FROM golang:1.10 AS build

COPY . $GOPATH/src/github.com/skycoin/skycoin

RUN cd $GOPATH/src/github.com/skycoin/skycoin && \
  CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./...

# skycoin image
FROM alpine:3.7

ENV COIN="skycoin" \
    RPC_ADDR="0.0.0.0:6430" \
    DATA_DIR="/data/.$COIN" \
    WALLET_DIR="/wallet" \
    WALLET_NAME="$COIN_cli.wlt"

# copy binaries
COPY --from=build /go/bin/* /usr/bin/

# copy assets
COPY --from=build /go/src/github.com/skycoin/skycoin/src/gui/static /usr/local/skycoin/src/gui/static

# copy launcher
COPY docker_launcher.sh /usr/local/bin

# volumes
VOLUME $WALLET_DIR
VOLUME $DATA_DIR

EXPOSE 6000 6420 6430

WORKDIR /usr/local/skycoin

CMD ["skycoin", "--web-interface-addr=0.0.0.0"]
