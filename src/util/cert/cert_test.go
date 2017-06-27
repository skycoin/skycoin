package cert

import (
	"crypto/tls"
	"os"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCert(t *testing.T) {
	defer os.Remove("certtest.pem")
	defer os.Remove("keytest.pem")
	err := GenerateCert("certtest.pem", "keytest.pem", "127.0.0.1", "org",
		2048, false, utc.Now(), time.Hour*24)
	assert.Nil(t, err)
	_, err = tls.LoadX509KeyPair("certtest.pem", "keytest.pem")
	assert.Nil(t, err)
}
