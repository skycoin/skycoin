package api

import (
	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

//go:generate mockery -name Gatewayer -case underscore -inpkg -testonly

// Gateway is the api gateway
type Gateway struct {
	Device *skyWallet.Device
}

// NewGateway creates a Gateway
func NewGateway(device *skyWallet.Device) *Gateway {
	return &Gateway{
		device,
	}
}

// Gatewayer interface for Gateway methods
type Gatewayer interface {
	skyWallet.Devicer
}
