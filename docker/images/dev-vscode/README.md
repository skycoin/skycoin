# Supported tags and respective `Dockerfile` links

## Simple Tags

-	[`develop` (*docker/images/dev-vscode/Dockerfile*)](https://github.com/SkycoinProject/skycoin/tree/develop/docker/images/dev-vscode/Dockerfile)
-	[`vscode` (*docker/images/dev-vscode/Dockerfile*)](https://github.com/SkycoinProject/skycoin/tree/develop/docker/images/dev-vscode/Dockerfile)

# Skycoin Docker image for development with [VS Code](https://code.visualstudio.com/) IDE

This image has the necessary tools to build, test, edit, lint and version the Skycoin
source code. It comes with [Visual Studio Code](https://code.visualstudio.com/) installed and some extensions
ot speed up workspace setup for Skycoin developers.

# How use this image

## Initialize your development environment.

0. Make sure you're on a system running [X](https://en.wikipedia.org/wiki/X_Window_System).
  - *GNU/Linux* users should be ready to go
  - *Mac OS* users can follow the following steps
    * `brew install xquartz`
    * **Important**  Log out and log back into OS/X
1. Disable X access control (don't do this on a public-facing machine):
  - *GNU/Linux* : `$ xhost +` or `$ xhost +local:docker`
  - *Mac OS* : XQuartz users should follow these steps
    * `open -a Xquartz`
    * With `xterm` active, open up `XQuartz` in menu bar => `Preferences` => `Security`. There make sure the `Allow connections from network clients` is checked `on`.
2. `$ cd` to a path where you want to write some code (e.g. a working copy of [`SkycoinProject/skycoin`](https://github.com/SkycoinProject/skycoin) )
3. Since Visual Studio Code inside docker container runs as user `skydev`, it's necessary apply permissions to files.
    ```sh
    $ sudo chown -R 777 .
    ```
4. Run docker image, either `SkycoinProject/skycoindev-vscode:develop` or `SkycoinProject/skycoindev-vscode:dind`
  - *GNU/Linux*
    ```sh
    $ docker run --rm -it -v /tmp/.X11-unix:/tmp/.X11-unix \
            -v $(pwd):$GOPATH/src/github.com/SkycoinProject/skycoin \
            -w $GOPATH/src/github.com/SkycoinProject/skycoin \
            -e DISPLAY=$DISPLAY \
            SkycoinProject/skycoindev-vscode:develop
    ```
  - *Mac OS* users running XQuartz should launch `socat` for Docker to be able to connect to the X server. Assuming `en0` is your primary network interface
    ```sh
    $ brew install socat
    $ socat TCP-LISTEN:6000,reuseaddr,fork UNIX-CLIENT:\"$DISPLAY\"
    $ export IP=$(ifconfig en0 | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p')
    $ docker run --rm -it -v /tmp/.X11-unix:/tmp/.X11-unix \
            -v $(pwd):$GOPATH/src/github.com/SkycoinProject/skycoin \
            -w $GOPATH/src/github.com/SkycoinProject/skycoin \
            -e DISPLAY=$IP:0 \
            SkycoinProject/skycoindev-vscode:develop
    ```
5. You should see vscode pop up.
6. Have fun. Write some code. Close VS Code IDE window when you're done, and press `Ctrl+C` to shut down the container. Your files will be in the host machine at the same path chosen in step `2` above.
7. __Reenable X access control:__ `$ xhost -`

For the sake of brevity, the examples that follow only include the invocation of `docker` command for GNU/Linux. Beware of the fact that there will be differences in running the Docker images on other operating systems. The hints provided above will still be valid, though.

## Add more VS Code extensions

If you want add more extensions, you must define `VS_EXTENSIONS` environment variable to the command-line with extensions you prefer.

```sh
    $ docker run --rm -it -v /tmp/.X11-unix:/tmp/.X11-unix 
            -v $PWD:/go/src/github.com/SkycoinProject/skycoin \
            -w $GOPATH/src/github.com/SkycoinProject/skycoin \
            -e DISPLAY=$DISPLAY \
            -e VS_EXTENSIONS="ms-python.python rebornix.Ruby" \
            SkycoinProject/skycoindev-vscode:dind
```

This downloads the skycoin source to src/SkycoinProject/skycoin and changes the owner
to your user. This is necessary, because all processes inside the container run
as root and the files created by it are therefore owned by root.

If you already have a Go development environment installed, you just need to
mount the src directory from your `$GOPATH` at path `/go/src` inside the
container.

```sh
    $ docker run --rm -it -v /tmp/.X11-unix:/tmp/.X11-unix 
            -v $GOPATH/src:$GOPATH/src \
            -w $GOPATH/src/github.com/SkycoinProject/skycoin \
            -e DISPLAY=$DISPLAY \
            -e VS_EXTENSIONS="ms-python.python rebornix.Ruby" \
            SkycoinProject/skycoindev-vscode:dind
```

# Build your own images

The following arguments influence the Docker build process.

- `SOURCE_COMMIT`: the SHA1 hash of the commit being tested.
- `IMAGE_NAME`: the name and tag of the Docker repository being built.
- `DOCKERFILE_PATH`: the dockerfile currently being built.

For instance, the following commands can be executed in order to build this VS Code dev image using `skycoindev-cli:develop` as base image.

```sh
$ git clone https://github.com/SkycoinProject/skycoin
$ cd skycoin
$ SOURCE_COMMIT=$(git rev-parse HEAD)
$ IMAGE_NAME=SkycoinProject/skycoindev-vscode:develop
$ DOCKERFILE_PATH=docker/images/dev-vscode/Dockerfile
$ docker build --build-arg BDATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
               --build-arg SCOMMIT=$SOURCE_COMMIT \
               --build-arg VS_EXTENSIONS="ms-vscode.Go windmilleng.vscode-go-autotest defaltd.go-coverage-viewer" \
               -f $DOCKERFILE_PATH \
               -t "$IMAGE_NAME" .
```

Or, if a decision has been made for including a Docker daemon then specify `skycoindev-cli:dind` instead and run:

```sh
$ git clone https://github.com/SkycoinProject/skycoin
$ cd skycoin
$ SOURCE_COMMIT=$(git rev-parse HEAD)
$ IMAGE_NAME=SkycoinProject/skycoindev-vscode:dind
$ DOCKERFILE_PATH=docker/images/dev-vscode/Dockerfile
$ docker build --build-arg IMAGE_FROM="SkycoinProject/skycoindev-cli:dind" \
               --build-arg BDATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
               --build-arg SCOMMIT=$SOURCE_COMMIT \
               --build-arg VS_EXTENSIONS="ms-python.python rebornix.Ruby"
               -f $DOCKERFILE_PATH \
               -t "$IMAGE_NAME" .
```

As a result of following theses steps two new images will be obtained:

`SkycoinProject/skycoindev-vscode:develop` based on [SkycoinProject/skycoindev-cli:develop](skycoin/docker/images/dev-cli) 
`SkycoinProject/skycoindev-vscode:dind` based on [SkycoinProject/skycoindev-cli:dind](skycoin/docker/images/dev-docker)

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

### Editing code with terminal

Comman line tools are still available . For instance it's possible to run `vim`

```sh
$ docker run --rm \
    -v src:/go/src SkycoinProject/skycoindev-cli \
    vim
```

## How to use docker in docker image

### Start a daemon instance

```sh
$ docker run --privileged --name some-name -d SkycoinProject/skycoindev-vscode:dind
```

### Where to store data

Create a data directory on the host system (outside the container) and mount this to a directory visible from inside the container.

The downside is that you need to make sure that the directory exists, and that e.g. directory permissions and other security mechanisms on the host system are set up correctly.

1. Create a data directory on a suitable volume on your host system, e.g. /my/own/var-lib-docker.
2. Start your docker container like this:

```sh
$ docker run --privileged --name some-name -v /my/own/var-lib-docker:/var/lib/docker \ 
-d SkycoinProject/skycoindev-vscode:dind
```

## Additional tools and packages installed

### Packages

- dep
- tig
- swig

### Vim plugins

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

