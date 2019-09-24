# Supported tags and respective `Dockerfile` links

## Simple Tags

-	[`latest` (*docker/images/dev-cli/Dockerfile*)](https://github.com/SkycoinProject/skycoin/tree/develop/docker/images/dev-cli/Dockerfile)

# Skycoin development image

This image has the necessary tools to build, test, edit, lint and version the Skycoin
source code.  It comes with the Vim editor installed, along with some plugins
to ease go development and version control with git.

_Plase note that there is also a sister image with ["docker in docker"](https://github.com/SkycoinProject/skycoin/tree/develop/docker/images/dev-docker/) feature on it._

# How to use this image

## Initialize your development environment.

```sh
$ mkdir src
$ docker run --rm \
    -v src:/go/src SkycoinProject/skycoindev-cli \
    go get github.com/SkycoinProject/skycoin
$ sudo chown -R `whoami` src
```

This downloads the skycoin source to src/SkycoinProject/skycoin and changes the owner
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
    -v src:/go/src SkycoinProject/skycoindev-cli \
    sh -c "cd skycoin; make test"
```

### Running lint

```sh
$ docker run --rm \
    -v src:/go/src SkycoinProject/skycoindev-cli \
    sh -c "cd skycoin; make lint"
```

### Editing code

```sh
$ docker run --rm \
    -v src:/go/src SkycoinProject/skycoindev-cli \
    vim
```

## Additional tools and packages installed

### Packages

- dep
- tig
- swig

### Vim's plugins

- Ale
- tig-explorer

## Automated builds

Docker Cloud is configured to build images from `develop` branch on every push.
The same process is triggered for all feature branches matching the pattern
`/^([^_]+)_t([0-9]+)_.*docker.*/`. The tag generated for those images will be of the form
`feature-{\1}-{\2}`.

