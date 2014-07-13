package linux

import (
	//"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute(t *testing.T) {
	route := NewRoute()

	// IsInstalled
	if !route.IsInstalled() {
		t.Skip("skipping test, program not installed")
	}
}
