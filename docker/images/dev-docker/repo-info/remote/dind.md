# `skycoin/skycoindev-cli:dind`

```console
$ docker pull skycoin/skycoindev-cli@sha256:059f9228a6dfe11c08e475b383cd851edbbb2d11ee766cc4067e329f3b6ce5c2
```

- Manifest MIME: `application/vnd.docker.distribution.manifest.v2+json`

- Platform: 
	- linux, amd64

- Layers:
	- sha256:f715ed19c28b66943ac8bc12dbfb828e8394de2530bbaf1ecce906e748e4fdff
	- sha256:8bb25f9cdc41e7d085033af15a522973b44086d6eedd24c11cc61c9232324f77
	- sha256:08a01612ffca33483a1847c909836610610ce523fb7e1aca880140ee84df23e9
	- sha256:1191b3f5862aa9231858809b7ac8b91c0b727ce85c9b3279932f0baacc92967d
	- sha256:52da4d4dcf59a2a01a4c4928516bf922944775862547dad6142bc477bbbdfa47
	- sha256:757fc57926be167de1ddcffd80c1cfafb412f460ece41fd3204e903c662538fd
	- sha256:7a7147543a13ba6e9fd6c49b7fdb9ee2599833c0cb5065fe10c993f64203a9e0
	- sha256:f314cb931d880c2f734ba06dce3f86ed17983d9e54c3d03b7ea1b0a4351b77c4
	- sha256:7602206fe527a9fecc5fde419e40c43b54c8f3699478a1a730b857be75da5b97
	- sha256:5f8d1d941448905f6e71e0b5cd7acc59433277b4d6bab80215bd08815000ca30
	- sha256:c1536c09a721f4e3ba6a4a007303dfbc0725780f701527c56f976050aab7444b
	- sha256:728023b74bb899c7f7ae2aa14e2916e9bc197c23ae84f4fb2741b47ff420fc45
	- sha256:e56e2b829b439d85f22e58909a7a41e93fc2e0ffbb590855f3b02b7bf2d23d51
	- sha256:1a71e5504bbbbc69c81032dc2553f998105b726babcfa536b987042eb74f4fa4
	- sha256:672cd8f534572da2590e9b4bc4b532b043780a72ccf9c97361b7e971e640bc2c
	- sha256:43e02ea201e8c01c0987efdabd7339f9cc18564d5feba145db37e6507169c021
	- sha256:37670bea79f6e9fc267f189bc1e97fd16832991a8c701e391f2c7991704964b6
	- sha256:17f311ecbf47425c4759f4d7c9d9e849815918533011e3bf6dd7912ac5bfb123
	- sha256:8ef5a41e28f1f051b38c6d80afda404b36a3cf56a78400a6d64d8aa79d5711e4
	- sha256:bca28110d5dae1ebe5eeb85275a61ed9e676cb6b65f6ddb3cbd0148488bcaae5
	- sha256:97f47526eb1c11e1eef0e46e4ac3110ea951916667d2bf1ad382dc5c7195523a
	- sha256:6a20cf4f1fc294c9ee8b995b352007adf480de9f47cf284f306cf36563bd772b
	- sha256:8b6cbd1ad0315808fc430f8069948693747ff897e71e1a24f27493d43d04df38
	- sha256:3b5a388c7799e4c4c04b65af3ba73cea70ff3c5308e855ee22a5cc1d8a1cce66
	- sha256:569c0813afae170750483916ec6a39d0612b5c2d4a1c9d277d4ac8519785dbe7
	- sha256:1e6d11189cc735a30839e4fba869f0ef41de25c363a3590de0b950db020b6a0c
	- sha256:acbee7f1aebed815da498a47ff8da4dbb1862f72ed474db8c2fe66fb306206dd
	- sha256:35c42c1d542d3599f7e9d9eb8fc725a7ee23d53472ce5d2fb367526b33164ab3

- Exposed Ports:
	- 2375/tcp

```dockerfile
# 2018-10-15T23:24:20.7838109Z
ADD file:b3598c18dc395846ab2c5e4e8422c4a5dad7bc3b5b08c09ebceee80989904641 in / 
# 2018-10-15T23:24:21.111611143Z
 CMD ["bash"]
# 2018-10-16T01:00:49.945915631Z
/bin/sh -c apt-get update && apt-get install -y --no-install-recommends ca-certificates curl netbase wget && rm -rf /var/lib/apt/lists/*
# 2018-10-16T01:00:57.258682859Z
/bin/sh -c set -ex; if ! command -v gpg > /dev/null; then apt-get update; apt-get install -y --no-install-recommends gnupg dirmngr ; rm -rf /var/lib/apt/lists/*; fi
# 2018-10-16T01:01:29.418111504Z
/bin/sh -c apt-get update && apt-get install -y --no-install-recommends bzr git mercurial openssh-client subversion procps && rm -rf /var/lib/apt/lists/*
# 2018-10-16T08:42:27.8302524Z
/bin/sh -c apt-get update && apt-get install -y --no-install-recommends g++ gcc libc6-dev make pkg-config && rm -rf /var/lib/apt/lists/*
# 2018-10-16T08:42:28.211440772Z
 ENV GOLANG_VERSION=1.11.1
# 2018-10-16T08:42:41.25124204Z
/bin/sh -c set -eux; dpkgArch="$(dpkg --print-architecture)"; case "${dpkgArch##*-}" in amd64) goRelArch='linux-amd64'; goRelSha256='2871270d8ff0c8c69f161aaae42f9f28739855ff5c5204752a8d92a1c9f63993' ;; armhf) goRelArch='linux-armv6l'; goRelSha256='bc601e428f458da6028671d66581b026092742baf6d3124748bb044c82497d42' ;; arm64) goRelArch='linux-arm64'; goRelSha256='25e1a281b937022c70571ac5a538c9402dd74bceb71c2526377a7e5747df5522' ;; i386) goRelArch='linux-386'; goRelSha256='52935db83719739d84a389a8f3b14544874fba803a316250b8d596313283aadf' ;; ppc64el) goRelArch='linux-ppc64le'; goRelSha256='f929d434d6db09fc4c6b67b03951596e576af5d02ff009633ca3c5be1c832bdd' ;; s390x) goRelArch='linux-s390x'; goRelSha256='93afc048ad72fa2a0e5ec56bcdcd8a34213eb262aee6f39a7e4dfeeb7e564c9d' ;; *) goRelArch='src'; goRelSha256='558f8c169ae215e25b81421596e8de7572bd3ba824b79add22fba6e284db1117'; echo >&2; echo >&2 "warning: current architecture ($dpkgArch) does not have a corresponding Go binary release; will be building from source"; echo >&2 ;; esac; url="https://golang.org/dl/go${GOLANG_VERSION}.${goRelArch}.tar.gz"; wget -O go.tgz "$url"; echo "${goRelSha256} *go.tgz" | sha256sum -c -; tar -C /usr/local -xzf go.tgz; rm go.tgz; if [ "$goRelArch" = 'src' ]; then echo >&2; echo >&2 'error: UNIMPLEMENTED'; echo >&2 'TODO install golang-any from jessie-backports for GOROOT_BOOTSTRAP (and uninstall after build)'; echo >&2; exit 1; fi; export PATH="/usr/local/go/bin:$PATH"; go version
# 2018-10-16T08:42:41.931278427Z
 ENV GOPATH=/go
# 2018-10-16T08:42:42.218913506Z
 ENV PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
# 2018-10-16T08:42:43.147661979Z
/bin/sh -c mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
# 2018-10-16T08:42:43.432359416Z
WORKDIR /go
# 2018-10-17T07:37:02.17339676Z
/bin/sh -c set -ex; apt-get update; apt-get install -y --no-install-recommends autoconf automake bzip2 dpkg-dev file g++ gcc imagemagick libbz2-dev libc6-dev libcurl4-openssl-dev libdb-dev libevent-dev libffi-dev libgdbm-dev libgeoip-dev libglib2.0-dev libjpeg-dev libkrb5-dev liblzma-dev libmagickcore-dev libmagickwand-dev libncurses5-dev libncursesw5-dev libpng-dev libpq-dev libreadline-dev libsqlite3-dev libssl-dev libtool libwebp-dev libxml2-dev libxslt-dev libyaml-dev make patch xz-utils zlib1g-dev build-essential ; apt-get clean; rm -rf /var/lib/apt/lists/*
# 2018-10-17T07:37:04.040629881Z
/bin/sh -c groupadd --gid 2000 node && useradd --uid 2000 --gid node --shell /bin/bash --create-home node
# 2018-10-17T07:38:11.439532764Z
/bin/sh -c set -ex && for key in 94AE36675C464D64BAFA68DD7434390BDBE9B9C5 FD3A5288F042B6850C66B31F09FE44734EB7990E 71DCFD284A79C3B38668286BC97EC7A07EDE3FC1 DD8F2338BAE7501E3DD5AC78C273792F7D83545D C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8 B9AE9905FFD7803F25714661B63B535A4C206CA9 56730D5401028683275BD23C23EFEFE93C4CFFFE 77984A986EBC2AA786BC0F66B01FBB92821C587A ; do gpg --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys "$key" || gpg --keyserver hkp://ipv4.pool.sks-keyservers.net --recv-keys "$key" || gpg --keyserver hkp://pgp.mit.edu:80 --recv-keys "$key" ; done
# 2018-10-17T07:38:11.63067098Z
 ENV NODE_VERSION=10.2.1
# 2018-10-17T07:38:20.699035146Z
/bin/sh -c ARCH= && dpkgArch="$(dpkg --print-architecture)" && case "${dpkgArch##*-}" in amd64) ARCH='x64';; ppc64el) ARCH='ppc64le';; s390x) ARCH='s390x';; arm64) ARCH='arm64';; armhf) ARCH='armv7l';; i386) ARCH='x86';; *) echo "unsupported architecture"; exit 1 ;; esac && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-$ARCH.tar.xz" && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" && gpg --batch --decrypt --output SHASUMS256.txt SHASUMS256.txt.asc && grep " node-v$NODE_VERSION-linux-$ARCH.tar.xz\$" SHASUMS256.txt | sha256sum -c - && tar -xJf "node-v$NODE_VERSION-linux-$ARCH.tar.xz" -C /usr/local --strip-components=1 --no-same-owner && rm "node-v$NODE_VERSION-linux-$ARCH.tar.xz" SHASUMS256.txt.asc SHASUMS256.txt && ln -s /usr/local/bin/node /usr/local/bin/nodejs
# 2018-10-17T07:38:21.153369203Z
 ENV YARN_VERSION=1.7.0
# 2018-10-17T07:38:25.479495596Z
/bin/sh -c set -ex && for key in 6A010C5166006599AA17F08146C2130DFD2497F5 ; do gpg --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys "$key" || gpg --keyserver hkp://ipv4.pool.sks-keyservers.net --recv-keys "$key" || gpg --keyserver hkp://pgp.mit.edu:80 --recv-keys "$key" ; done && curl -fsSLO --compressed "https://yarnpkg.com/downloads/$YARN_VERSION/yarn-v$YARN_VERSION.tar.gz" && curl -fsSLO --compressed "https://yarnpkg.com/downloads/$YARN_VERSION/yarn-v$YARN_VERSION.tar.gz.asc" && gpg --batch --verify yarn-v$YARN_VERSION.tar.gz.asc yarn-v$YARN_VERSION.tar.gz && mkdir -p /opt && tar -xzf yarn-v$YARN_VERSION.tar.gz -C /opt/ && ln -s /opt/yarn-v$YARN_VERSION/bin/yarn /usr/local/bin/yarn && ln -s /opt/yarn-v$YARN_VERSION/bin/yarnpkg /usr/local/bin/yarnpkg && rm yarn-v$YARN_VERSION.tar.gz.asc yarn-v$YARN_VERSION.tar.gz
# 2018-10-17T07:41:44.706504732Z
/bin/sh -c set -ex ; apt-get update ; apt-get install -y --no-install-recommends cmake libpcre3-dev gdbserver gdb vim less ctags vim-scripts screen sudo doxygen valgrind bsdmainutils texlive-latex-base ; apt-get clean ; rm -rf /var/lib/apt/lists/* ; npm install moxygen -g ; echo 'Installing Criterion ...' ; git clone --recurse-submodules -j8 https://github.com/skycoin/Criterion /go/Criterion ; cd /go/Criterion ; cmake . ; make install ; rm -r /go/Criterion ; echo 'Success nstalling Criterion ...'
# 2018-10-17T07:44:41.078027675Z
/bin/sh -c go get -u github.com/derekparker/delve/cmd/dlv && go get -u github.com/FiloSottile/vendorcheck && go get -u github.com/alecthomas/gometalinter && gometalinter --vendored-linters --install && go get -u github.com/zmb3/gogetdoc && go get -u golang.org/x/tools/cmd/guru && go get -u github.com/davidrjenni/reftools/cmd/fillstruct && go get -u github.com/rogpeppe/godef && go get -u github.com/fatih/motion && go get -u github.com/nsf/gocode && go get -u github.com/jstemmer/gotags && go get -u github.com/josharian/impl && go get -u github.com/fatih/gomodifytags && go get -u github.com/dominikh/go-tools/cmd/keyify && go get -u golang.org/x/tools/cmd/gorename && go get -u github.com/klauspost/asmfmt/cmd/asmfmt && go get -u github.com/vektra/mockery/.../ && go get -u github.com/wadey/gocovmerge && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# 2018-10-17T07:44:51.741341123Z
/bin/sh -c git clone https://github.com/fatih/vim-go /usr/share/vim/vim80/pack/dev/start/vim-go && git clone https://github.com/tpope/vim-fugitive /usr/share/vim/vim80/pack/dev/start/vim-fugitive && git clone https://github.com/Shougo/vimshell.vim /usr/share/vim/vim80/pack/dev/start/0vimshell && git clone https://github.com/Shougo/vimproc.vim /usr/share/vim/vim80/pack/dev/start/0vimproc && git clone https://github.com/w0rp/ale.git /usr/share/vim/vim80/pack/dev/start/ale && cd /usr/share/vim/vim80/pack/dev/start/0vimproc && make ; git clone https://github.com/iberianpig/tig-explorer.vim.git /tmp/tig-explorer; cp /tmp/tig-explorer/autoload/tig_explorer.vim /usr/share/vim/vim80/autoload; cp /tmp/tig-explorer/plugin/tig_explorer.vim /usr/share/vim/vim80/plugin; rm -rf /tmp/tig-explorer
# 2018-10-17T07:46:41.184461422Z
/bin/sh -c cd /tmp/; wget http://prdownloads.sourceforge.net/swig/swig-3.0.12.tar.gz && tar -zxf swig-3.0.12.tar.gz ; cd swig-3.0.12 ; ./configure --prefix=/usr && make && make install && rm -rf /tmp/swig-*
# 2018-10-17T07:46:41.378468069Z
 ENV GOLANGCI_LINT=1.10.2
# 2018-10-17T07:46:44.541327017Z
/bin/sh -c curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $GOPATH/bin v$GOLANGCI_LINT
# 2018-10-17T07:46:44.761903317Z
WORKDIR /go/src/github.com/skycoin
# 2018-10-17T07:46:45.9222658Z
 VOLUME [/go/src/]
# 2018-10-17T07:46:47.385077348Z
 ENV LD_LIBRARY_PATH=/usr/local/lib
# 2018-10-21T05:42:03.135569523Z
/bin/sh -c curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
# 2018-10-21T05:42:06.775092982Z
/bin/sh -c apt-key fingerprint 0EBFCD88
# 2018-10-21T05:42:24.441063936Z
/bin/sh -c set -ex; apt-get update; apt-get install -y --no-install-recommends lsb-release software-properties-common apt-transport-https
# 2018-10-21T05:42:27.881882246Z
/bin/sh -c add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
# 2018-10-21T05:42:55.252560207Z
/bin/sh -c set -ex; apt-get update; apt-get install -y --no-install-recommends btrfs-progs e2fsprogs iptables xfsprogs ca-certificates gnupg2 software-properties-common pigz docker-ce ; apt-get clean; rm -rf /var/lib/apt/lists/*
# 2018-10-21T05:42:55.825119189Z
COPY file:7542b8556b602563f72a47a7602958d8aaa7570b8d1decd67f014edcd7840d56 in /usr/local/bin/modprobe 
# 2018-10-21T05:42:57.910836374Z
/bin/sh -c set -x && groupadd dockremap && useradd -g dockremap dockremap && echo 'dockremap:165536:65536' >> /etc/subuid && echo 'dockremap:165536:65536' >> /etc/subgid
# 2018-10-21T05:42:59.349476646Z
 ENV DIND_COMMIT=3b5fac462d21ca164b3778647420016315289034
# 2018-10-21T05:43:01.720664356Z
/bin/sh -c set -ex; wget -O /usr/local/bin/dind "https://raw.githubusercontent.com/docker/docker/${DIND_COMMIT}/hack/dind"; chmod +x /usr/local/bin/dind;
# 2018-10-21T05:43:03.764273308Z
COPY file:8c7efafc9ff2ddd0b88764e8647f34ab4c5bca079a93694700343d4603f7f8a6 in /usr/local/bin/ 
# 2018-10-21T05:43:05.718390432Z
chmod +x /usr/local/bin/dockerd-entrypoint.sh /usr/local/bin/modprobe
# 2018-10-21T05:43:07.15905469Z
 VOLUME [/var/lib/docker]
# 2018-10-21T05:43:08.561961432Z
 EXPOSE 2375
# 2018-10-21T05:43:09.989421422Z
 ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh"]
# 2018-10-21T05:43:11.428595401Z
 CMD []
```

