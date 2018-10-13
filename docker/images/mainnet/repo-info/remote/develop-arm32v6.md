## `skycoin/skycoin:develop-arm32v6`

```console
$ docker pull skycoin/skycoin@sha256:8b5312f0f7c5096bbaaf3c0c028cc2c64a698e01abe5ac601c32498e856302d5
```

- Manifest MIME: `application/vnd.docker.distribution.manifest.list.v2+json`
- Platforms:
	- linux; amd64

- Layers:
    - sha256:3de7e3894033b9df2faed9182c17e7376f2081a1b0a62a55779e3fa0b3d6d49d
    - sha256:07e79c5a00e3b11d0b59e40df64d3a22696d4e633a53116eba30154876c89ea6
    - sha256:d8231cb38bd19e8578fa4cc2045999bdc14b3fedd4727be1488d5f9a5181596b
    - sha256:22e32bd4288e8a0f0c4cfc93b9ab8f87fe45986cd69c30887aa5094ca8aaa92c

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
