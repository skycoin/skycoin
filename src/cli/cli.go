package cli

import (
    "flag"
    "log"

    "../util/"
    "github.com/op/go-logging"
)

func ParseArgs() {
    RegisterArgs()
    flag.Parse()
    DataDirectory = util.InitDataDir(DataDirectory)
    ll, err := logging.LogLevel(logLevel)
    if err != nil {
        log.Panic("Invalid -log-level %s: %v\n", logLevel, err)
    }
    LogLevel = ll
}
