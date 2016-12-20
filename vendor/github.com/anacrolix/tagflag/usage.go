package tagflag

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/anacrolix/missinggo"
	"github.com/anacrolix/missinggo/slices"
)

func (p *parser) printUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage:\n  %s", p.program)
	if p.hasOptions() {
		fmt.Fprintf(w, " [OPTIONS...]")
	}
	for _, arg := range p.posArgs {
		fs := func() string {
			switch arg.arity {
			case arity{0, 1}:
				return "[%s]"
			case arity{1, infArity}:
				return "%s..."
			case arity{0, infArity}:
				return "[%s...]"
			default:
				return "<%s>"
			}
		}()
		// if arg.arity != arity{1,1} {
		fmt.Fprintf(w, " "+fs, arg.name)
		// }
		// if arg.arity > 1 {
		//  for range iter.N(int(arg.arity - 1)) {
		//      fmt.Fprintf(w, " "+fs, arg.name)
		//  }
		// }
	}
	fmt.Fprintf(w, "\n")
	if p.description != "" {
		fmt.Fprintf(w, "\n%s\n", missinggo.Unchomp(p.description))
	}
	if awd := p.posWithHelp(); len(awd) != 0 {
		fmt.Fprintf(w, "Arguments:\n")
		tw := newUsageTabwriter(w)
		for _, a := range awd {
			fmt.Fprintf(tw, "  %s\t(%s)\t%s\n", a.name, a.value.Type(), a.help)
		}
		tw.Flush()
	}
	var opts []arg
	for _, v := range p.flags {
		opts = append(opts, v)
	}
	slices.Sort(opts, func(left, right arg) bool {
		return left.name < right.name
	})
	writeOptionUsage(w, opts)
}

func newUsageTabwriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 8, 2, 3, ' ', 0)
}

func writeOptionUsage(w io.Writer, flags []arg) {
	if len(flags) == 0 {
		return
	}
	fmt.Fprintf(w, "Options:\n")
	tw := newUsageTabwriter(w)
	for _, f := range flags {
		fmt.Fprint(tw, "  ")
		fmt.Fprintf(tw, "%s%s", flagPrefix, f.name)
		fmt.Fprintf(tw, "\t(%s)\t%s\n", f.value.Type(), f.help)
	}
	tw.Flush()
}
