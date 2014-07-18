package linux

import (
	"errors"
	"os/exec"
	"strings"
	"time"
)

// Wrapper for linux utility: nmcli (network-manager)
// If networkmanager is installed and running, it blocks iwconfig
// and possibly more, so in those cases use networkmanager.
//
// nmcli does not need sudo (except the service and main NetworkManager tool)
// so if no sudo permissions we can default to use nmcli.
//
// The connection profiles are stored in:
// /etc/NetworkManager/system-connections/*  with owner set to root
//
// Also see:
//	NetworkManager (case sensitive, requires sudo)
//	/etc/NetworkManager/NetworkManager.conf
//	nmcli
//	nm-tool
//	nm-applet
//	nm-online
//	nm-settings
//	nm-connection-editor
//
type NetworkManager struct{}

func NewNetworkManager() NetworkManager {
	return NetworkManager{}
}

// Checks if the program nmcli exists using PATH environment variable
func (self NetworkManager) IsInstalled() bool {
	_, err := exec.LookPath("nmcli")
	if err != nil {
		return false
	}
	return true
}

func (self NetworkManager) NetworkManagerID(interfaceName string) string {
	return "darknet_" + interfaceName
}

// Connect to an access point
// only certain versions of NetworkManager support this...
func (self NetworkManager) Connect(interfaceName string, ssid string,
	secProtocol string, secKey string) error {
	logger.Debug("NetworkManager: Connecting to an access point")

	result, errp := self.ProfileExistsById(interfaceName)
	if result == true && errp == nil {
		self.DeleteById(interfaceName)
	}

	var cmd *exec.Cmd
	if secProtocol == "none" {
		cmd = exec.Command("nmcli", "dev", "wifi", "con", ssid,
			"name", self.NetworkManagerID(interfaceName),
			"iface", interfaceName)
	} else {
		cmd = exec.Command("nmcli", "dev", "wifi", "con", ssid,
			"name", self.NetworkManagerID(interfaceName),
			"password", secKey, "iface", interfaceName)
	}
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Delete a connection profile by NetworkManager id
func (self NetworkManager) ProfileExistsById(interfaceName string) (bool, error) {
	logger.Debug("NetworkManager: Check profile exists by nm id")

	networkmanagerID := self.NetworkManagerID(interfaceName)

	// nmcli dev disconnect iface wlan0
	cmd := exec.Command("nmcli", "-m", "multiline", "con", "list")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return false, err
	}
	logger.Debug("Command Return: %v", limitText(out))

	nmID, _ := self.parseProfiles(string(out))
	if networkmanagerID == nmID {
		return true, nil
	}

	return false, nil
}

func (self NetworkManager) parseProfiles(text string) (string, string) {
	nmID := ""
	nmUUID := ""
	for _, line := range strings.Split(text, "\n") {
		fs := strings.SplitN(line, " ", 2)
		fsLeft := strings.TrimSpace(fs[0])
		fsRight := ""
		if len(fs) >= 2 {
			fsRight = strings.TrimSpace(fs[1])
		}
		if fsLeft == "NAME:" {
			nmID = fsRight
		}
		if fsLeft == "UUID:" {
			nmUUID = fsRight
		}
	}
	return nmID, nmUUID
}

// Delete a connection profile by NetworkManager id
func (self NetworkManager) DeleteById(interfaceName string) {
	logger.Debug("NetworkManager: Delete a connection profile by nm id")

	networkmanagerID := self.NetworkManagerID(interfaceName)

	// nmcli dev disconnect iface wlan0
	cmd := exec.Command("nmcli", "con", "delete", "id", networkmanagerID)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
	}
	logger.Debug("Command Return: %v", limitText(out))
}

// Deactivate a connection on an interface by NetworkManager id
func (self NetworkManager) DeactivateById(interfaceName string) error {
	logger.Debug("NetworkManager: Deactivating an interface by nm id")

	networkmanagerID := self.NetworkManagerID(interfaceName)

	// nmcli dev disconnect iface wlan0
	cmd := exec.Command("nmcli", "con", "down", "id", networkmanagerID)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Disconnect all connections on an interface
func (self NetworkManager) DisconnectAll(interfaceName string) error {
	logger.Debug("NetworkManager: Disconnecting all connections on an interface")

	// nmcli dev disconnect iface wlan0
	cmd := exec.Command("nmcli", "dev", "disconnect", "iface", interfaceName)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Checks if network-manager service is running
func (self NetworkManager) ServiceIsRunning() (bool, error) {
	logger.Debug("NetworkManager: Checking if service running")

	if !self.IsInstalled() {
		return false, errors.New("service not installed")
	}

	if !authorized() {
		return false, ErrAuthRequired
	}

	// you must type sudo for services even if superuser
	cmd := exec.Command("sudo", "service", "network-manager", "status")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return false, err
	}
	logger.Debug("Command Return: %v", limitText(out))

	if strings.Contains(string(out), ",") {
		return true, nil
	}
	return false, nil
}

// Stop the network-manager service. Superuser authentication is required.
func (self NetworkManager) ServiceStop() error {
	logger.Debug("NetworkManager: Service stopping")

	if result, _ := self.ServiceIsRunning(); result == false {
		logger.Debug("NetworkManager: Service already stopped")
		return nil
	}

	if !authorized() {
		return ErrAuthRequired
	}

	// you must type sudo for services even if superuser
	cmd := exec.Command("sudo", "service", "network-manager", "stop")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	if result, _ := self.ServiceIsRunning(); result == true {
		logger.Debug("NetworkManager: Service failed to stop")
		return errors.New("service failed to stop")
	} else {
		logger.Debug("NetworkManager: Service stopped successfully")
	}

	return nil
}

// Start the network-manager service. Superuser authentication is required.
func (self NetworkManager) ServiceStart() error {
	logger.Debug("NetworkManager: Service starting")

	if result, _ := self.ServiceIsRunning(); result == true {
		logger.Debug("NetworkManager: Service already started")
		return nil
	}

	if !authorized() {
		return ErrAuthRequired
	}

	// you must type sudo for services even if superuser
	cmd := exec.Command("sudo", "service", "network-manager", "start")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
	}
	logger.Debug("Command Return: %v", limitText(out))

	// Wait a few seconds for it to really load
	time.Sleep(3 * time.Second)

	if result, _ := self.ServiceIsRunning(); result == false {
		logger.Debug("NetworkManager: Service failed to start")
		return errors.New("service failed to start")
	} else {
		logger.Debug("NetworkManager: Service started successfully")
	}

	return nil
}
