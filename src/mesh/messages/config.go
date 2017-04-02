package messages

import (
	"fmt"
	"gopkg.in/gcfg.v1"
	"time"
)

const (
	DEBUG = iota
	INFO
)

type ConfigStruct struct {
	SendInterval      time.Duration
	TimeUnit          time.Duration
	StartPort         uint32
	VPNSubnet         string
	SimulateDelay     bool
	MaxSimulatedDelay int
	MaxPacketSize     int
	SocksPacketSize   int
	AppTimeout        uint32
	ConnectionTimeout uint32
	TransportTimeout  uint32
	RetransmitLimit   int
	LogLevel          uint8
	MaxBuffer         uint64
}

type ConfigFromFile struct {
	General struct {
		LogLevel        string
		StartPort       uint32
		AppTimeout      uint32
		VPNSubnet       string
		SocksPacketSize int
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
}

var config = &ConfigStruct{ // default values
	LogLevel:          INFO,
	StartPort:         6000,
	AppTimeout:        1000000,
	VPNSubnet:         "192.168.11.",
	SocksPacketSize:   16384,
	TransportTimeout:  100000,
	RetransmitLimit:   10,
	SimulateDelay:     false,
	MaxSimulatedDelay: 500,
	ConnectionTimeout: 1000000,
	MaxPacketSize:     16384,
	MaxBuffer:         8192,
	SendInterval:      1500 * time.Microsecond,
	TimeUnit:          10 * time.Microsecond,
}

func init() {
	cfgFromFile := &ConfigFromFile{}
	err := gcfg.ReadFileInto(cfgFromFile, "/etc/meshnet.cfg")
	if err != nil {
		fmt.Println("Cannot read settings from file, applying defaults. Error:", err)
		panic(err)
		//return
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
	config.TimeUnit = time.Duration(cfgFromFile.Congestion.TimeUnit) * time.Microsecond
}

func GetConfig() *ConfigStruct {
	return config
}

func SetLogLevel(loglevel uint8) {
	if loglevel == DEBUG || loglevel == INFO {
		config.LogLevel = loglevel
	}
}

func SetDebugLogLevel() {
	SetLogLevel(DEBUG)
}

func SetInfoLogLevel() {
	SetLogLevel(INFO)
}

func IsDebug() bool {
	return config.LogLevel == DEBUG
}
