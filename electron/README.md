Build system
------------

The GUI client is an Electron (http://electron.atom.io/) app.

It cross compiles for osx, linux and windows 64 bit systems.

32bit windows is blocked on a go cross-compiler bug.

Requirements
------------

gox (go cross compiler), node and npm.

To install gox:

```
go get github.com/mitchellh/gox
```

Node and npm installation is system dependent.

Updating NPM
-----

```
sudo apt-get install npm
sudo apt-get install nodejs-legacy
sudo npm cache clean -f

sudo npm install -g n
sudo n stable

node -v
npm -v
```

Use gulp-electron to build app
-----

Install gulp-cli tool:

```
sudo npm install --global gulp-cli
```

Build
-----

```
./build-without-builder.sh
```

Use electron-builder to pack and create app installer
-----

Install electron-builder
-----

```
sudo npm install --global electron-builder
```

Install electron-download
-----

```
sudo npm install --global electron-download
```

For macOS
-----
Use brew to install required packages.

To build app for Windows on macOS:

```
brew install wine --without-x11
brew install mono
```

To build app for Linux on macOS:

```
brew install gnu-tar graphicsmagick xz
```

Code signing
-----

Set the CSC_IDENTITY_AUTO_DISCOVERY environment variable to false if you don't want to do code signing,
otherwise, you can create a certificate in login.keychain for testing purpose.

Create new certificate:
```
Keychain Access -> Certificate Assistant -> Create a Certificate...
```

Set certificate name and select `Code Signing` as `Certificate Type`.

Once you generated the certificate, you can use it by setting your environment variable:

```
export CSC_NAME="Certificate Name"
```

Now, when you run electron-builder, it will choose the name and sign the app with the certificate.


For Linux
-----
To build app in distributable format for Linux:

```
sudo apt-get install --no-install-recommends -y icnsutils graphicsmagick xz-utils
```

To build app for Windows on Linux:

* Install Wine (1.8+ is required):

```
sudo apt-get install software-properties-common
sudo add-apt-repository ppa:ubuntu-wine/ppa -y
sudo apt-get update
sudo apt-get install --no-install-recommends -y wine1.8
```

* Install Mono (4.2+ is required):

```
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 3FA7E0328081BFF6A14DA29AA6A19B38D3D831EF
echo "deb http://download.mono-project.com/repo/debian wheezy main" | sudo tee /etc/apt/sources.list.d/mono-xamarin.list
sudo apt-get update
sudo apt-get install --no-install-recommends -y mono-devel ca-certificates-mono
```

To build app in 32 bit from a machine with 64 bit:
-----

```
sudo apt-get install --no-install-recommends -y gcc-multilib g++-multilib
```

Manually download electron files
-----

To speed up download speed in China, use Electron Mirror of China:

```
export ELECTRON_MIRROR="https://npm.taobao.org/mirrors/electron/"
```

Download:

```
./electron-downloader.sh
```

Setup
-----

Once requirements are installed, node dependencies must be downloaded.

```
npm install
```

A folder `node_modules/` should now exist.

Building
--------

```
./build.sh
```

* compiles the skycoin app with gox (in parallel for all targets),
* creates the base electron app
* copies the skycoin binaries and static assets into the electron app
* compresses the electron app

Final results are placed in the `release/` folder.
