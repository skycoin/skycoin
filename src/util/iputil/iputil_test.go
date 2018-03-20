package iputil

import (
	"fmt"
	"testing"
)

func TestSplitAddr(t *testing.T) {
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
	}

	for _, test := range testData {
		actual := IsLocalhost(test.host)

		if test.expected != actual {
			t.Errorf("Expected %t is not equal to actual %t for host %s",
				test.expected, actual, test.host)
		}
	}
}

func TestIsLocalhost(t *testing.T) {
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
			err:   nil,
		},
		{
			input: "0.0.0.0:",
			host:  "0.0.0.0",
			err:   fmt.Errorf("Invalid port in %s", "0.0.0.0:"),
		},
		{
			input: ":9999",
			port:  9999,
			err:   fmt.Errorf("Invalid port in %s", ":9999"),
		},
		{
			input: "127.0.0.1",
			host:  "127.0.0.1",
			err:   fmt.Errorf("Invalid addr %s", "127.0.0.1"),
		},
	}

	for _, test := range testData {
		addr, port, err := SplitAddr(test.input)

		if test.err == nil && err != nil {
			t.Errorf("Unexpected error %v", err)
			return
		}

		if test.err != nil && err != nil && err.Error() != test.err.Error() {
			t.Errorf("Expected error %v, actual %v", test.err, err)
			return
		}

		if addr != test.host {
			t.Errorf("Wrong host, expected %s actual %s", test.host, addr)
			return
		}

		if port != test.port {
			t.Errorf("Wrong port, expected %d actual %d", test.port, port)
			return
		}
	}
}
