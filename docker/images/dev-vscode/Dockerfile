ARG IMAGE_FROM=SkycoinProject/skycoindev-cli:develop
FROM $IMAGE_FROM

ARG BDATE
ARG SCOMMIT
ARG VS_EXTENSIONS

# Image labels (see ./hooks/build for ARGS)
LABEL "org.label-schema.name"="skycoindev-cli" \
      "org.label-schema.description"="Docker image with go, node, dev tools and Visual Studio Code for Skycoin developers" \
      "org.label-schema.vendor"="Skycoin project" \
      "org.label-schema.url"="skycoin.com" \
      "org.label-schema.version"="1.0.0-rc.1" \
      "org.label-schema.schema-version"="1.0" \
      "org.label-schema.build-date"=$BDATE \
      "org.label-schema.vcs-url"="https://github.com/SkycoinProject/skycoin.git" \
      "org.label-schema.vcs-ref"=$SCOMMIT \
      "org.label-schema.usage"="https://github.com/SkycoinProject/skycoin/blob/"$SCOMMIT"/docker/images/dev-vscode/README.md" \
      "org.label-schema.docker.cmd"="xhost +; cd src; docker run --rm -it -v /tmp/.X11-unix:/tmp/.X11-unix -v $PWD:/go/src -w /go/src -e DISPLAY=$DISPLAY SkycoinProject/skycoindev-vscode:develop"

# Tell debconf to run in non-interactive mode
ENV DEBIAN_FRONTEND noninteractive

# Create a diferent user to run VS Code
ENV HOME /home/skydev
RUN useradd --create-home --home-dir $HOME skydev \
	&& chown -R skydev:skydev $HOME

# Install dependencies for vs code
# Create and assign permissions to `user` folders
# Install golang and npm necessaries dependencies to VS Code extensions
# Add the vscode debian repo
# Install VS Code extensions passed on build arg VS_EXTENSIONS
RUN apt-get update \
    && apt-get install -y \
	    apt-transport-https \
	    ca-certificates \
	    curl \
	    gnupg \
	    apt-utils \
	    libasound2 \
	    libatk1.0-0 \
      libcairo2 \
      libcups2 \
      libexpat1 \
      libfontconfig1 \
      libfreetype6 \
      libgtk2.0-0 \
      libpango-1.0-0 \
      libx11-xcb1 \
      libxcomposite1 \
      libxcursor1 \
      libxdamage1 \
      libxext6 \
      libxfixes3 \
      libxi6 \
      libxrandr2 \
      libxrender1 \
      libxss1 \
      libxtst6 \
      openssh-client \
      xdg-utils \
      dconf-editor \
      dbus-x11 \
      libfile-mimeinfo-perl \
      xdg-user-dirs \
      xsel \
	    --no-install-recommends \
	&& mkdir -p $HOME/.cache/dconf \
  && mkdir -p $HOME/.config/dconf \
  && chown skydev:skydev -R $HOME/.config \
  && chown skydev:skydev -R $HOME/.cache \
  && go get -v github.com/ramya-rao-a/go-outline \
  && go get -v github.com/uudashr/gopkgs/cmd/gopkgs \
  && go get -v github.com/acroca/go-symbols \
  && go get -v github.com/stamblerre/gocode \
  && go get -v github.com/ianthehat/godef \
  && go get -v github.com/sqs/goreturns \
  && ln -s /go/bin/gocode /go/bin/gocode-gomod \
  && ln -s /go/bin/godef /go/bin/godef-gomod \
  && npm install -g tslint typescript \
  && curl -sSL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor | apt-key add - \
  && echo "deb [arch=amd64] https://packages.microsoft.com/repos/vscode stable main" > /etc/apt/sources.list.d/vscode.list \
  && apt-get update \
  && apt-get -y install code \
  && for ext in $VS_EXTENSIONS; do code --user-data-dir $HOME --install-extension $ext; done \
	&& apt clean \
	&& rm -rf /var/lib/apt/lists/*

# Change to `skydev` and generate default user folders to avoid config problems in future
USER skydev
RUN xdg-user-dirs-update --force

# Back to root user
USER root

# Copy start.sh script to use it as our Docker ENTRYPOINT
COPY ./start.sh /usr/local/bin/start.sh
# backwards compat
RUN ln -s usr/local/bin/start.sh /

WORKDIR $GOPATH/src/github.com/skycoin/

ENTRYPOINT ["start.sh"]

#CMD [ "su", "skydev", "-p", "-c", "/usr/share/code/code" ]
