package cli

import (
    "../util/"
    "flag"
    "github.com/op/go-logging"
    "log"
    "strings"
)

func parseLogLevel(level string) logging.Level {
    switch strings.ToLower(level) {
    case "critical":
        return logging.CRITICAL
    case "error":
        return logging.ERROR
    case "warning":
        return logging.WARNING
    case "notice":
        return logging.NOTICE
    case "info":
        return logging.INFO
    case "debug":
        return logging.DEBUG
    default:
        log.Fatal("Unknown log level: %s\n", level)
        return logging.INFO
    }
    return logging.INFO
}

func ParseArgs() {
    RegisterArgs()
    flag.Parse()
    DataDirectory = util.InitDataDir(DataDirectory)
    LogLevel = parseLogLevel(logLevel)
}
