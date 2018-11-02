
# Install gcc6 (6.4.0-2 on Mac OS) for Travis builds

if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
  sudo apt-get install -qq g++-6 && sudo update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-6 90;
fi

if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
  brew update
  echo 'Available versions (gcc)'
  brew list --versions gcc
  echo 'Installing gcc@64 formula'
  cd "$(brew --repository)/Library/Taps/homebrew/homebrew-core"
  git show 42d31bba7772fb01f9ba442d9ee98b33a6e7a055:Formula/gcc\@6.rb > Formula/gcc\@64.rb
  sed -i '' -e 's/GccAT6/GccAT64/g' Formula/gcc\@64.rb
  brew install gcc@64 || brew link --overwrite gcc\@64
fi

cd $TRAVIS_BUILD_DIR

