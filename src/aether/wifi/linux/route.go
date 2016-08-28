package linux

import (
	"os/exec"
)

// Wrapper for linux utility: route
//
// To verify, can type "route -n" to see all current routes
//
type Route struct{}

func NewRoute() Route {
	return Route{}
}

// Checks if the program route exists using PATH environment variable
func (self Route) IsInstalled() bool {
	_, err := exec.LookPath("route")
	if err != nil {
		return false
	}
	return true
}

// Sets the system to add default gateway. Typically the IP of the Access Point
func (self Route) AddDefaultGateway(ipAddress string) error {
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
