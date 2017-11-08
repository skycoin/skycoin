#!/usr/bin/env bash
BRANCH="release"

# Get build version from package.json
APP_VERSION=`grep version electron/package.json | sed  's/[,\", ]//g'| awk '{split($0,a,":");print a[2]}'`

# Are we on the right branch?
if [ "$TRAVIS_BRANCH" = "$BRANCH" ]; then
  
  # Is this not a Pull Request?
  if [ "$TRAVIS_PULL_REQUEST" = false ]; then
    
    # Is this not a build which was triggered by setting a new tag?
    if [ -z "$TRAVIS_TAG" ]; then
      echo -e "Starting to tag commit.\n"

      git config --global user.email "heaven.mao@gmail.com"
      git config --global user.name "iketheadore"

      # Add tag and push to master.
      git tag -a v${APP_VERSION} -m "Travis build $APP_VERSION pushed a tag."
      git push origin --tags
      git fetch origin

      echo -e "Done magic with tags.\n"
  fi
  fi
fi