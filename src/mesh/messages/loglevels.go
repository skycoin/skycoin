package messages

const (
	DEBUG = iota
	INFO
)

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
