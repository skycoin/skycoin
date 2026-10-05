package electrum

// DefaultServers are public Bitcoin Electrum servers, tried in this order. The
// wallet's node settings offer the same list. Entries come from Electrum's own
// list (https://github.com/spesmilo/electrum/blob/master/electrum/chains/servers.json)
// and each answered server.version when last checked, on 2026-10-04.
var DefaultServers = []string{
	"ssl://electrum.blockstream.info:50002",
	"ssl://fortress.qtornado.com:443",
	"ssl://electrum.emzy.de:50002",
}
