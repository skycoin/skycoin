package linux

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

// RFKill Wrapper for linux utility: rfkill
// Checks the status of kill switches. If either is set, the device will be disabled.
// Soft = Software (set by software)
// Hard = Hardware (physical on/off switch on the device)
// Identifiers = all, wifi, wlan, bluetooth, uwb, ultrawideband, wimax, wwan, gps, fm
// See: http://wireless.kernel.org/en/users/Documentation/rfkill
type RFKill struct{}

// NewRFKill creates new RFKill instance
func NewRFKill() RFKill {
	return RFKill{}
}

// RFKillResult represents the result of RFKill
type RFKillResult struct {
	Index          int
	IdentifierType string
	Description    string
	SoftBlocked    bool
	HardBlocked    bool
}

// IsInstalled checks if the program rfkill exists using PATH environment variable
func (rfk RFKill) IsInstalled() bool {
	_, err := exec.LookPath("rfkill")
	if err != nil {
		return false
	}
	return true
}

// ListAll returns a list of rfkill results for every identifier type
func (rfk RFKill) ListAll() ([]RFKillResult, error) {
	rfkrs := []RFKillResult{}
	rfkr := RFKillResult{}
	fq := rfk.fileQuery

	// instead of parsing "rfkill list", query the filesystem
	dirInfos, err := ioutil.ReadDir("/sys/class/rfkill/")
	if err != nil {
		return nil, fmt.Errorf(
			"RFKill: Error reading directory '/sys/class/rfkill/': %v", err)
	}

	for _, dirInfo := range dirInfos {
		// directory starts with "rfkill"
		if len(dirInfo.Name()) > 6 && dirInfo.Name()[0:6] == "rfkill" {
			qp := "/sys/class/rfkill/" + dirInfo.Name()

			rfkr.Index, _ = strconv.Atoi(fq(qp + "/index"))
			rfkr.IdentifierType = fq(qp + "/type")
			rfkr.Description = fq(qp + "/name")
			rfkr.SoftBlocked = false
			rfkr.HardBlocked = false
			if fq(qp+"/soft") == "1" {
				rfkr.SoftBlocked = true
			}
			if fq(qp+"/hard") == "1" {
				rfkr.HardBlocked = true
			}

			rfkrs = append(rfkrs, rfkr)
			rfkr = RFKillResult{}
		}
	}
	return rfkrs, nil
}

// SoftBlock sets a software block on an identifier
func (rfk RFKill) SoftBlock(identifier string) error {
	logger.Debug("RFKill: Soft Blocking %v", identifier)

	cmd := exec.Command("rfkill", "block", identifier)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// SoftUnblock Removes a software block on an identifier
func (rfk RFKill) SoftUnblock(identifier string) error {
	logger.Debug("RFKill: Soft Unblocking %v", identifier)

	cmd := exec.Command("rfkill", "unblock", identifier)
	logger.Debug("Command Start: %v", cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Command Error: %v : %v", err, limitText(out))
		return err
	}
	logger.Debug("Command Return: %v", limitText(out))

	return nil
}

// IsBlocked Checks if an identifier has a software or hardware block
func (rfk RFKill) IsBlocked(identifier string) bool {
	rfkrs, _ := rfk.ListAll()
	for _, rfkr := range rfkrs {
		if rfk.checkThis(rfkr, identifier) {
			if rfkr.SoftBlocked || rfkr.HardBlocked {
				return true
			}
		}
	}
	return false
}

// IsSoftBlocked checks if an identifier has a software block
func (rfk RFKill) IsSoftBlocked(identifier string) bool {
	rfkrs, _ := rfk.ListAll()
	for _, rfkr := range rfkrs {
		if rfk.checkThis(rfkr, identifier) {
			if rfkr.SoftBlocked {
				return true
			}
		}
	}
	return false
}

// IsHardBlocked checks if an identifier has a hardware block
func (rfk RFKill) IsHardBlocked(identifier string) bool {
	rfkrs, _ := rfk.ListAll()
	for _, rfkr := range rfkrs {
		if rfk.checkThis(rfkr, identifier) {
			if rfkr.HardBlocked {
				return true
			}
		}
	}
	return false
}

// IsBlockedAfterUnblocking Checks if an identifier has a software or hardware block after
// removing a software block if it exists
func (rfk RFKill) IsBlockedAfterUnblocking(identifier string) bool {
	if rfk.IsBlocked(identifier) {
		rfk.SoftUnblock(identifier)
		if rfk.IsBlocked(identifier) {
			return true
		}
	}

	return false
}

func (rfk RFKill) checkThis(rfkr RFKillResult, identifier string) bool {
	switch identifier {
	case "":
		return true
	case "all":
		return true
	case rfkr.IdentifierType:
		return true
	}
	return false
}

func (rfk RFKill) fileQuery(queryFile string) string {
	out, _ := ioutil.ReadFile(queryFile)
	outs := string(out)
	outs = strings.TrimSpace(outs)
	return outs
}
