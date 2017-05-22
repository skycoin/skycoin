package linux

import (
	"os/exec"
	"strconv"
	"strings"
)

// IWList Wrapper for linux utility: iwlist
type IWList struct{}

// NewIWList creates iwlist
func NewIWList() IWList {
	return IWList{}
}

// IWListNetwork records the iwlist network info
type IWListNetwork struct {
	Address             string
	ESSID               string
	Protocol            string
	Mode                string
	Frequency           string
	EncryptionKeyStatus string
	BitRates            string
	InformationElement  []IWListNetworkIE
	QualityLevel        int
	SignalLevel         int
	NoiseLevel          int
}

// IWListNetworkIE records iwlist networkie
type IWListNetworkIE struct {
	Protocol        string
	GroupCipher     string
	PairwiseCiphers string
	AuthSuites      string
	Extra           string
}

// IsInstalled checks if the program iwlist exists using PATH environment variable
func (iwl IWList) IsInstalled() bool {
	_, err := exec.LookPath("iwlist")
	if err != nil {
		return false
	}
	return true
}

// ScanCached returns a cached list of wireless networks found with an interface.
// No Superuser authentication is required.
func (iwl IWList) ScanCached(interfaceName string) ([]IWListNetwork, error) {
	logger.Debug("IWList: Scanning for wifi networks (cached)")

	cmd := exec.Command("iwlist", interfaceName, "scan")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return nil, err
	}
	logger.Debug("Command Return: %v", limitText(out))
	WifiNetworks := iwl.parse(string(out))

	return WifiNetworks, nil
}

// Scan returns a fresh list of wireless networks found with an interface.
// Superuser authentication is required.
func (iwl IWList) Scan(interfaceName string) ([]IWListNetwork, error) {
	logger.Debug("IWList: Scanning for wifi networks (live)")

	if !authorized() {
		return nil, ErrAuthRequired
	}

	cmd := exec.Command("sudo", "iwlist", interfaceName, "scan")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return nil, err
	}
	logger.Debug("Command Return: %v", limitText(out))
	WifiNetworks := iwl.parse(string(out))

	return WifiNetworks, err
}

func (iwl IWList) parse(content string) []IWListNetwork {
	wfns := []IWListNetwork{}
	wfn := IWListNetwork{}

	ies := []IWListNetworkIE{}
	ie := IWListNetworkIE{}

	for lineIndex, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)

		splitChar := ":"
		if !strings.Contains(line, ":") {
			splitChar = "="
		}
		cols := strings.SplitN(line, splitChar, 2)
		if len(cols) < 2 {
			continue
		}
		colLeft := strings.TrimSpace(cols[0])
		colRight := strings.TrimSpace(cols[1])

		switch {
		case strings.Contains(colLeft, "Cell"):
			if lineIndex > 1 {
				wfns = append(wfns, wfn)
				wfn = IWListNetwork{}
			}

			wfn.Address = colRight
		case colLeft == "ESSID":
			wfn.ESSID = strings.Trim(colRight, "\"")
		case colLeft == "Protocol":
			wfn.Protocol = colRight
		case colLeft == "Mode":
			wfn.Mode = colRight
		case colLeft == "Frequency":
			wfn.Frequency = colRight
		case colLeft == "Encryption key":
			wfn.EncryptionKeyStatus = colRight
		case colLeft == "Bit Rates":
			wfn.BitRates = colRight
		case colLeft == "IE":
			if len(ies) == 1 {
				ies = append(ies, ie)
				ie = IWListNetworkIE{}
			}
			ie.Protocol = colRight
		case strings.Contains(colLeft, "Group Cipher"):
			ie.GroupCipher = colRight
		case strings.Contains(colLeft, "Pairwise Ciphers"):
			ie.PairwiseCiphers = colRight
		case strings.Contains(colLeft, "Authentication Suites"):
			ie.AuthSuites = colRight
		case colLeft == "Extra":
			ie.Extra = colRight
		case colLeft == "Quality":
			if len(ies) == 0 {
				ies = append(ies, ie)
			}
			wfn.InformationElement = ies

			spaceCols := strings.Split(line, "  ") // two spaces
			for _, spaceCol := range spaceCols {
				spaceCol = strings.TrimSpace(spaceCol)

				spaceSplitChar := ":"
				if !strings.Contains(line, ":") {
					spaceSplitChar = "="
				}
				signalCols := strings.SplitN(spaceCol, spaceSplitChar, 2)
				if len(signalCols) < 2 {
					continue
				}
				signalColLeft := strings.TrimSpace(signalCols[0])
				signalColRight := strings.TrimSpace(signalCols[1])
				signalColRight = strings.SplitN(signalColRight, "/", 2)[0]
				switch signalColLeft {
				case "Quality":
					wfn.QualityLevel, _ = strconv.Atoi(signalColRight)
				case "Signal level":
					wfn.SignalLevel, _ = strconv.Atoi(signalColRight)
				case "Noise level":
					wfn.NoiseLevel, _ = strconv.Atoi(signalColRight)
				}
			}
		}
	}

	wfns = append(wfns, wfn)

	return wfns
}
