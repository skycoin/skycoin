// Package flags pkg/flags/helpall.go
//
// InstallHelp replaces cobra's default (auto-generated) help subcommand
// with a flag-aware version that consolidates three previously-separate
// concerns — printing every subcommand's help in one pass, rendering
// the subcommand tree, and regenerating markdown docs — into the same
// `help` command users already reach for. Three mutually-exclusive
// mode flags (-r / -t / -d) pick between the variants.
//
// Why not a separate `help-all`? `help` already exists at every level
// of a cobra tree, invocation is universally understood, and
// overloading it avoids cluttering the top-level command list with a
// near-duplicate.
package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// InstallHelp installs a flag-aware help command on root. Replaces
// cobra's auto-generated help command for that scope. Safe to call
// on multiple roots — each gets its own closure-captured flag state.
//
// Modes (mutually exclusive):
//
//	(default)     print root's help (or [command]'s if an arg is given)
//	-r/--recursive print root + every descendant's help
//	-t/--tree     print the subcommand tree only (no help text)
//	-d/--doc      print markdown documentation
//
// With an argument (`cmd help <subcmd>`), the flags apply to the named
// subcommand's subtree rather than the root's, matching the intuition
// of `help <subcmd>`.
func InstallHelp(root *cobra.Command) {
	var recursive, tree, doc bool
	helpCmd := &cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command",
		Long: `Print help for the given command, or for the root command when no argument is supplied.

Modes (mutually exclusive):
  (default)      print the command's own help
  -r, --recursive print the command's help plus every descendant
  -t, --tree     print the subcommand tree only (no help text)
  -d, --doc      print markdown documentation`,
		Hidden:                true,
		SilenceErrors:         true,
		SilenceUsage:          true,
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Run: func(_ *cobra.Command, args []string) {
			target := root
			if len(args) > 0 {
				if found, _, err := root.Find(args); err == nil && found != nil {
					target = found
				}
			}
			switch {
			case tree:
				PrintCommandTree(target)
			case doc:
				PrintCommandDocs(target)
			case recursive:
				PrintHelpAll(target, true)
			default:
				_ = target.Help() //nolint:errcheck,gosec
			}
		},
	}
	helpCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "print help for every descendant as well")
	helpCmd.Flags().BoolVarP(&tree, "tree", "t", false, "print the subcommand tree only")
	helpCmd.Flags().BoolVarP(&doc, "doc", "d", false, "print markdown documentation")
	helpCmd.MarkFlagsMutuallyExclusive("recursive", "tree", "doc")
	// SetHelpCommand alone stores helpCmd in root.helpCommand but
	// cobra's InitDefaultHelpCmd (which also adds the helpCommand to
	// root's commands slice so it resolves as `root help …`) runs
	// only on the top-level root at Execute time. Subcommand roots
	// (e.g. `skywire cli`) need the helpCommand wired up manually —
	// without the AddCommand, `cli help` would still print help, but
	// via a fallback path that never parses our -r/-t/-d flags.
	// Remove a previously-added "help" child first so repeated
	// InstallHelp calls on the same root are idempotent.
	if prev := findChild(root, "help"); prev != nil {
		root.RemoveCommand(prev)
	}
	root.SetHelpCommand(helpCmd)
	root.AddCommand(helpCmd)
}

// findChild returns the direct child of c matching name, or nil.
func findChild(c *cobra.Command, name string) *cobra.Command {
	for _, child := range c.Commands() {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

// PrintHelpAll writes root.Help() followed by each subcommand's Help().
// With recursive=true, every descendant is included.
func PrintHelpAll(root *cobra.Command, recursive bool) {
	printCmdHelp(root)
	for _, c := range root.Commands() {
		if !c.IsAvailableCommand() {
			continue
		}
		printCmdHelp(c)
		if recursive {
			printDescendantsHelp(c)
		}
	}
}

func printDescendantsHelp(c *cobra.Command) {
	for _, child := range c.Commands() {
		if !child.IsAvailableCommand() {
			continue
		}
		printCmdHelp(child)
		printDescendantsHelp(child)
	}
}

func printCmdHelp(c *cobra.Command) {
	fmt.Println(strings.Repeat("=", 72))
	fmt.Printf("%s\n", c.CommandPath())
	fmt.Println(strings.Repeat("=", 72))
	_ = c.Help() //nolint:errcheck,gosec
	fmt.Println()
}

// PrintCommandTree writes an ASCII tree of the command hierarchy rooted
// at root. Kept dependency-free (no pterm / tablewriter) so the
// utilities package stays lightweight.
func PrintCommandTree(root *cobra.Command) {
	fmt.Println(root.Name())
	printTreeChildren(root, "")
}

func printTreeChildren(c *cobra.Command, prefix string) {
	// Materialize visible children once so the last-child branch char
	// (`└──`) is picked correctly when some children are hidden.
	var visible []*cobra.Command
	for _, child := range c.Commands() {
		if !child.IsAvailableCommand() {
			continue
		}
		visible = append(visible, child)
	}
	for i, child := range visible {
		isLast := i == len(visible)-1
		branch := "├── "
		nextPrefix := prefix + "│   "
		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		}
		fmt.Printf("%s%s%s\n", prefix, branch, child.Name())
		printTreeChildren(child, nextPrefix)
	}
}

// PrintCommandDocs emits a Markdown document describing every command
// in the tree rooted at root. Headings deepen with nesting (root = H1,
// direct children = H2, etc.), capped at H6 to stay within CommonMark.
// The help text is wrapped in fenced code blocks.
func PrintCommandDocs(root *cobra.Command) {
	fmt.Printf("# %s\n\n", root.CommandPath())
	printCmdDoc(root, 1)
	walkCmdDocs(root, 2)
}

func walkCmdDocs(c *cobra.Command, depth int) {
	for _, child := range c.Commands() {
		if !child.IsAvailableCommand() {
			continue
		}
		printCmdDoc(child, depth)
		walkCmdDocs(child, depth+1)
	}
}

func printCmdDoc(c *cobra.Command, depth int) {
	if depth > 6 {
		depth = 6
	}
	fmt.Printf("%s %s\n\n", strings.Repeat("#", depth), c.CommandPath())
	fmt.Println("```")
	_ = c.Help() //nolint:errcheck,gosec
	fmt.Println("```")
	fmt.Println()
}
