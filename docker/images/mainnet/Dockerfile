# skycoin build
# reference https://github.com/skycoin/skycoin
FROM golang:1.9 AS build

ENV SKYCOIN_VERSION=0.22.0

# dirs
RUN mkdir -p $GOPATH/src/github.com/skycoin/skycoin

# clone
RUN wget -c https://github.com/skycoin/skycoin/archive/v${SKYCOIN_VERSION}.tar.gz

# uncompress
RUN tar zxf v${SKYCOIN_VERSION}.tar.gz -C $GOPATH/src/github.com/skycoin/skycoin/ --strip-components=1

# install cli
RUN cd $GOPATH/src/github.com/skycoin/skycoin/cmd/cli && \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/skycoin-cli .

# install skycoin
RUN cd $GOPATH/src/github.com/skycoin/skycoin/cmd/skycoin && \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/skycoin .

# install address_gen
RUN cd $GOPATH/src/github.com/skycoin/skycoin/cmd/address_gen && \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/address_gen .


# skycoin image
FROM alpine:3.7

ENV COIN="skycoin" \
    RPC_ADDR="0.0.0.0:6430" \
    DATA_DIR="/root/.$COIN" \
    WALLET_DIR="/wallet" \
    WALLET_NAME="$COIN_cli.wlt"

# copy all the binaries
COPY --from=build /go/bin/* /usr/bin/

# copy assets
COPY --from=build /go/src/github.com/skycoin/skycoin/src/gui/static /usr/local/skycoin/static

# volumes
VOLUME $WALLET_DIR
VOLUME $DATA_DIR

EXPOSE 6000 6420 6430

CMD ["/usr/bin/skycoin", "--gui-dir=/usr/local/skycoin/static", "--web-interface-addr=0.0.0.0"]
