## `skycoin/skycoin:develop`

```console
$ docker pull skycoin/skycoin@sha256:0207b712c3df38a66bda4d7b355ab7dfc45e960b7d46f44ef4f40f454860ff68
```

- Manifest MIME: `application/vnd.docker.distribution.manifest.list.v2+json`
- Platforms:
	- linux; amd64

- Layers:
    - sha256:cd6101adb6f7d5c6f64f1bfb40ef0a63eea38274ac16c8c04bb949769962897b
    - sha256:9d396ec3c9f752238b6ce3459984560c56c5ebf99975e82ee2b4b68726de463d
    - sha256:30199736bff8fdef5d4caacbcba3c4478cd08b2d3c3328b2e15a1a192d851f31
    - sha256:d1d21de26fcccf58cafa0470142e38def9ed327c193a0d7212614e22dbeced15

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
