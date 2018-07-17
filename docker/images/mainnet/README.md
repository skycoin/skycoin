# Supported tags

## Simple Tags

- latest
- latest-arm32v5
- latest-arm32v6
- latest-arm32v7
- latest-arm64v8
- develop
- develop-arm32v5
- develop-arm32v6
- develop-arm32v7
- develop-arm64v8
- release-v0.23.0
- release-v0.23.0-arm32v5
- release-v0.23.0-arm32v6
- release-v0.23.0-arm32v7
- release-v0.23.0-arm64v8
- release-v0.22.0

## Building your own images

This Dockerfile build your working copy by default, but if you pass the
SKYCOIN_VERSION build argument to the `docker build` command, it will checkout
to the branch, a tag or a commit you specify on that variable.

Example

```sh
$ git clone https://github.com/skycoin/skycoin
$ cd skycoin
$ SKYCOIN_VERSION=v0.23.0
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=SKYCOIN_VERSION=$SKYCOIN_VERSION \
  -t skycoin:$SKYCOIN_VERSION .
```

or just

```sh
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=SKYCOIN_VERSION=v0.23.0 \
  -t skycoin:v0.23.0 .
```

## ARM Architecture

Build arguments are provided to make it easy if you want to build for the ARM
architecture.

Example for ARMv5

```sh
$ git clone https://github.com/skycoin/skycoin
$ cd skycoin
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=ARCH=arm \
  --build-arg=GOARM=5 \
  --build-arg=IMAGE_FROM="arm32v5/alpine" \
  -t skycoin:latest-arm32v5 .
```

## How to use this images

### Create a Skycoin node

This command launch a skycoin(version 0.23.0) node in background on top of Docker

```sh
$ docker volume create skycoin0.23.0-data
$ docker volume create skycoin0.23.0-wallet
$ docker run --rm -d -v skycoin0.23.0-data:/data/.skycoin \
  -v skycoin0.23.0-wallet:/wallet \
  -p 6000:6000 -p 6420:6420 \
  --name skycoin-node skycoin:v0.23.0
```

If you want to stop it , just run

```sh
$ docker stop skycoin-node
```

You can pass parameters to skycoin process inside the container

```sh
 $ docker run --rm -d -v skycoin0.23.0-data:/data/.skycoin \
  -v skycoin0.23.0-wallet:/wallet \
  -p 6000:6000 -p 6420:6420 \
  --name skycoin-node skycoin:v0.23.0 -web-interface-addr 192.168.1.1
```

