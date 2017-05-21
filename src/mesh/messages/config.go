package messages

import (
	"fmt"
	"gopkg.in/gcfg.v1"
	"time"
)

type ConfigStruct struct {
	SendInterval      time.Duration
	SendIntervalNum   uint32
	TimeUnit          time.Duration
	TimeUnitNum       int
	StartPort         uint32
	VPNSubnet         string
	SimulateDelay     bool
	MaxSimulatedDelay int
	MaxPacketSize     int
	ProxyPacketSize   int
	ProxyTimeout      time.Duration
	AppTimeout        uint32
	ConnectionTimeout uint32
	TransportTimeout  uint32
	RetransmitLimit   int
	LogLevel          uint8
	MaxBuffer         uint64
	MsgSrvTimeout     uint32
}

type ConfigFromFile struct {
	General struct {
		LogLevel   string
		StartPort  uint32
		AppTimeout uint32
		VPNSubnet  string
	}
	Transport struct {
		TransportTimeout  uint32
		RetransmitLimit   int
		SimulateDelay     bool
		MaxSimulatedDelay int
	}
	Connection struct {
		ConnectionTimeout uint32
		MaxPacketSize     int
	}
	Congestion struct {
		MaxBuffer    uint64
		SendInterval int
		TimeUnit     int
	}
	Proxy struct {
		ProxyPacketSize int
		ProxyTimeout    time.Duration
	}
	MsgSrv struct {
		MsgSrvTimeout uint32
	}
}

var config = &ConfigStruct{ // default values
	LogLevel:          INFO,
	StartPort:         6000,
	AppTimeout:        10000,
	VPNSubnet:         "192.168.11.",
	ProxyPacketSize:   16384,
	ProxyTimeout:      10000 * time.Millisecond,
	TransportTimeout:  500,
	RetransmitLimit:   10,
	SimulateDelay:     false,
	MaxSimulatedDelay: 300,
	ConnectionTimeout: 5000,
	MaxPacketSize:     16384,
	MaxBuffer:         8192,
	SendInterval:      1500 * time.Microsecond,
	SendIntervalNum:   1500,
	TimeUnit:          10 * time.Microsecond,
	TimeUnitNum:       10,
	MsgSrvTimeout:     500,
}

func init() {
	cfgFromFile := &ConfigFromFile{}
	err := gcfg.ReadFileInto(cfgFromFile, "/etc/meshnet.cfg")
	if err != nil {
		fmt.Println("Cannot read settings from file, applying defaults. Error:", err)
		return
	}

	if cfgFromFile.General.LogLevel == "DEBUG" {
		config.LogLevel = DEBUG
	} else {
		config.LogLevel = INFO
	}
	config.StartPort = cfgFromFile.General.StartPort
	config.AppTimeout = cfgFromFile.General.AppTimeout
	config.VPNSubnet = cfgFromFile.General.VPNSubnet

	config.TransportTimeout = cfgFromFile.Transport.TransportTimeout
	config.RetransmitLimit = cfgFromFile.Transport.RetransmitLimit
	config.SimulateDelay = cfgFromFile.Transport.SimulateDelay
	config.MaxSimulatedDelay = cfgFromFile.Transport.MaxSimulatedDelay

	config.ConnectionTimeout = cfgFromFile.Connection.ConnectionTimeout
	config.MaxPacketSize = cfgFromFile.Connection.MaxPacketSize

	config.MaxBuffer = cfgFromFile.Congestion.MaxBuffer
	config.SendInterval = time.Duration(cfgFromFile.Congestion.SendInterval) * time.Microsecond
	config.SendIntervalNum = uint32(cfgFromFile.Congestion.SendInterval)
	config.TimeUnit = time.Duration(cfgFromFile.Congestion.TimeUnit) * time.Microsecond
	config.TimeUnitNum = cfgFromFile.Congestion.TimeUnit

	config.ProxyPacketSize = cfgFromFile.Proxy.ProxyPacketSize
	config.ProxyTimeout = time.Duration(cfgFromFile.Proxy.ProxyTimeout) * time.Millisecond

	config.MsgSrvTimeout = cfgFromFile.MsgSrv.MsgSrvTimeout
}

func GetConfig() *ConfigStruct {
	return config
}
