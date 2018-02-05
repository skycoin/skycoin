package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
)

/**
Current list of tested commands:

generateAddresses
verifyAddress
send
send -m (send-to-many)
broadcastTransaction
createRawTransaction
getWalletBalance
transaction
status
*/

/**
The minimal requirements for the cli are
- a wallet file
- two addresses ( to test send -m )
*/

const cliName = "skycoin-cli"

//@TODO We can put the commands to be tested and their arguments into an array and iterate over it to keep this more DRY
func main() {
	var (
		cmdOut    []byte
		err       error
		help      bool
		wltFile   string
		addresses string // comma separated addresses, to be used in send and send -m commands, min 2 required
		testAddrs []string
	)

	flag.BoolVar(&help, "help", false, "Show help")
	flag.StringVar(&wltFile, "wallet-file", wltFile, "wallet file used for testing cli commands")
	flag.StringVar(&addresses, "addrs", addresses, "destination addresses for sending coins")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if wltFile == "" {
		fmt.Fprint(os.Stderr, "no wallet file given")
		os.Exit(1)
	}

	if _, err := os.Stat(wltFile); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("wallet file %s does not exist", wltFile))
		os.Exit(1)
	}

	if addresses == "" {
		fmt.Fprintln(os.Stderr, "no test addresses given")
		os.Exit(1)
	}

	testAddrs = strings.Split(addresses, ",")
	if len(testAddrs) < 2 {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("minimum two test addresses required, given: %v", len(testAddrs)))
		os.Exit(1)
	}

	//@TODO add a check to make sure that the addresses are unique?
	for i := range testAddrs {
		_, err = cipher.DecodeBase58Address(testAddrs[i])
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("address %s is invalid: ", testAddrs[i]), err)
			os.Exit(1)
		}
	}

	cmdArgs := []string{"generateAddresses", "-f", wltFile}
	if cmdOut, err = exec.Command(cliName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	address := strings.TrimSpace(string(cmdOut))
	// verify that the generated address is correct
	_, err = cipher.DecodeBase58Address(address)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	// use the correct address from above to check verifyAddress
	cmdArgs = []string{"verifyAddress", address}
	if cmdOut, err = exec.Command(cliName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("There was an error running  %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	fmt.Println("All tests executed successfully")
}
