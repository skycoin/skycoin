package bip44

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/bip32"
	"github.com/skycoin/skycoin/src/cipher/bip39"
)

func mustDefaultSeed(t *testing.T) []byte {
	mnemonic := "dizzy cigar grant ramp inmate uniform gold success able payment faith practice"
	passphrase := ""
	seed, err := bip39.NewSeed(mnemonic, passphrase)
	require.NoError(t, err)
	require.Equal(t, "24e563fb095d766df3862c70432cc1b2210b24d232da69af7af09d2ec86d28782ce58035bae29994c84081836aebe36a9b46af1578262fefc53e37efbe94be57", hex.EncodeToString(seed))
	return seed
}

func TestNewCoin(t *testing.T) {
	// bad seed
	_, err := NewCoin(make([]byte, 3), CoinTypeBitcoin)
	require.Equal(t, err, bip32.ErrInvalidSeedLength)

	// bad coin_type
	_, err = NewCoin(mustDefaultSeed(t), CoinType(bip32.FirstHardenedChild))
	require.Equal(t, err, ErrInvalidCoinType)
	_, err = NewCoin(mustDefaultSeed(t), CoinType(1+bip32.FirstHardenedChild))
	require.Equal(t, err, ErrInvalidCoinType)

	c, err := NewCoin(mustDefaultSeed(t), CoinTypeBitcoin)
	require.NoError(t, err)

	account, err := c.Account(0)
	require.NoError(t, err)
	require.Equal(t, "xprv9yKAFQtFghZSe4mfdpdqFm1WWmGeQbYMB4MSGUB85zbKGQgSxty4duZb8k6hNoHVd2UR7Y3QhWU3rS9wox9ewgVG7gDLyYTL4yzEuqUCjvF", account.String())
	require.Equal(t, "xpub6CJWevR9X57jrYr8jrAqctxF4o78p4GCYHH34rajeL8J9D1bWSHKBht4yzwiTQ4FP4HyQpx99iLxvU54rbEbcxBUgxzTGGudBVXb1N2gcHF", account.PublicKey().String())

	account, err = c.Account(1)
	require.NoError(t, err)
	require.Equal(t, "xprv9yKAFQtFghZSgShGXkxHsYQfFaqMyutf3izng8tV4Tmp7gidQUPB8kCuv66yukidivM2oSaUvGus8ffnYvYKChB7DME2H2AvUq8LM2rXUzF", account.String())
	require.Equal(t, "xpub6CJWevR9X57jtvmjdnVJEgMPocfrPNcWQwvPUXJ6coJnzV3mx1hRgYXPmQJh5vLQvrVCY8LtJB5xLLiPJVmpSwBe2yhonQLoQuSsCF8YPLN", account.PublicKey().String())

	_, err = c.Account(0x80000000)
	require.Equal(t, err, ErrInvalidAccount)
	_, err = c.Account(0x80000001)
	require.Equal(t, err, ErrInvalidAccount)

	external, err := account.External()
	require.NoError(t, err)
	require.Equal(t, "xprv9zjsvjLiqSerDzbeRXPeXwz8tuQ7eRUABkgFAgLPHw1KzGKkgBhJhGaMYHM8j2KDXBZTCv4m19qjxrrD7gusrtdpZ7xzJywdXHaMZEjf3Uv", external.String())
	require.Equal(t, "xpub6DjELEscfpD9SUg7XYveu5vsSwEc3tC1Yybqy4jzrGYJs4euDj1ZF4tqPZYvViMn9cvBobHyubuuh69PZ1szaBBx5oxPiQzD492B6C4QDHe", external.PublicKey().String())

	external0, err := external.NewPublicChildKey(0)
	require.NoError(t, err)
	require.Equal(t, "034d36f3bcd74e19204e75b81b9c0726e41b799858b92bab73f4cd7498308c5c8b", hex.EncodeToString(external0.Key))

	external1, err := external.NewPublicChildKey(1)
	require.NoError(t, err)
	require.Equal(t, "02f7309e9f559d847ee9cc9ee144cfa490791e33e908fdbde2dade50a389408b01", hex.EncodeToString(external1.Key))

	change, err := account.Change()
	require.NoError(t, err)
	require.Equal(t, "xprv9zjsvjLiqSerGzJyBrpZgCaGpQCeFDnZEuAV714WigmFyHT4nFLhZLeuHzLNE19PgkZeQ5Uf2pjFZjQTHbkugDbmw5TAPAvgo2jsaTnZo2A", change.String())
	require.Equal(t, "xpub6DjELEscfpD9VUPSHtMa3LX1NS38egWQc865uPU8H2JEr5nDKnex78yP9GxhFr5cnCRgiQF1dkv7aR7moraPrv73KHwSkDaXdWookR1Sh9p", change.PublicKey().String())

	change0, err := change.NewPublicChildKey(0)
	require.NoError(t, err)
	require.Equal(t, "026d3eb891e81ecabedfa8560166af383457aedaf172af9d57d00508faa5f57c4c", hex.EncodeToString(change0.Key))

	change1, err := change.NewPublicChildKey(1)
	require.NoError(t, err)
	require.Equal(t, "02681b301293fdf0292cd679b37d60b92a71b389fd994b2b57c8daf99532bfb4a5", hex.EncodeToString(change1.Key))
}
