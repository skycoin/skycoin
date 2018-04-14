# Installing go

Skycoin supports go1.9+.  The preferred version is go1.10.

## For OSX
First you need to have `homebrew` installed, if you don't have it yet.

```sh
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

Then, let's install go's latest version.

```sh
brew install go
```

Lastly, let's install Mercurial and Bazaar

```sh
brew install mercurial bzr
```

## For linux
We need to install linux dependencies on the correct distribution.

#### Ubuntu and Debian
```sh
sudo apt-get update && sudo apt-get upgrade -y
sudo apt-get install -y curl git mercurial make binutils gcc bzr bison libgmp3-dev screen gcc build-essential
```

#### Centos and Fedora
```sh
sudo yum update -y && sudo yum upgrade -y
sudo yum install -y git curl make gcc mercurial binutils bzr bison screen
if [[ "$(cat /etc/redhat-release | grep -o CentOS)" == "CentOS" ]]; then sudo yum install -y build-essential libgmp3-dev; else sudo yum groupinstall -y "Development Tools" "Development Libraries" && sudo yum install -y gmp; fi;
```

### Install Go with Gvm
#### Install Gvm
`gvm` need to be installed.

```sh
curl -sSL https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer > gvm-installer && chmod a+x gvm-installer &&
source $HOME/.gvm/scripts/gvm
```

#### Install Go
In China, use `--source=https://github.com/golang/go` to bypass firewall when fetching golang source.

```sh
gvm install go1.4 --source=https://github.com/golang/go
gvm use go1.4

gvm install go1.9
gvm use go1.9 --default
```

#### Installation issues
If you open up new a terminal and the `go` command is not found then add this to `.bashrc`. GVM should add this automatically.

```sh
[[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"
gvm use go1.9 >/dev/null
```

## Install Go manually
### Install Go

Let's go to home directory and declare `go`'s version that you want to download.

```sh
cd ~
export GOV=1.10 # golang version. Could be any of the following versions 1.9, 1.10
```

After that, let's download and uncompress golang source.

```sh
curl -sS https://storage.googleapis.com/golang/go$GOV.linux-amd64.tar.gz > go$GOV.linux-amd64.tar.gz
tar xvf go$GOV.linux-amd64.tar.gz
rm go$GOV.linux-amd64.tar.gz
```

lastly, let's install `go`.

```sh
sudo mv go /usr/local/go
sudo ln -s /usr/local/go/bin/go /usr/local/bin/go
sudo ln -s /usr/local/go/bin/godoc /usr/local/bin/godoc
sudo ln -s /usr/local/go/bin/gofmt /usr/local/bin/gofmt
```

Note: Find any golang source version at [Go Website](https://golang.org/dl/)

### Setup your GOPATH
The $GOPATH environment variable specifies the location of your workspace. It defaults to a directory named `go` inside your home directory, so $HOME/go on Unix.

Create your workspace directory with it's respective inner folders:

```sh
mkdir -p $HOME/go
mkdir -p $HOME/go/bin
mkdir -p $HOME/go/src
mkdir -p $HOME/go/pkg
```

Setup $GOPATH variable, add it to ~/.bashrc. After editing, run `source ~/.bashrc` or open a new tab.

```sh
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
```

## Test your Go installation
Create and run the hello.go application described here: https://golang.org/doc/install#testing to check if your Go installation is working.
