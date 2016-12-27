package tagflag

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Struct fields after this one are considered positional arguments.
type StartPos struct{}

// Default help flag was provided, and should be handled.
var ErrDefaultHelp = errors.New("help flag")

// Parses given arguments, returning any error.
func ParseErr(cmd interface{}, args []string, opts ...parseOpt) (err error) {
	p, err := newParser(cmd, opts...)
	if err != nil {
		return
	}
	return p.parse(args)
}

// Parses the command-line arguments, exiting the process appropriately on
// errors or if usage is printed.
func Parse(cmd interface{}, opts ...parseOpt) {
	opts = append([]parseOpt{Program(filepath.Base(os.Args[0]))}, opts...)
	p, err := newParser(cmd, opts...)
	if err == nil {
		err = p.parse(os.Args[1:])
	}
	if err == ErrDefaultHelp {
		p.printUsage(os.Stderr)
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "tagflag: %s\n", err)
		if _, ok := err.(userError); ok {
			os.Exit(2)
		}
		os.Exit(1)
	}
}
