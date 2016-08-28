package main

import (
    "time"
)

import (
    "github.com/skycoin/skycoin/src/daemon/gnet"
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/skycoin/skycoin/src/mesh")

type Config struct {
    Node mesh.NodeConfig
    TCPConnections []TCPOutgoingConnectionConfig
    TCPConfig gnet.Config
    // Index into Routes in NodeConfig
    StdoutToRoute int
    IncomingToStdout bool
    Proxy ProxyConfig
}

type PhysicalConnectionConfig struct {
    PeerPubKey cipher.PubKey
    Type string
}

// Type: "tcp"
type TCPOutgoingConnectionConfig struct {
    PhysicalConnectionConfig
    Endpoint string
    RetryDelay time.Duration
}

