ARG IMAGE_FROM=skycoin/skycoindev-cli:develop

FROM $IMAGE_FROM

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
	    --no-install-recommends \
	&& apt clean \
	&& rm -rf /var/lib/apt/lists/*

# Add the vscode debian repo
RUN curl -sSL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor | apt-key add -
RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/vscode stable main" > /etc/apt/sources.list.d/vscode.list

# Install code
RUN apt-get update && apt-get -y install code \
	&& apt clean \
	&& rm -rf /var/lib/apt/lists/*

ENV HOME /home/user
RUN useradd --create-home --home-dir $HOME user \
	&& chown -R user:user $HOME

# Install and list vscode extensions
ENV VS_EXTENSIONS ms-vscode.Go windmilleng.vscode-go-autotest defaltd.go-coverage-viewer
RUN for ext in $VS_EXTENSIONS; do code --user-data-dir $HOME --install-extension $ext; done
RUN code --user-data-dir $HOME --list-extensions
ENV VS_EXTENSIONS ""

COPY start.sh /usr/local/bin/start.sh
RUN ln -s usr/local/bin/start.sh / # backwards compat

WORKDIR $HOME

ENTRYPOINT ["start.sh"]

CMD [ "su", "user", "-p", "-c", "/usr/share/code/code" ]
