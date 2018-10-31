package iputil

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLocalhost(t *testing.T) {
	testData := []struct {
		host     string
		expected bool
	}{
		{
			host:     "0:0:0:0:0:0:0:1",
			expected: true,
		},
		{
			host:     "localhost",
			expected: true,
		},
		{
			host:     "127.0.0.1",
			expected: true,
		},
		{
			host:     "localhost",
			expected: true,
		},
		{
			host:     "85.56.12.34",
			expected: false,
		},
		{
			host:     "::1",
			expected: true,
		},
		{
			host:     "::",
			expected: false,
		},
		{
			host:     "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected: false,
		},
		{
			host:     "",
			expected: false,
		},
	}

	for _, tc := range testData {
		t.Run(tc.host, func(t *testing.T) {
			actual := IsLocalhost(tc.host)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestSplitAddr(t *testing.T) {
	testData := []struct {
		input string
		host  string
		port  uint16
		err   error
	}{
		{
			input: "0.0.0.0:8888",
			host:  "0.0.0.0",
			port:  8888,
		},
		{
			input: "0.0.0.0:",
			err:   ErrInvalidPort,
		},
		{
			input: "0.0.0.0:x",
			err:   ErrInvalidPort,
		},
		{
			input: ":9999",
			err:   ErrMissingIP,
		},
		{
			input: "127.0.0.1",
			err: &net.AddrError{
				Err:  "missing port in address",
				Addr: "127.0.0.1",
			},
		},
		{
			input: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234",
			host:  "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			port:  1234,
		},
		{
			input: "[::]:1234",
			host:  "::",
			port:  1234,
		},
		{
			input: "[::]:x",
			err:   ErrInvalidPort,
		},
	}

	for _, tc := range testData {
		t.Run(tc.input, func(t *testing.T) {
			addr, port, err := SplitAddr(tc.input)

			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, tc.host, addr)
			require.Equal(t, tc.port, port)
		})
	}
}
