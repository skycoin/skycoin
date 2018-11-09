# Supported tags and respective `Dockerfile` links

## Simple Tags

-	[`latest` (*docker/images/dev-vscode/Dockerfile*)](https://github.com/skycoin/skycoin/tree/develop/docker/images/dev-vscode/Dockerfile)

# Skycoin development *vscode* image

This image has the necessary tools to build, test, edit, lint and version the Skycoin
source code. It comes with Visual Studio Code installed, along with some plugins
to ease go development and version control with git.

# Build

```sh
$ cd docker/images/dev-vscode/hoook
$ ./build.sh
```

When it finish you will have two new images:

`skycoin/skycoindev-vscode:develop` based on [skycoin/skycoindev-cli:develop](skycoin/docker/images/dev-cli) 
`skycoin/skycoindev-vscode:dind` based on [skycoin/skycoindev-cli:dind](skycoin/docker/images/dev-docker)

# How use this image

## Initialize your development environment.

0. Make sure you're on a system running [X](https://en.wikipedia.org/wiki/X_Window_System).
1. Disable X access control (don't do this on a public-facing machine): `$ xhost +` or `$ xhost +local:docker`
2. `$ cd` to a path where you want to write some code.
3. Run docker image
```sh
    $ docker run --rm -it -v /tmp/.X11-unix:/tmp/.X11-unix \
        -v $PWD:/go/src \
        -w /go/src \
        -e DISPLAY=$DISPLAY \
        skycoin/skycoindev-vscode:develop
```
5. You should see vscode pop up.
6. Have fun. Write some code. Close vscode when you're done, and ctrl+c to shut down the container. Your files will be in the path on the host where you started.
7. __Reenable X access control:__ `$ xhost -`

## Add more VS Code extensions

You must pass VS_EXTENSIONS environment variable to command-line with extensions of you prefer.

```sh
    $ docker run --rm -it -v /tmp/.X11-unix:/tmp/.X11-unix 
    -v $PWD:/go/src \
       -w /go/src \
       -e DISPLAY=$DISPLAY \
       -e VS_EXTENSIONS="ms-python.python rebornix.Ruby" \
       skycoin/skycoindev-vscode:dind
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

## Additional tools and packages installed

### Packages

- dep
- tig
- swig

### Vim's plugins

- Ale
- tig-explorer

### VS Code extensions

- [Go](https://marketplace.visualstudio.com/items?itemName=ms-vscode.Go)
- [Go Autotest](https://marketplace.visualstudio.com/items?itemName=windmilleng.vscode-go-autotest)
- [Go Coverage Viewer](https://marketplace.visualstudio.com/items?itemName=defaltd.go-coverage-viewer)

## Automated builds

Docker Cloud is configured to build images from `develop` branch on every push.
The same process is triggered for all feature branches matching the pattern
`/^([^_]+)_t([0-9]+)_.*vscode/`. The tag generated for those images will be of the form
`feature-{\1}-{\2}-vscode`.

