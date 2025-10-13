package daemon

import (
	"errors"
	"log"
	"strings"

	"github.com/skycoin/hardware-wallet-daemon/src/api"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"

	"github.com/skycoin/skycoin/src/util/file"
)

// Config records the daemon and build configuration
type Config struct {
	App   AppConfig
	Build api.BuildInfo
}

// AppConfig records the app's configuration
type AppConfig struct {
	// Remote web interface port
	WebInterfacePort int
	// Remote web interface address
	WebInterfaceAddr string

	// Enable CSRF check
	EnableCSRF bool

	// Disable Host, Origin and Referer header check in the wallet API
	DisableHeaderCheck bool
	// Comma separate list of hostnames to accept in the Host header, used to bypass the Host header check which only applies to localhost addresses
	HostWhitelist string
	hostWhitelist []string

	// Logging
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	LogLevel string
	// Enable logging to file
	LogToFile bool

	// Enable cpu profiling
	ProfileCPU bool
	// Where the file is written to
	ProfileCPUFile string
	// Enable HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool
	// Expose HTTP profiling on this interface
	HTTPProfHost string

	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory string

	// DaemonMode decides with what api is enabled, either wallet or emulator
	DaemonMode string
	daemonMode skyWallet.DeviceType
}

// NewAppConfig returns a new app config instance
func NewAppConfig(port int, datadir string) AppConfig {
	return AppConfig{
		WebInterfaceAddr: "127.0.0.1",
		WebInterfacePort: port,

		// Logging
		ColorLog:  true,
		LogLevel:  "INFO",
		LogToFile: false,

		// disable csrf by default
		EnableCSRF: false,

		// Enable cpu profiling
		ProfileCPU: false,
		// Where the file is written to
		ProfileCPUFile: "cpu.prof",
		// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
		HTTPProf:     false,
		HTTPProfHost: "localhost:6060",

		// Run daemon in wallet mode by default
		DaemonMode: skyWallet.DeviceTypeUSB.String(),

		DataDirectory: datadir,
	}
}

func (c *Config) postProcess() error {
	var err error
	home := file.UserHome()
	c.App.DataDirectory, err = file.InitDataDir(replaceHome(c.App.DataDirectory, home))
	panicIfError(err, "Invalid DataDirectory")

	if c.App.HostWhitelist != "" {
		if c.App.DisableHeaderCheck {
			return errors.New("host whitelist should be empty when header check is disabled")
		}
		c.App.hostWhitelist = strings.Split(c.App.HostWhitelist, ",")
	}

	c.App.daemonMode = skyWallet.DeviceTypeFromString(c.App.DaemonMode)
	if c.App.daemonMode == skyWallet.DeviceTypeInvalid {
		return errors.New("invalid device type")
	}

	return nil
}

func panicIfError(err error, msg string, args ...interface{}) { // nolint: unparam
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func replaceHome(path, home string) string {
	return strings.Replace(path, "$HOME", home, 1)
}
