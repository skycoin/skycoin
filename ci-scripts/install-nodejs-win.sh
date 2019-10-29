#!/usr/bin/env bash

set -e -o pipefail

# Install nodejs with choco
choco install nodejs --version=8.11.0 -y
echo 'export PATH="/c/Program Files/nodejs:${PATH}";' >> ~/.bashrc
