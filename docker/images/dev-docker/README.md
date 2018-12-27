# Supported tags and respective `Dockerfile` links

## Simple Tags

-	[`dind` (*docker/images/dev-docker/Dockerfile*)](https://github.com/skycoin/skycoin/tree/develop/docker/images/dev-docker/Dockerfile)

# Skycoin development image including [docker in docker](https://hub.docker.com/_/docker/)

This image has the necessary tools to build, test, edit, lint and version the Skycoin
source code.  It comes with the Vim editor installed, along with some plugins
to ease go development and version control with git, besides it comes with docker installed.

# How to use this image

## Initialize your development environment.

```sh
$ mkdir src
$ docker run --privileged --rm \
    -v src:/go/src skycoin/skycoindev-cli:dind \
    go get github.com/skycoin/skycoin
$ sudo chown -R `whoami` src
```

This downloads the skycoin source to src/skycoin/skycoin and changes the owner
to your user. This is necessary, because all processes inside the container run
as root and the files created by it are therefore owned by root.

If you already have a Go development environment installed, you just need to
mount the src directory from your $GOPATH in the /go/src volume of the
container.

## Running commands inside the container

You can run commands by just passing them to the image.  Everything is run
in a container and deleted when finished.

### Running tests

```sh
$ docker run --rm \
    -v src:/go/src skycoin/skycoindev-cli:dind \
    sh -c "cd skycoin; make test"
```

### Running lint

```sh
$ docker run --rm \
    -v src:/go/src skycoin/skycoindev-cli:dind \
    sh -c "cd skycoin; make lint"
```

### Editing code

```sh
$ docker run --rm \
    -v src:/go/src skycoin/skycoindev-cli:dind \
    vim
```

## How to use docker in docker image

### Start a daemon instance

```sh
$ docker run --privileged --name some-name -d skycoin/skycoindev-cli:dind
```

### Where to store data

Create a data directory on the host system (outside the container) and mount this to a directory visible from inside the container.

The downside is that you need to make sure that the directory exists, and that e.g. directory permissions and other security mechanisms on the host system are set up correctly.

1. Create a data directory on a suitable volume on your host system, e.g. /my/own/var-lib-docker.
2. Start your docker container like this:

```sh
$ docker run --privileged --name some-name -v /my/own/var-lib-docker:/var/lib/docker \ 
-d skycoin/skycoindev-cli:dind
```


