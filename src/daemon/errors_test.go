package daemon

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/daemon/gnet"
)

func TestDisconnectReasonCode(t *testing.T) {
	c := DisconnectReasonToCode(ErrDisconnectIdle)
	require.NotEqual(t, uint16(0), c)

	r := DisconnectCodeToReason(c)
	require.Equal(t, ErrDisconnectIdle, r)

	// unknown reason is fine
	c = DisconnectReasonToCode(gnet.DisconnectReason(errors.New("foo")))
	require.Equal(t, uint16(0), c)

	r = DisconnectCodeToReason(c)
	require.Equal(t, ErrDisconnectUnknownReason, r)

	// unknown code is fine
	r = DisconnectCodeToReason(999)
	require.Equal(t, ErrDisconnectUnknownReason, r)
}
