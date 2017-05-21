package linux

import (
	"net"
	"os/exec"
	"strings"
)

// ResolvConf Wrapper for linux utility: resolvconf
// This modifies /etc/resolv.conf because any manual changes
// not using resolvconf will be removed the next time something
// uses this tool. So we must use it too. Other tools that
// utilize resolvconf are: pppd, dhclient, ifup, ifdown, dnsmasq
//
// Files that might end up being modified to verify changes:
//	/etc/resolv.conf
//	/etc/resolv.conf.d
//	/etc/resolvconf/resolv.conf.d
//
// Networkmanager's resolvconf config file
//	/run/resolvconf/interface/NetworkManager
//
// If we use this tool as-is, we get a file like:
//	/run/resolvconf/interface/wlan0.darknet
//
type ResolvConf struct{}

// NewResolvConf creates ResolvConf instance
func NewResolvConf() ResolvConf {
	return ResolvConf{}
}

// IsInstalled checks if the program route exists using PATH environment variable
func (rc ResolvConf) IsInstalled() bool {
	_, err := exec.LookPath("resolvconf")
	if err != nil {
		return false
	}
	return true
}

// SetFromFile add or Overwrite the interface records using a file as input
func (rc ResolvConf) SetFromFile(interfaceName string, programName string,
	fileName string) error {
	logger.Debug("ResolveConf: Setting records from file")

	// cat filename | resolvconf -a [interface].[programname]
	cmd1 := exec.Command("cat", "filename", fileName)
	cmd2 := exec.Command("resolvconf", "-a", interfaceName+"."+programName)
	cmd2.Stdin, _ = cmd1.StdoutPipe()

	logger.Debug("Command Start: %v | %v", cmd1.Args, cmd2.Args)
	out, err := cmd2.Output()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Set add or Overwrite the interface records using array of nameserver IPs
func (rc ResolvConf) Set(interfaceName string, programName string,
	nameserversList []net.IP) error {
	logger.Debug("ResolveConf: Setting records from list")

	resolvText := ""
	for _, nameserverIP := range nameserversList {
		resolvText += "nameserver " + nameserverIP.String() + "\n"
	}

	// echo nameserver [address] | resolvconf -a [interface].[programname]
	cmd := exec.Command("resolvconf", "-a", interfaceName+"."+programName)
	cmd.Stdin = strings.NewReader(resolvText)
	logger.Debug("Command Start: \n Piping \n %v into %v", resolvText, cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Delete the interface records
func (rc ResolvConf) Delete(interfaceName string, programName string) error {
	logger.Debug("ResolveConf: Deleting records")

	// resolvconf -d [interface].[programname]
	cmd := exec.Command("resolvconf", "-d", interfaceName+"."+programName)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// Update refresh (Update) the interface records
func (rc ResolvConf) Update() error {
	logger.Debug("ResolveConf: Updating records")

	// resolvconf -u
	cmd := exec.Command("resolvconf", "-u")
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}
