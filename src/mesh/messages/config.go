package messages

const (
	DEBUG = iota
	INFO
)

type ConfigStruct struct {
	StartPort         uint32
	MaxSimulatedDelay int
	MaxPacketSize     int
	ConnectionTimeout uint32
	TransportTimeout  uint32
	RetransmitLimit   int
	LogLevel          uint8
}

var config = &ConfigStruct{
	StartPort:         6000,
	MaxSimulatedDelay: 500,
	MaxPacketSize:     256,
	ConnectionTimeout: 10000,
	TransportTimeout:  1000,
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
