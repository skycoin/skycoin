#!/bin/bash
IMAGE_NAME='skycoin/skycoindev-vscode'
docker push $IMAGE_NAME:develop
docker push $IMAGE_NAME:dind

