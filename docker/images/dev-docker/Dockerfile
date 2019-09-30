# Creates an image for skycoin development with docker in docker
FROM SkycoinProject/skycoindev-cli:develop
ARG BDATE
ARG SCOMMIT

# Image labels
LABEL "org.label-schema.name"="skycoindev-cli:dind" \
      "org.label-schema.description"="Skycoin cli develop image with docker in docker support" \
      "org.label-schema.vcs-url"="https://github.com/SkycoinProject/skycoin/tree/develop/docker/images/dev-docker" \
      "org.label-schema.vendor"="Skycoin project" \
      "org.label-schema.url"="skycoin.com" \
      "org.label-schema.schema-version"="1.0" \
      "org.label-schema.build-date"=$BDATE \
      "org.label-schema.vcs-ref"=$SCOMMIT \
      "org.label-schema.version"="1.0.0-rc.1" \
      "org.label-schema.usage"="https://github.com/SkycoinProject/skycoin/blob/"$SCOMMIT"/docker/images/dev-docker/README.md" \
      "org.label-schema.docker.cmd"="mkdir src; docker run --privileged --rm -v src:/go/src SkycoinProject/skycoindev-cli:dind go get github.com/SkycoinProject/skycoin; sudo chown -R `whoami` src"

# Install docker

RUN curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
RUN apt-key fingerprint 0EBFCD88
RUN set -ex; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
    lsb-release \
    software-properties-common \
    apt-transport-https ;\
    apt clean

RUN add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"

RUN set -ex; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
    btrfs-progs \
    e2fsprogs \
    iptables \
    xfsprogs \
    ca-certificates \
    gnupg2 \
    software-properties-common \
    pigz \
    docker-ce ;\
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*

COPY modprobe.sh /usr/local/bin/modprobe

RUN set -x \
	&& groupadd dockremap \
	&& useradd -g dockremap dockremap \
	&& echo 'dockremap:165536:65536' >> /etc/subuid \
	&& echo 'dockremap:165536:65536' >> /etc/subgid

ENV DIND_COMMIT 3b5fac462d21ca164b3778647420016315289034

RUN set -ex; \
	wget -O /usr/local/bin/dind "https://raw.githubusercontent.com/docker/docker/${DIND_COMMIT}/hack/dind"; \
	chmod +x /usr/local/bin/dind;

COPY dockerd-entrypoint.sh /usr/local/bin/

RUN ["chmod", "+x", "/usr/local/bin/dockerd-entrypoint.sh","/usr/local/bin/modprobe"]

VOLUME /var/lib/docker

EXPOSE 2375

#WORKDIR $GOPATH/src/github.com/skycoin
#VOLUME $GOPATH/src/

#ENV LD_LIBRARY_PATH=/usr/local/lib

ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh"]
CMD []
