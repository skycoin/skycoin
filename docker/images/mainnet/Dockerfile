# skycoin build
# reference https://github.com/skycoin/skycoin
ARG IMAGE_FROM=busybox
FROM golang:1.10.1-stretch AS build
ARG ARCH=amd64
ARG GOARM
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
    COMMIT=$(git rev-parse HEAD) BRANCH=$(git rev-parse -—abbrev-ref HEAD) \
    GOARCH=$ARCH GOARM=$GOARM CGO_ENABLED=0 GOOS=linux \
    GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}" \
    go install -ldflags "${GOLDFLAGS}" ./cmd/... && \
    sh -c "if test -d $GOPATH/bin/linux_arm ; then mv $GOPATH/bin/linux_arm/* $GOPATH/bin/; fi; \
           if test -d $GOPATH/bin/linux_arm64 ; then mv $GOPATH/bin/linux_arm64/* $GOPATH/bin/; fi"

# skycoin image
FROM $IMAGE_FROM

ENV COIN="skycoin"
ENV RPC_ADDR="http://0.0.0.0:6420" \
    DATA_DIR="/data/.$COIN" \
    WALLET_DIR="/wallet" \
    USE_CSRF="1" \
    WALLET_NAME="$COIN_cli.wlt"

# copy all the binaries
COPY --from=build /go/bin/* /usr/bin/

# copy assets
COPY --from=build /go/src/github.com/skycoin/skycoin/src/gui/static /usr/local/skycoin/src/gui/static

# copy launcher
COPY docker_launcher.sh /usr/local/bin/docker_launcher.sh

# volumes
VOLUME $WALLET_DIR
VOLUME $DATA_DIR

EXPOSE 6000 6420

ENTRYPOINT ["docker_launcher.sh", "--web-interface-addr=0.0.0.0", "--gui-dir=/usr/local/skycoin/src/gui/static"]
