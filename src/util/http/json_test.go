package httphelper

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/testutil"
)

func TestFromDuration(t *testing.T) {
	dur := 10 * time.Second
	require.Equal(t, FromDuration(dur).Duration, dur)
}

func TestDurationMarshalJSON(t *testing.T) {
	d := Duration{time.Second * 10}
	data, err := d.MarshalJSON()

	require.NoError(t, err)
	require.Equal(t, `"10s"`, string(data))
}

func TestDurationUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		s        string
		expected time.Duration
		err      string
	}{
		{
			name:     "valid duration",
			s:        "1m",
			expected: time.Minute,
		},
		{
			name: "invalid duration",
			s:    "foo",
			err:  "time: invalid duration foo",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var d Duration
			err := d.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, tc.s)))

			if tc.err != "" {
				require.Equal(t, errors.New(tc.err), err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, d.Duration)
			}

		})
	}

	var d Duration
	err := d.UnmarshalJSON([]byte("invalidjson"))
	testutil.RequireError(t, err, "invalid character 'i' looking for beginning of value")
}

func TestAddressMarshalJSON(t *testing.T) {
	addrStr := "2bfYafFtdkCRNcCyuDvsATV66GvBR9xfvjy"
	addrInner, err := cipher.DecodeBase58Address(addrStr)
	require.NoError(t, err)

	addr := Address{addrInner}

	data, err := addr.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"2bfYafFtdkCRNcCyuDvsATV66GvBR9xfvjy"`, string(data))
}

func TestAddressUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		addr string
		err  string
	}{
		{
			name: "empty address",
			addr: "",
			err:  "invalid address: Invalid base58 string",
		},
		{
			name: "short address",
			addr: "xxx",
			err:  "invalid address: Invalid address length",
		},
		{
			name: "invalid base58 address",
			addr: "2blYafFtdkCRNcCyuDvsATV66GvBR9xfvjy",
			err:  "invalid address: Invalid base58 character",
		},
		{
			name: "valid address",
			addr: "2bfYafFtdkCRNcCyuDvsATV66GvBR9xfvjy",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var a Address
			err := a.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, tc.addr)))
			if tc.err != "" {
				require.Equal(t, errors.New(tc.err), err)
			} else {
				require.NoError(t, err)
				addr, err := cipher.DecodeBase58Address(tc.addr)
				require.NoError(t, err)

				require.Equal(t, addr, a.Address)
			}
		})
	}

	var a Address
	err := a.UnmarshalJSON([]byte("invalidjson"))
	testutil.RequireError(t, err, "invalid character 'i' looking for beginning of value")
}

func TestCoinsMarshalJSON(t *testing.T) {
	c := Coins(111)

	data, err := c.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"0.000111"`, string(data))

	c = Coins(math.MaxUint64)
	_, err = c.MarshalJSON()
	testutil.RequireError(t, err, "Droplet string conversion failed: Value is too large")
}

func TestCoinsUnmarshalJSON(t *testing.T) {
	cases := []struct {
		c        string
		expected uint64
		err      string
	}{
		{
			c:   "-1",
			err: "Droplet string conversion failed: Negative balance",
		},
		{
			c:        "0",
			expected: 0,
		},
		{
			c:   "0.1234567",
			err: "Droplet string conversion failed: Too many decimal places",
		},
		{
			c:        "1.001",
			expected: 1001e3,
		},
		{
			c:        "1.234567",
			expected: 1234567,
		},
		{
			c:        "9",
			expected: 9e6,
		},
		{
			c:   ".",
			err: "can't convert . to decimal",
		},
		{
			c:   "inf",
			err: "can't convert inf to decimal",
		},
		{
			c:        "9223372036854.775807",
			expected: uint64(math.MaxInt64),
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprint(tc.c), func(t *testing.T) {
			var c Coins
			err := c.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, tc.c)))
			if tc.err != "" {
				require.Equal(t, errors.New(tc.err), err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, c.Value())
			}
		})
	}

	var c Coins
	err := c.UnmarshalJSON([]byte("invalidjson"))
	testutil.RequireError(t, err, "invalid character 'i' looking for beginning of value")
}

func TestHoursMarshalJSON(t *testing.T) {
	c := Hours(111)

	data, err := c.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"111"`, string(data))
}

func TestHoursUnmarshalJSON(t *testing.T) {
	cases := []struct {
		c        string
		expected uint64
		err      string
	}{
		{
			c:   "-1",
			err: "invalid hours value: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			c:        "0",
			expected: 0,
		},
		{
			c:   "0.1",
			err: "invalid hours value: strconv.ParseUint: parsing \"0.1\": invalid syntax",
		},
		{
			c:        "9",
			expected: 9,
		},
		{
			c:   ".",
			err: "invalid hours value: strconv.ParseUint: parsing \".\": invalid syntax",
		},
		{
			c:   "inf",
			err: "invalid hours value: strconv.ParseUint: parsing \"inf\": invalid syntax",
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprint(tc.c), func(t *testing.T) {
			var c Hours
			err := c.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, tc.c)))
			if tc.err != "" {
				require.Equal(t, errors.New(tc.err), err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, c.Value())
			}
		})
	}

	var c Hours
	err := c.UnmarshalJSON([]byte("invalidjson"))
	testutil.RequireError(t, err, "invalid character 'i' looking for beginning of value")
}

func TestSHA256MarshalJSON(t *testing.T) {
	hash := "97dd062820314c46da0fc18c8c6c10bfab1d5da80c30adc79bbe72e90bfab11d"

	c := SHA256{cipher.MustSHA256FromHex(hash)}

	data, err := c.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"97dd062820314c46da0fc18c8c6c10bfab1d5da80c30adc79bbe72e90bfab11d"`, string(data))
}

func TestSHA256UnmarshalJSON(t *testing.T) {
	cases := []struct {
		c        string
		expected cipher.SHA256
		err      string
	}{
		{
			c:   "",
			err: "invalid SHA256 hash: Invalid hex length",
		},

		{
			c:   "foo",
			err: "invalid SHA256 hash: encoding/hex: invalid byte: U+006F 'o'",
		},

		{
			c:   "97dd0628",
			err: "invalid SHA256 hash: Invalid hex length",
		},

		{
			c:   "97dd062820314c46da0fc18c8c6c10bfab1d5da80c30adc79bbe72e90bfab11d97",
			err: "invalid SHA256 hash: Invalid hex length",
		},

		{
			c:   "97dd062820314c46da0fc18c8c6c10bfab1d5da80c30adc79bbe72e90bfab11Q",
			err: "invalid SHA256 hash: encoding/hex: invalid byte: U+0051 'Q'",
		},

		{
			c:        "97dd062820314c46da0fc18c8c6c10bfab1d5da80c30adc79bbe72e90bfab11d",
			expected: cipher.MustSHA256FromHex("97dd062820314c46da0fc18c8c6c10bfab1d5da80c30adc79bbe72e90bfab11d"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.c, func(t *testing.T) {
			var c SHA256
			err := c.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, tc.c)))
			if tc.err != "" {
				require.Equal(t, errors.New(tc.err), err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, c.SHA256)
			}
		})
	}
}
