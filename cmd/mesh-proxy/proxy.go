package main

import(
	"net"
)

import(
    "github.com/skycoin/skycoin/src/mesh")

type PortRange struct {
	SourceIP net.IP
	Minimum uint16
	Maximum uint16
}

type ProxyConfig struct {
	SourcePortRanges []PortRange
	ClientSourcePortLimit uint32
}

type Proxy struct {
	Config ProxyConfig
}

func NewProxy(config ProxyConfig) *Proxy {
    ret := &Proxy{}
    ret.Config = config
    ret.MeshNode = node
    return ret
}

func (self *Proxy) Run() {
	
}