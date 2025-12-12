### Docker image

```sh
$ docker volume create skycoin-data
$ docker volume create skycoin-wallet
$ docker run -ti --rm \
    -v skycoin-data:/data/.skycoin \
    -v skycoin-wallet:/wallet \
    -p 6000:6000 \
    -p 6420:6420 \
    skycoin/skycoin
```

This image has a `skycoin` user for the skycoin daemon to run, with UID and GID 10000.
When you mount the volumes, the container will change their owner, so you
must be aware that if you are mounting an existing host folder any content you
have there will be own by 10000.

The container will run with some default options, but you can change them
by just appending flags at the end of the `docker run` command. The following
example will show you the available options.

```sh
$ docker run --rm skycoin/skycoin -help
```

Access the API: [http://localhost:6420/version](http://localhost:6420/version).

### Building your own images

[Building your own images](docker/images/mainnet/README.md).

### Development image

The [skycoin/skycoindev-cli docker image](docker/images/dev-cli/README.md) is provided in order to make
easy to start developing Skycoin. It comes with the compiler, linters, debugger
and the vim editor among other tools.

The [skycoin/skycoindev-dind docker image](docker/images/dev-docker/README.md) comes with docker installed
and all tools available on `skycoin/skycoindev-cli:develop` docker image.

Also, the [skycoin/skycoindev-vscode docker image](docker/images/dev-vscode/README.md) is provided
to facilitate the setup of the development process with [Visual Studio Code](https://code.visualstudio.com)
and useful tools included in `skycoin/skycoindev-cli`.
