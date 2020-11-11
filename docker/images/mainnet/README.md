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
$ SKYCOIN_VERSION=v0.24.0
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=SKYCOIN_VERSION=$SKYCOIN_VERSION \
  -t skycoin:$SKYCOIN_VERSION .
```

or just

```sh
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=SKYCOIN_VERSION=v0.24.0 \
  -t skycoin:v0.24.0
```

## ARM Architecture

Build arguments are provided to make it easy if you want to build for the ARM
architecture.

Example for ARMv5.

```sh
$ git clone https://github.com/skycoin/skycoin
$ cd skycoin
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=ARCH=arm \
  --build-arg=GOARM=5 \
  --build-arg=IMAGE_FROM="arm32v5/alpine" \
  -t skycoin:$SKYCOIN_VERSION-arm32v5 .
```

## How to use this images

### Run a Skycoin node

This command pulls latest stable image from Docker Hub, and launches a node inside a Docker container that runs as a service daemon in the background. It is possible to use the tags listed above to run another version of the node

```sh
$ docker volume create skycoin-data
$ docker volume create skycoin-wallet
$ docker run -d -v skycoin-data:/data/.skycoin \
  -v skycoin-wallet:/wallet \
  -p 6000:6000 -p 6420:6420 \
  --name skycoin-node-stable skycoin/skycoin
```

When invoking the container this way the options of the skycoin command are set to their respective default values , except the following

| Parameter  | Value |
| ------------- | ------------- |
| web-interface-addr | 0.0.0.0  |
| gui-dir | /usr/local/skycoin/src/gui/static |

In order to stop the container , just run

```sh
$ docker stop skycoin-node-stable
```

Restart it once again by executing

```sh
$ docker start skycoin-node-stable
```

### Customizing node server with parameters

The container accepts parameters in order to customize the execution of the skycoin node. For instance, in order to run the bleeding edge development image and listen for REST API requests at a non-standard port (e.g. `6421`) it is possible to execute the following command.

```sh
 $ docker run --rm -d -v skycoin-data:/data/.skycoin \
  -v skycoin-wallet:/wallet \
  -p 6000:6000 -p 6421:6421 \
  --name skycoin-node-develop skycoin/skycoin:develop -web-interface-port 6421
```

Notice that the value of node parameter (e.g. `-web-interface-port`) affects the execution context inside the container. Therefore, in this particular case, the port mapping should be updated accordingly.

To get a full list of skycoin's parameters, just run

```sh
 $ docker run --rm skycoin/skycoin:develop -help
```

To run multiple nodes concurrently in the same host, it is highly recommended to create separate volumes for each node. For example, in order to run a block publisher node along with the one launched above, it is necessary to execute

```sh
$ docker volume create skycoin-block-publisher-data
$ docker volume create skycoin-block-publisher-wallet
$ docker run -d -v skycoin-block-publisher-data:/data/.skycoin \
  -v skycoin-block-publisher-wallet:/wallet \
  -p 6001:6000 -p 6421:6420 \
  --name skycoin-block-publisher-stable skycoin/skycoin -block-publisher
```

Notice that the host's port must be changed since collisions of two services listening at the same port are not allowed by the low-level operating system socket libraries.
