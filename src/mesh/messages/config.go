package messages

const (
	DEBUG = iota
	INFO
)

type ConfigStruct struct {
	StartPort         uint32
	SimulateDelay     bool
	MaxSimulatedDelay int
	MaxPacketSize     int
	ConnectionTimeout uint32
	TransportTimeout  uint32
	RetransmitLimit   int
	LogLevel          uint8
}

var config = &ConfigStruct{
	StartPort:         6000,
	SimulateDelay:     false,
	MaxSimulatedDelay: 500,
	MaxPacketSize:     512,
	ConnectionTimeout: 100000,
	TransportTimeout:  10000,
	RetransmitLimit:   10,
	LogLevel:          DEBUG,
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
