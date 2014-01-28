#!/usr/bin/env bash

# Install gvm and go1.2
bash < <(curl -s https://raw.github.com/moovweb/gvm/master/binscripts/gvm-installer)
gvm install go1.2
grep -q 'gvm use' ~/.bashrc;
if [[ $? -ne 0 ]]; then
    echo "gvm use go1.2 >/dev/null" >> ~/.bashrc
    gvm use go1.2
fi
