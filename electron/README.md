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
