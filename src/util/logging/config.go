package logging

// Configure package loggers
func ConfigPkgLogging(logger *MasterLogger, config []PkgLogConfig) {
	for _, pkgcfg := range config {
		logger.PkgConfig[pkgcfg.PkgName] = pkgcfg
	}
}
