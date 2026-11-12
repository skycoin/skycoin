// Package skywallet src/hardware-wallet/skywallet/transaction_sign_test.go c1-app-hw
package skywallet

import (
	"reflect"
	"testing"

	"github.com/gogo/protobuf/proto"

	"github.com/skycoin/skycoin/src/hardware-wallet/messages"
)

// The firmware signs an input with the key at address_n[0], and skips one
// whose address_n is empty (tiny-firmware/firmware/fsm_impl.c msgTxAckImpl).
func TestTxAckTypesAddressN(t *testing.T) {
	ins, outs := txAckTypes(
		[]*messages.SkycoinTransactionInput{
			{HashIn: proto.String("a"), Index: proto.Uint32(0)},
			{HashIn: proto.String("b"), Index: proto.Uint32(2)},
			{HashIn: proto.String("c")},
		},
		[]*messages.SkycoinTransactionOutput{
			{Address: proto.String("x"), Coin: proto.Uint64(1), Hour: proto.Uint64(2)},
			{Address: proto.String("y"), Coin: proto.Uint64(3), Hour: proto.Uint64(4), AddressIndex: proto.Uint32(5)},
		},
	)
	var gotIn [][]uint32
	for _, in := range ins {
		gotIn = append(gotIn, in.AddressN)
	}
	if want := [][]uint32{{0}, {2}, nil}; !reflect.DeepEqual(gotIn, want) {
		t.Errorf("inputs address_n = %v, want %v", gotIn, want)
	}
	if outs[0].AddressN != nil || !reflect.DeepEqual(outs[1].AddressN, []uint32{5}) {
		t.Errorf("outputs address_n = %v, %v, want none and [5]", outs[0].AddressN, outs[1].AddressN)
	}
}
