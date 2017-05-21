package linux

import (
	"os/exec"
)

// Route Wrapper for linux utility: route
//
// To verify, can type "route -n" to see all current routes
//
type Route struct{}

// NewRoute create route instance
func NewRoute() Route {
	return Route{}
}

// IsInstalled checks if the program route exists using PATH environment variable
func (rt Route) IsInstalled() bool {
	_, err := exec.LookPath("route")
	if err != nil {
		return false
	}
	return true
}

// AddDefaultGateway sets the system to add default gateway. Typically the IP of the Access Point
func (rt Route) AddDefaultGateway(ipAddress string) error {
	logger.Debug("Route: Adding default gateway")

	cmd := exec.Command("route", "add", "default", "gw", ipAddress)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}
