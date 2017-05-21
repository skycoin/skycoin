package linux

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"

	logging "github.com/op/go-logging"
)

var (
	// ErrAuthRequired authentication required error
	ErrAuthRequired = errors.New("sudo authentication required")
	logger          = logging.MustGetLogger("darknet.network")
)

func authorized() bool {
	if authorizedCheck() {
		return true
	}

	logger.Debug("Requesting superuser authentication")

	var cmd *exec.Cmd
	var sudoTool = "sudo"
	/*if !authorizedIsCheck() {
		_, err := exec.LookPath("gksudo")
		if err == nil {
			return "gksudo"
		}
	}*/
	if sudoTool == "sudo" {
		cmd = exec.Command("sudo", "echo")
	} else {
		cmd = exec.Command("gksudo", "darknet")
	}
	_, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	if authorizedCheck() {
		return true
	}

	return false
}

func authorizedCheck() bool {
	if os.Geteuid() == 0 {
		return true
	}

	cmd := exec.Command("sudo", "-n", "true")
	if err := cmd.Start(); err != nil {
		logger.Debug("sudo check error: %v")
		return false
	}

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
				_ = exitStatus
				return false
			}
		}
	}
	return true
}

func limitText(text []byte) string {
	t := strings.TrimSpace(string(text))
	if len(t) > 150 {
		t = t[0:150] + "..."
	}
	return "[" + t + "]"
}
