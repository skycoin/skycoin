package messages

const (
	DEBUG = iota
	INFO
)

type ConfigStruct struct {
	MaxSimulatedDelay int
	ConnectionTimeout uint32
	TransportTimeout  uint32
	RetransmitLimit   int
	LogLevel          uint8
}

var config = &ConfigStruct{
	MaxSimulatedDelay: 500,
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
