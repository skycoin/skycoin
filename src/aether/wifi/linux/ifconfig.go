package linux

import (
	"net"
	"os/exec"
)

// IFConfig Wrapper for linux utility: ifconfig
type IFConfig struct{}

// NewIFConfig create ifconfig
func NewIFConfig() IFConfig {
	return IFConfig{}
}

/*
type IFConfigInfo struct {
	// wlan0
	InterfaceName string
	// Local Loopback|Ethernet
	LinkEncap string
	// 00:16:3E:E1:A2:A3
	HWAddr string
	// 127.0.0.1|192.168.1.123
	INetAddr string
	// 192.168.1.255
	Bcast string
	// 255.0.0.0|255.255.255.0
	Mask string
	// ::1/128|fd0c:172b:6750:0fae::/64
	INet6Addr string
	// Host|Link
	Scope string
	// 65536|1500
	MTU int
	// 1
	Metric int
	// #
	RxPackets int
	// #
	RxPacketsErrors int
	// #
	RxPacketsDropped int
	// #
	RxPacketsOverruns int
	// #
	RxPacketsFrame int
	// #
	TxPackets int
	// #
	TxPacketsErrors int
	// #
	TxPacketsDropped int
	// #
	TxPacketsOverruns int
	// #
	TxPacketsCarrier int
	// #
	Collisions int
	// #
	TxQueueLen int
	// #
	RxBytes int
	// #
	TxBytes int
}


func (self IFConfig) Info(interfaceName string) (IFConfigInfo, error) {
	ifconfigInfoResults, err := self.InfoList()

	for _, ifconfigInfoResult := range ifconfigInfoResults {
		if ifconfigInfoResult.InterfaceName == interfaceName {
			return ifconfigInfoResult, err
		}
	}

	return IFConfigInfo{}, err
}

func (self IFConfig) InfoList() ([]IFConfigInfo, error) {
	cmd := exec.Command("ifconfig")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%v", err)
	}
	ifconfigInfoResults, err := self.parse(string(out))

	return ifconfigInfoResults, err
}

func (self IFConfig) parse(content string) ([]IFConfigInfo, error) {
	ifrs := []IFConfigInfo{}
	ifr := IFConfigInfo{}

	for _, line := range strings.Split(content, "\n") {
		ifaceLine := 0

		spaceCols := strings.Split(line, "  ") // two spaces
		for spaceColIndex, spaceCol := range spaceCols {
			spaceCol = strings.TrimSpace(spaceCol)

			if len(spaceCol) == 0 {
				continue
			}

			if strings.Contains(line, "error fetching interface information") {
				return nil, errors.New("error fetching interface information")
			}

			if spaceColIndex == 0 { // interface's first line
				ifr.InterfaceName = spaceCol
				ifaceLine = 1
			}

			fs := strings.SplitN(spaceCol, ":", 2)
			fsL := strings.TrimSpace(fs[0])
			fsR := ""
			if len(fs) >= 2 {
				fsr = strings.TrimSpace(fs[1])
			}

			if ifaceLine == 1 {
				if fsl == "Link encap" {
					ifr.LinkEncap = fsr
				} else if strings.Contains(fsl, "Nickname") {
					ifr.NickName = fsr
				} else {
					ifr.Protocol = spaceCol
				}
			} else {
				switch fsl {
				case "Mode":
					ifr.Mode = fsr
				case "Frequency":
					ifr.Frequency = fsr
				case "Access Point":
					ifr.AccessPoint = fsr
				case "Bit Rate":
					ifr.BitRate = fsr
				case "Sensitivity":
					ifr.Sensitivity, _ = strconv.Atoi(fsr)
				case "Retry":
					ifr.Retry = fsr
				case "RTS thr":
					ifr.RTSThr = fsr
				case "Fragment thr":
					ifr.FragementThr = fsr
				case "Power Management":
					ifr.PowerManagement = fsr
				case "Link Quality":
					ifr.LinkQuality, _ = strconv.Atoi(fsr)
				case "Signal level":
					ifr.SignalLevel, _ = strconv.Atoi(fsr)
				case "Noise level":
					ifr.NoiseLevel, _ = strconv.Atoi(fsr)
				case "Rx invalid nwid":
					ifr.RxInvalidNWID, _ = strconv.Atoi(fsr)
				case "Rx invalid crypt":
					ifr.RxInvalidCrypt, _ = strconv.Atoi(fsr)
				case "Rx invalid frag":
					ifr.RxInvalidFrag, _ = strconv.Atoi(fsr)
				case "Tx excessive retries":
					ifr.TxExcessiveRetries, _ = strconv.Atoi(fsr)
				case "Invalid misc":
					ifr.InvalidMisc, _ = strconv.Atoi(fsr)
				case "Missed beacon":
					ifr.MissedBeacon, _ = strconv.Atoi(fsr)

					// last thing before new record
					ifrs = append(ifrs, ifr)
					ifr = IFConfigInfo{}
				}
			}
		}

		if ifaceFirstLine {
			ifaceFirstLine = false
		}
	}

	return ifrs, nil
}
*/

// IsInstalled checks if the program ifconfig exists using PATH environment variable
func (ifc IFConfig) IsInstalled() bool {
	_, err := exec.LookPath("ifconfig")
	if err != nil {
		return false
	}
	return true
}

// Up brings the interface Up. Superuser authentication is required.
func (ifc IFConfig) Up(interfaceName string) error {
	logger.Debug("IFConfig: Bringing interface up")

	if !authorized() {
		return ErrAuthRequired
	}

	cmd := exec.Command("sudo", "ifconfig", interfaceName, "up")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Down brings the interface Down. Superuser authentication is required.
func (ifc IFConfig) Down(interfaceName string) error {
	logger.Debug("IFConfig: Bringing interface down")

	if !authorized() {
		return ErrAuthRequired
	}

	cmd := exec.Command("sudo", "ifconfig", interfaceName, "down")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// SetIP sets the IP for the interface. Superuser authentication is required.
func (ifc IFConfig) SetIP(interfaceName string,
	ipAddress net.IP, netMask net.IPMask) error {
	logger.Debug("IFConfig: Setting the IP and Netmask")

	if !authorized() {
		return ErrAuthRequired
	}

	cmd := exec.Command("sudo", "ifconfig", interfaceName, ipAddress.String(),
		"netmask", netMask.String())
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}
