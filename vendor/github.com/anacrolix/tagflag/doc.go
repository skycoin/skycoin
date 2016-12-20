// Package tagflag uses reflection to derive flags and positional arguments to a
// program, and parses and sets them from a slice of arguments.
//
// For example:
//  var opts struct {
//      Mmap           bool           `help:"memory-map torrent data"`
//      TestPeer       []*net.TCPAddr `help:"addresses of some starting peers"`
//      tagflag.StartPos              // Marks beginning of positional arguments.
//      Torrent        []string       `arity:"+" help:"torrent file path or magnet uri"`
//  }
//  tagflag.Parse(&opts)
//
// Supported tags include:
//  help: a line of text to show after the option
//  arity: defaults to 1. the number of arguments a field requires, or ? for one
//         optional argument, + for one or more, or * for zero or more.
//
// MarshalArgs is called on fields that implement ArgsMarshaler. A number of
// arguments matching the arity of the field are passed if possible.
//
// Slices will collect successive values, within the provided arity constraints.
//
// A few helpful types have builtin marshallers, for example Bytes,
// *net.TCPAddr, *url.URL, time.Duration, and net.IP.
//
// Flags are strictly passed with the form -K or -K=V. No space between -K and
// the value is allowed. This allows positional arguments to be mixed in with
// flags, and prevents any confusion due to some flags occasionally not taking
// values. A `--` will terminate flag parsing, and treat all further arguments
// as positional.
//
// A builtin help and usage printer are provided, and activated when passing
// -h or -help.
//
// Flag and positional argument names are automatically munged to fit the
// standard scheme within tagflag.
package tagflag
