#!/bin/bash

# This are the git ENV vars present at this point for the build process.
# more on that in https://docs.docker.com/docker-cloud/builds/advanced
#
# SOURCE_BRANCH: the name of the branch or the tag that is currently being tested.
# SOURCE_COMMIT: the SHA1 hash of the commit being tested.
# COMMIT_MSG: the message from the commit being tested and built.
# DOCKERFILE_PATH: the dockerfile currently being built.
# DOCKER_REPO: the name of the Docker repository being built.
# CACHE_TAG: the Docker repository tag being built.
# IMAGE_NAME: the name and tag of the Docker repository being built.
#             (This variable is a combination of DOCKER_REPO:CACHE_TAG.)

DOCKERFILE_PATH='Dockerfile'
IMAGE_NAME='skycoin/skycoindev-vscode'

echo "Build hook running"

cd ../

# Build skycoin/skycoindev-vscode:develop
docker build --build-arg BDATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
             --build-arg SHEAD=`git rev-parse HEAD` \
             --build-arg SCOMMIT=$SOURCE_COMMIT \
             -f $DOCKERFILE_PATH \
             -t $IMAGE_NAME:develop .

# Build skycoin/skycoindev-vscode:dind
docker build --build-arg IMAGE_FROM="skycoin/skycoindev-cli:dind" \
             --build-arg BDATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
             --build-arg SHEAD=`git rev-parse HEAD` \
             --build-arg SCOMMIT=$SOURCE_COMMIT \
             -f $DOCKERFILE_PATH \
             -t $IMAGE_NAME:dind .
