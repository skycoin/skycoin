package main

import (
	"flag" // cli parsing
	"fmt"
	"log"
)

import (
	"coin" //skycoin daemon
	"common"
	"gui"
)

var disable_gui bool = false
var disable_coind bool = false
var disable_coin_tests bool = false

func parse_args() {
	_gui_port := flag.Uint64("gui-port", uint64(sb_gui.GuiServerPort), "port to run gui server")
	flag.BoolVar(&disable_coind, "disable-coind", disable_coind, "disable the coin daemon")
	flag.BoolVar(&disable_gui, "disable-gui", disable_gui, "disable the gui server")
	flag.Parse()
	sb_gui.GuiServerPort = uint32(*_gui_port)
}

func tests() {

	pub, private := sb.GenerateKeyPair()
	fmt.Printf("pub: %s \n", pub)

	_ = private

	var BC *sb_coin.BlockChain = sb_coin.NewBlockChain()

	if false {
		fmt.Printf("l= %v\n", len(BC.Blocks))
	}

	B := BC.NewBlock()

	var T sb_coin.Transaction

	T.UpdateHeader() //sets hash

	/*
		Need to add input
	*/
	err := BC.AppendTransaction(B, &T)

	if err != nil {
		log.Panic(err)
	}

}

func main() {

	tests()
	return

	if false {
		log.Panic()
	}

	parse_args()

	var n_servers int = 2

	if disable_gui {
		n_servers--
	}
	if disable_coind {
		n_servers--
	}

	if n_servers <= 0 {
		fmt.Printf("Nothing to run. You disabled everything\n")
		return
	}

	errchan := make(chan error, n_servers)

	if !disable_gui {

	}

	if !disable_coind {

	}

	for err := range errchan {
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			break
		}
	}

	fmt.Println("Goodbye")
}

/*
   Junk
*/
