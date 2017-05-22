package network

// Wifi Interface Control Library
//
// NetworkManager can do almost everything. However some lightweight linux
// distributions might not have it or a user has disabled NetworkManager
// because it can conflict with other networking software or desired manual control.
//
// Because of that, we use this operation preference:
//    If NetworkManager installed and running, use it because it blocks other tools
//    If NetworkManager installed and not running, dont use it, use old tools.
//    If NetworkManager not installed, use old tools (ifconfig, iwconfig, etc)
//    If NetworkManager not installed and want WPA or WPA2, then require wpa_supplicant
//
// Wireless modes to be supported:
//    managed         connect to access point
//    ad-hoc          (not yet supported) p2p networking
//    master          (not yet supported) act as access point
//    monitor         (not yet supported) RFMON mode
//
// Tools used in this wifi library:
//    dhclient         DHCP client
//    ifconfig         Set interface up, down, IP
//    iwconfig         Set wireless mode, channel, essid, no security, wep security
//    iwlist           Get detailed list of wireless networks
//    resolvconf       Set nameserver information
//    rfkill           Set, unset, detect software or hardware device locks
//    route            Set gateway
//    Sysfs            Linux virtual filesystem used to easily query info
//    udevadm          Get detailed information for an interface
//    wpa_supplicant   Manages wpa, wpa2 security (wpa_cli)
//    NetworkManager   Manages all of the above (nmcli)
//
//    Because command line tools are used. We set LC_ALL=C to set the locale
//    to ASCII (not UTF8 or other langs) to stabilize parsing results. (todo)
//
// Tools deemed unusable:
//    iw/ifup/ifdown/ifquery	Requires cards with nl80211 drivers
//    lshw -C network			Slow
//
// Research later:
//    ip link
//
import (
	"fmt"
	"net"
)

/*
const (
	WIFI_SECURITY_NONE = iota
	WIFI_SECURITY_WEP
	WIFI_SECURITY_WPA_PERSONAL
)
*/

type WifiConnection struct {
	ConnectionName string
	InterfaceName  string
	//
	Mode             string
	SSID             string
	Channel          string
	Frequency        string
	SecurityProtocol string // [NONE, WEP, WPA]
	SecurityKey      string
	DHCPEnabled      bool
	Addresses        []Address
	Routes           []Route
	Nameservers      []net.IP
	DefaultGateway   net.IP
}

type WifiStats struct {
	Carrier     string
	OperState   string
	Address     string
	LinkQuality int
	SignalLevel int
	NoiseLevel  int
	RxBytes     int
	TxBytes     int
	RxPackets   int
	TxPackets   int
}

type WifiNetwork struct {
	Address             string `json:"address"`
	ESSID               string `json:"essid"`
	Protocol            string `json:"protocol"`
	Mode                string `json:"mode"`
	Frequency           string `json:"frequency"`
	EncryptionKeyStatus string `json:"encryption_key_status"`
	BitRates            string `json:"bit_rates"`
	InformationElement  []WifiNetworkIE
	QualityLevel        int `json:"quality_level"`
	SignalLevel         int `json:"signal_level"`
	NoiseLevel          int `json:"noise_level"`
}

type WifiNetworkIE struct {
	Protocol        string `json:"protocol"`
	GroupCipher     string `json:"group_cipher"`
	PairwiseCiphers string `json:"pairwise_ciphers"`
	AuthSuites      string `json:"auth_suites"`
	Extra           string `json:"extra"`
}

// Extends net.Interface for managing wifi interfaces
// See: http://golang.org/src/pkg/net/interface.go
type WifiInterface struct {
	net.Interface
	Model      string
	Driver     string
	Vendor     string
	Connection WifiConnection
	Statistics WifiStats
	Networks   WifiNetwork
}

// Returns a list of wifi interfaces installed
func NewWifiInterfaces() ([]WifiInterface, error) {
	logger.Debug("Wifi: Gathering interfaces")

	ifacesOut := []WifiInterface{}

	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	// Add if wireless interface
	for _, iface := range ifaces {
		fqs := sysfs.Run(iface.Name)
		if fqs.WirelessDirectoryExists {
			ifaceOut := NewWifiInterface(iface)
			logger.Info("Wifi: Found wifi interface %v %v %v %v",
				ifaceOut.Name, ifaceOut.Model, ifaceOut.Driver, ifaceOut.Vendor)

			ifacesOut = append(ifacesOut, ifaceOut)
		}
	}

	return ifacesOut, nil
}

// Return an initialized WifiInterface type
func NewWifiInterface(ifaceNet net.Interface) WifiInterface {
	iface := WifiInterface{Interface: ifaceNet}

	if udevadm.IsInstalled() {
		// Use udevadm if installed
		udevInfo, _ := udevadm.Run(iface.Name)
		iface.Model = udevInfo.IDModel
		iface.Driver = udevInfo.IDUSBDriver
		iface.Vendor = udevInfo.IDVendor
	} else {
		// Use sysfs as backup, vendor will sometimes return blank
		sysfsInfo := sysfs.Run(iface.Name)
		iface.Model = sysfsInfo.Device.Interface
		iface.Driver = sysfsInfo.Other.DriverName
		iface.Vendor = sysfsInfo.Device.Vendor
	}

	return iface
}

func (self *WifiInterface) Scan() ([]WifiNetwork, error) {
	logger.Info("Wifi: Scanning for networks using %v", self.Name)
	wnetsIn, err := iwlist.Scan(self.Name)

	// Clone the values from iwlist results
	// They are similar now, but that could change
	wnets := []WifiNetwork{}
	for _, wnetIn := range wnetsIn {
		wnet := WifiNetwork{}

		wnet.Address = wnetIn.Address
		wnet.ESSID = wnetIn.ESSID
		wnet.Protocol = wnetIn.Protocol
		wnet.Mode = wnetIn.Mode
		wnet.Frequency = wnetIn.Frequency
		wnet.EncryptionKeyStatus = wnetIn.EncryptionKeyStatus
		wnet.BitRates = wnetIn.BitRates

		for _, ieIn := range wnetIn.InformationElement {
			ie := WifiNetworkIE{}

			ie.Protocol = ieIn.Protocol
			ie.GroupCipher = ieIn.Protocol
			ie.PairwiseCiphers = ieIn.Protocol
			ie.AuthSuites = ieIn.Protocol
			ie.Extra = ieIn.Protocol

			wnet.InformationElement = append(wnet.InformationElement, ie)
		}

		wnet.QualityLevel = wnetIn.QualityLevel
		wnet.SignalLevel = wnetIn.SignalLevel
		wnet.NoiseLevel = wnetIn.NoiseLevel

		wnets = append(wnets, wnet)
	}

	return wnets, err
}

// Start interface config
func (self *WifiInterface) Start() {
	logger.Info("Wifi: Starting interface %v", self.Name)

	if !prerequisites() {
		return
	}

	// Now start the interface
	if softwareType == "networkmanager" {
		networkmanager.Connect(self.Name,
			self.Connection.SSID,
			self.Connection.SecurityProtocol,
			self.Connection.SecurityKey)
	} else {
		// Before starting, make sure everything is stopped
		//self.Stop()

		switch self.Connection.Mode {
		case "adhoc":
			if self.isUp() {
				self.down()
			}
			self.mode()
			self.security()
			self.channel()
			self.frequency()
			self.essid()

			// ad-hoc needs settings first, then interface up
			self.up()
			// iwlist wlan0 scan
		case "managed":
			self.up()
			self.mode()
			self.security()
			self.channel()
			self.frequency()
			self.essid()
			self.addressing("start")
			//check we are connected
		}
	}
}

// Stop interface config
func (self *WifiInterface) Stop() {
	logger.Info("Wifi: Stopping interface %#v", self.Name)

	if !prerequisites() {
		return
	}

	switch self.Connection.Mode {
	case "adhoc":
		self.down()
	case "managed":
		// addressing goes down before interface goes down
		self.addressing("stop")
		self.down()

		if self.Connection.SecurityProtocol == "WPA" {
			wpasupplicant.DaemonShutdown()
		}
	}
}

// Updates interface statistics
func (self *WifiInterface) Stats() {
	logger.Info("Wifi: Getting statistics")

	si := sysfs.Run(self.Name)

	stats := &self.Statistics
	stats.Carrier = si.Main.Carrier
	stats.OperState = si.Main.OperState
	stats.Address = si.Main.Address
	stats.RxBytes = si.Statistics.RxBytes
	stats.TxBytes = si.Statistics.TxBytes
	stats.LinkQuality = si.ProcNetWireless.LinkQuality
	stats.SignalLevel = si.ProcNetWireless.SignalLevel
	stats.NoiseLevel = si.ProcNetWireless.NoiseLevel
}

func (self *WifiInterface) mode() {
	logger.Info("Wifi: Setting mode")

	iwconfig.Mode(self.Name, self.Connection.Mode)
}

func (self *WifiInterface) channel() {
	logger.Info("Wifi: Setting channel")

	switch softwareType {
	case "networkmanager":
	case "legacy":
		if self.Connection.Mode == "managed" {
			if self.Connection.Channel == "" {
				// Channel not set, do nothing
				// In managed the access point automatically dictates the channel
			} else {
				// Channel is set, use it
				iwconfig.Channel(self.Name, self.Connection.Channel)
			}
		}
		if self.Connection.Mode == "adhoc" {
			// In ad-hoc channel or frequency must be set before initial cell
			// creation if we don't have one set, we'll put auto
			iwconfig.Channel(self.Name, "auto")
		}
	}
}

func (self *WifiInterface) frequency() {
	logger.Info("Wifi: Setting frequency")

	switch softwareType {
	case "networkmanager":
	case "legacy":
		if self.Connection.Mode == "managed" {
			if self.Connection.Frequency == "" {
				// Frequency not set, do nothing
			} else {
				// Frequency is set, use it
				iwconfig.Frequency(
					self.Name,
					self.Connection.Frequency)
			}
		}
		if self.Connection.Mode == "adhoc" {
		}
	}
}

func (self *WifiInterface) security() {
	logger.Info("Wifi: Setting security")

	if self.Connection.SecurityProtocol == "" {
		logger.Info("Wifi: Set network encryption - none")
	} else {
		logger.Info("Wifi: Set network encryption - %v",
			self.Connection.SecurityProtocol)
	}

	switch softwareType {
	case "networkmanager":
	case "legacy":
		switch self.Connection.SecurityProtocol {
		default:
			// No Encryption
			iwconfig.Key(self.Name, "off")
		case "WEP":
			// WEP, use iwconfig
			iwconfig.Key(self.Name, self.Connection.SecurityKey)
		case "WPA":
			wpasupplicant.ConfigWrite(self.Name,
				self.Connection.SSID,
				self.Connection.SecurityKey)
			wpasupplicant.DaemonStartup(self.Name)
			//wpasupplicant.Authenticate(self.Name, self.Connection.SSID,
			//	self.Connection.SecurityKey)
		}
	}
}

func (self *WifiInterface) essid() {
	logger.Info("Wifi: Setting SSID/Network Name to %v",
		self.Connection.SSID)

	switch softwareType {
	case "networkmanager":
	case "legacy":
		// WEP, use iwconfig
		iwconfig.ESSID(self.Name, self.Connection.SSID)
	}
}

func (self *WifiInterface) addressing(actionType string) {
	logger.Info("Wifi: Setting IP address configuration")

	switch actionType {
	case "start":
		if self.Connection.DHCPEnabled {
			dhclient.Startup(self.Name)
		} else {
			// IP and Subnetmask (from config)
			for _, address := range self.Connection.Addresses {
				// IP and Subnetmask
				ifconfig.SetIP(self.Name, address.IPNet.IP, address.IPNet.Mask)

				// Gateway if exists
				//address.Gateway.String() !=
			}

			// Default Gateway (from config if set)
			if self.Connection.DefaultGateway != nil {
				route.AddDefaultGateway(self.Connection.DefaultGateway.String())
			}

			// DNS Server (from config if set)
			if self.Connection.Nameservers != nil {
				resolvconf.Set(self.Name, "darknet", self.Connection.Nameservers)
				resolvconf.Update()
			}
		}
	case "stop":
		if self.Connection.DHCPEnabled {
			dhclient.Shutdown(self.Name)
		} else {
			// DNS Server
			if self.Connection.Nameservers != nil {
				resolvconf.Delete(self.Name, "darknet")
				resolvconf.Update()
			}
		}
	}
}

func (self *WifiInterface) down() {
	if self.isUp() {
		logger.Info("Wifi: Setting interface down")

		switch softwareType {
		case "networkmanager":
			//
		case "legacy":
			ifconfig.Down(self.Name)
		}
	} else {
		logger.Info("Wifi: %v already down, skipping", self.Name)
	}
}

func (self *WifiInterface) up() {
	if !self.isUp() {
		logger.Info("Wifi: Setting interface up")

		// Unblocking any rfkill switches
		self.unblock()

		switch softwareType {
		case "networkmanager":
			// nmcli automatically brings interface up
		case "legacy":
			ifconfig.Up(self.Name)
		}
	} else {
		logger.Info("Wifi: %v already up, skipping", self.Name)
	}
}

func (self *WifiInterface) isUp() bool {
	if self.Flags&net.FlagUp == net.FlagUp {
		return true
	}
	return false
}

// Unblock any rfkill software switches
func (self *WifiInterface) unblock() error {
	if rfkill.IsBlockedAfterUnblocking("wlan") {
		logger.Warning("Wifi: %v is blocked by rfkill", "wlan")
	}
	return nil
}

func prerequisites() bool {
	logger.Debug("Wifi: Checking software prerequisites")

	failed := false

	// Use NetworkManager if its installed & running, otherwise everything legacy
	running, err := networkmanager.ServiceIsRunning()
	if networkmanager.IsInstalled() && running && err == nil {
		softwareType = "networkmanager"
	} else {
		softwareType = "legacy"
	}

	var installList = map[string]bool{
		"dhclient":       dhclient.IsInstalled(),
		"ifconfig":       ifconfig.IsInstalled(),
		"iwconfig":       iwconfig.IsInstalled(),
		"iwlist":         iwlist.IsInstalled(),
		"networkmanager": networkmanager.IsInstalled(),
		"resolvconf":     resolvconf.IsInstalled(),
		"rfkill":         rfkill.IsInstalled(),
		"route":          route.IsInstalled(),
		"udevadm":        udevadm.IsInstalled(),
		"wpasupplicant":  wpasupplicant.IsInstalled(),
	}

	if failed {
		return false
	}

	logger.Debug("Wifi: Software chosen: %v", softwareType)
	logger.Debug("Wifi: Software installed: %v", installList)

	return true
}
