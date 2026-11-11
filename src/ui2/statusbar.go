package ui2

import "strconv"

// Bar is the strip under the top bar that warns about the node, as the
// Angular wallet's header does.
type Bar struct {
	Unreachable bool
	Syncing     bool
	Current     uint64
	Highest     uint64
	Percent     string
	NoPeers     bool
	Pending     uint64
	NoCSRF      bool
}

// Any reports whether the strip has anything to show.
func (b Bar) Any() bool {
	return b.Unreachable || b.Syncing || b.NoPeers || b.Pending > 0 || b.NoCSRF
}

// LoadBar reads the node's health and sync progress.
func LoadBar(n Node) Bar {
	h, err := n.Health()
	if err != nil {
		return Bar{Unreachable: true}
	}
	b := Bar{
		NoPeers: h.OpenConnections == 0,
		Pending: h.BlockchainMetadata.Unconfirmed,
		NoCSRF:  !h.CSRFEnabled,
	}
	if p, err := n.BlockchainProgress(); err == nil && p.Highest > p.Current {
		b.Syncing, b.Current, b.Highest = true, p.Current, p.Highest
		b.Percent = strconv.FormatFloat(float64(p.Current)*100/float64(p.Highest), 'f', 1, 64)
	}
	return b
}
