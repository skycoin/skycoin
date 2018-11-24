# Supported tags and respective `Dockerfile` links

## Simple Tags

-	[`develop` (*docker/images/dev-vscode/Dockerfile*)](https://github.com/skycoin/skycoin/tree/develop/docker/images/dev-vscode/Dockerfile)
-	[`vscode` (*docker/images/dev-vscode/Dockerfile*)](https://github.com/skycoin/skycoin/tree/develop/docker/images/dev-vscode/Dockerfile)

# Skycoin development *vscode* image

This image has the necessary tools to build, test, edit, lint and version the Skycoin
source code. It comes with Visual Studio Code installed, along with some plugins
to ease go development and version control with git.

# How use this image

## Initialize your development environment.

0. Make sure you're on a system running [X](https://en.wikipedia.org/wiki/X_Window_System).
1. Disable X access control (don't do this on a public-facing machine): `$ xhost +` or `$ xhost +local:docker`
2. `$ cd` to a path where you want to write some code.
3. Since Visual Studio Code inside docker container run as user `user`, it's necessary apply permissions to files.
    ```sh
    $ sudo chown -R 777 .
    ```
4. Run docker image.
    ```sh
    $ docker run --rm -it -v /tmp/.X11-unix:/tmp/.X11-unix \
            -v $PWD:/go/src/github.com/skycoin/skycoin \
            -w /go/src/github.com/skycoin/skycoin \
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
            -v $PWD:/go/src/github.com/skycoin/skycoin \
            -w /go/src/github.com/skycoin/skycoin \
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

# Build your own images

`SOURCE_COMMIT`: the SHA1 hash of the commit being tested.

`IMAGE_NAME`: the name and tag of the Docker repository being built.

`DOCKERFILE_PATH`: the dockerfile currently being built.

Build image from `skycoindev-cli:develop`.

```sh
$ git clone https://github.com/skycoin/skycoin
$ cd skycoin
$ SOURCE_COMMIT=$(git rev-parse HEAD)
$ IMAGE_NAME=skycoin/skycoindev-vscode:develop
$ DOCKERFILE_PATH=docker/images/dev-vscode/Dockerfile
$ docker build --build-arg BDATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
               --build-arg SCOMMIT=$SOURCE_COMMIT \
               --build-arg VS_EXTENSIONS="ms-vscode.Go windmilleng.vscode-go-autotest defaltd.go-coverage-viewer" \
               -f $DOCKERFILE_PATH \
               -t $IMAGE_NAME .
```

Or, if you prefer use `skycoindev-cli:dind`. Run:

```sh
$ git clone https://github.com/skycoin/skycoin
$ cd skycoin
$ SOURCE_COMMIT=$(git rev-parse HEAD)
$ IMAGE_NAME=skycoin/skycoindev-vscode:dind
$ DOCKERFILE_PATH=docker/images/dev-vscode/Dockerfile
$ docker build --build-arg IMAGE_FROM="skycoin/skycoindev-cli:dind" \
               --build-arg BDATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
               --build-arg SCOMMIT=$SOURCE_COMMIT \
               --build-arg VS_EXTENSIONS="ms-python.python rebornix.Ruby"
               -f $DOCKERFILE_PATH \
               -t IMAGE_NAME .
```

When it finish, you will have two new images:

`skycoin/skycoindev-vscode:develop` based on [skycoin/skycoindev-cli:develop](skycoin/docker/images/dev-cli) 
`skycoin/skycoindev-vscode:dind` based on [skycoin/skycoindev-cli:dind](skycoin/docker/images/dev-docker)

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

## How to use docker in docker image

### Start a daemon instance

```sh
$ docker run --privileged --name some-name -d skycoin/skycoindev-vscode:dind
```

### Where to store data

Create a data directory on the host system (outside the container) and mount this to a directory visible from inside the container.

The downside is that you need to make sure that the directory exists, and that e.g. directory permissions and other security mechanisms on the host system are set up correctly.

1. Create a data directory on a suitable volume on your host system, e.g. /my/own/var-lib-docker.
2. Start your docker container like this:

```sh
$ docker run --privileged --name some-name -v /my/own/var-lib-docker:/var/lib/docker \ 
-d skycoin/skycoindev-vscode:dind
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
- [TSLint](https://marketplace.visualstudio.com/items?itemName=eg2.tslint)
- [GitLens](https://marketplace.visualstudio.com/items?itemName=eamodio.gitlens)
- [Vim](https://marketplace.visualstudio.com/items?itemName=vscodevim.vim)

## Automated builds

Docker Cloud is configured to build images from `develop` branch on every push.
The same process is triggered for all feature branches matching the pattern
`/^([^_]+)_t([0-9]+)_.*vscode/`. The tag generated for those images will be of the form
`feature-{\1}-{\2}-vscode`.

