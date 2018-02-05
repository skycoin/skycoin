package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"encoding/hex"
	"encoding/json"

	"github.com/skycoin/skycoin/src/cipher"
)

/**
Current list of tested commands:

generateAddresses - done
verifyAddress - done
send - done
send -m (send-to-many) - done
broadcastTransaction - done
createRawTransaction - done
getWalletBalance - done
transaction done
status - done
*/

/**
The minimal requirements for the cli are
- a wallet file
- two addresses ( to test send -m )
- the wallet address with most coins has enough coinhours for 3 transactions
*/

const cliName = "skycoin-cli"

//@TODO We can put the commands to be tested and their arguments into an array and iterate over it to keep this more DRY
// @TODO can we fetch the error from the cli command and show the exact error?
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
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	cmdArgs = []string{"walletBalance", wltFile}
	if cmdOut, err = exec.Command(cliName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	if !json.Valid(cmdOut) {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command", cmdArgs[0]))
		os.Exit(1)
	}

	cmdArgs = []string{"status"}
	if cmdOut, err = exec.Command(cliName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	if !json.Valid(cmdOut) {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command", cmdArgs[0]))
		os.Exit(1)
	}

	cmdArgs = []string{"createRawTransaction", "-f", wltFile, testAddrs[0], "0.001"}
	if cmdOut, err = exec.Command(cliName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	// validate the rawTx
	rawTx := strings.TrimSpace(string(cmdOut))
	_, err = hex.DecodeString(rawTx)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	// use the valid rawTx from above to test broadcast transaction
	cmdArgs = []string{"broadcastTransaction", rawTx}
	if cmdOut, err = exec.Command(cliName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	txId := strings.TrimSpace(string(cmdOut))
	// validate the txId
	_, err = cipher.SHA256FromHex(txId)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	// use the txId from to test transaction
	cmdArgs = []string{"transaction", txId}
	if cmdOut, err = exec.Command(cliName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	if !json.Valid(cmdOut) {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command", cmdArgs[0]))
		os.Exit(1)
	}

	cmdArgs = []string{"send", "-f", wltFile, testAddrs[0], "0.001"}
	if cmdOut, err = exec.Command(cliName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	// check that response contains the substring `txid:`
	if !strings.Contains(string(cmdOut), "txid:") {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	// send many
	cmdArgs = []string{"send", "-f", wltFile, "-m", fmt.Sprintf("'[{\"addr\":\"%s\", \"coins\": \"0.001\"}, {\"addr\":\"%s\", \"coins\": \"0.001\"}]'", testAddrs[0], testAddrs[1])}
	if cmdOut, err = exec.Command(cliName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	// check that response contains the substring `txid:`
	if !strings.Contains(string(cmdOut), "txid:") {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cmdArgs[0]), err)
		os.Exit(1)
	}

	fmt.Println("All tests executed successfully")
}
