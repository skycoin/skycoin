package visor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/params"
)

func TestConfigVerifyBlockPublisher(t *testing.T) {
	pk, sk := cipher.GenerateKeyPair()
	_, otherSk := cipher.GenerateKeyPair()

	cases := []struct {
		name      string
		publisher bool
		pubkey    cipher.PubKey
		seckey    cipher.SecKey
		errSub    string
	}{
		{"publisher null seckey", true, pk, cipher.SecKey{}, "secret key is not set"},
		{"publisher mismatched seckey", true, pk, otherSk, "invalid seckey for pubkey"},
		{"publisher matching seckey", true, pk, sk, ""},
		{"non-publisher null seckey", false, pk, cipher.SecKey{}, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewConfig()
			c.Distribution = params.MainNetDistribution
			c.IsBlockPublisher = tc.publisher
			c.BlockchainPubkey = tc.pubkey
			c.BlockchainSeckey = tc.seckey

			err := c.Verify()
			if tc.errSub == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.errSub)
		})
	}
}
