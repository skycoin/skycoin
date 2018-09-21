package flag

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/util/collections"
)

// stringSetReader reads a string set from cli args
type stringSetReader struct {
	FlagName      string
	Values        *collections.StringSet
	Separator     string
	IgnoreInvalid bool
	ValidValues   collections.StringSet
}

// Set new values in this string set.
// Input string will be split using Separator string to add individual items to Values.
// If ValidValues is not empty only values in this set will be added, else any string will be accepted.
// If IgnoreInvalid is false then an error will be returned if unexpected values are provided.
// By setting this flag to true unexpected values will be silently ignored.
func (r *stringSetReader) Set(value string) error {
	values := strings.Split(value, r.Separator)
	if len(r.ValidValues) == 0 || (!r.ValidValues.Contains(values...) && !r.IgnoreInvalid) {
		return errors.New(fmt.Sprintf("Expected one of %s in -%s but got %s", r.ValidValues.String(), r.FlagName, value))
	}
	for _, value := range values {
		(*r.Values)[value] = struct{}{}
	}
	return nil
}

// String representation of string set
func (r *stringSetReader) String() string {
	return r.Values.String()
}

// StringSetVar command line argument
// Input string will be split using Separator string to add individual items to Values.
// If ValidValues is not empty only values in this set will be added, else any string will be accepted.
// If IgnoreInvalid is false then an error will be returned if unexpected values are provided.
// By setting this flag to true unexpected values will be silently ignored.
func StringSetVar(set *collections.StringSet, validValues []string, ignoreInvalid bool, name string, usage string) {
	cliVar := stringSetReader{
		Values:        set,
		Separator:     ",",
		IgnoreInvalid: ignoreInvalid,
		ValidValues:   collections.NewStringSet(validValues...),
		FlagName:      name,
	}
	flag.Var(&cliVar, name, usage)
}
