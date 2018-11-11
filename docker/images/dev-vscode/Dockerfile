ARG IMAGE_FROM=skycoin/skycoindev-cli:develop
FROM $IMAGE_FROM

ARG BDATE
ARG SCOMMIT
ARG VS_EXTENSIONS

# Image labels (see ./hooks/build for ARGS)
LABEL "org.label-schema.name"="skycoindev-cli" \
      "org.label-schema.description"="Docker image with go, node, dev tools and Visual Studio Code for Skycoin developers" \
      "org.label-schema.vendor"="Skycoin project" \
      "org.label-schema.url"="skycoin.net" \
      "org.label-schema.schema-version"="1.0.0-rc.1" \
      "org.label-schema.build-date"=$BDATE \
      "org.label-schema.vcs-url"="https://github.com/skycoin/skycoin.git" \
      "org.label-schema.vcs-ref"=$SCOMMIT \
      "org.label-schema.usage"="https://github.com/skycoin/skycoin/blob/"$SCOMMIT"/docker/images/dev-vscode/README.md" \
      "org.label-schema.docker.cmd"="xhost +; cd src; docker run --rm -it -v /tmp/.X11-unix:/tmp/.X11-unix \
                                    -v $PWD:/go/src \
                                    -w /go/src \
                                    -e DISPLAY=$DISPLAY \
                                    skycoin/skycoindev-vscode:develop"

# Tell debconf to run in non-interactive mode
ENV DEBIAN_FRONTEND noninteractive

# Install dependencies for vs code
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
	    --no-install-recommends \
	&& apt clean \
	&& rm -rf /var/lib/apt/lists/*

# Add the vscode debian repo
RUN curl -sSL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor | apt-key add -
RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/vscode stable main" > /etc/apt/sources.list.d/vscode.list

# Install Visual Studio Code
RUN apt-get update && apt-get -y install code \
	&& apt clean \
	&& rm -rf /var/lib/apt/lists/*

# Create a diferent user to run VS Code
ENV HOME /home/user
RUN useradd --create-home --home-dir $HOME user \
	&& chown -R user:user $HOME

# Generte default user folders to avoid .config problems in future
USER user
RUN xdg-user-dirs-update --force

# Install VS Code extensions passed on build arg VS_EXTENSIONS
RUN for ext in $VS_EXTENSIONS; do code --user-data-dir $HOME --install-extension $ext; done
RUN code --user-data-dir $HOME --list-extensions

# Back to root user
USER root

# Install golang necessaries dependencies to VS Code extensions
RUN go get -v github.com/ramya-rao-a/go-outline \
    && go get -v github.com/uudashr/gopkgs/cmd/gopkgs \
    && go get -v github.com/acroca/go-symbols \
    && go get -v github.com/stamblerre/gocode \
    && go get -v github.com/ianthehat/godef \
    && go get -v github.com/sqs/goreturns \
    && ln -s /go/bin/gocode /go/bin/gocode-gomod \
    && ln -s /go/bin/godef /go/bin/godef-gomod

COPY docker/images/dev-vscode/start.sh /usr/local/bin/start.sh
RUN ln -s usr/local/bin/start.sh / # backwards compat

WORKDIR $HOME

ENTRYPOINT ["start.sh"]

#CMD [ "su", "user", "-p", "-c", "/usr/share/code/code" ]
