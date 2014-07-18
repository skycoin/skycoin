package linux

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

// Wrapper for linux utility: rfkill
// Checks the status of kill switches. If either is set, the device will be disabled.
// Soft = Software (set by software)
// Hard = Hardware (physical on/off switch on the device)
// Identifiers = all, wifi, wlan, bluetooth, uwb, ultrawideband, wimax, wwan, gps, fm
// See: http://wireless.kernel.org/en/users/Documentation/rfkill
type RFKill struct{}

func NewRFKill() RFKill {
	return RFKill{}
}

type RFKillResult struct {
	Index          int
	IdentifierType string
	Description    string
	SoftBlocked    bool
	HardBlocked    bool
}

// Checks if the program rfkill exists using PATH environment variable
func (self RFKill) IsInstalled() bool {
	_, err := exec.LookPath("rfkill")
	if err != nil {
		return false
	}
	return true
}

// Returns a list of rfkill results for every identifier type
func (self RFKill) ListAll() ([]RFKillResult, error) {
	rfks := []RFKillResult{}
	rfk := RFKillResult{}
	fq := self.fileQuery

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

			rfk.Index, _ = strconv.Atoi(fq(qp + "/index"))
			rfk.IdentifierType = fq(qp + "/type")
			rfk.Description = fq(qp + "/name")
			rfk.SoftBlocked = false
			rfk.HardBlocked = false
			if fq(qp+"/soft") == "1" {
				rfk.SoftBlocked = true
			}
			if fq(qp+"/hard") == "1" {
				rfk.HardBlocked = true
			}

			rfks = append(rfks, rfk)
			rfk = RFKillResult{}
		}
	}
	return rfks, nil
}

// Sets a software block on an identifier
func (self RFKill) SoftBlock(identifier string) error {
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

// Removes a software block on an identifier
func (self RFKill) SoftUnblock(identifier string) error {
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

// Checks if an identifier has a software or hardware block
func (self RFKill) IsBlocked(identifier string) bool {
	rfks, _ := self.ListAll()
	for _, rfk := range rfks {
		if self.checkThis(rfk, identifier) {
			if rfk.SoftBlocked || rfk.HardBlocked {
				return true
			}
		}
	}
	return false
}

// Checks if an identifier has a software block
func (self RFKill) IsSoftBlocked(identifier string) bool {
	rfks, _ := self.ListAll()
	for _, rfk := range rfks {
		if self.checkThis(rfk, identifier) {
			if rfk.SoftBlocked {
				return true
			}
		}
	}
	return false
}

// Checks if an identifier has a hardware block
func (self RFKill) IsHardBlocked(identifier string) bool {
	rfks, _ := self.ListAll()
	for _, rfk := range rfks {
		if self.checkThis(rfk, identifier) {
			if rfk.HardBlocked {
				return true
			}
		}
	}
	return false
}

// Checks if an identifier has a software or hardware block after
// removing a software block if it exists
func (self RFKill) IsBlockedAfterUnblocking(identifier string) bool {
	if self.IsBlocked(identifier) {
		self.SoftUnblock(identifier)
		if self.IsBlocked(identifier) {
			return true
		}
	}

	return false
}

func (self RFKill) checkThis(rfk RFKillResult, identifier string) bool {
	switch identifier {
	case "":
		return true
	case "all":
		return true
	case rfk.IdentifierType:
		return true
	}
	return false
}

func (self RFKill) fileQuery(queryFile string) string {
	out, _ := ioutil.ReadFile(queryFile)
	outs := string(out)
	outs = strings.TrimSpace(outs)
	return outs
}
