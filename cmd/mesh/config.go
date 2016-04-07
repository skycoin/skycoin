package main

import (
    "time"
)

import (
    "github.com/skycoin/skycoin/src/daemon/gnet"
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/skycoin/skycoin/src/mesh")

type Config struct {
    Node mesh.Config
    TCPConnections []TCPOutgoingConnectionConfig
    TCPConfig gnet.Config
    Routes []RouteConfig
    RouteToPipe int
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

type RouteConfig struct {
    PeerPubKeys []cipher.PubKey
}
