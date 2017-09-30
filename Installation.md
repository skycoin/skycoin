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
We need to install linux dependencies on the correct distribution.
Currently we have support for the following linux distributions:
 * Ubuntu 16.04
 * Debian 9.1
 * Centos 7
 * Fedora 26

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

### Install go with gvm
#### Install gvm
<strong>gvm</strong> need to be installed. It's used to install go programming language as easy as posible.

```sh
curl -sSL https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer > gvm-installer && chmod a+x gvm-installer && 
source $HOME/.gvm/scripts/gvm
```

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

### Install go manually.
Let's download and uncompress golang.

```sh
export GOV=1.9 # golang version. Could be 1.1, 1.2, 1.3, 1.4, ..., 1.9
curl -sS https://storage.googleapis.com/golang/go$GOV.linux-amd64.tar.gz > go$GOV.linux-amd64.tar.gz
tar xvf go$GOV.linux-amd64.tar.gz
rm go$GOV.linux-amd64.tar.gz
```

After that, let's install go.

```sh
mv go /usr/local/go
ln -s /usr/local/go/bin/go /usr/local/bin/go
ln -s /usr/local/go/bin/godoc /usr/local/bin/godoc
ln -s /usr/local/go/bin/gofmt /usr/local/bin/gofmt
```


## Setup your GOPATH
The <strong>$GOPATH</strong> environment variable specifies the location of your workspace. It defaults to a directory named <strong>go</strong> inside your home directory, so <strong>$HOME/go</strong> on Unix.

Create your workspace directory with it's respective inner folders:

```sh
mkdir -p $HOME/go
mkdir -p $HOME/go/bin
mkdir -p $HOME/go/src
mkdir -p $HOME/go/pkg
```

Setup <strong>$GOPATH</strong> variable, add it to ~/.bashrc. After editing, run `source ~/.bashrc` or open a new tab.

```sh
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
```

## Test your Go installation
Create and run the hello.go application described here: https://golang.org/doc/install#testing to check if your Go installation is working.
