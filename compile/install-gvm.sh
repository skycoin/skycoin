#!/usr/bin/env bash

# Install gvm and go1.2
bash < <(curl -s https://raw.github.com/moovweb/gvm/master/binscripts/gvm-installer)
gvm install go1.2
echo "gvm use go1.2" >> ~/.bashrc
source ~/.bashrc
