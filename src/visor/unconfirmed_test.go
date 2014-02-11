package visor

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestUnconfirmedTxnHash(t *testing.T) {
    utx := createUnconfirmedTxn()
    assert.Equal(t, utx.Hash(), utx.Txn.Header.Hash)
}

func TestNewUnconfirmedTxnPool(t *testing.T) {
    ut := NewUnconfirmedTxnPool()
    assert.NotNil(t, ut.Txns)
    assert.Equal(t, len(ut.Txns), 0)
}

func TestSetAnnounced(t *testing.T) {
    ut := NewUnconfirmedTxnPool()
    assert.Equal(t, len(ut.Txns), 0)
    // Unknown should be safe and a noop
    assert.NotPanics(t, func() {
        ut.SetAnnounced(coin.SHA256{}, util.Now())
    })
    assert.Equal(t, len(ut.Txns), 0)
    utx := createUnconfirmedTxn()
    assert.True(t, utx.Announced.IsZero())
    ut.Txns[utx.Hash()] = utx
    now := util.Now()
    ut.SetAnnounced(utx.Hash(), now)
    assert.Equal(t, ut.Txns[utx.Hash()].Announced, now)
}
