package linux

import (
	"os/exec"
	"strings"
)

// Wrapper for linux utility: udevadm
type UDevAdm struct{}

func NewUDevAdm() UDevAdm {
	return UDevAdm{}
}

type UDevAdmInfo struct {
	DevPath                   string // /devices/pci0000:00/...
	DevType                   string // wlan
	IdBus                     string // usb
	IdMMCandidate             string // 1
	IdModel                   string // Model (useful)
	IdModelEnc                string // Model Encoded
	IdModelId                 string // Model Id
	IdNetNameMac              string //
	IdNetNamePath             string //
	IdOUIFromDatabase         string // Corporation Name
	IdRevision                string //
	IdSerial                  string //
	IdSerialShort             string //
	IdType                    string // generic
	IdUSBClassFromDatabase    string // Miscellaneous Device
	IdUSBDriver               string // Driver Name (useful)
	IdUSBInterfaces           string //
	IdUSBInterfaceNum         string //
	IdUSBProtocolFromDatabase string // Interface Association
	IdVendor                  string // Manufacturer name (useful)
	IdVendorEnc               string // Manufacturer name encoded
	IdVendorFromDatabase      string // Manufacturer name long
	IdVendorID                string //
	IfIndex                   string
	InterfaceName             string // wlan0
	Subsystem                 string // net
	USecInitialized           string //
}

// Checks if the program udevadm exists using PATH environment variable
func (self UDevAdm) IsInstalled() bool {
	_, err := exec.LookPath("udevadm")
	if err != nil {
		return false
	}
	return true
}

// Return information on an interface
func (self UDevAdm) Run(interfaceName string) (UDevAdmInfo, error) {
	logger.Debug("UDevAdm: Getting device information")

	cmd := exec.Command("udevadm", "info", "-p", "/sys/class/net/"+interfaceName)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return UDevAdmInfo{}, err
	}
	logger.Debug("Command Return: %v", limitText(out))

	content := string(out)
	uInf, err := self.udevadmParse(content)

	return uInf, err
}

func (self UDevAdm) udevadmParse(content string) (UDevAdmInfo, error) {
	uInf := UDevAdmInfo{}

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)

		cols := strings.SplitN(line, "=", 2)
		if len(cols) < 2 {
			continue
		}
		colLeft := strings.TrimSpace(cols[0])
		colLeftSplit := strings.SplitN(colLeft, " ", 2)
		if len(colLeftSplit) < 2 {
			continue
		}
		colLeft = colLeftSplit[1]
		colRight := strings.TrimSpace(cols[1])

		switch {
		case colLeft == "DEVPATH":
			uInf.DevPath = colRight
		case colLeft == "DEVTYPE":
			uInf.DevType = colRight
		case colLeft == "ID_BUS":
			uInf.IdBus = colRight
		case colLeft == "ID_MM_CANDIDATE":
			uInf.IdMMCandidate = colRight
		case colLeft == "ID_MODEL":
			uInf.IdModel = colRight
		case colLeft == "ID_MODEL_ENC":
			uInf.IdModelEnc = colRight
		case colLeft == "ID_MODEL_ID":
			uInf.IdModelId = colRight
		case colLeft == "ID_NET_NAME_MAC":
			uInf.IdNetNameMac = colRight
		case colLeft == "ID_NET_NAME_PATH":
			uInf.IdNetNamePath = colRight
		case colLeft == "ID_OUI_FROM_DATABASE":
			uInf.IdOUIFromDatabase = colRight
		case colLeft == "ID_REVISION":
			uInf.IdRevision = colRight
		case colLeft == "ID_SERIAL":
			uInf.IdSerial = colRight
		case colLeft == "ID_SERIAL_SHORT":
			uInf.IdSerialShort = colRight
		case colLeft == "ID_TYPE":
			uInf.IdType = colRight
		case colLeft == "ID_USB_CLASS_FROM_DATABASE":
			uInf.IdUSBClassFromDatabase = colRight
		case colLeft == "ID_USB_DRIVER":
			uInf.IdUSBDriver = colRight
		case colLeft == "ID_USB_INTERFACES":
			uInf.IdUSBInterfaces = colRight
		case colLeft == "ID_USB_INTERFACE_NUM":
			uInf.IdUSBInterfaceNum = colRight
		case colLeft == "ID_USB_PROTOCOL_FROM_DATABASE":
			uInf.IdUSBProtocolFromDatabase = colRight
		case colLeft == "ID_VENDOR":
			uInf.IdVendor = colRight
		case colLeft == "ID_VENDOR_ENC":
			uInf.IdVendorEnc = colRight
		case colLeft == "ID_VENDOR_FROM_DATABASE":
			uInf.IdVendorFromDatabase = colRight
		case colLeft == "ID_VENDOR_ID":
			uInf.IdVendorID = colRight
		case colLeft == "IFINDEX":
			uInf.IfIndex = colRight
		case colLeft == "INTERFACE":
			uInf.InterfaceName = colRight
		case colLeft == "SUBSYSTEM":
			uInf.Subsystem = colRight
		case colLeft == "USEC_INITIALIZED":
			uInf.USecInitialized = colRight
		}
	}

	return uInf, nil
}
