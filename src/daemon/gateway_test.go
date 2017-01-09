package daemon

import (
	"reflect"
	"testing"

	"github.com/skycoin/skycoin/src/visor"
)

func TestFbyAddresses(t *testing.T) {
	tests := []struct {
		name    string
		addrs   []string
		outputs []visor.ReadableOutput
		want    []visor.ReadableOutput
	}{
		// TODO: Add test cases.
		{
			"filter with one address",
			[]string{"abc"},
			[]visor.ReadableOutput{
				{
					Address: "abc",
				},
				{
					Address: "cde",
				},
			},
			[]visor.ReadableOutput{
				{
					Address: "abc",
				},
			},
		},
		{
			"filter with multiple addresses",
			[]string{"abc", "cde"},
			[]visor.ReadableOutput{
				{
					Address: "abc",
				},
				{
					Address: "cde",
				},
				{
					Address: "efg",
				},
			},
			[]visor.ReadableOutput{
				{
					Address: "abc",
				},
				{
					Address: "cde",
				},
			},
		},
	}
	for _, tt := range tests {
		outs := FbyAddresses(tt.addrs)(tt.outputs)
		if !reflect.DeepEqual(outs, tt.want) {
			t.Errorf("%q. FbyAddresses() = %v, want %v", tt.name, outs, tt.want)
		}
	}
}

func TestFbyHashes(t *testing.T) {
	type args struct {
		hashes []string
	}
	tests := []struct {
		name    string
		hashes  []string
		outputs []visor.ReadableOutput
		want    []visor.ReadableOutput
	}{
		// TODO: Add test cases.
		{
			"filter with one hash",
			[]string{"abc"},
			[]visor.ReadableOutput{
				{
					Hash: "abc",
				},
				{
					Hash: "cde",
				},
			},
			[]visor.ReadableOutput{
				{
					Hash: "abc",
				},
			},
		},
		{
			"filter with multiple hash",
			[]string{"abc", "cde"},
			[]visor.ReadableOutput{
				{
					Hash: "abc",
				},
				{
					Hash: "cde",
				},
				{
					Hash: "efg",
				},
			},
			[]visor.ReadableOutput{
				{
					Hash: "abc",
				},
				{
					Hash: "cde",
				},
			},
		},
	}
	for _, tt := range tests {
		outs := FbyHashes(tt.hashes)(tt.outputs)
		if !reflect.DeepEqual(outs, tt.want) {
			t.Errorf("%q. FbyHashes() = %v, want %v", tt.name, outs, tt.want)
		}
	}
}
