
## Getting started

### Build image

```
$ docker build -t skycoin .
```

### Running

```
$ docker run -ti --rm \
    -p 6000:6000 \
    -p 6420:6420 \
    -p 6430:6430 \
    skycoin
```

Access the dashboard: [http://localhost:6420](http://localhost:6420).

Access the API: [http://localhost:6420/version](http://localhost:6420/version).

### Data persistency

```
$ docker volume create skycoin-data
$ docker volume create skycoin-wallet
$ docker run -ti --rm \
    -v skycoin-data:/root/.skycoin \
    -v skycoin-wallet:/wallet \
    -p 6000:6000 \
    -p 6420:6420 \
    -p 6430:6430 \
    skycoin
```

### API

https://github.com/skycoin/skycoin/blob/develop/src/gui/README.md

https://github.com/skycoin/skycoin/blob/v0.21.1/src/api/webrpc/README.md
