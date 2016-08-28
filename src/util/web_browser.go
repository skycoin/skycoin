package util

import (
	"github.com/toqueteos/webbrowser"
)

func OpenBrowser(url string) error {
	return webbrowser.Open(url)
}
