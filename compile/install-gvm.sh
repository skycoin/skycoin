#!/usr/bin/env bash

# Install gvm and go1.2
bash < <(curl -s https://raw.github.com/moovweb/gvm/master/binscripts/gvm-installer)
gvm install go1.2
grep -q "scripts/gvm" ~/.bashrc;
if [[ $? -ne 0 ]]; then
    echo "" >> ~/.bashrc
    echo "# Setup gvm on shell start" >> ~/.bashrc
    echo '[[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"' >> ~/.bashrc
    [[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"
fi
grep -q 'gvm use' ~/.bashrc;
if [[ $? -ne 0 ]]; then
    echo "" >> ~/.bashrc
    echo "# Use go1.2 on shell start" >> ~/.bashrc
    echo "gvm use go1.2 >/dev/null" >> ~/.bashrc
    gvm use go1.2
fi
