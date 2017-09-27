# Installing go

## For OSX

Install `homebrew`, if you don't have it yet.

```sh
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

Install the latest version of go

```sh
brew install go
```

Install Mercurial and Bazaar

```sh
brew install mercurial bzr
```

## For linux

```sh
sudo apt-get update && apt-get install curl git mercurial make binutils gcc bzr bison libgmp3-dev screen gcc build-essential -y
```

### Install gvm

```sh
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)

source $HOME/.gvm/scripts/gvm
```

### Install go with gvm

In China, use `--source=https://github.com/golang/go` to bypass firewall when fetching golang source.

```sh
gvm install go1.4 --source=https://github.com/golang/go
gvm use go1.4
gvm install go1.9
gvm use go1.9 --default
```

If you open up new terminal and the go command is not found then add this to `.bashrc`. GVM should add this automatically.

```sh
[[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"
gvm use go1.9 >/dev/null
```

## Setup your GOPATH

The GOPATH environment variable specifies the location of your workspace. It defaults to a directory named go inside your home directory, so $HOME/go on Unix.

Create your workspace directory:

```sh
mkdir -p $HOME/go
```

Setup $GOPATH variable, add it to ~/.bashrc. After editing, run `source ~/.bashrc` or open a new tab.

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

## Test your Go installation

Create and run the hello.go application described here: [https://golang.org/doc/install#testing](https://golang.org/doc/install#testing) to check if your Go installation is working.
