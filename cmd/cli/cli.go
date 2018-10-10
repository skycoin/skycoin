/*
cli is a command line client for interacting with a skycoin node and offline wallet management
*/
package main

import (
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/cli"
)

func setupApp() *cli.App {
	cfg, err := cli.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app, err := cli.NewApp(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return app
}

func run(app *cli.App) {
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	run(setupApp())
}
