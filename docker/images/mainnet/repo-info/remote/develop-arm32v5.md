## `skycoin/skycoin:develop-arm32v5`

```console
$ docker pull skycoin/skycoin@sha256:2e95c6d7a10cc4f69821c79d3dee8ac1c1fb329949796c2b04f81faddef8b4f0
```

- Manifest MIME: `application/vnd.docker.distribution.manifest.list.v2+json`
- Platforms:
	- linux; amd64

- Layers:
    - sha256:b4f45fb2b1d01c55a010dc4585efa7c58022437c15734b329df6afdd7670005b
    - sha256:b3a3ece1ae8b8bd197c940a7acde1956f3f9138008ca948a1fc25b92991f02dd
    - sha256:52546602f55a7fcd444a5942e0429cdfd30e7f718ff422a7b2856c37302a206a
    - sha256:03169c874cfd04b8604b67c33ee55e1b9b69d95e19bf83f4869ba59692fa4b89

- Expose Ports:
	- 6000
    - 6420

```dockerfile
# 2018-07-22T07:53:54.681656183Z
EXPOSE 6000 6420
# 2018-07-22T07:53:54.450638675Z
VOLUME [/data/.skycoin]
# 2018-07-22T07:53:54.450638675Z
VOLUME [/wallet]
# 2018-07-20T01:31:36.608478714Z
COPY file:6ac857b94e8b21cfa7f4c9a4d19387c91ec0b0eeb0faf318a16758e7c280e791 in /usr/local/bin/docker_launcher.sh
# 2018-07-20T01:31:35.779110731Z
COPY dir:0a4f98c7af3e020a45ac06413d1f1cb6409bd9ef2ba1546d2a4970fb73bc8c31 in /usr/local/skycoin/src/gui/static
# 2018-07-16T22:19:41.841251284Z
COPY multi:d033726808550b3bf4ec4dc28a2156e03a05e265d8e928b8762a8d0ad1f2583e in /usr/bin/
# 2018-07-16T22:19:41.841251284Z
ENV RPC_ADDR=http://0.0.0.0:6420 DATA_DIR=/data/.skycoin WALLET_DIR=/wallet WALLET_NAME=.wlt
# 2018-07-16T22:19:41.841251284Z
ENV COIN=skycoin
# 2018-07-16T22:19:41.841251284Z
ADD file:2a4c44bdcb743a52ffa1c4b07ce471d8735a5d59cb45da2e6bfe0c2b5311ca90 in /
```
