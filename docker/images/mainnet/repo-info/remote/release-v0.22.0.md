## `skycoin/skycoin:develop`

```console
$ docker pull skycoin/skycoin@sha256:4b7e56145a7584d14273de62495e21622712fd07908ec2972fabaaf1fb0e7ca8
```

- Manifest MIME: `application/vnd.docker.distribution.manifest.list.v2+json`
- Platforms:
	- linux; amd64

- Layers:
    - sha256:ff3a5c916c92643ff77519ffa742d3ec61b7f591b6b7504599d95a4a41134e28
    - sha256:4ea2ebce37a98ced054ce64a83b8e2c38d4c63995799c76e43008b627bec12bf
    - sha256:bfe7dc706c0d6054d556e7a8635d22f58a67ee121c59b4edd9ea6ba611688893
    - sha256:516441122b6918902d3d3b085b312d430d16fca130415b6f26b2d361a43c4ba8

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
ENV RPC_ADDR=http://0.0.0.0:6420 DATA_DIR=/data/.skycoin WALLET_DIR=/wallet USE_CSRF=1 WALLET_NAME=.wlt
# 2018-07-16T22:19:41.841251284Z
ENV COIN=skycoin
# 2018-07-16T22:19:41.841251284Z
ADD file:2a4c44bdcb743a52ffa1c4b07ce471d8735a5d59cb45da2e6bfe0c2b5311ca90 in /
```
