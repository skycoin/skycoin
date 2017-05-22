package linux

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Sysfs wrapper for linux virtual system filesystem.
// Contains driver, device and kernal info.
type Sysfs struct{}

// NewSysfs creates sysfs instance
func NewSysfs() Sysfs {
	return Sysfs{}
}

// SysfsInfo records the sysfs info
type SysfsInfo struct {
	Main            SysfsInfoMain
	Statistics      SysfsInfoStatistics
	Device          SysfsInfoDevice
	Power           SysfsInfoPower
	ProcNetWireless SysfsInfoProcNetWireless
	Other           SysfsInfoOther
	// Does this interface have a wireless directory
	WirelessDirectoryExists bool
}

// SysfsInfoMain Information from /sys/class/net/{interface}/
type SysfsInfoMain struct {
	// If Hardware (MAC) Address assigned randomly value=1, else 0
	AddrAssignType string
	// Hardware (MAC) Address
	Address string
	// Hardware (MAC) Address Length
	AddrLen string
	// Broadcast Address
	Broadcast string
	// Broken or unplugged physical link/cable value=0, possible with operstate="up",
	// else 1.  For Wifi if signal loss, try setting channel=auto, or access point
	// went down. If operstate="down", value="Invalid argument"
	Carrier string
	// Device id setting used for interrupt handlers
	DevID string
	// If operstate="down", value="Invalid argument"
	Dormant     string
	Flags       string
	IfAlias     string
	IfIndex     string
	IfLink      string
	LinkMode    string
	MTU         string
	NetDevGroup string
	OperState   string
	// If operstate="down", value="Invalid argument"
	Speed      string
	TxQueueLen string
	Type       string
	// Values related to interface hotplugging events.
	// These values can be found individually elsewhere.
	Uevent string
}

// SysfsInfoStatistics Information from /sys/class/net/{interface}/statistics/
type SysfsInfoStatistics struct {
	Collisions int
	Multicast  int
	// Number of bytes received
	RxBytes      int
	RxCompressed int
	RxCrcErrors  int
	// Number of packets dropped while receiving due to no buffer space
	RxDropped int
	// Number of packets received with errors
	RxErrors       int
	RxFifoErrors   int
	RxFrameErrors  int
	RxLengthErrors int
	RxMissedErrors int
	RxOverErrors   int
	// Number of packets received
	RxPackets       int
	TxAbortedErrors int
	// Number of bytes transmitted
	TxBytes int
	// Number of transmission errors due to carrier problems
	TxCarrierErrors int
	TxCompressed    int
	// Number of packets dropped while transmitting due to no buffer space
	TxDropped int
	// Number of packets transmitted with errors
	TxErrors          int
	TxFifoErrors      int
	TxHeartbeatErrors int
	// Number of packets transmitted
	TxPackets      int
	TxWindowErrors int
}

// SysfsInfoDevice Information from /sys/class/net/{interface}/device/
type SysfsInfoDevice struct {
	// Description of the interface, ex: 802.11g WLAN Adapter
	Interface string
	// Formatted form of hardware info that can be referenced
	// with /lib/modules/*/modules.alias to extract info such as
	// vendor and device.
	ModAlias string
	// If interface can automatically go into a suspended state
	// value=0 for false, value=1 for true
	SupportsAutoSuspend string
	// Values related to interface hotplugging events.
	// These values can be found individually elsewhere.
	Uevent string
	// Supposed to contain the vendor but apparently doesn't always exist
	// so might come back blank
	Vendor string
}

// SysfsInfoPower Information from /sys/class/net/{interface}/power/
type SysfsInfoPower struct {
	// allows async operation of interface's suspend and resume
	// callbacks during system power states (hibernation, suspend)
	Async string
	// value=auto allows interface to be power managed at runtime
	// value=on prevents this
	Control              string
	RuntimeActiveKids    string
	RuntimeActiveTime    string
	RuntimeEnabled       string
	RuntimeStatus        string
	RuntimeSuspendedTime string
	RuntimeUsage         string
}

// SysfsInfoProcNetWireless Information from /proc/net/wireless/
type SysfsInfoProcNetWireless struct {
	LinkQuality int
	SignalLevel int
	NoiseLevel  int
}

// SysfsInfoOther Information from things that don't fit anywhere else
type SysfsInfoOther struct {
	// Driver name for the interface
	DriverName string
}

// Run returns wireless information from the file system
func (sfs Sysfs) Run(interfaceName string) SysfsInfo {
	logger.Debug("Sysfs: Querying virtual filesystem for %v", interfaceName)

	inf := SysfsInfo{}
	fq := sfs.sysfsQuery
	qp := ""

	qp = "/sys/class/net/" + interfaceName
	inf.Main.AddrAssignType = fq(qp + "/addr_assign_type")
	inf.Main.Address = fq(qp + "/address")
	inf.Main.AddrLen = fq(qp + "/addr_len")
	inf.Main.Broadcast = fq(qp + "/broadcast")
	inf.Main.Carrier = fq(qp + "/carrier")
	inf.Main.DevID = fq(qp + "/dev_id")
	inf.Main.Dormant = fq(qp + "/dormant")
	inf.Main.Flags = fq(qp + "/flags")
	inf.Main.IfAlias = fq(qp + "/ifalias")
	inf.Main.IfIndex = fq(qp + "/ifindex")
	inf.Main.IfLink = fq(qp + "/iflink")
	inf.Main.LinkMode = fq(qp + "/link_mode")
	inf.Main.MTU = fq(qp + "/mtu")
	inf.Main.NetDevGroup = fq(qp + "/netdev_group")
	inf.Main.OperState = fq(qp + "/operstate")
	inf.Main.Speed = fq(qp + "/speed")
	inf.Main.TxQueueLen = fq(qp + "/tx_queue_len")
	inf.Main.Type = fq(qp + "/type")
	inf.Main.Uevent = fq(qp + "/uevent")

	qp = "/sys/class/net/" + interfaceName + "/device/"
	inf.Device.Interface = fq(qp + "interface")
	inf.Device.ModAlias = fq(qp + "modalias")
	inf.Device.SupportsAutoSuspend = fq(qp + "supports_autosuspend")
	inf.Device.Uevent = fq(qp + "uevent")
	inf.Device.Vendor = fq(qp + "vendor")

	qp = "/sys/class/net/" + interfaceName + "/statistics/"
	inf.Statistics.Collisions, _ = strconv.Atoi(fq(qp + "collisions"))
	inf.Statistics.Multicast, _ = strconv.Atoi(fq(qp + "multicast"))
	inf.Statistics.RxBytes, _ = strconv.Atoi(fq(qp + "rx_bytes"))
	inf.Statistics.RxCompressed, _ = strconv.Atoi(fq(qp + "rx_compressed"))
	inf.Statistics.RxCrcErrors, _ = strconv.Atoi(fq(qp + "rx_crc_errors"))
	inf.Statistics.RxDropped, _ = strconv.Atoi(fq(qp + "rx_dropped"))
	inf.Statistics.RxErrors, _ = strconv.Atoi(fq(qp + "rx_errors"))
	inf.Statistics.RxFifoErrors, _ = strconv.Atoi(fq(qp + "rx_fifo_errors"))
	inf.Statistics.RxFrameErrors, _ = strconv.Atoi(fq(qp + "rx_frame_errors"))
	inf.Statistics.RxLengthErrors, _ = strconv.Atoi(fq(qp + "rx_length_errors"))
	inf.Statistics.RxMissedErrors, _ = strconv.Atoi(fq(qp + "rx_missed_errors"))
	inf.Statistics.RxOverErrors, _ = strconv.Atoi(fq(qp + "rx_over_errors"))
	inf.Statistics.RxPackets, _ = strconv.Atoi(fq(qp + "rx_packets"))
	inf.Statistics.TxAbortedErrors, _ = strconv.Atoi(fq(qp + "tx_aborted_errors"))
	inf.Statistics.TxBytes, _ = strconv.Atoi(fq(qp + "tx_bytes"))
	inf.Statistics.TxCarrierErrors, _ = strconv.Atoi(fq(qp + "tx_carrier_errors"))
	inf.Statistics.TxCompressed, _ = strconv.Atoi(fq(qp + "tx_compressed"))
	inf.Statistics.TxDropped, _ = strconv.Atoi(fq(qp + "tx_dropped"))
	inf.Statistics.TxErrors, _ = strconv.Atoi(fq(qp + "tx_errors"))
	inf.Statistics.TxFifoErrors, _ = strconv.Atoi(fq(qp + "tx_fifo_errors"))
	inf.Statistics.TxHeartbeatErrors, _ = strconv.Atoi(fq(qp + "tx_heartbeat_errors"))
	inf.Statistics.TxPackets, _ = strconv.Atoi(fq(qp + "tx_packets"))
	inf.Statistics.TxWindowErrors, _ = strconv.Atoi(fq(qp + "tx_window_errors"))

	qp = "/sys/class/net/" + interfaceName + "/power/"
	inf.Power.Async = fq(qp + "async")
	inf.Power.Control = fq(qp + "control")
	inf.Power.RuntimeActiveKids = fq(qp + "runtime_active_kids")
	inf.Power.RuntimeActiveTime = fq(qp + "runtime_active_time")
	inf.Power.RuntimeEnabled = fq(qp + "runtime_enabled")
	inf.Power.RuntimeStatus = fq(qp + "runtime_status")
	inf.Power.RuntimeSuspendedTime = fq(qp + "runtime_suspended_time")
	inf.Power.RuntimeUsage = fq(qp + "runtime_usage")

	// supposedly in newer kernel's these values do not exist anymore
	// link quality: /sys/class/net/{interface}/wireless/link
	// signal level: /sys/class/net/{interface}/wireless/level (signed integer)
	// so instead we parse the legacy procfs
	pnwStats := fq("/proc/net/wireless")
	for lineIndex, line := range strings.Split(pnwStats, "\n") {
		if lineIndex >= 2 {
			pnwInterfaceName := ""

			fields := strings.Fields(line)
			for fieldIndex, field := range fields {
				field = strings.Trim(field, ":")
				field = strings.Trim(field, ".")

				if fieldIndex == 0 {
					pnwInterfaceName = field
				}

				if pnwInterfaceName == interfaceName {
					switch fieldIndex {
					case 2:
						inf.ProcNetWireless.LinkQuality, _ = strconv.Atoi(field)
					case 3:
						inf.ProcNetWireless.SignalLevel, _ = strconv.Atoi(field)
					case 4:
						inf.ProcNetWireless.NoiseLevel, _ = strconv.Atoi(field)
					}
				}
			}
		}
	}

	// Get the driver name from the first directory found in this path
	qp = "/sys/class/net/" + interfaceName + "/device/driver/module/drivers/"
	dirInfos, err := ioutil.ReadDir(qp)
	if err == nil && len(dirInfos) > 0 {
		inf.Other.DriverName = dirInfos[0].Name()
	}

	// Check for existence of wireless path
	qp = "/sys/class/nintentedet/" + interfaceName + "/wireless"
	inf.WirelessDirectoryExists = sfs.dirExists(qp)

	return inf
}

func (sfs Sysfs) sysfsQuery(queryFile string) string {
	out, _ := ioutil.ReadFile(queryFile)
	outs := string(out)
	outs = strings.TrimSpace(outs)
	return outs
}

func (sfs Sysfs) dirExists(dirPath string) bool {
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
