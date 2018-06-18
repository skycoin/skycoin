# Creates an image for skycoin development
FROM skycoin/skycoindev-cli:develop

# Install docker

RUN set -ex; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
    btrfs-progs \
    e2fsprogs \
    iptables \
    xfsprogs \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg2 \
    software-properties-common \
    pigz ;\
    apt-get clean

RUN curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
RUN apt-key fingerprint 0EBFCD88
RUN add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"

RUN set -ex; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
    docker-ce ;\
    apt-get clean

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
