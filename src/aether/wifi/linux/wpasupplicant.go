package linux

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// WPASupplicant Wrapper for linux utility: wpa_supplicant
//
// Areas for wpa_supplicant:
//	/etc/wpa_supplicant/wpa_supplicant.conf
//	Config file with network connection info and authentication keys
//
//	wpa_supplicant
//	main program that communicates with interface
//
//	wpa_cli
//	high-level wrapper around wpa_supplicant for easier use
//
//	wpa_passphrase
//	helps create wpa_supplicant.conf with encrypted keys
//
type WPASupplicant struct{}

// NewWPASupplicant creates instance
func NewWPASupplicant() WPASupplicant {
	return WPASupplicant{}
}

// IsInstalled checks if the program wpa_supplicant exists using PATH environment variable
func (ws WPASupplicant) IsInstalled() bool {
	_, err := exec.LookPath("wpa_supplicant")
	if err != nil {
		return false
	}
	return true
}

// ConfigName returns the config file name
func (ws WPASupplicant) ConfigName() string {
	return os.TempDir() + "/" + "darknet_wpa_supplicant.conf"
}

// ConfigWrite config write
func (ws WPASupplicant) ConfigWrite(interfaceName string, ssid string,
	password string) error {
	logger.Debug("WPASupplicant: Writing config file %v", ws.ConfigName())

	if !authorized() {
		return ErrAuthRequired
	}

	text := ""
	text += "ctrl_interface=/var/run/wpa_supplicant" + "\n"
	text += "ctrl_interface_group=0" + "\n"
	text += "eapol_version=1" + "\n"
	text += "ap_scan=1" + "\n"
	text += "fast_reauth=1" + "\n"
	text += "" + "\n"

	authKey, errA := ws.ConfigGenerateAuth(ssid, password)
	if errA == nil {
		text += "network={" + "\n"
		text += "     ssid=\"" + ssid + "\"" + "\n"
		text += "     scan_ssid=1" + "\n"
		text += "     key_mgmt=WPA-PSK" + "\n"
		//text += "     pairwise=CCMP TKIP" + "\n"
		//text += "     group=CCMP TKIP" + "\n"
		//text += "     proto=RSN" + "\n"
		text += "     psk=" + authKey + "\n"
		text += "}" + "\n"
	}

	err := ioutil.WriteFile(ws.ConfigName(), []byte(text), 0644)
	if err != nil {
		logger.Error("WPASupplicant: Writing config file failed.")
		return err
	}

	return nil
}

// ConfigRemove removes config
func (ws WPASupplicant) ConfigRemove() {
	logger.Debug("WPASupplicant: Removing config file %v", ws.ConfigName())
	os.Remove(ws.ConfigName())
}

// ConfigGenerateAuth generate auth
func (ws WPASupplicant) ConfigGenerateAuth(ssid string,
	password string) (string, error) {
	logger.Debug("WPASupplicant: Config Generate Auth")

	cmd := exec.Command("wpa_passphrase", ssid, password)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return "", err
	}
	logger.Debug("Command Return: %v", limitText(out))

	for _, line := range strings.Split(string(out), "\n") {
		fieldSplit := strings.SplitN(line, "=", 2)
		fieldSplitLeft := strings.TrimSpace(fieldSplit[0])
		fieldSplitRight := ""
		if len(fieldSplit) >= 2 {
			fieldSplitRight = strings.TrimSpace(fieldSplit[1])
		}
		if strings.Contains(fieldSplitLeft, "psk") &&
			!strings.Contains(fieldSplitLeft, "#psk") {
			return fieldSplitRight, nil
		}
	}

	return "", errors.New("generate auth key failed")
}

// DaemonStartup start service
func (ws WPASupplicant) DaemonStartup(interfaceName string) error {
	logger.Debug("WPASupplicant: Service starting")

	if !authorized() {
		return ErrAuthRequired
	}

	running, runerr := ws.DaemonIsRunning()
	if running && runerr == nil {
		logger.Debug("WPASupplicant: Service already started")
		ws.DaemonConfigReload()
		return nil
	}

	// wpa_supplicant -B -iwlan0 -Dwext -c/tmp/darknet_wpa_supplicant.conf
	// -B is run daemon in background, -d is increase debugging verbosity
	pInterface := "-i" + interfaceName

	// use wext driver, its like catch all
	pDriver := "-D" + "wext"
	//pDriver := "-D" + "nl80211"

	pConfigPath := "-c" + ws.ConfigName()

	cmd := exec.Command("sudo", "wpa_supplicant", "-B", "-d",
		pInterface, pDriver, pConfigPath)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// DaemonShutdown shutdowns the service
func (ws WPASupplicant) DaemonShutdown() error {
	logger.Debug("WPASupplicant: Daemon Shutdown")

	if !authorized() {
		return ErrAuthRequired
	}

	cmd := exec.Command("sudo", "wpa_cli", "terminate")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// DaemonConfigReload reload daemon config
func (ws WPASupplicant) DaemonConfigReload() error {
	logger.Debug("WPASupplicant: Daemon Config Reload")

	if !authorized() {
		return ErrAuthRequired
	}

	cmd := exec.Command("sudo", "wpa_cli", "reconfigure")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// DaemonIsRunning check if running
func (ws WPASupplicant) DaemonIsRunning() (bool, error) {
	logger.Debug("WPASupplicant: Checking if running")

	if !authorized() {
		return false, ErrAuthRequired
	}

	cmd := exec.Command("sudo", "wpa_cli", "status")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return false, err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return true, nil
}

// Authenticate authenticates
func (ws WPASupplicant) Authenticate(interfaceName string,
	ssid string, password string) error {
	logger.Debug("WPASupplicant: Authenticate")

	if !authorized() {
		return ErrAuthRequired
	}

	// sudo wpa_cli identity "" otp "" (otp is one time password)
	cmd := exec.Command("sudo", "wpa_cli", "interface", interfaceName,
		"identity", ssid, "otp", password)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}
