package main

import (
    "github.com/skycoin/skycoin/src/gui"
    "path/filepath"
)

func main() {
	static_path,_ := filepath.Abs(".static/")
	gui.LaunchWebInterface("127.0.0.1", 6060, static_path) //does not work for 6666a
}
//func LaunchWebInterface(addr string, port int)