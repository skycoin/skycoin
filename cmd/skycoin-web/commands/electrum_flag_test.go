package commands

import (
	"reflect"
	"testing"

	"github.com/skycoin/skycoin/src/electrum"
)

func TestElectrumServerList(t *testing.T) {
	for _, tc := range []struct {
		flag string
		want []string
	}{
		{"default", electrum.DefaultServers},
		{"none", nil},
		{"", nil},
		{" ssl://a:1 , tcp://b:2 ", []string{"ssl://a:1", "tcp://b:2"}},
	} {
		if got := electrumServerList(tc.flag); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("electrumServerList(%q) = %v, want %v", tc.flag, got, tc.want)
		}
	}
}

func TestElectrumFlagDefaultsToBuiltInList(t *testing.T) {
	if got := RootCmd.Flags().Lookup("btc-electrum-url").DefValue; got != "default" {
		t.Fatalf("--btc-electrum-url default = %q, want default", got)
	}
}
