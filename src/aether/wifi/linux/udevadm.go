package linux

import (
	"os/exec"
	"strings"
)

// UDevAdm Wrapper for linux utility: udevadm
type UDevAdm struct{}

// NewUDevAdm creates instance
func NewUDevAdm() UDevAdm {
	return UDevAdm{}
}

// UDevAdmInfo records udevadm info
type UDevAdmInfo struct {
	DevPath                   string // /devices/pci0000:00/...
	DevType                   string // wlan
	IDBus                     string // usb
	IDMMCandidate             string // 1
	IDModel                   string // Model (useful)
	IDModelEnc                string // Model Encoded
	IDModelID                 string // Model Id
	IDNetNameMac              string //
	IDNetNamePath             string //
	IDOUIFromDatabase         string // Corporation Name
	IDRevision                string //
	IDSerial                  string //
	IDSerialShort             string //
	IDType                    string // generic
	IDUSBClassFromDatabase    string // Miscellaneous Device
	IDUSBDriver               string // Driver Name (useful)
	IDUSBInterfaces           string //
	IDUSBInterfaceNum         string //
	IDUSBProtocolFromDatabase string // Interface Association
	IDVendor                  string // Manufacturer name (useful)
	IDVendorEnc               string // Manufacturer name encoded
	IDVendorFromDatabase      string // Manufacturer name long
	IDVendorID                string //
	IfIndex                   string
	InterfaceName             string // wlan0
	Subsystem                 string // net
	USecInitialized           string //
}

// IsInstalled checks if the program udevadm exists using PATH environment variable
func (uda UDevAdm) IsInstalled() bool {
	_, err := exec.LookPath("udevadm")
	if err != nil {
		return false
	}
	return true
}

// Run returns information on an interface
func (uda UDevAdm) Run(interfaceName string) (UDevAdmInfo, error) {
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
	uInf, err := uda.udevadmParse(content)

	return uInf, err
}

func (uda UDevAdm) udevadmParse(content string) (UDevAdmInfo, error) {
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
			uInf.IDBus = colRight
		case colLeft == "ID_MM_CANDIDATE":
			uInf.IDMMCandidate = colRight
		case colLeft == "ID_MODEL":
			uInf.IDModel = colRight
		case colLeft == "ID_MODEL_ENC":
			uInf.IDModelEnc = colRight
		case colLeft == "ID_MODEL_ID":
			uInf.IDModelID = colRight
		case colLeft == "ID_NET_NAME_MAC":
			uInf.IDNetNameMac = colRight
		case colLeft == "ID_NET_NAME_PATH":
			uInf.IDNetNamePath = colRight
		case colLeft == "ID_OUI_FROM_DATABASE":
			uInf.IDOUIFromDatabase = colRight
		case colLeft == "ID_REVISION":
			uInf.IDRevision = colRight
		case colLeft == "ID_SERIAL":
			uInf.IDSerial = colRight
		case colLeft == "ID_SERIAL_SHORT":
			uInf.IDSerialShort = colRight
		case colLeft == "ID_TYPE":
			uInf.IDType = colRight
		case colLeft == "ID_USB_CLASS_FROM_DATABASE":
			uInf.IDUSBClassFromDatabase = colRight
		case colLeft == "ID_USB_DRIVER":
			uInf.IDUSBDriver = colRight
		case colLeft == "ID_USB_INTERFACES":
			uInf.IDUSBInterfaces = colRight
		case colLeft == "ID_USB_INTERFACE_NUM":
			uInf.IDUSBInterfaceNum = colRight
		case colLeft == "ID_USB_PROTOCOL_FROM_DATABASE":
			uInf.IDUSBProtocolFromDatabase = colRight
		case colLeft == "ID_VENDOR":
			uInf.IDVendor = colRight
		case colLeft == "ID_VENDOR_ENC":
			uInf.IDVendorEnc = colRight
		case colLeft == "ID_VENDOR_FROM_DATABASE":
			uInf.IDVendorFromDatabase = colRight
		case colLeft == "ID_VENDOR_ID":
			uInf.IDVendorID = colRight
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
