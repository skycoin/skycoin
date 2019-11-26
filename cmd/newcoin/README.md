# Fiber Coin Creator CLI Documention
This tool can be used to create a new fiber coin easily from a config file.
- [Install](#install)
 - [Usage](#usage)
   - [Create New Coin](#create-new-coin)
     - [Example](#example)

## Install

```bash
$ cd $GOPATH/src/github.com/SkycoinProject/skycoin/cmd/newcoin
$ go install ./...
```

## Usage

After the installation, you can run `newcoin` to see the usage:

```
$ newcoin

NAME:
   newcoin - newcoin is a helper tool for creating new fiber coins

USAGE:
   newcoin [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
     createcoin  Create a new coin from a template file
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Create New Coin
When using the `newcoin` command, you should run it from the `$GOPATH/src/github.com/SkycoinProject/skycoin` folder to utilise the built in default templates.

```bash
$ cd $GOPATH/src/github.com/SkycoinProject/skycoin
$ newcoin createcoin [command options]
```

```
OPTIONS:
   --coin value                             name of the coin to create (default: "skycoin")
   --template-dir value, --td value         template directory path (default: "./template")
   --coin-template-file value, --ct value   coin template file (default: "coin.template")
   --visor-template-file value, --vt value  visor template file (default: "visor_parameters.template")
   --config-dir value, --cd value           config directory path (default: "./")
   --config-file value, --cf value          config file path (default: "fiber.toml")
```

#### Example
Create a test coin using application defaults.

```bash
$ cd $GOPATH/src/github.com/SkycoinProject/skycoin
$ newcoin createcoin --coin testcoin
```

This will create a new directory, `testcoin`, in `cmd` folder and a `testcoin.go` file inside that folder.
It will also use the built-in defaul options (specified above) and draw template configuration from `$GOPATH/src/github.com/SkycoinProject/skycoin/template`

This file can be used to run a "testcoin" node.
