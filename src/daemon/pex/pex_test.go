package pex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAddress(t *testing.T) {
	// empty string
	assert.False(t, validateAddress("", false, 7200))
	// doubled ip:port
	assert.False(t, validateAddress("112.32.32.14:100112.32.32.14:101", false, 7200))
	// requires port
	assert.False(t, validateAddress("112.32.32.14", false, 7200))
	// not ip
	assert.False(t, validateAddress("112", false, 7200))
	assert.False(t, validateAddress("112.32", false, 7200))
	assert.False(t, validateAddress("112.32.32", false, 7200))
	// bad part
	assert.False(t, validateAddress("112.32.32.14000", false, 7200))
	// large port
	assert.False(t, validateAddress("112.32.32.14:66666", false, 7200))
	// unspecified
	assert.False(t, validateAddress("0.0.0.0:8888", false, 7200))
	// no ip
	assert.False(t, validateAddress(":8888", false, 7200))
	// multicast
	assert.False(t, validateAddress("224.1.1.1:8888", false, 7200))
	// invalid ports
	assert.False(t, validateAddress("112.32.32.14:0", false, 7200))
	assert.False(t, validateAddress("112.32.32.14:1", false, 7200))
	assert.False(t, validateAddress("112.32.32.14:10", false, 7200))
	assert.False(t, validateAddress("112.32.32.14:100", false, 7200))
	assert.False(t, validateAddress("112.32.32.14:1000", false, 7200))
	assert.False(t, validateAddress("112.32.32.14:1023", false, 7200))
	assert.False(t, validateAddress("112.32.32.14:65536", false, 7200))
	assert.False(t, validateAddress("112.32.32.14:7200", false, 7201))
	// valid ones
	assert.True(t, validateAddress("112.32.32.14:1024", false, 1024))
	assert.True(t, validateAddress("112.32.32.14:10000", false, 10000))
	assert.True(t, validateAddress("112.32.32.14:65535", false, 65535))
	// localhost is allowed
	assert.True(t, validateAddress("127.0.0.1:8888", true, 8888))
	// localhost is not allowed
	assert.False(t, validateAddress("127.0.0.1:8888", false, 7200))
}
