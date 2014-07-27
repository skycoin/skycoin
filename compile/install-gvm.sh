#!/usr/bin/env bash

# Install gvm
bash < <(curl -s https://raw.github.com/moovweb/gvm/master/binscripts/gvm-installer)
grep -q "scripts/gvm" ~/.bashrc;
if [[ $? -ne 0 ]]; then
    echo "" >> ~/.bashrc
    echo "# Setup gvm on shell start" >> ~/.bashrc
    echo '[[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"' >> ~/.bashrc
    [[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"
fi

# Install go1.3 with gvm
gvm install go1.3
grep -q 'gvm use' ~/.bashrc;
if [[ $? -ne 0 ]]; then
    echo "" >> ~/.bashrc
    echo "# Use go1.3 on shell start" >> ~/.bashrc
    echo "gvm use go1.3 >/dev/null" >> ~/.bashrc
    gvm use go1.3
fi
