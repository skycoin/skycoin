#!/bin/bash

DOCKERFILE_PATH='Dockerfile'
IMAGE_NAME='skycoin/skycoindev-vscode'
cd ../
docker build -f $DOCKERFILE_PATH -t $IMAGE_NAME:develop .
docker build --build-arg=IMAGE_FROM="skycoin/skycoindev-cli:dind" -f $DOCKERFILE_PATH -t $IMAGE_NAME:dind .
