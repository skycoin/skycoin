package daemon

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
)

func TestSerializeRejectMessageReplyingIntroduction(t *testing.T) {
	// TODO: Test fixture
	msgcfg := []MessageConfig{
		NewMessageConfig("INTR", &IntroductionMessage{}),
		NewMessageConfig("RJCT", &RejectMessage{})}
	for i, _ := range msgcfg {
		to := reflect.TypeOf(msgcfg[i].Message)
		_, succ := gnet.MessageIDMap[to]
		if !succ {
			gnet.RegisterMessage(msgcfg[i].Prefix, msgcfg[i].Message)
			fmt.Println("Register %v", msgcfg[i].Prefix)
		}
	}

	originalMessage := NewIntroductionMessage(1234, 2, 6000)

	peers := make([]IPAddr, 0)
	addr, _ := NewIPAddr("192.168.1.1:6001")
	peers = append(peers, addr)
	addr, _ = NewIPAddr("192.168.1.2:6002")
	peers = append(peers, addr)
	addr, _ = NewIPAddr("192.168.1.3:6003")
	peers = append(peers, addr)
	addr, _ = NewIPAddr("192.168.1.4:6004")
	peers = append(peers, addr)

	// TODO: Expected message length
	bLen := encoder.SerializeAtomic(uint32(0))
	errCode := GetErrorCode(pex.ErrPeerlistFull)
	prefix := gnet.MessagePrefixFromString("RJCT")

	b := make([]byte, 0)
	// Expected message length
	b = append(b, bLen...)
	// Message prefix
	b = append(b, prefix[:]...)
	// Rejected message prefix
	prefix = gnet.MessagePrefixFromString("INTR")
	b = append(b, prefix[:]...)
	// Error code for peer list overflow
	b = append(b, encoder.SerializeAtomic(errCode)...)
	// Reason length
	b = append(b, encoder.SerializeAtomic(uint32(0))...)
	// Reason string
	// Addresses
	b = append(b, encoder.Serialize(peers)...)

	msg := NewRejectMessage(originalMessage, pex.ErrPeerlistFull, "", peers)
	realBytes := gnet.EncodeMessage(msg)

	require.Equal(t, b, realBytes)
}
