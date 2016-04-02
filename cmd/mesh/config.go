package main

import (
    "github.com/skycoin/skycoin/src/daemon/gnet"
    "github.com/skycoin/skycoin/src/cipher")

type Config struct {
    MyPubKey cipher.PubKey
    Connections []PhysicalConnectionConfig
    TCPServer* gnet.Config
    RouteToPipe* RouteConfig
}

type PhysicalConnectionConfig struct {
    PeerPubKey cipher.PubKey
    Type string
}

// Type: "tcp"
type TCPOutgoingConnectionConfig struct {
    PhysicalConnectionConfig

    Address string
    Port int
}

type RouteConfig struct {
    PeerPubKeys []cipher.PubKey
    FullDuplex bool
}
