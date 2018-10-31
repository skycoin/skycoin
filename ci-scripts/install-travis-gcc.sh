
# Install gcc6 (6.4.0-2 on Mac OS) for Travis builds

if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
  sudo apt-get install -qq g++-6 && sudo update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-6 90;
fi

if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
  brew update
  echo 'Available versions (gcc)' && brew list --versions gcc
  brew list gcc@6 &>/dev/null || (echo "Install gcc6 deps" && (brew deps -n gcc@6 | xargs -I NAME bash -c 'brew list NAME &> /dev/null || brew install NAME') && echo "Install gcc@6 from sources" && HOMEBREW_BUILD_FROM_SOURCE=1 brew install gcc@6 )
fi

cd $TRAVIS_BUILD_DIR

