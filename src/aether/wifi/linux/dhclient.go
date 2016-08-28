package linux

import (
	"os/exec"
)

// Wrapper for linux utility: dhclient
type DHClient struct{}

func NewDHClient() DHClient {
	return DHClient{}
}

// Checks if the program dhclient exists using PATH environment variable
func (self DHClient) IsInstalled() bool {
	_, err := exec.LookPath("dhclient")
	if err != nil {
		return false
	}
	return true
}

// Start DHCP client and get settings for IP, subnetmask, DNS server, and
// default gateway from the connected Access Point or Router
func (self DHClient) Startup(interfaceName string) error {
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

// Start DHCP client and get settings for IP, subnetmask, DNS server, and
// default gateway from the connected Access Point or Router
// and immediately return
func (self DHClient) StartupFast(interfaceName string) error {
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

// Stop DHCP client and release the current lease
func (self DHClient) Shutdown(interfaceName string) error {
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

// Stop DHCP client and don't bother releasing the current lease
func (self DHClient) ShutdownFast(interfaceName string) error {
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
