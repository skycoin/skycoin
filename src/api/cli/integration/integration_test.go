package cli_integration_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

const (
	binaryName = "skycoin-cli"

	testModeStable           = "stable"
	testModeLive             = "live"
	testModeDisableWalletApi = "disable-wallet-api"

	// Number of random transactions of live transaction test.
	randomLiveTransactionNum = 500

	testFixturesDir = "test-fixtures"

	stableWalletName = "integration-test.wlt"
)

var (
	binaryPath string

	update     = flag.Bool("update", false, "update golden files")
	liveTxFull = flag.Bool("live-tx-full", false, "run live transaction test against full blockchain")
	testWallet = flag.Bool("test-wallet", false, "run wallet tests")
)

type TestData struct {
	actual   interface{}
	expected interface{}
}

func init() {
	rand.Seed(time.Now().Unix())
}

// Do setup and teardown here.
func TestMain(m *testing.M) {
	abs, err := filepath.Abs(binaryName)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("get binary name absolute path failed: %v\n", err))
		os.Exit(1)
	}

	binaryPath = abs

	// Build cli binary file.
	args := []string{"build", "-o", binaryPath, "../../../../cmd/cli/cli.go"}
	if err := exec.Command("go", args...).Run(); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Make %v binary failed: %v\n", binaryName, err))
		os.Exit(1)
	}

	ret := m.Run()

	// Remove the generated cli binary file.
	if err := os.Remove(binaryPath); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Delete %v failed: %v", binaryName, err))
		os.Exit(1)
	}

	os.Exit(ret)
}

func nodeAddress() string {
	addr := os.Getenv("SKYCOIN_NODE_HOST")
	if addr == "" {
		return "http://127.0.0.1:6420"
	}
	return addr
}

// createTempWalletFile creates a temporary dir, and copy the 'from' file to dir.
// returns the temporary wallet path, cleanup callback function, and error if any.
func createTempWalletFile(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", "wallet-data-dir")
	require.NoError(t, err)

	// Copy the testdata/$stableWalletName to the temporary dir.
	walletPath := filepath.Join(dir, stableWalletName)
	f, err := os.Create(walletPath)
	require.NoError(t, err)

	defer f.Close()

	rf, err := os.Open(filepath.Join(testFixturesDir, stableWalletName))
	require.NoError(t, err)

	defer rf.Close()
	io.Copy(f, rf)

	os.Setenv("WALLET_DIR", dir)
	os.Setenv("WALLET_NAME", stableWalletName)

	fun := func() {
		os.Setenv("WALLET_DIR", "")
		os.Setenv("WALLET_NAME", "")

		// Delete the temporary dir
		os.RemoveAll(dir)
	}

	return walletPath, fun
}

// createTempWalletDir creates a temporary wallet dir,
// sets the WALLET_DIR environment variable.
// Returns wallet dir path and callback function to clean up the dir.
func createTempWalletDir(t *testing.T) func() {
	dir, err := ioutil.TempDir("", "wallet-data-dir")
	require.NoError(t, err)
	os.Setenv("WALLET_DIR", dir)

	return func() {
		os.Setenv("WALLET_DIR", "")
		os.RemoveAll(dir)
	}
}

func loadJSON(t *testing.T, filename string, obj interface{}) {
	f, err := os.Open(filename)
	require.NoError(t, err)
	defer f.Close()

	err = json.NewDecoder(f).Decode(obj)
	require.NoError(t, err)
}

func loadGoldenFile(t *testing.T, filename string, testData TestData) {
	require.NotEmpty(t, filename, "loadGoldenFile golden filename missing")

	goldenFile := filepath.Join(testFixturesDir, filename)

	if *update {
		updateGoldenFile(t, goldenFile, testData.actual)
	}

	f, err := os.Open(goldenFile)
	require.NoError(t, err)
	defer f.Close()

	err = json.NewDecoder(f).Decode(testData.expected)
	require.NoError(t, err, filename)
}

func updateGoldenFile(t *testing.T, filename string, content interface{}) {
	contentJSON, err := json.MarshalIndent(content, "", "\t")
	require.NoError(t, err)
	contentJSON = append(contentJSON, '\n')
	err = ioutil.WriteFile(filename, contentJSON, 0644)
	require.NoError(t, err)
}

func writeJSON(t *testing.T, filename string, obj interface{}) {
	f, err := os.Create(filename)
	require.NoError(t, err)
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	require.NoError(t, enc.Encode(obj))
}

func mode(t *testing.T) string {
	mode := os.Getenv("SKYCOIN_INTEGRATION_TEST_MODE")
	switch mode {
	case "":
		mode = testModeStable
	case testModeLive, testModeStable, testModeDisableWalletApi:
	default:
		t.Fatal("Invalid test mode, must be stable, live or disable-wallet-api")
	}
	return mode
}

func enabled() bool {
	return os.Getenv("SKYCOIN_INTEGRATION_TESTS") == "1"
}

func doStable(t *testing.T) bool {
	if enabled() && mode(t) == testModeStable {
		return true
	}

	t.Skip("Stable tests disabled")
	return false
}

func doLive(t *testing.T) bool {
	if enabled() && mode(t) == testModeLive {
		return true
	}

	t.Skip("Live tests disabled")
	return false
}

func doWallet(t *testing.T) bool {
	if *testWallet {
		return true
	}

	t.Skip("Wallet tests disabled")
	return false
}

// doLiveEnvCheck checks if the WALLET_DIR and WALLET_NAME environment value do exist
func doLiveEnvCheck(t *testing.T) {
	t.Helper()
	walletDir := os.Getenv("WALLET_DIR")
	if walletDir == "" {
		t.Fatal("missing WALLET_DIR environment value")
	}

	walletName := os.Getenv("WALLET_NAME")
	if walletName == "" {
		t.Fatal("missing WALLET_NAME environment value")
	}
}

//  getWalletPathFromEnv gets wallet file path from environment variables
func getWalletPathFromEnv(t *testing.T) (string, string) {
	walletDir := os.Getenv("WALLET_DIR")
	if walletDir == "" {
		t.Fatal("missing WALLET_DIR environment value")
	}

	walletName := os.Getenv("WALLET_NAME")
	if walletName == "" {
		t.Fatal("missing WALLET_NAME environment value")
	}

	return walletDir, walletName
}

func doLiveOrStable(t *testing.T) bool {
	if enabled() {
		switch mode(t) {
		case testModeStable, testModeLive:
			return true
		}
	}

	t.Skip("Live and stable tests disabled")
	return false
}

func rpcAddress() string {
	rpcAddr := os.Getenv("RPC_ADDR")
	if rpcAddr == "" {
		rpcAddr = "127.0.0.1:6430"
	}

	return rpcAddr
}

func TestStableGenerateAddresses(t *testing.T) {
	if !doStable(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	tt := []struct {
		name         string
		args         []string
		expectOutput []byte
		goldenFile   string
	}{
		{
			"generateAddresses",
			[]string{"generateAddresses"},
			[]byte("7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\n"),
			"generate-addresses.golden",
		},
		{
			"generateAddresses -n 2 -j",
			[]string{"generateAddresses", "-n", "2", "-j"},
			[]byte("{\n    \"addresses\": [\n        \"7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\",\n        \"2EDapDfn1VC6P2hx4nTH2cRUkboGAE16evV\"\n    ]\n}\n"),
			"generate-addresses-2.golden",
		},
		{
			"generateAddresses -n -2 -j",
			[]string{"generateAddresses", "-n", "-2", "-j"},
			[]byte("Error: invalid value \"-2\" for flag -n: strconv.ParseUint: parsing \"-2\": invalid syntax"),
			"generate-addresses-2.golden",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			walletPath, clean := createTempWalletFile(t)
			defer clean()

			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			require.NoError(t, err)
			if bytes.Contains(output, []byte("Error: ")) {
				require.True(t, bytes.Contains(output, tc.expectOutput))
				return
			}

			require.Equal(t, string(tc.expectOutput), string(output))

			var w wallet.ReadableWallet
			loadJSON(t, walletPath, &w)

			// Use loadJSON instead of loadGoldenFile because this golden file
			// should not use the *update flag
			goldenFile := filepath.Join(testFixturesDir, tc.goldenFile)
			var expect wallet.ReadableWallet
			loadJSON(t, goldenFile, &expect)
			require.Equal(t, expect, w)
		})
	}
}

func TestVerifyAddress(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	tt := []struct {
		name   string
		addr   string
		err    error
		errMsg string
	}{
		{
			"valid skycoin address",
			"2Kg3eRXUhY6hrDZvNGB99DKahtrPDQ1W9vN",
			nil,
			"",
		},
		{
			"invalid skycoin address",
			"2KG9eRXUhx6hrDZvNGB99DKahtrPDQ1W9vn",
			errors.New("exit status 1"),
			"Invalid version",
		},
		{
			"invalid bitcoin address",
			"1Dcb9gpaZpBKmjqjCsiBsP3sBW1md2kEM2",
			errors.New("exit status 1"),
			"Invalid version",
		},
	}

	for _, tc := range tt {
		output, err := exec.Command(binaryPath, "verifyAddress", tc.addr).CombinedOutput()
		if err != nil {
			require.Equal(t, tc.err.Error(), err.Error())
			require.Equal(t, tc.errMsg, strings.Trim(string(output), "\n"))
			return
		}

		require.Empty(t, output)
	}
}

func TestDecodeRawTransaction(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	tt := []struct {
		name       string
		rawTx      string
		goldenFile string
		errMsg     []byte
	}{
		{
			name:       "success",
			rawTx:      "2601000000a1d3345ac47f897f24084b1c6b9bd6e03fc92887050d0748bdab5e639c1fdcd401000000a2a10f07e0e06cf6ba3e793b3186388a126591ee230b3f387617f1ccb6376a3f18e094bd3f7719aa8191c00764f323872f5192da393852bd85dab70b13409d2b01010000004d78de698a33abcfff22391c043b57a56bb0efbdc4a5b975bf8e7889668896bc0400000000bae12bbf671abeb1181fc85f1c01cdfee55deb97980c9c0a00000000543600000000000000373bb3675cbf3880bba3f3de7eb078925b8a72ad0095ba0a000000001c12000000000000008829025fe45b48f29795893a642bdaa89b2bb40e40d2df03000000001c12000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec09c0fc9b01000000001b12000000000000",
			goldenFile: "decode-raw-transaction.golden",
		},
		{
			name:   "invalid raw transaction",
			rawTx:  "2601000000a1d",
			errMsg: []byte("invalid raw transaction: encoding/hex: odd length hex string\nencoding/hex: odd length hex string\n"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, "decodeRawTransaction", tc.rawTx).CombinedOutput()
			if err != nil {
				require.Error(t, err, "exit status 1")
				require.Equal(t, tc.errMsg, output)
				return
			}

			require.NoError(t, err)
			if bytes.Contains(output, []byte("Error: ")) {
				require.Equal(t, tc.errMsg, string(output))
				return
			}

			var txn visor.TransactionJSON
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&txn)
			require.NoError(t, err)

			var expect visor.TransactionJSON
			loadGoldenFile(t, tc.goldenFile, TestData{txn, &expect})
			require.Equal(t, expect, txn)
		})
	}

}

func TestAddressGen(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	tt := []struct {
		name        string
		args        []string
		outputCheck func(t *testing.T, output []byte)
	}{
		{
			"addressGen",
			[]string{"addressGen"},
			func(t *testing.T, v []byte) {
				var w wallet.ReadableWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the wallet type is skycoin
				require.Equal(t, "skycoin", w.Meta["coin"])

				// Confirms that the seed is consisted of 12 words
				seed := w.Meta["seed"]
				require.NotEmpty(t, seed)
				ss := strings.Split(seed, " ")
				require.Len(t, ss, 12)
			},
		},
		{
			"addressGen --count 2",
			[]string{"addressGen", "--count", "2"},
			func(t *testing.T, v []byte) {
				var w wallet.ReadableWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the wallet type is skycoin
				require.Equal(t, "skycoin", w.Meta["coin"])

				// Confirms that the seed is consisted of 12 words
				seed := w.Meta["seed"]
				require.NotEmpty(t, seed)
				ss := strings.Split(seed, " ")
				require.Len(t, ss, 12)

				// Confirms that the wallet have 2 address
				require.Len(t, w.Entries, 2)

				// Confirms the addresses are generated from the seed
				_, keys := cipher.GenerateDeterministicKeyPairsSeed([]byte(seed), 2)
				for i, key := range keys {
					pk := cipher.PubKeyFromSecKey(key)
					addr := cipher.AddressFromSecKey(key)
					require.Equal(t, addr.String(), w.Entries[i].Address)
					require.Equal(t, pk.Hex(), w.Entries[i].Public)
					require.Equal(t, key.Hex(), w.Entries[i].Secret)
				}
			},
		},
		{
			"addressGen -c 2",
			[]string{"addressGen", "-c", "2"},
			func(t *testing.T, v []byte) {
				var w wallet.ReadableWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the wallet type is skycoin
				require.Equal(t, "skycoin", w.Meta["coin"])

				// Confirms that the seed is consisted of 12 words
				seed := w.Meta["seed"]
				require.NotEmpty(t, seed)
				ss := strings.Split(seed, " ")
				require.Len(t, ss, 12)

				// Confirms that the wallet have 2 address
				require.Len(t, w.Entries, 2)

				// Confirms the addresses are generated from the seed
				_, keys := cipher.GenerateDeterministicKeyPairsSeed([]byte(seed), 2)
				for i, key := range keys {
					pk := cipher.PubKeyFromSecKey(key)
					addr := cipher.AddressFromSecKey(key)
					require.Equal(t, addr.String(), w.Entries[i].Address)
					require.Equal(t, pk.Hex(), w.Entries[i].Public)
					require.Equal(t, key.Hex(), w.Entries[i].Secret)
				}
			},
		},
		{
			"addressGen --hide-secret -c 2",
			[]string{"addressGen", "--hide-secret", "-c", "2"},
			func(t *testing.T, v []byte) {
				var w wallet.ReadableWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the wallet type is skycoin
				require.Equal(t, "skycoin", w.Meta["coin"])

				// Confirms the secrets in Entries are hidden
				require.Len(t, w.Entries, 2)
				for _, e := range w.Entries {
					require.Equal(t, e.Secret, "")
				}
			},
		},
		{
			"addressGen -s -c 2",
			[]string{"addressGen", "-s", "-c", "2"},
			func(t *testing.T, v []byte) {
				var w wallet.ReadableWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the wallet type is skycoin
				require.Equal(t, "skycoin", w.Meta["coin"])

				// Confirms the secrets in Entries are hidden
				require.Len(t, w.Entries, 2)
				for _, e := range w.Entries {
					require.Equal(t, e.Secret, "")
				}
			},
		},
		{
			"addressGen --bitcoin -c 2",
			[]string{"addressGen", "--bitcoin", "-c", "2"},
			func(t *testing.T, v []byte) {
				var w wallet.ReadableWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the wallet type is skycoin
				require.Equal(t, "bitcoin", w.Meta["coin"])

				require.Len(t, w.Entries, 2)

				// Confirms the addresses are bitcoin addresses that generated from the seed
				seed := w.Meta["seed"]
				_, keys := cipher.GenerateDeterministicKeyPairsSeed([]byte(seed), 2)
				for i, key := range keys {
					pk := cipher.PubKeyFromSecKey(key)
					sk := cipher.BitcoinWalletImportFormatFromSeckey(key)
					address := cipher.BitcoinAddressFromPubkey(pk)
					require.Equal(t, address, w.Entries[i].Address)
					require.Equal(t, pk.Hex(), w.Entries[i].Public)
					require.Equal(t, sk, w.Entries[i].Secret)
				}
			},
		},
		{
			"addressGen -b -c 2",
			[]string{"addressGen", "-b", "-c", "2"},
			func(t *testing.T, v []byte) {
				var w wallet.ReadableWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the wallet type is skycoin
				require.Equal(t, "bitcoin", w.Meta["coin"])

				require.Len(t, w.Entries, 2)

				// Confirms the addresses are bitcoin addresses that generated from the seed
				seed := w.Meta["seed"]
				_, keys := cipher.GenerateDeterministicKeyPairsSeed([]byte(seed), 2)
				for i, key := range keys {
					pk := cipher.PubKeyFromSecKey(key)
					sk := cipher.BitcoinWalletImportFormatFromSeckey(key)
					address := cipher.BitcoinAddressFromPubkey(pk)
					require.Equal(t, address, w.Entries[i].Address)
					require.Equal(t, pk.Hex(), w.Entries[i].Public)
					require.Equal(t, sk, w.Entries[i].Secret)
				}
			},
		},
		{
			"addressGen --hex",
			[]string{"addressGen", "--hex"},
			func(t *testing.T, v []byte) {
				var w wallet.ReadableWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the seed is a valid hex string
				_, err = hex.DecodeString(w.Meta["seed"])
				require.NoError(t, err)
			},
		},
		{
			"addressGen -x",
			[]string{"addressGen", "-x"},
			func(t *testing.T, v []byte) {
				var w wallet.ReadableWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the seed is a valid hex string
				_, err = hex.DecodeString(w.Meta["seed"])
				require.NoError(t, err)
			},
		},
		{
			"addressGen --only-addr",
			[]string{"addressGen", "--only-addr"},
			func(t *testing.T, v []byte) {
				// Confirms that only addresses are returned
				v = bytes.Trim(v, "\n")
				_, err := cipher.DecodeBase58Address(string(v))
				require.NoError(t, err)
			},
		},
		{
			"addressGen --oa",
			[]string{"addressGen", "--oa"},
			func(t *testing.T, v []byte) {
				// Confirms that only addresses are returned
				v = bytes.Trim(v, "\n")
				_, err := cipher.DecodeBase58Address(string(v))
				require.NoError(t, err)
			},
		},
		{
			"addressGen --seed=123",
			[]string{"addressGen", "--seed", "123"},
			func(t *testing.T, v []byte) {
				var w wallet.ReadableWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				pk, sk := cipher.GenerateDeterministicKeyPair([]byte("123"))
				addr := cipher.AddressFromPubKey(pk)
				require.Len(t, w.Entries, 1)
				require.Equal(t, addr.String(), w.Entries[0].Address)
				require.Equal(t, pk.Hex(), w.Entries[0].Public)
				require.Equal(t, sk.Hex(), w.Entries[0].Secret)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			require.NoError(t, err)
			tc.outputCheck(t, output)
		})
	}
}

func TestStableListWallets(t *testing.T) {
	if !doStable(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	_, clean := createTempWalletFile(t)
	defer clean()

	output, err := exec.Command(binaryPath, "listWallets").CombinedOutput()
	require.NoError(t, err)

	var wlts struct {
		Wallets []cli.WalletEntry `json:"wallets"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wlts)
	require.NoError(t, err)

	var expect struct {
		Wallets []cli.WalletEntry `json:"wallets"`
	}
	loadGoldenFile(t, "list-wallets.golden", TestData{wlts, &expect})
	require.Equal(t, expect, wlts)
}

func TestLiveListWallets(t *testing.T) {
	if !doLive(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	doLiveEnvCheck(t)

	output, err := exec.Command(binaryPath, "listWallets").CombinedOutput()
	require.NoError(t, err)

	var wlts struct {
		Wallets []cli.WalletEntry `json:"wallets"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wlts)
	require.NoError(t, err)
}

func TestStableListAddress(t *testing.T) {
	if !doStable(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	_, clean := createTempWalletFile(t)
	defer clean()

	output, err := exec.Command(binaryPath, "listAddresses").CombinedOutput()
	require.NoError(t, err)

	var wltAddresses struct {
		Addresses []string `json:"addresses"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltAddresses)
	require.NoError(t, err)

	var expect struct {
		Addresses []string `json:"addresses"`
	}
	loadGoldenFile(t, "list-addresses.golden", TestData{wltAddresses, &expect})
	require.Equal(t, expect, wltAddresses)
}

func TestLiveListAddresses(t *testing.T) {
	if !doLive(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	doLiveEnvCheck(t)

	output, err := exec.Command(binaryPath, "listAddresses").CombinedOutput()
	require.NoError(t, err)

	var wltAddresses struct {
		Addresses []string `json:"addresses"`
	}

	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltAddresses)
	require.NoError(t, err)
}

func TestStableAddressBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "addressBalance", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt").CombinedOutput()
	require.NoError(t, err)

	var addrBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrBalance)
	require.NoError(t, err)

	var expect cli.BalanceResult
	loadGoldenFile(t, "address-balance.golden", TestData{addrBalance, &expect})
	require.Equal(t, expect, addrBalance)
}

func TestLiveAddressBalance(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := exec.Command(binaryPath, "addressBalance", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt").CombinedOutput()
	require.NoError(t, err)

	var addrBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrBalance)
	require.NoError(t, err)
}

func TestStableWalletBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	_, clean := createTempWalletFile(t)
	defer clean()

	output, err := exec.Command(binaryPath, "walletBalance").CombinedOutput()
	require.NoError(t, err)

	var wltBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltBalance)
	require.NoError(t, err)

	var expect cli.BalanceResult
	loadGoldenFile(t, "wallet-balance.golden", TestData{wltBalance, &expect})
	require.Equal(t, expect, wltBalance)
}

func TestLiveWalletBalance(t *testing.T) {
	if !doLive(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	doLiveEnvCheck(t)

	output, err := exec.Command(binaryPath, "walletBalance").CombinedOutput()
	require.NoError(t, err)

	var wltBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltBalance)
	require.NoError(t, err)
}

func TestStableWalletOutputs(t *testing.T) {
	if !doStable(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	_, clean := createTempWalletFile(t)
	defer clean()

	output, err := exec.Command(binaryPath, "walletOutputs").CombinedOutput()
	require.NoError(t, err)

	var wltOutput webrpc.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltOutput)
	require.NoError(t, err)

	var expect webrpc.OutputsResult
	loadGoldenFile(t, "wallet-outputs.golden", TestData{wltOutput, &expect})
	require.Equal(t, expect, wltOutput)
}

func TestLiveWalletOutputs(t *testing.T) {
	if !doLive(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	doLiveEnvCheck(t)

	output, err := exec.Command(binaryPath, "walletOutputs").CombinedOutput()
	require.NoError(t, err)

	var wltOutput webrpc.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltOutput)
	require.NoError(t, err)
}

func TestStableAddressOutputs(t *testing.T) {
	if !doStable(t) {
		return
	}

	tt := []struct {
		name       string
		args       []string
		goldenFile string
	}{
		{
			"addressOutputs one address",
			[]string{"addressOutputs", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
			"address-outputs.golden",
		},
		{
			"addressOutputs two address",
			[]string{"addressOutputs", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt", "ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm"},
			"two-addresses-outputs.golden",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			require.NoError(t, err)

			var addrOutputs webrpc.OutputsResult
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrOutputs)
			require.NoError(t, err)

			var expect webrpc.OutputsResult
			loadGoldenFile(t, tc.goldenFile, TestData{addrOutputs, &expect})
			require.Equal(t, expect, addrOutputs)

		})
	}
}

func TestLiveAddressOutputs(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := exec.Command(binaryPath, "addressOutputs", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt").CombinedOutput()
	require.NoError(t, err)

	var addrOutputs webrpc.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrOutputs)
	require.NoError(t, err)
}

func TestStableStatus(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "status").CombinedOutput()
	require.NoError(t, err)
	var ret struct {
		webrpc.StatusResult
		RPCAddress string `json:"webrpc_address"`
	}

	err = json.NewDecoder(bytes.NewReader(output)).Decode(&ret)
	require.NoError(t, err)

	// TimeSinceLastBlock is not stable
	ret.TimeSinceLastBlock = ""

	var expect struct {
		webrpc.StatusResult
		RPCAddress string `json:"webrpc_address"`
	}

	loadGoldenFile(t, "status.golden", TestData{ret, &expect})
	require.Equal(t, expect, ret)
}

func TestLiveStatus(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := exec.Command(binaryPath, "status").CombinedOutput()
	require.NoError(t, err)

	var ret struct {
		webrpc.StatusResult
		RPCAddress string `json:"webrpc_address"`
	}

	err = json.NewDecoder(bytes.NewReader(output)).Decode(&ret)
	require.NoError(t, err)
	require.True(t, ret.Running)
	require.Equal(t, ret.RPCAddress, rpcAddress())
}

func TestStableTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	tt := []struct {
		name       string
		args       []string
		err        error
		errMsg     string
		goldenFile string
	}{
		{
			"invalid txid",
			[]string{"abcd"},
			errors.New("exit status 1"),
			"invalid txid\n",
			"",
		},
		{
			"not exist",
			[]string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			errors.New("exit status 1"),
			"transaction doesn't exist [code: -32600]\n",
			"",
		},
		{
			"empty txid",
			[]string{""},
			errors.New("exit status 1"),
			"txid is empty\n",
			"",
		},
		{
			"genesis transaction",
			[]string{"d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add"},
			nil,
			"",
			"genesis-transaction.golden",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"transaction"}, tc.args...)
			o, err := exec.Command(binaryPath, args...).CombinedOutput()
			if err != nil {
				require.Equal(t, tc.err.Error(), err.Error())
				require.Equal(t, tc.errMsg, string(o))
				return
			}

			// Decode the output into visor.TransactionJSON
			var tx webrpc.TxnResult
			err = json.NewDecoder(bytes.NewReader(o)).Decode(&tx)
			require.NoError(t, err)

			var expect webrpc.TxnResult
			loadGoldenFile(t, tc.goldenFile, TestData{tx, &expect})

			require.Equal(t, expect, tx)
		})
	}

	scanTransactions(t, true)
}

func TestLiveTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	o, err := exec.Command(binaryPath, "transaction", "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add").CombinedOutput()
	require.NoError(t, err)
	var tx webrpc.TxnResult
	err = json.NewDecoder(bytes.NewReader(o)).Decode(&tx)
	require.NoError(t, err)

	var expect webrpc.TxnResult

	loadGoldenFile(t, "genesis-transaction.golden", TestData{tx, &expect})
	require.Equal(t, expect.Transaction.Transaction, tx.Transaction.Transaction)

	scanTransactions(t, *liveTxFull)

	// scan pending transactions
	scanPendingTransactions(t)
}

// cli doesn't have command to querying pending transactions yet.
func scanPendingTransactions(t *testing.T) {
}

// scanTransactions scans transactions against blockchain.
// If fullTest is true, scan the whole blockchain, and test every transactions,
// otherwise just test random transactions.
func scanTransactions(t *testing.T, fullTest bool) {
	// Gets blockchain height through "status" command
	output, err := exec.Command(binaryPath, "status").CombinedOutput()
	require.NoError(t, err)
	var status struct {
		webrpc.StatusResult
		RPCAddress string `json:"webrpc_address"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&status)
	require.NoError(t, err)

	txids := getTxids(t, status.BlockNum)

	l := len(txids)
	if !fullTest && l > randomLiveTransactionNum {
		txidMap := make(map[string]struct{})
		var ids []string
		for len(txidMap) < randomLiveTransactionNum {
			// get random txid
			txid := txids[rand.Intn(l)]
			if _, ok := txidMap[txid]; !ok {
				ids = append(ids, txid)
				txidMap[txid] = struct{}{}
			}
		}

		// reassign the txids
		txids = ids
	}

	checkTransactions(t, txids)
}

func checkTransactions(t *testing.T, txids []string) {
	// Start goroutines to check transactions
	var wg sync.WaitGroup
	txC := make(chan string, 500)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case txid, ok := <-txC:
					if !ok {
						return
					}

					t.Run(fmt.Sprintf("%v", txid), func(t *testing.T) {
						o, err := exec.Command(binaryPath, "transaction", txid).CombinedOutput()
						require.NoError(t, err)
						var txRlt webrpc.TxnResult
						err = json.NewDecoder(bytes.NewReader(o)).Decode(&txRlt)
						require.NoError(t, err)
						require.Equal(t, txid, txRlt.Transaction.Transaction.Hash)
						require.True(t, txRlt.Transaction.Status.Confirmed)
					})
				}
			}
		}()
	}

	for _, txid := range txids {
		txC <- txid
	}
	close(txC)

	wg.Wait()
}

func getTxids(t *testing.T, blockNum uint64) []string {
	// p represents the number of blocks that each time we query,
	// do not get all blocks in one query, which might run out of
	// memory when blockchain becomes very huge.
	p := 500
	n := int(blockNum / uint64(p))

	// Collects all transactions' id
	var txids []string
	for i := 0; i < int(n); i++ {
		txids = append(txids, getTxidsInBlocks(t, i*p+1, (i+1)*p)...)
	}

	if (blockNum % uint64(p)) > 0 {
		txids = append(txids, getTxidsInBlocks(t, n*p+1, int(blockNum)-1)...)
	}

	return txids
}

func getTxidsInBlocks(t *testing.T, start, end int) []string {
	s := strconv.Itoa(start)
	e := strconv.Itoa(end)
	o, err := exec.Command(binaryPath, "blocks", s, e).CombinedOutput()
	require.NoError(t, err)
	var blocks visor.ReadableBlocks
	err = json.NewDecoder(bytes.NewReader(o)).Decode(&blocks)
	require.NoError(t, err)
	require.Len(t, blocks.Blocks, end-start+1)

	var txids []string
	for _, b := range blocks.Blocks {
		for _, tx := range b.Body.Transactions {
			txids = append(txids, tx.Hash)
		}
	}
	return txids
}

func TestStableBlocks(t *testing.T) {
	if !doStable(t) {
		return
	}

	testKnownBlocks(t)

	// Tests blocks 180~181, should only return block 180.
	output, err := exec.Command(binaryPath, "blocks", "180", "181").CombinedOutput()
	require.NoError(t, err)

	var blocks visor.ReadableBlocks
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&blocks)
	require.NoError(t, err)

	var expect visor.ReadableBlocks
	loadGoldenFile(t, "blocks180.golden", TestData{blocks, &expect})
	require.Equal(t, expect, blocks)
}

func TestLiveBlocks(t *testing.T) {
	if !doLive(t) {
		return
	}

	testKnownBlocks(t)

	// These blocks were affected by the coinhour overflow issue, make sure that they can be queried
	blockSeqs := []int{11685, 11707, 11710, 11709, 11705, 11708, 11711, 11706, 11699}

	for _, seq := range blockSeqs {
		output, err := exec.Command(binaryPath, "blocks", strconv.Itoa(seq)).CombinedOutput()
		require.NoError(t, err)
		var blocks visor.ReadableBlocks
		err = json.NewDecoder(bytes.NewReader(output)).Decode(&blocks)
		require.NoError(t, err)
	}
}

func testKnownBlocks(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		goldenFile string
	}{
		{
			"blocks 0",
			[]string{"blocks", "0"},
			"block0.golden",
		},
		{
			"blocks 0 5",
			[]string{"blocks", "0", "5"},
			"blocks0-5.golden",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			require.NoError(t, err)

			var blocks visor.ReadableBlocks
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&blocks)
			require.NoError(t, err)

			var expect visor.ReadableBlocks
			loadGoldenFile(t, tc.goldenFile, TestData{blocks, &expect})
			require.Equal(t, expect, blocks)
		})
	}

	scanBlocks(t, "0", "180")
}

func scanBlocks(t *testing.T, s, e string) {
	outputs, err := exec.Command(binaryPath, "blocks", s, e).CombinedOutput()
	require.NoError(t, err)

	var blocks visor.ReadableBlocks
	err = json.NewDecoder(bytes.NewReader(outputs)).Decode(&blocks)
	require.NoError(t, err)

	var preBlocks visor.ReadableBlock
	preBlocks.Head.BlockHash = "0000000000000000000000000000000000000000000000000000000000000000"
	for _, b := range blocks.Blocks {
		require.Equal(t, b.Head.PreviousBlockHash, preBlocks.Head.BlockHash)
		preBlocks = b
	}
}

func TestStableLastBlocks(t *testing.T) {
	if !doStable(t) {
		return
	}

	tt := []struct {
		name       string
		args       []string
		goldenFile string
		errMsg     []byte
	}{
		{
			name:       "lastBlocks 0",
			args:       []string{"lastBlocks", "0"},
			goldenFile: "last-blocks0.golden",
		},
		{
			name:       "lastBlocks 1",
			args:       []string{"lastBlocks", "1"},
			goldenFile: "last-blocks1.golden",
		},
		{
			name:       "lastBlocks 2",
			args:       []string{"lastBlocks", "2"},
			goldenFile: "last-blocks2.golden",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()

			if bytes.Contains(output, []byte("Error: ")) {
				fmt.Println(string(output))
				require.Equal(t, string(tc.errMsg), string(output))
				return
			}

			require.NoError(t, err)

			var blocks visor.ReadableBlocks
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&blocks)
			require.NoError(t, err)

			var expect visor.ReadableBlocks
			loadGoldenFile(t, tc.goldenFile, TestData{blocks, &expect})
			require.Equal(t, expect, blocks)
		})
	}
}

func TestLiveLastBlocks(t *testing.T) {
	if !doLive(t) {
		return
	}

	tt := []struct {
		name string
		args []string
	}{
		{
			"lastBlocks 0",
			[]string{"lastBlocks", "0"},
		},
		{
			"lastBlocks 1",
			[]string{"lastBlocks", "1"},
		},
		{
			"lastBlocks 2",
			[]string{"lastBlocks", "2"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			require.NoError(t, err)

			var blocks visor.ReadableBlocks
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&blocks)
			require.NoError(t, err)
		})
	}
}

func TestStableWalletDir(t *testing.T) {
	if !doStable(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	walletPath, clean := createTempWalletFile(t)
	defer clean()

	dir := filepath.Dir(walletPath)
	output, err := exec.Command(binaryPath, "walletDir").CombinedOutput()
	require.NoError(t, err)
	require.Equal(t, dir, strings.TrimRight(string(output), "\n"))
}

func TestLiveWalletDir(t *testing.T) {
	if !doLive(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	doLiveEnvCheck(t)

	walletDir := os.Getenv("WALLET_DIR")
	output, err := exec.Command(binaryPath, "walletDir").CombinedOutput()
	require.NoError(t, err)

	require.Equal(t, walletDir, strings.Trim(string(output), "\n"))
}

// TestLiveSend sends coin from specific wallet file, user should manually specify the
// wallet file by setting the enviroment variables: WALLET_DIR and WALLET_NAME. The WALLET_DIR
// points to the directory of the wallet, and WALLET_NAME represents the wallet file name.
//
// Note:
// 1. This test might modify the wallet file, in order to avoid losing coins, we don't send coins to
// addresses that are not belong to the wallet, when addresses in the wallet are not sufficient, we
// will automatically generate enough addresses as coin recipient.
// 2. The wallet must must have at least 2 coins and 16 coinhours.
func TestLiveSend(t *testing.T) {
	if !doLive(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	// prepares wallet and confirms the wallet has at least 2 coins and 16 coin hours.
	w, totalCoins, _ := prepareAndCheckWallet(t, 2e6, 16)

	tt := []struct {
		name    string
		args    func() []string
		errMsg  []byte
		checkTx func(t *testing.T, txid string)
	}{
		{
			// Send all coins to the first address to one output.
			name: "send all coins to the first address",
			args: func() []string {
				coins, err := droplet.ToString(totalCoins)
				require.NoError(t, err)
				return []string{"send", w.Entries[0].Address.String(), coins}
			},
			checkTx: func(t *testing.T, txid string) {
				// Confirms all coins are in the first address in one output
				tx := getTransaction(t, txid)
				require.Len(t, tx.Transaction.Transaction.Out, 1)
				c, err := droplet.FromString(tx.Transaction.Transaction.Out[0].Coins)
				require.NoError(t, err)
				require.Equal(t, totalCoins, c)
			},
		},
		{
			// Send 0.5 coin to the second address.
			// Send 0.5 coin to the third address.
			// After sending, the first address should have at least 1 coin left.
			name: "send to multiple address with -m option",
			args: func() []string {
				addrCoins := []struct {
					Addr  string `json:"addr"`
					Coins string `json:"coins"`
				}{
					{
						w.Entries[1].Address.String(),
						"0.5",
					},
					{
						w.Entries[2].Address.String(),
						"0.5",
					},
				}

				v, err := json.Marshal(addrCoins)
				require.NoError(t, err)

				return []string{"send", "-m", string(v)}
			},
			checkTx: func(t *testing.T, txid string) {
				tx := getTransaction(t, txid)
				// Confirms the second address receives 0.5 coin and 1 coinhour in this transaction
				checkCoinsAndCoinhours(t, tx, w.Entries[1].Address.String(), 5e5, 1)
				// Confirms the third address receives 0.5 coin and 1 coinhour in this transaction
				checkCoinsAndCoinhours(t, tx, w.Entries[2].Address.String(), 5e5, 1)
				// Confirms the first address has at least 1 coin left.
				coins, _ := getAddressBalance(t, w.Entries[0].Address.String())
				require.True(t, coins >= 1e6)
			},
		},
		{
			// Send 0.001 coin from the third address to the second address.
			// Set the second as change address, so the 0.499 change coin will also be sent to the second address.
			// After sending, the second address should have 1 coin and 1 coin hour.
			name: "send with -c(change address) -a(from address) options",
			args: func() []string {
				return []string{"send", "-c", w.Entries[1].Address.String(),
					"-a", w.Entries[2].Address.String(), w.Entries[1].Address.String(), "0.001"}
			},
			checkTx: func(t *testing.T, txid string) {
				tx := getTransaction(t, txid)
				// Confirms the second address receives 0.5 coin and 0 coinhour in this transaction
				checkCoinsAndCoinhours(t, tx, w.Entries[1].Address.String(), 5e5, 0)
				// Confirms the second address have 1 coin and 1 coin hour
				coins, hours := getAddressBalance(t, w.Entries[1].Address.String())
				require.Equal(t, uint64(1e6), coins)
				require.Equal(t, uint64(1), hours)
			},
		},
		{
			// Send 1 coin from second to the the third address, this will spend three outputs(0.2, 0.3. 0.5 coin),
			// and burn out the remaining 1 coin hour.
			name: "send to burn all coin hour",
			args: func() []string {
				return []string{"send", "-a", w.Entries[1].Address.String(),
					w.Entries[2].Address.String(), "1"}
			},
			checkTx: func(t *testing.T, txid string) {
				// Confirms that the third address has 1 coin and 0 coin hour
				coins, hours := getAddressBalance(t, w.Entries[2].Address.String())
				require.Equal(t, uint64(1e6), coins)
				require.Equal(t, uint64(0), hours)
			},
		},
		{
			// Send with 0 coin hour, this test should fail.
			name: "send 0 coin hour",
			args: func() []string {
				return []string{"send", "-a", w.Entries[2].Address.String(),
					w.Entries[1].Address.String(), "1"}
			},
			errMsg:  []byte("ERROR: Transaction has zero coinhour fee. See 'skycoin-cli send --help'"),
			checkTx: func(t *testing.T, txid string) {},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args()...).CombinedOutput()
			if err != nil {
				t.Fatalf("err: %v, output: %v", err, string(output))
				return
			}
			require.NoError(t, err)
			output = bytes.TrimRight(output, "\n")
			if bytes.Contains(output, []byte("ERROR:")) {
				require.Equal(t, tc.errMsg, output)
				return
			}

			// output: "txid:$txid_string"
			// split the output to get txid value
			v := bytes.Split(output, []byte(":"))
			require.Len(t, v, 2)
			txid := string(v[1])
			fmt.Println("txid:", txid)
			_, err = cipher.SHA256FromHex(txid)
			require.NoError(t, err)

			// Wait untill transaction is confirmed.
			tk := time.NewTicker(time.Second)
		loop:
			for {
				select {
				case <-time.After(30 * time.Second):
					t.Fatal("Wait tx confirmation timeout")
				case <-tk.C:
					if isTxConfirmed(t, txid) {
						break loop
					}
				}
			}

			tc.checkTx(t, txid)
		})
	}

	// Send with too small decimal value
	// CLI send is a litte bit slow, almost 300ms each. so we only test 20 invalid decimal coin.
	errMsg := []byte("ERROR: invalid amount, too many decimal places. See 'skycoin-cli send --help'")
	for i := uint64(1); i < uint64(20); i++ {
		v, err := droplet.ToString(i)
		require.NoError(t, err)
		name := fmt.Sprintf("send %v", v)
		t.Run(name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, "send", w.Entries[0].Address.String(), v).CombinedOutput()
			require.NoError(t, err)
			output = bytes.Trim(output, "\n")
			require.Equal(t, errMsg, output)
		})
	}
}

// TestLiveCreateAndBroadcastRawTransaction does almost the same procedure as TestLiveSend.
// Create raw transaction with command arguments the same as TestLiveSend, then broadcast the
// created raw transaction. After the transaction is confirmed, run the same transaction check
// function like in TestLiveSend.
func TestLiveCreateAndBroadcastRawTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	// prepares wallet and confirms the wallet has at least 2 coins and 2 coin hours.
	w, totalCoins, _ := prepareAndCheckWallet(t, 2e6, 2)

	tt := []struct {
		name    string
		args    func() []string
		errMsg  []byte
		checkTx func(t *testing.T, txid string)
	}{
		{
			// Send all coins to the first address to one output.
			name: "send all coins to the first address",
			args: func() []string {
				coins, err := droplet.ToString(totalCoins)
				require.NoError(t, err)
				return []string{"createRawTransaction", w.Entries[0].Address.String(), coins}
			},
			checkTx: func(t *testing.T, txid string) {
				// Confirms all coins are in the first address in one output
				tx := getTransaction(t, txid)
				require.Len(t, tx.Transaction.Transaction.Out, 1)
				c, err := droplet.FromString(tx.Transaction.Transaction.Out[0].Coins)
				require.NoError(t, err)
				require.Equal(t, totalCoins, c)
			},
		},
		{
			// Send 0.5 coin to the second address.
			// Send 0.5 coin to the third address.
			// After sending, the first address should have at least 1 coin left.
			name: "send to multiple address with -m option",
			args: func() []string {
				addrCoins := []struct {
					Addr  string `json:"addr"`
					Coins string `json:"coins"`
				}{
					{
						w.Entries[1].Address.String(),
						"0.5",
					},
					{
						w.Entries[2].Address.String(),
						"0.5",
					},
				}

				v, err := json.Marshal(addrCoins)
				require.NoError(t, err)

				return []string{"createRawTransaction", "-m", string(v)}
			},
			checkTx: func(t *testing.T, txid string) {
				// Confirms the first address has at least 1 coin left.
				coins, _ := getAddressBalance(t, w.Entries[0].Address.String())
				require.True(t, coins >= 1e6)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Create raw transaction first
			output, err := exec.Command(binaryPath, tc.args()...).CombinedOutput()
			if err != nil {
				t.Fatalf("err: %v, output: %v", err, string(output))
				return
			}
			require.NoError(t, err)
			output = bytes.TrimRight(output, "\n")
			if bytes.Contains(output, []byte("ERROR:")) {
				require.Equal(t, tc.errMsg, output)
				return
			}

			// Broadcast transaction
			output, err = exec.Command(binaryPath, "broadcastTransaction", string(output)).CombinedOutput()
			require.NoError(t, err)

			txid := string(bytes.TrimRight(output, "\n"))
			fmt.Println("txid:", txid)
			_, err = cipher.SHA256FromHex(txid)
			require.NoError(t, err)

			// Wait untill transaction is confirmed.
			tk := time.NewTicker(time.Second)
		loop:
			for {
				select {
				case <-time.After(30 * time.Second):
					t.Fatal("Wait tx confirmation timeout")
				case <-tk.C:
					if isTxConfirmed(t, txid) {
						break loop
					}
				}
			}

			tc.checkTx(t, txid)
		})
	}

	// Send with too small decimal value
	errMsg := []byte("ERROR: invalid amount, too many decimal places. See 'skycoin-cli createRawTransaction --help'")
	for i := uint64(1); i < uint64(20); i++ {
		v, err := droplet.ToString(i)
		require.NoError(t, err)
		name := fmt.Sprintf("send %v", v)
		t.Run(name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, "createRawTransaction", w.Entries[0].Address.String(), v).CombinedOutput()
			require.NoError(t, err)
			output = bytes.Trim(output, "\n")
			require.Equal(t, errMsg, output)
		})
	}
}

func getTransaction(t *testing.T, txid string) *webrpc.TxnResult {
	output, err := exec.Command(binaryPath, "transaction", txid).CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return &webrpc.TxnResult{}
	}
	require.NoError(t, err)

	var tx webrpc.TxnResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&tx)
	require.NoError(t, err)

	return &tx
}

func isTxConfirmed(t *testing.T, txid string) bool {
	tx := getTransaction(t, txid)
	return tx.Transaction.Status.Confirmed
}

// checkCoinhours checks if the address coinhours in transaction are correct
func checkCoinsAndCoinhours(t *testing.T, tx *webrpc.TxnResult, addr string, coins, coinhours uint64) {
	addrCoinhoursMap := make(map[string][]visor.ReadableTransactionOutput)
	for _, o := range tx.Transaction.Transaction.Out {
		addrCoinhoursMap[o.Address] = append(addrCoinhoursMap[o.Address], o)
	}

	os, ok := addrCoinhoursMap[addr]
	if !ok {
		t.Fatalf("transaction doesn't have receiver of address: %v", addr)
	}

	var totalCoins, totalHours uint64
	for _, o := range os {
		c, err := droplet.FromString(o.Coins)
		if err != nil {
			t.Fatalf("%v", err)
		}
		totalCoins += c
		totalHours += o.Hours
	}

	require.Equal(t, coins, totalCoins)
	require.Equal(t, coinhours, totalHours)
}

// prepareAndCheckWallet prepares wallet for live testing.
// Returns *wallet.Wallet, total coin, total hours.
// Confirms that the wallet meets the minimal requirements of coins and coinhours.
func prepareAndCheckWallet(t *testing.T, miniCoins, miniCoinHours uint64) (*wallet.Wallet, uint64, uint64) {
	walletDir, walletName := getWalletPathFromEnv(t)
	walletPath := filepath.Join(walletDir, walletName)
	// Checks if the wallet does exist
	if _, err := os.Stat(walletPath); os.IsNotExist(err) {
		t.Fatalf("Wallet file: %v does not exist", walletPath)
	}

	// Loads the wallet
	w, err := wallet.Load(walletPath)
	if err != nil {
		t.Fatalf("Load wallet failed: %v", err)
	}

	if len(w.Entries) < 3 {
		// Generates addresses
		_, err = w.GenerateAddresses(uint64(3 - len(w.Entries)))
		if err != nil {
			t.Fatalf("Wallet generateAddress failed: %v", err)
		}
	}

	outputs := getWalletOutputs(t, walletPath)
	// Confirms the wallet is not empty.
	if len(outputs) == 0 {
		t.Fatalf("Wallet %v has no coin", walletPath)
	}

	var totalCoins uint64
	var totalCoinhours uint64
	for _, output := range outputs {
		coins, err := droplet.FromString(output.Coins)
		if err != nil {
			t.Fatalf("%v", err)
		}

		totalCoins += coins
		totalCoinhours += output.CalculatedHours
	}

	// Confirms the coins meet minimal coins requirement
	if totalCoins < miniCoins {
		t.Fatalf("Wallet must have at least %v coins", miniCoins)
	}

	if totalCoinhours < miniCoinHours {
		t.Fatalf("Wallet must have at least %v coinhours", miniCoinHours)
	}

	if err := w.Save(walletDir); err != nil {
		t.Fatalf("%v", err)
	}
	return w, totalCoins, totalCoinhours
}

func getAddressBalance(t *testing.T, addr string) (uint64, uint64) {
	output, err := exec.Command(binaryPath, "addressBalance", addr).CombinedOutput()
	require.NoError(t, err)

	var addrBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrBalance)
	require.NoError(t, err)
	coins, err := droplet.FromString(addrBalance.Confirmed.Coins)
	require.NoError(t, err)

	hours, err := strconv.ParseUint(addrBalance.Confirmed.Hours, 10, 64)
	require.NoError(t, err)
	return coins, hours
}

func getWalletOutputs(t *testing.T, walletPath string) visor.ReadableOutputs {
	output, err := exec.Command(binaryPath, "walletOutputs", walletPath).CombinedOutput()
	require.NoError(t, err)

	var wltOutput webrpc.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltOutput)
	require.NoError(t, err)

	return wltOutput.Outputs.HeadOutputs
}

func TestStableWalletHistory(t *testing.T) {
	if !doStable(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	_, clean := createTempWalletFile(t)
	defer clean()

	output, err := exec.Command(binaryPath, "walletHistory").CombinedOutput()
	require.NoError(t, err)

	var history []cli.AddrHistory
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&history)
	require.NoError(t, err)

	var expect []cli.AddrHistory
	loadGoldenFile(t, "wallet-history.golden", TestData{history, &expect})
	require.Equal(t, expect, history)
}

func TestLiveWalletHistory(t *testing.T) {
	if !doLive(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	doLiveEnvCheck(t)

	output, err := exec.Command(binaryPath, "walletHistory").CombinedOutput()
	require.NoError(t, err)
	var his []cli.AddrHistory
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&his)
	require.NoError(t, err)
}

func TestStableCheckDB(t *testing.T) {
	if !doStable(t) {
		return
	}

	tt := []struct {
		name   string
		dbPath string
		result []byte
		errMsg []byte
	}{
		{
			name:   "no signature",
			dbPath: "../../../visor/testdata/data.db.nosig",
			errMsg: []byte("checkdb failed: find no signature of block: seq=1000\n"),
		},
		{
			name:   "invalid database",
			dbPath: "../../../visor/testdata/data.db.garbage",
			errMsg: []byte("open db failed: invalid database\n"),
		},
		{
			name:   "valid database",
			dbPath: "../../../gui/integration/test-fixtures/blockchain-180.db",
			result: []byte("check db success\n"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, "checkdb", tc.dbPath).CombinedOutput()
			if err != nil {
				fmt.Println(string(output))
				require.Equal(t, tc.errMsg, output)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.result, output)
		})
	}
}

func TestVersion(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	// Gets version in json format.
	output, err := exec.Command(binaryPath, "version", "-j").CombinedOutput()
	require.NoError(t, err)

	var ver = struct {
		Skycoin string `json:"skycoin"`
		Cli     string `json:"cli"`
		RPC     string `json:"rpc"`
		Wallet  string `json:"wallet"`
	}{}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&ver)
	require.NoError(t, err)
	require.True(t, ver.Skycoin != "")
	require.True(t, ver.Cli != "")
	require.True(t, ver.RPC != "")
	require.True(t, ver.Wallet != "")

	// Gets version without json format.
	output, err = exec.Command(binaryPath, "version").CombinedOutput()
	require.NoError(t, err)

	// Confirms the result contains 4 version componments
	output = bytes.TrimRight(output, "\n")
	vers := bytes.Split(output, []byte("\n"))
	require.Len(t, vers, 4)
}

func TestStableGenerateWallet(t *testing.T) {
	if !doStable(t) {
		return
	}

	if !doWallet(t) {
		return
	}

	tt := []struct {
		name        string
		args        []string
		setup       func(t *testing.T) func()
		errMsg      []byte
		checkWallet func(t *testing.T, w *wallet.Wallet)
	}{
		{
			name:  "generate wallet with -r option",
			args:  []string{"-r"},
			setup: createTempWalletDir,
			checkWallet: func(t *testing.T, w *wallet.Wallet) {
				// Confirms the default wallet name is skycoin_cli.wlt
				require.Equal(t, "skycoin_cli.wlt", w.GetFilename())

				// Confirms the seed is a valid hex string
				_, err := hex.DecodeString(w.Meta["seed"])
				require.NoError(t, err)

				// Confirms the label is empty
				require.Empty(t, w.Meta["label"])
			},
		},
		{
			name:  "generate wallet with --rd option",
			args:  []string{"--rd"},
			setup: createTempWalletDir,
			checkWallet: func(t *testing.T, w *wallet.Wallet) {
				// Confirms the default wallet name is skycoin_cli.wlt
				require.Equal(t, "skycoin_cli.wlt", w.GetFilename())

				// Confirms the seed is consisited of 12 words
				seed := w.Meta["seed"]
				words := strings.Split(seed, " ")
				require.Len(t, words, 12)

				// Confirms the label is empty
				require.Empty(t, w.Meta["label"])
			},
		},
		{
			name:  "generate wallet with -s option",
			args:  []string{"-s", "great duck trophy inhale dad pluck include maze smart mechanic ring merge"},
			setup: createTempWalletDir,
			checkWallet: func(t *testing.T, w *wallet.Wallet) {
				// Confirms the default wallet name is skycoin_cli.wlt
				require.Equal(t, "skycoin_cli.wlt", w.GetFilename())
				// Confirms the label is empty
				require.Empty(t, w.Meta["label"])

				require.Equal(t, "great duck trophy inhale dad pluck include maze smart mechanic ring merge", w.Meta["seed"])
				require.Equal(t, "2amA8sxKJhNRp3wfWrE5JfTEUjr9S3C2BaU", w.Entries[0].Address.String())
				require.Equal(t, "02b4a4b63f2f8ba56f9508712815eca3c088693333715eaf7a73275d8928e1be5a", w.Entries[0].Public.Hex())
				require.Equal(t, "f4a281d094a6e9e95a84c23701a7d01a0e413c838758e94ad86a10b9b83e0434", w.Entries[0].Secret.Hex())
			},
		},
		{
			name:  "generate wallet with -n option",
			args:  []string{"-n", "5"},
			setup: createTempWalletDir,
			checkWallet: func(t *testing.T, w *wallet.Wallet) {
				// Confirms the default wallet name is skycoin_cli.wlt
				require.Equal(t, "skycoin_cli.wlt", w.GetFilename())
				// Confirms the label is empty
				require.Empty(t, w.Meta["label"])
				// Confirms wallet has 5 address entries
				require.Len(t, w.Entries, 5)
			},
		},
		{
			name:  "generate wallet with -f option",
			args:  []string{"-f", "integration-cli.wlt"},
			setup: createTempWalletDir,
			checkWallet: func(t *testing.T, w *wallet.Wallet) {
				// Confirms the default wallet name is skycoin_cli.wlt
				require.Equal(t, "integration-cli.wlt", w.GetFilename())
				// Confirms the label is empty
				require.Empty(t, w.Meta["label"])
			},
		},
		{
			name:  "generate wallet with -l option",
			args:  []string{"-l", "integration-cli"},
			setup: createTempWalletDir,
			checkWallet: func(t *testing.T, w *wallet.Wallet) {
				// Confirms the default wallet name is skycoin_cli.wlt
				require.Equal(t, "skycoin_cli.wlt", w.GetFilename())
				label, ok := w.Meta["label"]
				require.True(t, ok)
				require.Equal(t, "integration-cli", label)
			},
		},
		{
			name: "generate wallet with duplicate wallet name",
			args: []string{},
			setup: func(t *testing.T) func() {
				_, clean := createTempWalletFile(t)
				return clean
			},
			errMsg: []byte("ERROR: integration-test.wlt already exist. See 'skycoin-cli generateWallet --help'"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			clean := tc.setup(t)
			defer clean()

			// Run command with arguments
			args := append([]string{"generateWallet"}, tc.args...)
			output, err := exec.Command(binaryPath, args...).CombinedOutput()
			require.NoError(t, err)
			// Trims the suffix "\n"
			output = bytes.TrimRight(output, "\n")

			// Checks if the output is start with "Error: ",
			// confirms the error message are matched.
			if bytes.Contains(output, []byte("ERROR: ")) {
				require.Equal(t, tc.errMsg, output)
				return
			}

			var rw wallet.ReadableWallet
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&rw)
			require.NoError(t, err)

			// Converts to wallet.Wallet
			w, err := rw.ToWallet()
			require.NoError(t, err)

			// Validate the wallet
			err = w.Validate()
			require.NoError(t, err)

			// Confirms all entries and lastSeed are derived from seed.
			checkWalletEntriesAndLastSeed(t, &w)

			// Checks the wallet with provided checking method.
			tc.checkWallet(t, &w)
		})
	}
}

// checkWalletEntriesAndLastSeed confirms the wallet entries and lastSeed are derivied
// from the seed.
func checkWalletEntriesAndLastSeed(t *testing.T, w *wallet.Wallet) {
	seed, ok := w.Meta["seed"]
	require.True(t, ok)
	newSeed, seckeys := cipher.GenerateDeterministicKeyPairsSeed([]byte(seed), len(w.Entries))
	require.Len(t, seckeys, len(w.Entries))
	for i, sk := range seckeys {
		require.Equal(t, w.Entries[i].Secret, sk)
		pk := cipher.PubKeyFromSecKey(sk)
		require.Equal(t, w.Entries[i].Public, pk)
	}
	lastSeed, ok := w.Meta["lastSeed"]
	require.True(t, ok)
	require.Equal(t, lastSeed, hex.EncodeToString(newSeed))
}

// TestLiveGUIInjectTransaction does almost the same procedure as TestCreateAndBroadcastRawTransaction.
// The only difference is we broadcast the raw transaction throught the gui /injectTransaction api.
func TestLiveGUIInjectTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	doLiveEnvCheck(t)

	if !doWallet(t) {
		return
	}

	c := gui.NewClient(nodeAddress())
	// prepares wallet and confirms the wallet has at least 2 coins and 2 coin hours.
	w, totalCoins, _ := prepareAndCheckWallet(t, 2e6, 2)

	tt := []struct {
		name    string
		args    func() []string
		errMsg  []byte
		checkTx func(t *testing.T, txid string)
	}{
		{
			// Send all coins to the first address to one output.
			name: "send all coins to the first address",
			args: func() []string {
				coins, err := droplet.ToString(totalCoins)
				require.NoError(t, err)
				return []string{"createRawTransaction", w.Entries[0].Address.String(), coins}
			},
			checkTx: func(t *testing.T, txid string) {
				// Confirms all coins are in the first address in one output
				tx := getTransaction(t, txid)
				require.Len(t, tx.Transaction.Transaction.Out, 1)
				c, err := droplet.FromString(tx.Transaction.Transaction.Out[0].Coins)
				require.NoError(t, err)
				require.Equal(t, totalCoins, c)
			},
		},
		{
			// Send 0.5 coin to the second address.
			// Send 0.5 coin to the third address.
			// After sending, the first address should have at least 1 coin left.
			name: "send to multiple address with -m option",
			args: func() []string {
				addrCoins := []struct {
					Addr  string `json:"addr"`
					Coins string `json:"coins"`
				}{
					{
						w.Entries[1].Address.String(),
						"0.5",
					},
					{
						w.Entries[2].Address.String(),
						"0.5",
					},
				}

				v, err := json.Marshal(addrCoins)
				require.NoError(t, err)

				return []string{"createRawTransaction", "-m", string(v)}
			},
			checkTx: func(t *testing.T, txid string) {
				// Confirms the first address has at least 1 coin left.
				coins, _ := getAddressBalance(t, w.Entries[0].Address.String())
				require.True(t, coins >= 1e6)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Create raw transaction first
			output, err := exec.Command(binaryPath, tc.args()...).CombinedOutput()
			if err != nil {
				t.Fatalf("err: %v, output: %v", err, string(output))
				return
			}
			require.NoError(t, err)
			output = bytes.TrimRight(output, "\n")
			if bytes.Contains(output, []byte("ERROR:")) {
				require.Equal(t, tc.errMsg, output)
				return
			}

			// Broadcast raw transaction with gui /injectTransaction
			txid, err := c.InjectTransaction(string(output))
			require.NoError(t, err)

			txid = strings.TrimRight(txid, "\n")
			fmt.Println("txid:", txid)
			_, err = cipher.SHA256FromHex(txid)
			require.NoError(t, err)

			// Wait untill transaction is confirmed.
			tk := time.NewTicker(time.Second)
		loop:
			for {
				select {
				case <-time.After(30 * time.Second):
					t.Fatal("Wait tx confirmation timeout")
				case <-tk.C:
					if isTxConfirmed(t, txid) {
						break loop
					}
				}
			}

			tc.checkTx(t, txid)
		})
	}
}
