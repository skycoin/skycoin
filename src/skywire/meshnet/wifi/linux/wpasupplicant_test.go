package linux

import (
	//"github.com/stretchr/testify/assert"
	"testing"
)

func TestWPASupplicant(t *testing.T) {
	wpasupplicant := NewWPASupplicant()

	// IsInstalled
	if !wpasupplicant.IsInstalled() {
		t.Skip("skipping test, program not installed")
	}
}
