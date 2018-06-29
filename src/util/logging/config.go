package logging

import (
	"errors"
	"os"
	"strings"
)

const (
	// PkgLogLevelVar OS env var for package-specific logging levels
	PkgLogLevelVar = "PKGLOGLEVEL"
)

func parsePkgLogLevelsFromEnv() ([]PkgLogConfig, error) {
	configStr := os.Getenv(PkgLogLevelVar)
	return ParsePackageLevels(configStr)
}

var (
	// ErrInvalidValue error representing invalid config value
	ErrInvalidValue = errors.New("Invalid config value")
)

// ParsePackageLevels parse a sequence of the form 'name1=INFO:name2=DEGUB:...'
func ParsePackageLevels(configStr string) ([]PkgLogConfig, error) {
	items := strings.Split(configStr, ":")
	configMap := make([]PkgLogConfig, len(items))
	for index, step := range items {
		item := strings.SplitN(step, "=", 2)
		if len(item) < 2 {
			return nil, ErrInvalidValue
		}
		level, _err := LevelFromString(item[1])
		if _err != nil {
			return nil, ErrInvalidValue
		}
		configMap[index] = PkgLogConfig{item[0], level}
	}
	return configMap, nil
}

// ConfigPkgLogging Configure package loggers
func ConfigPkgLogging(logger *MasterLogger, config []PkgLogConfig) {
	for _, pkgcfg := range config {
		logger.PkgConfig[pkgcfg.PkgName] = pkgcfg
	}
}
