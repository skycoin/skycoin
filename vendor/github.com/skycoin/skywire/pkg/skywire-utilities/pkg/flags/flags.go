// Package flags internal/flags/flags.go
package flags

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

// InitFlags is used to set command flags for help menu and styling with coloredcobra
func InitFlags(cmd *cobra.Command, usage bool) {
	var helpflag bool
	if !usage {
		cmd.SetUsageTemplate(help)
	} else {
		cmd.SetUsageTemplate(helpUsage)
	}
	cmd.PersistentFlags().BoolVarP(&helpflag, "help", "h", false, "show help menu")
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.PersistentFlags().MarkHidden("help") //nolint:errcheck,gosec

	// Set help template BEFORE cc.Init() so coloredcobra can colorize it
	if !usage {
		cmd.SetHelpTemplate(helpTemplateNoUsage)
	}

	cc.Init(&cc.Config{
		RootCmd:         cmd,
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
}

const helpTemplateNoUsage = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if .HasAvailableSubCommands}}{{HeadingStyle "Available Commands:"}}{{range .Commands}}{{if and (ne .Name "completion") .IsAvailableCommand}}
  {{rpad (CommandStyle .Name) (sum .NamePadding 12) }} {{CmdShortStyle .Short}}{{end}}{{end}}

{{end}}{{if .HasAvailableLocalFlags}}{{HeadingStyle "Flags:"}}
{{FlagStyle .LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

{{HeadingStyle "Global Flags:"}}
{{FlagStyle .InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`

// help is used as usage template (coloredcobra will add colors to it)
const help = `{{if gt (len .Aliases) 0}}{{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}Available Commands:{{range .Commands}}{{if and (ne .Name "completion") .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`

const helpUsage = `Usage:
  {{.UseLine}}

` + help
