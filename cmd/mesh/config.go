package main

import (
    "time"
)

import (
    "github.com/skycoin/skycoin/src/daemon/gnet"
    "github.com/skycoin/skycoin/src/cipher")

type Config struct {
    MyPubKey cipher.PubKey
    TCPConnections []TCPOutgoingConnectionConfig
    TCPConfig gnet.Config
    RouteToPipe* RouteConfig
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
