# skycoin build
# reference https://github.com/skycoin/skycoin
FROM golang:1.10 AS build

ARG SKYCOIN_VERSION

COPY . $GOPATH/src/github.com/skycoin/skycoin

# This code checks if SKYCOIN_VERSION is set and checkouts to that version if
# so. The git stash line prevents the build to fail if there are any uncommited
# changes in the working copy. It won't affect the host working copy.
RUN sh -c \
    'if test ${SKYCOIN_VERSION};then \
        echo "Revision is set to: "${SKYCOIN_VERSION}; \
        cd $GOPATH/src/github.com/skycoin/skycoin; \
        git stash; \
        git checkout ${SKYCOIN_VERSION}; \
     fi'

RUN cd $GOPATH/src/github.com/skycoin/skycoin && \
  CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./...

# skycoin image
FROM alpine:3.7

ENV COIN="skycoin"
ENV RPC_ADDR="0.0.0.0:6430" \
    DATA_DIR="/data/.$COIN" \
    WALLET_DIR="/wallet" \
    WALLET_NAME="$COIN_cli.wlt"

# copy all the binaries
COPY --from=build /go/bin/* /usr/bin/

# copy assets
COPY --from=build /go/src/github.com/skycoin/skycoin/src/gui/static /usr/local/skycoin/src/gui/static

# copy launcher
COPY docker_launcher.sh /usr/local/bin

# volumes
VOLUME $WALLET_DIR
VOLUME $DATA_DIR

EXPOSE 6000 6420 6430

ENTRYPOINT ["docker_launcher.sh", "--web-interface-addr=0.0.0.0", "--rpc-interface-addr=0.0.0.0", "--gui-dir=/usr/local/skycoin/src/gui/static"]
