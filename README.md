skycoin
=======


Setup
-----

* Clone this repo
* Install [gvm](https://github.com/moovweb/gvm) (or hack up $GOPATH yourself)
* Install `go1.2` with gvm, or manually.  Skycoin does not work with earlier releases of go.
* `./compile/getdeps.sh` - This installs go dependencies

Running
-------

### Command line 
Skycoin has three command line interfaces, `dev`, `daemon`, and `client`.

To run skycoin with any one of these, do

```
go run -tags $tag main.go
```

For the developer's convenience,

```
./run.sh
```

will run it in `dev` mode.


### GUI

To run the gui client, it must be built first.

```
./gui.sh build
```

Once it is built, you can run it with

```
./gui.sh
```

until the go source has changed and you need to rebuild.  
You do not need to rebuild if only modifying the GUI frontend code, 
located in `./static/`.

The GUI consists of a `node-webkit` binary and an `skycoin.nw` file which contains the frontend code and the skycoin binary.
When running the GUI, the `node-webkit` binary is executed, it unpacks the `skycoin.nw` file, forks skycoin which runs an http
server on `localhost:$randomport`, and the GUI's `index.html` is served from there.

If you are trying to run the skycoin GUI client on a platform that we are not targeting, you can run `node-webkit` with 
the `skycoin.nw` file produced by the build scripts in `compile/`.

Available Platforms
-------------------

The instructions for running the client apply to Linux, Windows and OSX.
Windows will need MingW.

Skycoin development is primarily done on Linux so Windows and OSX may break from time to time.

Please report any issues you have running skycoin on your system.

We will provide snapshot binary releases for Linux 32/64-bit, Windows 32-bit and OSX 32-bit once
the client is deemed ready for distribution.


Tests
-----

Skycoin tests can be run with 

```
./test.sh
```

At the moment, there are few tests for the core skycoin source.  Skycoin only recently stabilized enough
to be tested. [gnet](https://github.com/skycoin/gnet) is being tested next, and [pex](https://github.com/skycoin/pex)'s 
tests are complete.
