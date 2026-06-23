//go:build tinygo

package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// cobra renders its usage/help via text/template with a FuncMap. Validating a
// FuncMap calls reflect.Type.NumOut, which TinyGo's reflect does not implement,
// so the default templated help panics ("unimplemented: (reflect.Type).NumOut()")
// on --help, no-args, and usage-on-error. Installing plain-text Usage/Help funcs
// on the root command (inherited by every subcommand) avoids text/template
// entirely. This only affects TinyGo builds; standard builds keep the normal
// (coloredcobra) templated help.
func init() {
	RootCmd.SetUsageFunc(tinygoUsageFunc)
	RootCmd.SetHelpFunc(tinygoHelpFunc)
}

func tinygoUsageFunc(c *cobra.Command) error {
	w := c.OutOrStderr()
	fmt.Fprintf(w, "Usage:\n  %s\n", c.UseLine())

	if c.HasAvailableSubCommands() {
		fmt.Fprintln(w, "\nAvailable Commands:")
		for _, sub := range c.Commands() {
			if sub.IsAvailableCommand() {
				fmt.Fprintf(w, "  %-20s %s\n", sub.Name(), sub.Short)
			}
		}
	}
	if c.HasAvailableLocalFlags() {
		fmt.Fprintf(w, "\nFlags:\n%s", c.LocalFlags().FlagUsages())
	}
	if c.HasAvailableInheritedFlags() {
		fmt.Fprintf(w, "\nGlobal Flags:\n%s", c.InheritedFlags().FlagUsages())
	}
	if c.HasAvailableSubCommands() {
		fmt.Fprintf(w, "\nUse \"%s [command] --help\" for more information about a command.\n", c.CommandPath())
	}
	return nil
}

func tinygoHelpFunc(c *cobra.Command, _ []string) {
	if long := strings.TrimSpace(c.Long); long != "" {
		fmt.Fprintln(c.OutOrStdout(), long)
	} else if short := strings.TrimSpace(c.Short); short != "" {
		fmt.Fprintln(c.OutOrStdout(), short)
	}
	_ = c.Usage()
}
