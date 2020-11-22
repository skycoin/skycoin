package useragent

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataBuild(t *testing.T) {
	cases := []struct {
		name      string
		data      Data
		userAgent string
		err       error
	}{
		{
			name:      "without remark",
			userAgent: "Skycoin:0.24.1",
			data: Data{
				Coin:    "Skycoin",
				Version: "0.24.1",
			},
		},
		{
			name:      "with remark",
			userAgent: "Skycoin:0.24.1(remark; foo)",
			data: Data{
				Coin:    "Skycoin",
				Version: "0.24.1",
				Remark:  "remark; foo",
			},
		},
		{
			name: "invalid characters in coin",
			data: Data{
				Coin:    "foo<>",
				Version: "0.24.1",
			},
			err: ErrIllegalChars,
		},
		{
			name: "invalid characters in version",
			data: Data{
				Coin:    "foo",
				Version: "<0.24.1",
			},
			err: errors.New(`Invalid character(s) found in major number "<0"`),
		},
		{
			name: "invalid characters in remark",
			data: Data{
				Coin:    "foo",
				Version: "0.24.1",
				Remark:  "<>",
			},
			err: ErrIllegalChars,
		},
		{
			name: "missing coin",
			data: Data{
				Version: "0.24.1",
			},
			err: errors.New("missing coin name"),
		},
		{
			name: "missing version",
			data: Data{
				Coin: "Skycoin",
			},
			err: errors.New("missing version"),
		},
		{
			name: "version is not valid semver",
			data: Data{
				Coin:    "Skycoin",
				Version: "0.24",
			},
			err: errors.New("No Major.Minor.Patch elements found"),
		},
		{
			name: "invalid remark",
			data: Data{
				Coin:    "skycoin",
				Version: "0.24.1",
				Remark:  "\t",
			},
			err: ErrIllegalChars,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			userAgent, err := tc.data.Build()

			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.userAgent, userAgent)
		})
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		name      string
		userAgent string
		data      Data
		err       error
	}{
		{
			name:      "too long",
			userAgent: fmt.Sprintf("skycoin:0.24.1[abcdefg](%s)", strings.Repeat("a", 245)),
			err:       ErrTooLong,
		},
		{
			name:      "no tab chars allowed",
			userAgent: "skycoin:0.24.1(\t)",
			err:       ErrIllegalChars,
		},
		{
			name:      "no newlines allowed",
			userAgent: "skycoin:0.24.1(\n)",
			err:       ErrIllegalChars,
		},
		{
			name:      "valid",
			userAgent: "skycoin:0.25.0",
			data: Data{
				Coin:    "skycoin",
				Version: "0.25.0",
			},
		},
		{
			name:      "valid, version has suffix",
			userAgent: "skycoin:0.25.1",
			data: Data{
				Coin:    "skycoin",
				Version: "0.25.1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := Parse(tc.userAgent)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.data, d)
		})
	}
}

func TestDataJSON(t *testing.T) {
	d := Data{}

	x, err := json.Marshal(d)
	require.NoError(t, err)
	require.Equal(t, `""`, string(x))

	d.Coin = "skycoin"
	d.Version = "0.25.0"

	x, err = json.Marshal(d)
	require.NoError(t, err)
	require.Equal(t, `"skycoin:0.25.0"`, string(x))

	var e Data
	err = json.Unmarshal([]byte(x), &e)
	require.NoError(t, err)
	require.Equal(t, d, e)

	d.Remark = "foo; bar"

	x, err = json.Marshal(d)
	require.NoError(t, err)
	require.Equal(t, `"skycoin:0.25.0(foo; bar)"`, string(x))

	e = Data{}
	err = json.Unmarshal([]byte(x), &e)
	require.NoError(t, err)
	require.Equal(t, d, e)

	// Fails, does not parse to a string
	err = json.Unmarshal([]byte("{}"), &e)
	require.Error(t, err)

	// OK, empty string
	e = Data{}
	err = json.Unmarshal([]byte(`""`), &e)
	require.NoError(t, err)
	require.Equal(t, Data{}, e)

	// Fails, does not parse
	err = json.Unmarshal([]byte(`"skycoin:0.25.0(<>)"`), &e)
	require.Equal(t, ErrIllegalChars, err)
}

func TestSanitize(t *testing.T) {
	for i := 0; i < len(IllegalChars); i++ {
		x := "t" + IllegalChars[i:i+1]
		t.Run(x, func(t *testing.T) {
			require.Equal(t, "t", Sanitize(x))
		})
	}

	for i := 0; i < 256; i++ {
		j := byte(i)
		if j >= ' ' || j <= '~' {
			continue
		}

		v := []byte{'t', j}

		t.Run(fmt.Sprintf("%q", j), func(t *testing.T) {
			require.Equal(t, "t", Sanitize(string(v)))
		})
	}

	z := "dog\t\t\t\ncat\x01t\xE3\xE4t"
	require.Equal(t, "dogcattt", Sanitize(z))

	// Should not have anything stripped
	x := "Skycoin:0.25.0(foo; bar)"
	require.Equal(t, x, Sanitize(x))

	// Should not have anything stripped
	x = "Skycoin:0.25.1(foo; bar)"
	require.Equal(t, x, Sanitize(x))
}

func TestEmpty(t *testing.T) {
	var d Data
	require.True(t, d.Empty())

	d.Coin = "skycoin"
	d.Version = "0.25.0"
	require.False(t, d.Empty())
}

func TestMustParse(t *testing.T) {
	d := MustParse("skycoin:0.25.0")
	require.Equal(t, Data{
		Coin:    "skycoin",
		Version: "0.25.0",
	}, d)

	require.Panics(t, func() {
		MustParse("foo") //nolint:errcheck
	})
}

func TestMustBuild(t *testing.T) {
	d := Data{
		Version: "0",
	}
	require.Panics(t, func() {
		d.MustBuild() //nolint:errcheck
	})
}
