package main

import (
	"fmt"
	"os"

	skycli "github.com/skycoin/skycoin/src/api/cli"
	"github.com/urfave/cli"
)

func main() {
	var commandHelpTemplate = `USAGE:
		{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Category}}
		
CATEGORY:
		{{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
		{{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
		{{range .VisibleFlags}}{{.}}
		{{end}}{{end}}
	`
	cli.SubcommandHelpTemplate = commandHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help,h",
		Usage: "show help, can also be used to show subcommand help",
	}

	//   cli.NewApp().Run(os.Args)
	app := cli.NewApp()
	app.Usage = "the skycoin command line interface"
	app.Version = "0.1"
	app.Commands = skycli.Commands
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
