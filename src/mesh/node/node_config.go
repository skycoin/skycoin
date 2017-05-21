package node

type NodeConfig struct {
	ClientAddr  string   // address for talking with servers
	ServerAddrs []string // addresses of servers
	AppTalkPort int      // will establish connection with apps here
	Hostname    string   // can be empty, in this case node will respond only by pubkey
}
