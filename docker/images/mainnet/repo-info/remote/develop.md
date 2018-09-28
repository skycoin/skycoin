## `skycoin/skycoin:develop`

```console
$ docker pull skycoin/skycoin@sha256:940af212ee94d4da8b7dce2dbb4f2aece146f529c25e33908fbecc2bcfe2f425
```

- Manifest MIME: `application/vnd.docker.distribution.manifest.list.v2+json`
- Platforms:
	- linux; amd64

- Layers:
	- sha256:51c49aafd5a6f4bdf61a5a087ccb5e9cf616796d70dcc8ca62b7bf257473582e
    - sha256:c1b532ec52a9dd695a4a8dde5a56916f00433e1e8ccc8723a37fa50b2b3467f8
    - sha256:9f75ac52562ed718648d77a6705bd0b271fcc075fbec05085028ed34f1694f13
    - sha256:3b486c5a9b0fa0bb15ec4f952479279d9bbbe294b9ca834e97ccf4cd169f06db

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
