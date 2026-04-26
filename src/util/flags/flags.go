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
	// Install our flag-aware `help` command in place of cobra's
	// default. Supports -r/-t/-d modes — see InstallHelp. Replaces
	// the previous `&cobra.Command{Hidden: true}` stub that disabled
	// `cmd help` entirely.
	InstallHelp(cmd)
	cmd.PersistentFlags().MarkHidden("help") //nolint:errcheck,gosec

	// Set help template BEFORE cc.Init() so coloredcobra can colorize it
	if !usage {
		cmd.SetHelpTemplate(helpTemplateNoUsage)
	}

	initColoredCobra(cmd)
}

// InitStyle sets the usage + help templates and runs coloredcobra on
// cmd WITHOUT adding flags or a help command. Use this for subcommand
// roots (like scli.RootCmd when embedded under the skywire binary)
// that need colors but already inherit --help from their parent.
func InitStyle(cmd *cobra.Command) {
	cmd.SetUsageTemplate(helpUsage)
	initColoredCobra(cmd)
}

func initColoredCobra(cmd *cobra.Command) {
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

// Both templates below support cobra Command Groups (cobra ≥ 1.6):
// when a command has defined groups via AddGroup, subcommands render
// under their group's title; ungrouped subcommands fall through to an
// "Additional Commands:" block. When no groups are defined the old
// flat "Available Commands:" listing is produced, preserving the
// existing behavior for commands that haven't opted in.
const helpTemplateNoUsage = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if .HasAvailableSubCommands}}{{if eq (len .Groups) 0}}{{HeadingStyle "Available Commands:"}}{{range .Commands}}{{if and (ne .Name "completion") .IsAvailableCommand}}
  {{rpad (CommandStyle .Name) (sum .NamePadding 12) }} {{CmdShortStyle .Short}}{{end}}{{end}}
{{else}}{{$cmds := .Commands}}{{range $group := .Groups}}
{{HeadingStyle $group.Title}}{{range $cmds}}{{if and (eq .GroupID $group.ID) (ne .Name "completion") .IsAvailableCommand}}
  {{rpad (CommandStyle .Name) (sum .NamePadding 12) }} {{CmdShortStyle .Short}}{{end}}{{end}}
{{end}}{{end}}

{{end}}{{if .HasAvailableLocalFlags}}{{HeadingStyle "Flags:"}}
{{FlagStyle .LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

{{HeadingStyle "Global Flags:"}}
{{FlagStyle .InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`

// help is used as usage template (coloredcobra will add colors to it).
// IMPORTANT: the non-grouped path MUST use `range .Commands` (not a
// $variable) so coloredcobra's regex can find and wrap `.Name` /
// `.Short` in style functions. The grouped path uses explicit style
// function calls (HeadingStyle, CommandStyle, CmdShortStyle) because
// cc's regex can't match the `$cmds` variable or dynamic group
// titles. The functions are registered globally by cc.Init, so both
// paths render identically.
const help = `{{if gt (len .Aliases) 0}}{{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}{{if eq (len .Groups) 0}}Available Commands:{{range .Commands}}{{if and (ne .Name "completion") .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{$cmds := .Commands}}{{range $group := .Groups}}
{{HeadingStyle $group.Title}}{{range $cmds}}{{if and (eq .GroupID $group.ID) (ne .Name "completion") .IsAvailableCommand}}
  {{rpad (CommandStyle .Name) (sum .NamePadding 12)}} {{CmdShortStyle .Short}}{{end}}{{end}}
{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`

const helpUsage = `Usage:
  {{.UseLine}}

` + help
