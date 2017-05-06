package node

type NodeConfig struct {
	ClientAddr  string   // address for talking with servers
	ServerAddrs []string // addresses of servers
	AppTalkAddr string   // will establish connection with apps here
}
