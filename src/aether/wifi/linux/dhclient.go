package linux

import (
	"os/exec"
)

// DHClient Wrapper for linux utility: dhclient
type DHClient struct{}

// NewDHClient creates dhclient
func NewDHClient() DHClient {
	return DHClient{}
}

// IsInstalled checks if the program dhclient exists using PATH environment variable
func (dhc DHClient) IsInstalled() bool {
	_, err := exec.LookPath("dhclient")
	if err != nil {
		return false
	}
	return true
}

// Startup starts DHCP client and get settings for IP, subnetmask, DNS server, and
// default gateway from the connected Access Point or Router
func (dhc DHClient) Startup(interfaceName string) error {
	logger.Debug("DHClient: Starting DHCP client and retreiving settings")

	if !authorized() {
		return ErrAuthRequired
	}

	// This can take few seconds or over a minute
	cmd := exec.Command("sudo", "dhclient", "-v", "-nw", interfaceName)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// StartupFast start DHCP client and get settings for IP, subnetmask, DNS server, and
// default gateway from the connected Access Point or Router
// and immediately return
func (dhc DHClient) StartupFast(interfaceName string) error {
	logger.Debug("DHClient: Starting DHCP client (fast) and retreiving settings")

	if !authorized() {
		return ErrAuthRequired
	}

	// -nw is no wait so after we will have to verify our ip config updated
	cmd := exec.Command("sudo", "dhclient", "-v", "-nw", interfaceName)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Shutdown stop DHCP client and release the current lease
func (dhc DHClient) Shutdown(interfaceName string) error {
	logger.Debug("DHClient: Stopping DHCP client and releasing lease")

	if !authorized() {
		return ErrAuthRequired
	}

	// -r shutdown the client and waits until it released the current lease
	cmd := exec.Command("sudo", "dhclient", "-v", "-x", interfaceName)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// ShutdownFast stop DHCP client and don't bother releasing the current lease
func (dhc DHClient) ShutdownFast(interfaceName string) error {
	logger.Debug("DHClient: Stopping DHCP client and releasing lease")

	if !authorized() {
		return ErrAuthRequired
	}

	// -x immediately shutdown client
	cmd := exec.Command("sudo", "dhclient", "-v", "-x", interfaceName)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}
