package linux

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// IWConfig Wrapper for linux utility: iwconfig
type IWConfig struct{}

// NewIWConfig creates iwconfig
func NewIWConfig() IWConfig {
	return IWConfig{}
}

// IWConfigInfo records iwconfig info
type IWConfigInfo struct {
	// wlan0
	InterfaceName string
	// IEEE 802.11bgn
	Protocol string
	// MyHomeWifi
	ESSID string
	// <WIFI@MANUFACTURER>
	NickName string
	// Managed
	Mode string
	// 2.412 Ghz
	Frequency string
	// 00:16:3E:00:30:BF
	AccessPoint string
	// 68 Mb/s
	BitRate string
	// # out of 100
	Sensitivity int
	// off
	Retry string
	// off
	RTSThr string
	// off
	FragementThr string
	// off
	PowerManagement string
	// # out of 100
	LinkQuality int
	// # out of 100
	SignalLevel int
	// # out of 100
	NoiseLevel int
	// 0
	RxInvalidNWID int
	// 0
	RxInvalidCrypt int
	// 0
	RxInvalidFrag int
	// 0
	TxExcessiveRetries int
	// 0
	InvalidMisc int
	// 0
	MissedBeacon int
}

// IsInstalled checks if the program iwconfig exists using PATH environment variable
func (iwc IWConfig) IsInstalled() bool {
	_, err := exec.LookPath("iwconfig")
	if err != nil {
		return false
	}
	return true
}

// Info returns information on a wireless interface
func (iwc IWConfig) Info(interfaceName string) (IWConfigInfo, error) {
	results, err := iwc.InfoList()

	for _, result := range results {
		if result.InterfaceName == interfaceName {
			return result, err
		}
	}

	return IWConfigInfo{}, err
}

// InfoList returns information on all wireless interfaces
func (iwc IWConfig) InfoList() ([]IWConfigInfo, error) {
	logger.Debug("IWConfig: Getting wifi information")

	cmd := exec.Command("iwconfig")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return nil, err
	}
	results, err := iwc.parseInfo(string(out))

	return results, err
}

// Mode sets the interface operating mode. Superuser authentication is required.
func (iwc IWConfig) Mode(interfaceName string, interfaceMode string) error {
	logger.Debug("IWConfig: Setting mode")

	if !authorized() {
		return ErrAuthRequired
	}

	cmd := exec.Command("sudo", "iwconfig", interfaceName, "mode", interfaceMode)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Frequency sets the interface operating frequency. Superuser authentication is required.
func (iwc IWConfig) Frequency(interfaceName string, interfaceFrequency string) error {
	logger.Debug("IWConfig: Setting frequency")

	if !authorized() {
		return ErrAuthRequired
	}

	cmd := exec.Command("sudo", "iwconfig", interfaceName, "freq", interfaceFrequency)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Channel sets the interface operating channel. Superuser authentication is required.
func (iwc IWConfig) Channel(interfaceName string, channel string) error {
	logger.Debug("IWConfig: Setting channel")

	if !authorized() {
		return ErrAuthRequired
	}

	cmd := exec.Command("sudo", "iwconfig", interfaceName, "channel", channel)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Key sets the WEP key used for an access point.
func (iwc IWConfig) Key(interfaceName string, key string) error {
	logger.Debug("IWConfig: Setting Key")

	if !authorized() {
		return ErrAuthRequired
	}

	cmd := exec.Command("sudo", "iwconfig", interfaceName, "key", key)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// ESSID Sets the ESSID of the access point to connect to. Superuser authentication is required.
func (iwc IWConfig) ESSID(interfaceName string, essid string) error {
	logger.Debug("IWConfig: Setting ESSID")

	if !authorized() {
		return ErrAuthRequired
	}

	cmd := exec.Command("sudo", "iwconfig", interfaceName, "essid", essid)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Parse iwconfig information
func (iwc IWConfig) parseInfo(content string) ([]IWConfigInfo, error) {
	iwrs := []IWConfigInfo{}
	iwr := IWConfigInfo{}

	for _, line := range strings.Split(content, "\n") {
		ifaceFirstLine := false

		spaceCols := strings.Split(line, "  ") // two spaces
		for spaceColIndex, spaceCol := range spaceCols {
			spaceCol = strings.TrimSpace(spaceCol)

			if len(spaceCol) == 0 {
				continue
			}

			if strings.Contains(line, "no wireless extensions") {
				continue
			}

			if strings.Contains(line, "no such device") {
				return nil, errors.New("no such device")
			}

			if spaceColIndex == 0 { // interface's first line
				iwr.InterfaceName = spaceCol
				ifaceFirstLine = true
			}

			splitChar := ":"
			if !strings.Contains(spaceCol, ":") {
				splitChar = "="
			}
			fieldSplit := strings.SplitN(spaceCol, splitChar, 2)
			fieldSplitLeft := strings.TrimSpace(fieldSplit[0])
			fieldSplitRight := ""
			if len(fieldSplit) >= 2 {
				fieldSplitRight = strings.TrimSpace(fieldSplit[1])
			}

			if ifaceFirstLine {
				if strings.Contains(fieldSplitLeft, "ESSID") {
					iwr.ESSID = strings.Trim(fieldSplitRight, "\"")
				} else if strings.Contains(fieldSplitLeft, "Nickname") {
					iwr.NickName = strings.Trim(fieldSplitRight, "\"")
				} else {
					iwr.Protocol = spaceCol
				}
			} else {
				switch fieldSplitLeft {
				case "Mode":
					iwr.Mode = fieldSplitRight
				case "Frequency":
					iwr.Frequency = fieldSplitRight
				case "Access Point":
					iwr.AccessPoint = fieldSplitRight
				case "Bit Rate":
					iwr.BitRate = fieldSplitRight
				case "Sensitivity":
					iwr.Sensitivity, _ = strconv.Atoi(fieldSplitRight)
				case "Retry":
					iwr.Retry = fieldSplitRight
				case "RTS thr":
					iwr.RTSThr = fieldSplitRight
				case "Fragment thr":
					iwr.FragementThr = fieldSplitRight
				case "Power Management":
					iwr.PowerManagement = fieldSplitRight
				case "Link Quality":
					iwr.LinkQuality, _ = strconv.Atoi(strings.SplitN(fieldSplitRight,
						"/", 2)[0])
				case "Signal level":
					iwr.SignalLevel, _ = strconv.Atoi(strings.SplitN(fieldSplitRight,
						"/", 2)[0])
				case "Noise level":
					iwr.NoiseLevel, _ = strconv.Atoi(strings.SplitN(fieldSplitRight,
						"/", 2)[0])
				case "Rx invalid nwid":
					iwr.RxInvalidNWID, _ = strconv.Atoi(fieldSplitRight)
				case "Rx invalid crypt":
					iwr.RxInvalidCrypt, _ = strconv.Atoi(fieldSplitRight)
				case "Rx invalid frag":
					iwr.RxInvalidFrag, _ = strconv.Atoi(fieldSplitRight)
				case "Tx excessive retries":
					iwr.TxExcessiveRetries, _ = strconv.Atoi(fieldSplitRight)
				case "Invalid misc":
					iwr.InvalidMisc, _ = strconv.Atoi(fieldSplitRight)
				case "Missed beacon":
					iwr.MissedBeacon, _ = strconv.Atoi(fieldSplitRight)

					// last thing before new record
					iwrs = append(iwrs, iwr)
					iwr = IWConfigInfo{}
				}
			}
		}

		if ifaceFirstLine {
			ifaceFirstLine = false
		}
	}

	return iwrs, nil
}
