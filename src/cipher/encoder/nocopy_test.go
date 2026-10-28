package encoder

import (
	"bytes"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type noCopyCase struct {
	ID     uint32
	Name   string `enc:",maxlen=16"`
	Blob   []byte
	Nested []Contained
	Tags   map[string][]byte
	Tail   []byte `enc:",omitempty"`
}

// checkNoCopyMatches decodes in with every variant into new values of the
// type of proto and requires identical results, and in left unchanged.
func checkNoCopyMatches(t *testing.T, in []byte, proto interface{}) {
	t.Helper()
	typ := reflect.TypeOf(proto)
	orig := bytes.Clone(in)

	want := reflect.New(typ)
	wn, werr := DeserializeRaw(bytes.Clone(in), want.Interface())

	got := reflect.New(typ)
	gn, gerr := DeserializeRawNoCopy(in, got.Interface())
	require.Equal(t, werr, gerr)
	require.Equal(t, wn, gn)
	if werr == nil {
		require.Equal(t, want.Interface(), got.Interface())
	}

	wantV := reflect.New(typ)
	wvn, wverr := DeserializeRawToValue(bytes.Clone(in), wantV)
	gotV := reflect.New(typ)
	gvn, gverr := DeserializeRawToValueNoCopy(in, gotV)
	require.Equal(t, wverr, gverr)
	require.Equal(t, wvn, gvn)
	if wverr == nil {
		require.Equal(t, wantV.Interface(), gotV.Interface())
	}

	require.Equal(t, orig, in, "the no-copy decoders must not write to in")
}

func noCopyInputs(t *testing.T) map[string]struct {
	in    []byte
	proto interface{}
} {
	golden, err := os.ReadFile("testdata/encode-every-type.golden")
	require.NoError(t, err)
	c := noCopyCase{
		ID:     7,
		Name:   "short",
		Blob:   bytes.Repeat([]byte{0xab}, 300),
		Nested: []Contained{{X: 1, Y: 2, Bytes: []uint8{3, 4}, Ints: []uint16{5}}, {}},
		Tags:   map[string][]byte{"k": []byte("v")},
		Tail:   []byte("tail"),
	}
	noTail := c
	noTail.Tail = nil
	long := c
	long.Name = "this name is longer than sixteen"
	return map[string]struct {
		in    []byte
		proto interface{}
	}{
		"every type": {golden, hasEveryType{}},
		"case":       {Serialize(c), noCopyCase{}},
		"no tail":    {Serialize(noTail), noCopyCase{}},
		"maxlen":     {Serialize(long), noCopyCase{}},
		"empty":      {Serialize(noCopyCase{}), noCopyCase{}},
	}
}

// The no-copy decoders give the same values, byte counts and errors as
// the copying ones, on whole, truncated, extended and corrupted inputs.
func TestDeserializeRawNoCopyMatches(t *testing.T) {
	r := rand.New(rand.NewSource(1)) //nolint:gosec
	for name, tc := range noCopyInputs(t) {
		t.Run(name, func(t *testing.T) {
			for i := 0; i <= len(tc.in); i++ {
				checkNoCopyMatches(t, bytes.Clone(tc.in[:i]), tc.proto)
			}
			checkNoCopyMatches(t, append(bytes.Clone(tc.in), 1, 2, 3), tc.proto)
			for i := 0; i < 500; i++ {
				in := bytes.Clone(tc.in)
				for j := r.Intn(4); j >= 0 && len(in) > 0; j-- {
					in[r.Intn(len(in))] = byte(r.Intn(256))
				}
				checkNoCopyMatches(t, in, tc.proto)
			}
		})
	}
	_, err := DeserializeRawNoCopy(nil, noCopyCase{})
	require.Error(t, err)
	_, err = DeserializeRawToValueNoCopy(nil, reflect.ValueOf(noCopyCase{}))
	require.Error(t, err)
}

// Decoded byte slices point into the input instead of a copy.
func TestDeserializeRawNoCopyShares(t *testing.T) {
	in := Serialize(noCopyCase{Blob: []byte("shared")})
	var out noCopyCase
	_, err := DeserializeRawNoCopy(in, &out)
	require.NoError(t, err)
	in[bytes.Index(in, []byte("shared"))+5] = 'D'
	require.Equal(t, "shareD", string(out.Blob), "Blob aliases in")
}

func BenchmarkDeserializeRawBlob(b *testing.B) {
	in := Serialize(noCopyCase{Name: "leaf", Blob: make([]byte, 64<<10)})
	b.Run("copy", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out noCopyCase
			DeserializeRaw(in, &out) //nolint:errcheck,gosec
		}
	})
	b.Run("nocopy", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var out noCopyCase
			DeserializeRawNoCopy(in, &out) //nolint:errcheck,gosec
		}
	})
}
