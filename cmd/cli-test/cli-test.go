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
The minimal requirements for the cli are
- a wallet file
- two addresses ( to test send -m )
- the wallet address with most coins has enough coinhours for 3 transactions
- minimum of 3 unspents in the wallet
*/

const cliName = "skycoin-cli"

var (
	cliOut    []byte
	cliArgs   []string
	err       error
	help      bool
	wltFile   string
	addresses string // comma separated addresses, to be used in send and send -m commands, min 2 required
	testAddrs []string
)

func main() {
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
		fmt.Fprint(os.Stderr, fmt.Sprintf("wallet file %s does not exist", wltFile))
		os.Exit(1)
	}

	if addresses == "" {
		fmt.Fprint(os.Stderr, "no test addresses given")
		os.Exit(1)
	}

	testAddrs = strings.Split(addresses, ",")
	if len(testAddrs) < 2 {
		fmt.Fprint(os.Stderr, fmt.Sprintf("minimum two test addresses required, given: %v", len(testAddrs)))
		os.Exit(1)
	}

	for i := range testAddrs {
		_, err = cipher.DecodeBase58Address(testAddrs[i])
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("Address %s is invalid: ", testAddrs[i]), err)
			os.Exit(1)
		}
	}

	testCliAddressCommands()
	testCliWalletCommands()
	testCliStatusCommand()
	testCliTransactionCommands()
	testCliSendCommands()

	fmt.Println("All tests executed successfully")
}

func testCliAddressCommands() {
	cliArgs = []string{"generateAddresses", "-f", wltFile}
	if cliOut, err = exec.Command(cliName, cliArgs...).CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	address := strings.TrimSpace(string(cliOut))
	// verify that the generated address is correct
	_, err = cipher.DecodeBase58Address(address)
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	// use the correct address from above to check verifyAddress
	cliArgs = []string{"verifyAddress", address}
	if cliOut, err = exec.Command(cliName, cliArgs...).CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

}

func testCliWalletCommands() {
	cliArgs = []string{"walletBalance", wltFile}
	if cliOut, err = exec.Command(cliName, cliArgs...).CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	if !json.Valid(cliOut) {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command", cliArgs[0]))
		os.Exit(1)
	}
}

func testCliStatusCommand() {
	cliArgs = []string{"status"}
	if cliOut, err = exec.Command(cliName, cliArgs...).CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	if !json.Valid(cliOut) {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command", cliArgs[0]))
		os.Exit(1)
	}
}

func testCliTransactionCommands() {
	cliArgs = []string{"createRawTransaction", "-f", wltFile, testAddrs[0], "0.001"}
	if cliOut, err = exec.Command(cliName, cliArgs...).CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	// validate the rawTx
	rawTx := strings.TrimSpace(string(cliOut))
	_, err = hex.DecodeString(rawTx)
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	// use the valid rawTx from above to test broadcast transaction
	cliArgs = []string{"broadcastTransaction", rawTx}
	if cliOut, err = exec.Command(cliName, cliArgs...).CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	txId := strings.TrimSpace(string(cliOut))
	// validate the txId
	_, err = cipher.SHA256FromHex(txId)
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	// use the txId from to test transaction
	cliArgs = []string{"transaction", txId}
	if cliOut, err = exec.Command(cliName, cliArgs...).CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	if !json.Valid(cliOut) {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command", cliArgs[0]))
		os.Exit(1)
	}
}

func testCliSendCommands() {
	// send to a single address
	cliArgs = []string{"send", "-f", wltFile, testAddrs[0], "0.001"}
	if cliOut, err = exec.Command(cliName, cliArgs...).CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("there was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	// check that response contains the substring `txid:`
	if !strings.Contains(string(cliOut), "txid:") {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	// send many
	sendJsonMap := make([]map[string]string, len(testAddrs))
	for i := range testAddrs {
		sendJsonMap[i] = map[string]string{
			"addr":  testAddrs[i],
			"coins": "0.001",
		}
	}
	sendJson, err := json.Marshal(sendJsonMap)
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("Unable to marshal send many json string: %v", err))
	}

	cliArgs = []string{"send", "-f", wltFile, "-m", string(sendJson)}
	if cliOut, err = exec.Command(cliName, cliArgs...).CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v many command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}

	// check that response contains the substring `txid:`
	if !strings.Contains(string(cliOut), "txid:") {
		fmt.Fprint(os.Stderr, fmt.Sprintf("There was an error running %v many command: ", cliArgs[0]), string(cliOut))
		os.Exit(1)
	}
}
