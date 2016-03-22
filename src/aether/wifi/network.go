package network

import (
	"errors"
	"net"

	"github.com/skycoin/skycoin/src/aether/wifi/linux"
	"gopkg.in/op/go-logging.v1"
)

type Address struct {
	IPNet   net.IPNet
	Gateway net.IP
}

type Route struct {
	IPNet   net.IPNet
	Gateway net.IP
	Metric  string
}

var (
	softwareType    = ""
	ErrAuthRequired = errors.New("sudo authentication required")
	logger          = logging.MustGetLogger("darknet.network")

	dhclient       = linux.NewDHClient()
	ifconfig       = linux.NewIFConfig()
	iwconfig       = linux.NewIWConfig()
	iwlist         = linux.NewIWList()
	networkmanager = linux.NewNetworkManager()
	resolvconf     = linux.NewResolvConf()
	rfkill         = linux.NewRFKill()
	route          = linux.NewRoute()
	sysfs          = linux.NewSysfs()
	udevadm        = linux.NewUDevAdm()
	wpasupplicant  = linux.NewWPASupplicant()
)
