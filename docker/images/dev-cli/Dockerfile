# Creates an image for skycoin development
FROM golang:1.11-stretch

ARG BDATE="test_build"
ARG SCOMMIT="develop"

# Image labels (see ./hooks/build for ARGS)
LABEL "org.label-schema.name"="skycoindev-cli" \
      "org.label-schema.description"="Docker image with go, node and command line tools for Skycoin developers" \
      "org.label-schema.vendor"="Skycoin project" \
      "org.label-schema.url"="skycoin.com" \
      "org.label-schema.schema-version"="1.0" \
      "org.label-schema.build-date"=$BDATE \
      "org.label-schema.vcs-url"="https://github.com/SkycoinProject/skycoin.git" \
      "org.label-schema.vcs-ref"=$SCOMMIT \
      "org.label-schema.usage"="https://github.com/SkycoinProject/skycoin/blob/"$SCOMMIT"/docker/images/dev-cli/README.md" \
      "org.label-schema.version"="1.0.0-rc.1" \
      "org.label-schema.docker.cmd"="mkdir src; docker run --rm -v src:/go/src SkycoinProject/skycoindev-cli ; go get github.com/SkycoinProject/skycoin ; sudo chown -R `whoami` src"

# Installs nodejs and npm. Needed for moxygen.

# Packages installed in buildpack-deps:stretch
RUN set -ex; \
  apt-get update; \
  apt-get install -y --no-install-recommends \
    autoconf \
    automake \
    bzip2 \
    dpkg-dev \
    file \
    g++ \
    gcc \
    imagemagick \
    libbz2-dev \
    libc6-dev \
    libcurl4-openssl-dev \
    libdb-dev \
    libevent-dev \
    libffi-dev \
    libgdbm-dev \
    libgeoip-dev \
    libglib2.0-dev \
    libjpeg-dev \
    libkrb5-dev \
    liblzma-dev \
    libmagickcore-dev \
    libmagickwand-dev \
    libncurses5-dev \
    libncursesw5-dev \
    libpng-dev \
    libpq-dev \
    libreadline-dev \
    libsqlite3-dev \
    libssl-dev \
    libtool \
    libwebp-dev \
    libxml2-dev \
    libxslt-dev \
    libyaml-dev \
    make \
    patch \
    xz-utils \
    zlib1g-dev \
    build-essential \
    ruby \
    ruby-dev \
    \
# No need for MySQL client
#
# # https://lists.debian.org/debian-devel-announce/2016/09/msg00000.html
#     $( \
# # if we use just "apt-cache show" here, it returns zero because "Can't select versions from package 'libmysqlclient-dev' as it is purely virtual", hence the pipe to grep
#       if apt-cache show 'default-libmysqlclient-dev' 2>/dev/null | grep -q '^Version:'; then \
#         echo 'default-libmysqlclient-dev'; \
#       else \
#         echo 'libmysqlclient-dev'; \
#       fi \
#     ) \
   ; \
  apt-get clean; \
  rm -rf /var/lib/apt/lists/*

# Build steps in node:10 (uid=2000)
RUN groupadd --gid 2000 node \
  && useradd --uid 2000 --gid node --shell /bin/bash --create-home node

# gpg keys listed at https://github.com/nodejs/node#release-team
RUN set -ex \
  && for key in \
    94AE36675C464D64BAFA68DD7434390BDBE9B9C5 \
    FD3A5288F042B6850C66B31F09FE44734EB7990E \
    71DCFD284A79C3B38668286BC97EC7A07EDE3FC1 \
    DD8F2338BAE7501E3DD5AC78C273792F7D83545D \
    C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8 \
    B9AE9905FFD7803F25714661B63B535A4C206CA9 \
    56730D5401028683275BD23C23EFEFE93C4CFFFE \
    77984A986EBC2AA786BC0F66B01FBB92821C587A \
  ; do \
    gpg --no-tty --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys "$key" || \
    gpg --no-tty --keyserver hkp://ipv4.pool.sks-keyservers.net --recv-keys "$key" || \
    gpg --no-tty --keyserver hkp://pgp.mit.edu:80 --recv-keys "$key" ; \
  done

ENV NODE_VERSION 10.13.0

RUN ARCH= && dpkgArch="$(dpkg --print-architecture)" \
  && case "${dpkgArch##*-}" in \
    amd64) ARCH='x64';; \
    ppc64el) ARCH='ppc64le';; \
    s390x) ARCH='s390x';; \
    arm64) ARCH='arm64';; \
    armhf) ARCH='armv7l';; \
    i386) ARCH='x86';; \
    *) echo "unsupported architecture"; exit 1 ;; \
  esac \
  && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-$ARCH.tar.xz" \
  && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" \
  && gpg --batch --decrypt --output SHASUMS256.txt SHASUMS256.txt.asc \
  && grep " node-v$NODE_VERSION-linux-$ARCH.tar.xz\$" SHASUMS256.txt | sha256sum -c - \
  && tar -xJf "node-v$NODE_VERSION-linux-$ARCH.tar.xz" -C /usr/local --strip-components=1 --no-same-owner \
  && rm "node-v$NODE_VERSION-linux-$ARCH.tar.xz" SHASUMS256.txt.asc SHASUMS256.txt \
  && ln -s /usr/local/bin/node /usr/local/bin/nodejs

ENV YARN_VERSION 1.12.3

RUN set -ex \
  && for key in \
    6A010C5166006599AA17F08146C2130DFD2497F5 \
  ; do \
    gpg --no-tty --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys "$key" || \
    gpg --no-tty --keyserver hkp://ipv4.pool.sks-keyservers.net --recv-keys "$key" || \
    gpg --no-tty --keyserver hkp://pgp.mit.edu:80 --recv-keys "$key" ; \
  done \
  && curl -fsSLO --compressed "https://yarnpkg.com/downloads/$YARN_VERSION/yarn-v$YARN_VERSION.tar.gz" \
  && curl -fsSLO --compressed "https://yarnpkg.com/downloads/$YARN_VERSION/yarn-v$YARN_VERSION.tar.gz.asc" \
  && gpg --batch --verify yarn-v$YARN_VERSION.tar.gz.asc yarn-v$YARN_VERSION.tar.gz \
  && mkdir -p /opt \
  && tar -xzf yarn-v$YARN_VERSION.tar.gz -C /opt/ \
  && ln -s /opt/yarn-v$YARN_VERSION/bin/yarn /usr/local/bin/yarn \
  && ln -s /opt/yarn-v$YARN_VERSION/bin/yarnpkg /usr/local/bin/yarnpkg \
  && rm yarn-v$YARN_VERSION.tar.gz.asc yarn-v$YARN_VERSION.tar.gz

# Installs software
RUN set -ex ; \
    apt-get update ; \
    apt-get install -y --no-install-recommends \
    cmake \
    libpcre3-dev \
    gdbserver \
    gdb \
    vim \
    less \
    ctags \
    vim-scripts \
    screen \
    sudo \
    doxygen \
    valgrind \
    bsdmainutils \
    texlive-latex-base \
    ; \
    apt-get clean ; \
    rm -rf /var/lib/apt/lists/* ; \
    npm install moxygen -g ; \
    \
    \
    echo 'Installing Criterion ...' ; \
    git clone --recurse-submodules -j8 https://github.com/skycoin/Criterion /go/Criterion ; \
    cd /go/Criterion ; \
    cmake . ; \
    make install ; \
    rm -r /go/Criterion ; \
    echo 'Success nstalling Criterion ...'

# Installs go development tools
RUN go get -u github.com/derekparker/delve/cmd/dlv  && \
    go get -u github.com/FiloSottile/vendorcheck && \
    go get -u github.com/alecthomas/gometalinter && \
    gometalinter --vendored-linters --install && \
    go get -u github.com/zmb3/gogetdoc && \
    go get -u golang.org/x/tools/cmd/guru && \
    go get -u github.com/davidrjenni/reftools/cmd/fillstruct && \
    go get -u github.com/rogpeppe/godef && \
    go get -u github.com/fatih/motion && \
    go get -u github.com/nsf/gocode && \
    go get -u github.com/jstemmer/gotags && \
    go get -u github.com/josharian/impl && \
    go get -u github.com/fatih/gomodifytags && \
    go get -u github.com/dominikh/go-tools/cmd/keyify && \
    go get -u golang.org/x/tools/cmd/gorename && \
    go get -u github.com/klauspost/asmfmt/cmd/asmfmt && \
    go get -u github.com/vektra/mockery/.../ && \
    go get -u github.com/wadey/gocovmerge && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh


# Install vim-go development tools
RUN git clone https://github.com/fatih/vim-go /usr/share/vim/vim80/pack/dev/start/vim-go && \
    git clone https://github.com/tpope/vim-fugitive /usr/share/vim/vim80/pack/dev/start/vim-fugitive && \
    git clone https://github.com/Shougo/vimshell.vim /usr/share/vim/vim80/pack/dev/start/0vimshell && \
    git clone https://github.com/Shougo/vimproc.vim /usr/share/vim/vim80/pack/dev/start/0vimproc && \
    git clone https://github.com/w0rp/ale.git /usr/share/vim/vim80/pack/dev/start/ale && \
    cd /usr/share/vim/vim80/pack/dev/start/0vimproc && make ;\
    git clone https://github.com/iberianpig/tig-explorer.vim.git /tmp/tig-explorer;\
    cp /tmp/tig-explorer/autoload/tig_explorer.vim /usr/share/vim/vim80/autoload;\
    cp /tmp/tig-explorer/plugin/tig_explorer.vim /usr/share/vim/vim80/plugin;\
    rm -rf /tmp/tig-explorer


# Install SWIG-3.0.12
RUN cd /tmp/; \
    wget http://prdownloads.sourceforge.net/swig/swig-3.0.12.tar.gz && \
    tar -zxf swig-3.0.12.tar.gz ; \
    cd swig-3.0.12 ;\
    ./configure --prefix=/usr && make && make install && \
    rm -rf /tmp/swig-*
# Install Travis CLI
# Install golangci-lint
ENV GOLANGCI_LINT 1.12.3
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $GOPATH/bin v$GOLANGCI_LINT

# Install Travis CLI
RUN gem install travis

WORKDIR $GOPATH/src/github.com/skycoin
VOLUME $GOPATH/src/

ENV LD_LIBRARY_PATH=/usr/local/lib
