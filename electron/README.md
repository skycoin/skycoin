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

Install gulp
-----

```
npm rm --global gulp
npm install --global gulp-cli
```

Setup
-----

Once requirements are installed, node dependencies must be downloaded.

```
npm install
```

A folder `node_modules/` should now exist.

Manually download electron release 1.2.2
-----

 The gulp-electron now can't download the electron release of version 1.2.2
 automatically. We need to download them manually from the following links:

 * [electron-v1.2.2-darwin-x64.zip](https://github.com/electron/electron/releases/download/v1.2.2/electron-v1.2.2-darwin-x64.zip)
 * [electron-v1.2.2-linux-x64.zip](https://github.com/electron/electron/releases/download/v1.2.2/electron-v1.2.2-linux-x64.zip)
 * [electron-v1.2.2-win32-x64.zip](https://github.com/electron/electron/releases/download/v1.2.2/electron-v1.2.2-win32-x64.zip)

Copy the downloaded zip files into `electron/.electron_cache` folder.

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
