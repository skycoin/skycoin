package util

import (
	"github.com/toqueteos/webbrowser"
)

// OpenBrowser opens browser
func OpenBrowser(url string) error {
	return webbrowser.Open(url)
}
