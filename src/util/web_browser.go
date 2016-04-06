package util

import (
	"github.com/toqueteos/webbrowser" //open webbrowser
)

func OpenBrowser(url string) {
	webbrowser.Open(url)
}
