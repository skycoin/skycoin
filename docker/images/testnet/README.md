# Skycoin Docker Testnet

The Skycoin testnet provides an environment to test the interaction between
several skycoin nodes.

## Requirements

- Docker v17.x
- docker-compose v1.19.0

## Using the testnet

### Running the testnet

```sh
cd $GOPATH/src/github.com/skycoin/skycoin/docker/images/testnet
go run testnet.go
```

By default it will run 5 headless skycoin nodes using the source code as the
context. This behavior can be changed using the flags `--nodes` and
`--buildcontext`. For more information `go run testnet.go -h`

### Querying the logs

On start-up the testnet will display the following message:

```
Compose files will be copied to /tmp/skycointest<number>
```

It means that all the files, including the docker-compose.yml, will be at that
location. In order to see the logs stored in the oklog container, do the
following.

```sh
cd /tmp/skycointest<number>
docker-compose exec oklog oklog query
```
