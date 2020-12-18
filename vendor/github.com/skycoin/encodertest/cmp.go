package encodertest

import (
	"unicode"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
)

// IgnoreAllUnexported returns an Option that only ignores the all unexported
// fields of a struct and its descendents, including anonymous fields of unexported types,
func IgnoreAllUnexported() cmp.Option {
	return cmp.FilterPath(allUnexportedFilter, cmp.Ignore())
}

func allUnexportedFilter(p cmp.Path) bool {
	sf, ok := p.Index(-1).(cmp.StructField)
	if !ok {
		return false
	}
	return !isExported(sf.Name())
}

// isExported reports whether the identifier is exported.
func isExported(id string) bool {
	r, _ := utf8.DecodeRuneInString(id)
	return unicode.IsUpper(r)
}
