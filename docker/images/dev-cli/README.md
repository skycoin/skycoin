# Supported tags and respective `Dockerfile` links

## Simple Tags

-	[`latest` (*docker/images/skycoindev-cli/Dockerfile*)](https://github.com/skycoin/skycoin/tree/develop/docker/images/skycoindev-cli/Dockerfile)

# Skycoin development image

This image has the necessary tools to build, test, edit, lint and version the Skycoin
source code.  It comes with the Vim editor installed, along with some plugins
to ease go development and version control with git.

# How to use this image

## Initialize your development environment.

```sh
$ mkdir src
$ docker run --rm \
    -v src:/go/src skycoin/skycoindev-cli \
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

You can run commands by just passing the them to the image.  Everything is run
in a container and deleted when finished.

### Running tests

```sh
$ docker run --rm \
    -v src:/go/src skycoin/skycoindev-cli \
    sh -c "cd skycoin; make test"
```

### Running lint

```sh
$ docker run --rm \
    -v src:/go/src skycoin/skycoindev-cli \
    sh -c "cd skycoin; make lint"
```

### Editing code

```sh
$ docker run --rm \
    -v src:/go/src skycoin/skycoindev-cli \
    vim
```


