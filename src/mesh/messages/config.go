package messages

import "time"

const (
	DEBUG = iota
	INFO
)

type ConfigStruct struct {
	SendInterval      time.Duration
	TimeUnit          time.Duration
	StartPort         uint32
	SimulateDelay     bool
	MaxSimulatedDelay int
	MaxPacketSize     int
	AppTimeout        uint32
	ConnectionTimeout uint32
	TransportTimeout  uint32
	RetransmitLimit   int
	LogLevel          uint8
	MaxBuffer         uint64
}

var config = &ConfigStruct{
	SendInterval:      1500 * time.Microsecond,
	TimeUnit:          10 * time.Microsecond,
	StartPort:         6000,
	SimulateDelay:     false,
	MaxSimulatedDelay: 500,
	MaxPacketSize:     16384,
	AppTimeout:        1000000,
	ConnectionTimeout: 1000000,
	TransportTimeout:  100000,
	RetransmitLimit:   10,
	LogLevel:          DEBUG,
	MaxBuffer:         8192,
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
