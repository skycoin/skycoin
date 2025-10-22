// Package main provides the skycoin-cli tool.
// package main cmd/newcoin/newcoin.go
/*
newcoin generates a new coin cmd from a toml configuration file
*/

package main

import (
	"log"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/cmd/skycoin-cli/commands"
)

func init() {
	var helpflag bool
	commands.RootCmd.SetUsageTemplate(help)
	commands.RootCmd.PersistentFlags().BoolVarP(&helpflag, "help", "h", false, "help menu")
	commands.RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	commands.RootCmd.PersistentFlags().MarkHidden("help") //nolint
}

func main() {
	cc.Init(&cc.Config{
		RootCmd:         commands.RootCmd,
		Headings:        cc.HiBlue + cc.Bold,
		Commands:        cc.HiBlue + cc.Bold,
		CmdShortDescr:   cc.HiBlue,
		Example:         cc.HiBlue + cc.Italic,
		ExecName:        cc.HiBlue + cc.Bold,
		Flags:           cc.HiBlue + cc.Bold,
		FlagsDescr:      cc.HiBlue,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})
	err := commands.RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

const help = "{{if .HasAvailableSubCommands}}{{end}} {{if gt (len .Aliases) 0}}\r\n\r\n" +
	"{{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}" +
	"Available Commands:{{range .Commands}}  {{if and (ne .Name \"completion\") .IsAvailableCommand}}\r\n  " +
	"{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}\r\n\r\n" +
	"Flags:\r\n" +
	"{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}\r\n\r\n" +
	"Global Flags:\r\n" +
	"{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}\r\n\r\n"
