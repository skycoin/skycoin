// Package integration_test implements CLI integration tests
package integration_test

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/andreyvit/diff"
	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/api"
	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/cli"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/readable"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/util/droplet"
	wh "github.com/SkycoinProject/skycoin/src/util/http"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/SkycoinProject/skycoin/src/wallet/deterministic"

	// register wallets
	_ "github.com/SkycoinProject/skycoin/src/wallet/bip44wallet"
	_ "github.com/SkycoinProject/skycoin/src/wallet/collection"
	_ "github.com/SkycoinProject/skycoin/src/wallet/xpubwallet"
)

const (
	binaryName = "skycoin-cli.test"

	testModeStable = "stable"
	testModeLive   = "live"

	// Number of random transactions of live transaction test.
	randomLiveTransactionNum = 500

	testFixturesDir = "testdata"
)

var (
	binaryPath string

	update         = flag.Bool("update", false, "update golden files")
	liveTxFull     = flag.Bool("live-tx-full", false, "run live transaction test against full blockchain")
	testLiveWallet = flag.Bool("test-live-wallet", false, "run live wallet tests, requires wallet envvars set")

	cryptoTypes = []crypto.CryptoType{crypto.CryptoTypeScryptChacha20poly1305, crypto.CryptoTypeSha256Xor}

	validNameRegexp     = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	stripCoverageReport = regexp.MustCompile(`PASS\ncoverage: [\d\.]+% of statements in github.com/SkycoinProject/skycoin/\.\.\.\n$`)
)

type TestData struct {
	actual   interface{}
	expected interface{}
}

func init() {
	rand.Seed(time.Now().Unix())
}

func sanitizeName(s string) (string, error) {
	if s == "" {
		return "", errors.New("sanitizeName name empty")
	}
	s = strings.Replace(s, " ", "-", -1)
	if !validNameRegexp.MatchString(s) {
		return "", errors.New("sanitizeName name has invalid characters")
	}
	return s, nil
}

// coverprofileNames manages names for unique coverprofile outputs from invocations of the cli test binary
type coverprofileNames struct {
	names map[string]struct{}
	sync.Mutex
}

func newCoverprofileNames() *coverprofileNames {
	return &coverprofileNames{
		names: make(map[string]struct{}),
	}
}

func (n *coverprofileNames) makeName(name string) (string, error) {
	name, err := sanitizeName(name)
	if err != nil {
		return "", err
	}
	coverprofile := fmt.Sprintf("cli-%s", name)

	coverprofile = n.add(coverprofile)

	return coverprofile, nil
}

func (n *coverprofileNames) add(s string) string {
	n.Lock()
	defer n.Unlock()

	i := 1
	_, ok := n.names[fmt.Sprintf("%s.coverage.out", s)]
	for ok {
		x := fmt.Sprintf("%s-%d.coverage.out", s, i)
		_, ok = n.names[x]
		if !ok {
			s = x
		}
		i++
	}

	s = fmt.Sprintf("%s.coverage.out", s)
	n.names[s] = struct{}{}
	return s
}

var (
	cpNames = newCoverprofileNames()
)

func execCommand(args ...string) *exec.Cmd {
	// Add test flags to arguments to generate a coverage report
	coverprofile, err := cpNames.makeName(args[0])
	if err != nil {
		panic(err)
	}
	args = append(args, []string{fmt.Sprintf("--test.coverprofile=../../../coverage/%s", coverprofile)}...)
	return exec.Command(binaryPath, args...)
}

func execCommandCombinedOutput(args ...string) ([]byte, error) {
	cmd := execCommand(args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	// Remove the trailing coverage statements that the test cli binary produces due to coverage mode, e.g.
	// PASS
	// coverage: 8.1% of statements in github.com/SkycoinProject/skycoin/...
	output = stripCoverageReport.ReplaceAll(output, nil)
	return output, nil
}

func TestMain(m *testing.M) {
	if !enabled() {
		return
	}

	abs, err := filepath.Abs(binaryName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get binary name absolute path failed: %v\n", err)
		os.Exit(1)
	}

	binaryPath = abs

	// Build cli binary file.
	// Args to build the cli binary without coverage:
	// args := []string{"build", "-o", binaryPath, "../../../cmd/skycoin-cli/skycoin-cli.go"}
	// Compile the binary with test flags enabled to get a coverage report from the binary
	args := []string{"test", "-c", "-tags", "testrunmain", "-o", binaryPath, "-coverpkg=github.com/SkycoinProject/skycoin/...", "../../../cmd/skycoin-cli/"}
	if err := exec.Command("go", args...).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Make %v binary failed: %v\n", binaryName, err)
		os.Exit(1)
	}

	ret := m.Run()

	// Remove the generated cli binary file.
	if err := os.Remove(binaryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Delete %v failed: %v", binaryName, err)
		os.Exit(1)
	}

	os.Exit(ret)
}

// createTempWallet creates a temporary dir, and if encrypt is true, copy
// the testdata/$stableEncryptedWalletName file to the dir. If it's false, then
// copy the testdata/$stableWalletName file to the dir
// returns the temporary wallet filename, cleanup callback function, and error if any.
func createTempWallet(t *testing.T, label, seed string, encrypt bool, password []byte) *api.WalletResponse {
	c := newClient()
	wlt, err := c.CreateWallet(api.CreateWalletOptions{
		Seed:     seed,
		Type:     wallet.WalletTypeDeterministic,
		Label:    label,
		Encrypt:  encrypt,
		Password: string(password),
	})
	require.NoError(t, err)
	return wlt
}

type readableDeterministicEntry struct {
	Address string `json:"address"`
	Public  string `json:"public_key"`
	Secret  string `json:"secret_key"`
}

type readableDeterministicWallet struct {
	Meta    wallet.Meta                  `json:"meta"`
	Entries []readableDeterministicEntry `json:"entries"`
}

// createTempWalletDir creates a temporary wallet dir,
// Returns wallet dir path and callback function to clean up the dir.
func createTempWalletDir(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", "wallet-data-dir")
	require.NoError(t, err)

	return dir, func() {
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

func checkGoldenFile(t *testing.T, goldenFile string, td TestData) {
	loadGoldenFile(t, goldenFile, td)
	require.Equal(t, reflect.Indirect(reflect.ValueOf(td.expected)).Interface(), td.actual)
	checkGoldenFileObjectChanges(t, goldenFile, td)
}

func checkGoldenFileObjectChanges(t *testing.T, goldenFile string, td TestData) {
	// Serialize expected to JSON and compare to the goldenFile's contents
	// This will detect field changes that could be missed otherwise
	b, err := json.MarshalIndent(td.expected, "", "\t")
	require.NoError(t, err)

	goldenFile = filepath.Join(testFixturesDir, goldenFile)

	f, err := os.Open(goldenFile)
	require.NoError(t, err)
	defer f.Close()

	c, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	sc := string(c)
	sb := string(b) + "\n"

	require.Equal(t, sc, sb, "JSON struct output differs from golden file, was a field added to the struct?\nDiff:\n"+diff.LineDiff(sc, sb))
}

func mode(t *testing.T) string {
	mode := os.Getenv("SKYCOIN_INTEGRATION_TEST_MODE")
	switch mode {
	case "":
		mode = testModeStable
	case testModeLive, testModeStable:
	default:
		t.Fatal("Invalid test mode, must be stable or live")
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

func doLiveWallet(t *testing.T) bool {
	if *testLiveWallet {
		return true
	}

	t.Skip("Tests requiring wallet envvars are disabled")
	return false
}

func dbNoUnconfirmed(t *testing.T) bool {
	x := os.Getenv("DB_NO_UNCONFIRMED")
	if x == "" {
		return false
	}

	v, err := strconv.ParseBool(x)
	require.NoError(t, err)
	return v
}

// requireWalletEnv checks if live wallet tests are enabled and the CLI_WALLET_FILE environment variable is set
func requireWalletEnv(t *testing.T) string {
	if !doLiveWallet(t) {
		t.Fatal("not doing live wallet tests, should have skipped")
	}

	walletFile := os.Getenv("CLI_WALLET_FILE")
	if walletFile == "" {
		t.Fatal("missing CLI_WALLET_FILE environment variable")
	}

	return walletFile
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
		rpcAddr = "http://127.0.0.1:6420"
	}

	return rpcAddr
}

func newClient() *api.Client {
	c := api.NewClient(rpcAddress())
	c.SetAuth(os.Getenv("RPC_USER"), os.Getenv("RPC_PASS"))
	return c
}

func useCSRF(t *testing.T) bool {
	x := os.Getenv("USE_CSRF")
	if x == "" {
		return false
	}

	useCSRF, err := strconv.ParseBool(x)
	require.NoError(t, err)
	return useCSRF
}

func doHeaderCheck(t *testing.T) bool {
	x := os.Getenv("HEADER_CHECK")
	if x == "" {
		return false
	}

	doHeaderCheck, err := strconv.ParseBool(x)
	require.NoError(t, err)
	return doHeaderCheck

}

func TestWalletAddAddresses(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	tt := []struct {
		name         string
		encrypted    bool
		seed         string
		args         []string
		isUsageErr   bool
		expectOutput []byte
		goldenFile   string
	}{
		{
			name:         "walletAddAddresses",
			encrypted:    false,
			seed:         "exchange stage green marine palm tobacco decline shadow cereal chapter lamp copy",
			expectOutput: []byte("7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\n"),
			goldenFile:   "generate-addresses.golden",
		},
		{
			name:         "walletAddAddresses -n 2 -j",
			encrypted:    false,
			seed:         "visual ancient fancy body choose trigger drama window toward resource enough another",
			args:         []string{"-n", "2", "-j"},
			expectOutput: []byte("{\n    \"addresses\": [\n        \"buDFq2kR9JLJcPoirZbiEL5DJGGBgpbXaU\",\n        \"Xzm3BCV8XCWUgCuM7rtdZ1RUTZnPqKcvw1\"\n    ]\n}\n"),
			goldenFile:   "generate-addresses-2.golden",
		},
		{
			name:         "walletAddAddresses -n -2 -j",
			encrypted:    false,
			seed:         "bronze nut vehicle book vehicle matter curve amused jaguar fall finger fade",
			args:         []string{"-n", "-2", "-j"},
			isUsageErr:   true,
			expectOutput: []byte("Error: invalid value \"-2\" for flag -n: strconv.ParseUint: parsing \"-2\": invalid syntax"),
		},
		{
			name:         "walletAddAddresses in encrypted wallet",
			encrypted:    true,
			seed:         "chunk tortoise solid extra casual lend merry tooth captain inform alpha zebra",
			args:         []string{"-p", "pwd", "-j"},
			expectOutput: []byte("{\n    \"addresses\": [\n        \"2c3Dr4YdHSyc9HAPrjnHcXLQEKnHEitHUn2\"\n    ]\n}\n"),
			goldenFile:   "generate-addresses-encrypted.golden",
		},
		{
			name:         "walletAddAddresses in encrypted wallet with invalid password",
			encrypted:    true,
			seed:         "lazy poverty prepare mad pen celery come panel animal approve cattle already",
			args:         []string{"-p", "invalid password", "-j"},
			expectOutput: []byte("invalid password\n"),
			isUsageErr:   true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var password []byte
			if tc.encrypted {
				password = []byte("pwd")
			}
			wlt := createTempWallet(t, "test", tc.seed, tc.encrypted, password)
			id := wlt.Meta.Filename

			args := append([]string{"walletAddAddresses", id}, tc.args...)
			output, err := execCommandCombinedOutput(args...)
			if err != nil {
				require.EqualError(t, err, "exit status 1")
				return
			}

			if tc.isUsageErr {
				require.True(t, bytes.Contains(output, tc.expectOutput))
				return
			}

			require.Equal(t, tc.expectOutput, output)

			c := newClient()
			wlt, err = c.Wallet(id)
			require.NoError(t, err)

			require.Equal(t, tc.encrypted, wlt.Meta.Encrypted)

			var addrs struct {
				Addresses []string
			}
			for _, e := range wlt.Entries {
				addrs.Addresses = append(addrs.Addresses, e.Address)
			}

			var expect struct {
				Addresses []string
			}
			checkGoldenFile(t, tc.goldenFile, TestData{addrs, &expect})
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
			"Error: Invalid checksum",
		},
		{
			"invalid bitcoin address",
			"1Dcb9gpaZpBKmjqjCsiBsP3sBW1md2kEM2",
			errors.New("exit status 1"),
			"Error: Invalid checksum",
		},
	}

	for _, tc := range tt {
		output, err := execCommandCombinedOutput("verifyAddress", tc.addr)
		if err != nil {
			require.Error(t, tc.err, "%v", err)
			require.Equal(t, tc.err.Error(), err.Error())
			require.Equal(t, tc.errMsg, strings.Trim(string(output), "\n"))
			return
		}

		require.Empty(t, output, string(output))
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
			errMsg: []byte("Error: invalid raw transaction: encoding/hex: odd length hex string\n"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := execCommandCombinedOutput("decodeRawTransaction", tc.rawTx)
			if err != nil {
				require.Error(t, err, "exit status 1")
				require.Equal(t, tc.errMsg, output)
				return
			}

			require.NoError(t, err)

			var txn readable.Transaction
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&txn)
			require.NoError(t, err)

			var expect readable.Transaction
			checkGoldenFile(t, tc.goldenFile, TestData{txn, &expect})
		})
	}

}

func TestEncodeJSONTransaction(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	pathToGoldenFile := func(filename string) string {
		return filepath.Join(testFixturesDir, filename)
	}

	tt := []struct {
		name   string
		rawTx  string
		errMsg []byte
		args   []string
	}{
		{
			name:  "encode success",
			rawTx: "2601000000a1d3345ac47f897f24084b1c6b9bd6e03fc92887050d0748bdab5e639c1fdcd401000000a2a10f07e0e06cf6ba3e793b3186388a126591ee230b3f387617f1ccb6376a3f18e094bd3f7719aa8191c00764f323872f5192da393852bd85dab70b13409d2b01010000004d78de698a33abcfff22391c043b57a56bb0efbdc4a5b975bf8e7889668896bc0400000000bae12bbf671abeb1181fc85f1c01cdfee55deb97980c9c0a00000000543600000000000000373bb3675cbf3880bba3f3de7eb078925b8a72ad0095ba0a000000001c12000000000000008829025fe45b48f29795893a642bdaa89b2bb40e40d2df03000000001c12000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec09c0fc9b01000000001b12000000000000",
			args:  []string{"encodeJsonTransaction", pathToGoldenFile("decode-raw-transaction.golden")},
		},
		{
			name:  "encode recompute success",
			rawTx: "b7000000006b1a69b76b2412314b2b928ad5e97c31c034be5734f9fa77f31f11b6b933b97601000000ddf4bd79f66ea9c7849c5240a27d9a4745ad4661bdac2179184447088512bb3e62c89efa4fd2cee05980c59b38ef23ddbd09bb77e54e94a0f9123d968a090d420101000000f7c183c1823266aff172928f8d06aa65531643456f97ccca6bd34d15e92fac7d01000000000421aa1694a5b04955781d91880d16c0e9c4227ae8030000000000001600000000000000",
			args:  []string{"encodeJsonTransaction", "--fix", pathToGoldenFile("recompute-transaction-hash.golden")},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := execCommandCombinedOutput(tc.args...)
			if err != nil {
				require.Error(t, err, "exit status 1")
				require.Equal(t, tc.errMsg, output)
				return
			}
			require.NoError(t, err)

			output = bytes.Trim(output, "\n")

			require.Equal(t, tc.rawTx, string(output))
		})
	}
}

func TestAddressGen(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	tt := []struct {
		name  string
		args  []string
		err   error
		check func(t *testing.T, output []byte)
	}{
		{
			name: "addressGen --mode=wallet",
			args: []string{"addressGen", "--mode=wallet"},
			check: func(t *testing.T, v []byte) {
				//var w wallet.ReadableDeterministicWallet
				//err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				//require.NoError(t, err)
				w := loadDeterministicWalletFromBytes(t, v)

				// Confirms the wallet type is skycoin
				require.Equal(t, wallet.CoinTypeSkycoin, wallet.CoinType(w.Meta["coin"]))

				// Confirms that the seed is consisted of 12 words
				seed := w.Seed()
				require.NotEmpty(t, seed)
				ss := strings.Split(seed, " ")
				require.Len(t, ss, 12)
			},
		},
		{
			name: "addressGen  --mode=wallet --num 2",
			args: []string{"addressGen", "--mode=wallet", "--num", "2"},
			check: func(t *testing.T, v []byte) {
				//var w wallet.ReadableDeterministicWallet
				//err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				//require.NoError(t, err)
				w := loadDeterministicWalletFromBytes(t, v)

				// Confirms the wallet type is skycoin
				require.Equal(t, wallet.CoinTypeSkycoin, w.Coin())

				// Confirms that the seed is consisted of 12 words
				seed := w.Seed()
				require.NotEmpty(t, seed)
				ss := strings.Split(seed, " ")
				require.Len(t, ss, 12)

				// Confirms that the wallet have 2 address
				entries, err := w.GetEntries()
				require.NoError(t, err)
				require.Len(t, entries, 2)

				// Confirms the addresses are generated from the seed
				_, keys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 2)
				for i, key := range keys {
					pk := cipher.MustPubKeyFromSecKey(key)
					addr := cipher.MustAddressFromSecKey(key)
					require.Equal(t, addr.String(), entries[i].Address.String())
					require.Equal(t, pk.Hex(), entries[i].Public.Hex())
					require.Equal(t, key.Hex(), entries[i].Secret.Hex())
				}
			},
		},
		{
			name: "addressGen --mode=wallet -n 2",
			args: []string{"addressGen", "--mode=wallet", "-n", "2"},
			check: func(t *testing.T, v []byte) {
				//var w wallet.ReadableDeterministicWallet
				//err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				//require.NoError(t, err)
				w := loadDeterministicWalletFromBytes(t, v)

				// Confirms the wallet type is skycoin
				require.Equal(t, wallet.CoinTypeSkycoin, w.Coin())

				// Confirms that the seed is consisted of 12 words
				seed := w.Seed()
				require.NotEmpty(t, seed)
				ss := strings.Split(seed, " ")
				require.Len(t, ss, 12)

				// Confirms that the wallet have 2 address
				entries, err := w.GetEntries()
				require.NoError(t, err)
				require.Len(t, entries, 2)

				// Confirms the addresses are generated from the seed
				_, keys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 2)
				for i, key := range keys {
					pk := cipher.MustPubKeyFromSecKey(key)
					addr := cipher.MustAddressFromSecKey(key)
					require.Equal(t, addr.String(), entries[i].Address.String())
					require.Equal(t, pk.Hex(), entries[i].Public.Hex())
					require.Equal(t, key.Hex(), entries[i].Secret.Hex())
				}
			},
		},
		{
			name: "addressGen  --mode=wallet --hide-secrets -n 2",
			args: []string{"addressGen", "--mode=wallet", "--hide-secrets", "-n", "2"},
			check: func(t *testing.T, v []byte) {
				var w readableDeterministicWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the wallet type is skycoin
				require.Equal(t, wallet.CoinTypeSkycoin, wallet.CoinType(w.Meta["coin"]))

				// Confirms the secrets in Entries are hidden
				require.Len(t, w.Entries, 2)
				for _, e := range w.Entries {
					require.Equal(t, e.Secret, "")
				}
			},
		},
		{
			name: "addressGen -m=wallet -i -n 2",
			args: []string{"addressGen", "-m=wallet", "-i", "-n", "2"},
			check: func(t *testing.T, v []byte) {
				var w readableDeterministicWallet
				err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				require.NoError(t, err)

				// Confirms the wallet type is skycoin
				require.Equal(t, wallet.CoinTypeSkycoin, wallet.CoinType(w.Meta["coin"]))

				// Confirms the secrets in Entries are hidden
				require.Len(t, w.Entries, 2)
				for _, e := range w.Entries {
					require.Equal(t, e.Secret, "")
				}
			},
		},
		{
			name: "addressGen --mode=wallet--coin=bitcoin -n 2",
			args: []string{"addressGen", "--mode=wallet", "--coin=bitcoin", "-n", "2"},
			check: func(t *testing.T, v []byte) {
				//var w wallet.ReadableDeterministicWallet
				//err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				//require.NoError(t, err)
				w := loadDeterministicWalletFromBytes(t, v)

				// Confirms the wallet type is skycoin
				require.Equal(t, wallet.CoinTypeBitcoin, wallet.CoinType(w.Meta["coin"]))

				entries, err := w.GetEntries()
				require.NoError(t, err)
				require.Len(t, entries, 2)

				// Confirms the addresses are bitcoin addresses that generated from the seed
				seed := w.Seed()
				_, keys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 2)
				for i, key := range keys {
					pk := cipher.MustPubKeyFromSecKey(key)
					address := cipher.BitcoinAddressFromPubKey(pk)
					require.Equal(t, address.String(), entries[i].Address.String())
					require.Equal(t, pk.Hex(), entries[i].Public.Hex())
					require.Equal(t, key, entries[i].Secret)
				}
			},
		},
		{
			name: "addressGen --mode=wallet -c=btc -n 2",
			args: []string{"addressGen", "--mode=wallet", "-c=btc", "-n", "2"},
			check: func(t *testing.T, v []byte) {
				//var w wallet.ReadableDeterministicWallet
				//err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				//require.NoError(t, err)
				w := loadDeterministicWalletFromBytes(t, v)

				// Confirms the wallet type is skycoin
				require.Equal(t, wallet.CoinTypeBitcoin, wallet.CoinType(w.Meta["coin"]))

				entries, err := w.GetEntries()
				require.NoError(t, err)
				require.Equal(t, len(entries), 2)

				// Confirms the addresses are bitcoin addresses that generated from the seed
				seed := w.Seed()
				_, keys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 2)
				for i, key := range keys {
					pk := cipher.MustPubKeyFromSecKey(key)
					address := cipher.BitcoinAddressFromPubKey(pk)
					require.Equal(t, address.String(), entries[i].Address.String())
					require.Equal(t, pk.Hex(), entries[i].Public.Hex())
					require.Equal(t, key, entries[i].Secret)
				}
			},
		},
		{
			name: "addressGen --mode=wallet --hex",
			args: []string{"addressGen", "--mode=wallet", "--hex"},
			check: func(t *testing.T, v []byte) {
				//var w wallet.ReadableDeterministicWallet
				//err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				//require.NoError(t, err)
				w := loadDeterministicWalletFromBytes(t, v)

				// Confirms the seed is a valid hex string
				_, err := hex.DecodeString(w.Seed())
				require.NoError(t, err)
			},
		},
		{
			name: "addressGen --mode=addresses",
			args: []string{"addressGen", "--mode=addresses"},
			check: func(t *testing.T, v []byte) {
				// Confirms that only addresses are returned
				v = bytes.Trim(v, "\n")
				_, err := cipher.DecodeBase58Address(string(v))
				require.NoError(t, err)
			},
		},
		{
			name: "addressGen --mode=addresses --entropy=256",
			args: []string{"addressGen", "--mode=addresses", "--entropy=256"},
			check: func(t *testing.T, v []byte) {
				// Confirms that only addresses are returned
				v = bytes.Trim(v, "\n")
				_, err := cipher.DecodeBase58Address(string(v))
				require.NoError(t, err)
			},
		},
		{
			name: "addressGen --entropy=9",
			args: []string{"addressGen", "--entropy=9"},
			check: func(t *testing.T, v []byte) {
				require.Equal(t, "Error: entropy must be 128 or 256\n", string(v))
			},
			err: errors.New("exit status 1"),
		},
		{
			name: "addressGen --mode=wallet --seed 123",
			args: []string{"addressGen", "--mode=wallet", "--seed", "123"},
			check: func(t *testing.T, v []byte) {
				//var w wallet.ReadableDeterministicWallet
				//err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				//require.NoError(t, err)
				w := loadDeterministicWalletFromBytes(t, v)

				pk, sk := cipher.MustGenerateDeterministicKeyPair([]byte("123"))
				addr := cipher.AddressFromPubKey(pk)
				entries, err := w.GetEntries()
				require.NoError(t, err)
				require.Len(t, entries, 1)
				require.Equal(t, addr.String(), entries[0].Address.String())
				require.Equal(t, pk.Hex(), entries[0].Public.Hex())
				require.Equal(t, sk.Hex(), entries[0].Secret.Hex())
			},
		},
		{
			name: "addressGen --mode=wallet -s 123",
			args: []string{"addressGen", "--mode=wallet", "-s", "123"},
			check: func(t *testing.T, v []byte) {
				//var w wallet.ReadableDeterministicWallet
				//err := json.NewDecoder(bytes.NewReader(v)).Decode(&w)
				//require.NoError(t, err)
				w := loadDeterministicWalletFromBytes(t, v)

				pk, sk := cipher.MustGenerateDeterministicKeyPair([]byte("123"))
				addr := cipher.AddressFromPubKey(pk)
				entries, err := w.GetEntries()
				require.NoError(t, err)
				require.Len(t, entries, 1)
				require.Equal(t, addr.String(), entries[0].Address.String())
				require.Equal(t, pk.Hex(), entries[0].Public.Hex())
				require.Equal(t, sk.Hex(), entries[0].Secret.Hex())
			},
		},
		{
			name: "addressGen --hide-secrets --mode=secrets",
			args: []string{"addressGen", "--mode=secrets", "--hide-secrets"},
			check: func(t *testing.T, v []byte) {
				require.Equal(t, "Error: secrets mode selected but hideSecrets enabled\n", string(v))
			},
			err: errors.New("exit status 1"),
		},
		{
			name: "addressGen --mode=secrets",
			args: []string{"addressGen", "--mode=secrets"},
			check: func(t *testing.T, v []byte) {
				// Confirms that only secret keys are returned
				v = bytes.Trim(v, "\n")
				_, err := cipher.SecKeyFromHex(string(v))
				require.NoError(t, err)
			},
		},
		{
			name: "addressGen --mode=secrets --coin=bitcoin",
			args: []string{"addressGen", "--mode=secrets", "--coin=bitcoin"},
			check: func(t *testing.T, v []byte) {
				// Confirms that only secret keys are returned
				v = bytes.Trim(v, "\n")
				_, err := cipher.SecKeyFromBitcoinWalletImportFormat(string(v))
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := execCommandCombinedOutput(tc.args...)
			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			tc.check(t, output)
		})
	}
}

func TestFiberAddressGen(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	checkAddrsFile := func(t *testing.T, fn string, n int) []string {
		b, err := ioutil.ReadFile(fn)
		require.NoError(t, err)

		addrs := strings.Split(strings.TrimSpace(string(b)), "\n")
		require.Equal(t, n, len(addrs))

		rx := regexp.MustCompile(`"([a-zA-Z0-9]+)",`)

		addrsMap := make(map[string]struct{}, len(addrs))
		out := make([]string, len(addrs))
		for i, a := range addrs {
			_, ok := addrsMap[a]
			require.False(t, ok)
			addrsMap[a] = struct{}{}

			matches := rx.FindStringSubmatch(a)
			require.Len(t, matches, 2)
			addr := matches[1]

			_, err := cipher.DecodeBase58Address(addr)
			require.NoError(t, err)

			out[i] = addr
		}

		return out
	}

	checkSeedsFile := func(t *testing.T, fn string, entropy int, addrs []string) {
		f, err := os.Open(fn)
		require.NoError(t, err)
		defer f.Close()

		r := csv.NewReader(f)
		records, err := r.ReadAll()
		require.NoError(t, err)
		require.Len(t, records, len(addrs))

		seedsMap := make(map[string]struct{}, len(records))

		for i, x := range records {
			require.Len(t, x, 2)

			// addr is valid and matches the addr written to the addrs file
			addr := x[0]
			_, err := cipher.DecodeBase58Address(addr)
			require.NoError(t, err)
			require.Equal(t, addrs[i], addr)

			seed := x[1]

			// no duplicate seeds
			_, ok := seedsMap[seed]
			require.False(t, ok)
			seedsMap[seed] = struct{}{}

			// seed is a valid mnemonic
			err = bip39.ValidateMnemonic(seed)
			require.NoError(t, err)

			// seed entropy is as expected
			switch entropy {
			case 128:
				require.Len(t, strings.Split(seed, " "), 12)
			case 256:
				require.Len(t, strings.Split(seed, " "), 24)
			default:
				t.Fatalf("entropy must be 128 or 256")
			}

			// seed generates the correct address
			pk, _ := cipher.MustGenerateDeterministicKeyPair([]byte(seed))
			regenAddr := cipher.AddressFromPubKey(pk)
			require.Equal(t, addr, regenAddr.String())
		}
	}

	touch := func(t *testing.T, fn string) {
		f, err := os.Create(fn)
		require.NoError(t, err)
		defer f.Close()
		err = f.Close()
		require.NoError(t, err)
	}

	addrsFilename := "addresses.txt"
	seedsFilename := "seeds.csv"

	cases := []struct {
		name  string
		args  []string
		err   error
		setup func(t *testing.T)
		check func(t *testing.T, v []byte)
	}{
		{
			name: "fiberAddressGen",
			args: []string{"fiberAddressGen"},
			setup: func(t *testing.T) {
				testutil.RequireFileNotExists(t, addrsFilename)
				testutil.RequireFileNotExists(t, seedsFilename)
			},
			check: func(t *testing.T, v []byte) {
				defer os.Remove(addrsFilename)
				defer os.Remove(seedsFilename)
				testutil.RequireFileExists(t, addrsFilename)
				testutil.RequireFileExists(t, seedsFilename)
				addrs := checkAddrsFile(t, addrsFilename, 100)
				checkSeedsFile(t, seedsFilename, 128, addrs)
			},
		},
		{
			name: "fiberAddressGen --entropy=256",
			args: []string{"fiberAddressGen", "--entropy=256"},
			setup: func(t *testing.T) {
				testutil.RequireFileNotExists(t, addrsFilename)
				testutil.RequireFileNotExists(t, seedsFilename)
			},
			check: func(t *testing.T, v []byte) {
				defer os.Remove(addrsFilename)
				defer os.Remove(seedsFilename)
				testutil.RequireFileExists(t, addrsFilename)
				testutil.RequireFileExists(t, seedsFilename)
				addrs := checkAddrsFile(t, addrsFilename, 100)
				checkSeedsFile(t, seedsFilename, 256, addrs)
			},
		},
		{
			name: "fiberAddressGen -n=1",
			args: []string{"fiberAddressGen", "-n=1"},
			setup: func(t *testing.T) {
				testutil.RequireFileNotExists(t, addrsFilename)
				testutil.RequireFileNotExists(t, seedsFilename)
			},
			check: func(t *testing.T, v []byte) {
				defer os.Remove(addrsFilename)
				defer os.Remove(seedsFilename)
				testutil.RequireFileExists(t, addrsFilename)
				testutil.RequireFileExists(t, seedsFilename)
				addrs := checkAddrsFile(t, addrsFilename, 1)
				checkSeedsFile(t, seedsFilename, 128, addrs)
			},
		},
		{
			name: "fiberAddressGen can't overwrite addrs file",
			args: []string{"fiberAddressGen"},
			setup: func(t *testing.T) {
				testutil.RequireFileNotExists(t, addrsFilename)
				testutil.RequireFileNotExists(t, seedsFilename)
				touch(t, addrsFilename)
				testutil.RequireFileExists(t, addrsFilename)
			},
			check: func(t *testing.T, v []byte) {
				defer os.Remove(addrsFilename)
				defer os.Remove(seedsFilename)
				testutil.RequireFileNotExists(t, seedsFilename)
				require.Equal(t, "Error: -addrs-file \"addresses.txt\" already exists. Use -overwrite to force writing\n", string(v))
			},
			err: errors.New("exit status 1"),
		},
		{
			name: "fiberAddressGen can't overwrite seeds file",
			args: []string{"fiberAddressGen"},
			setup: func(t *testing.T) {
				testutil.RequireFileNotExists(t, addrsFilename)
				testutil.RequireFileNotExists(t, seedsFilename)
				touch(t, seedsFilename)
				testutil.RequireFileExists(t, seedsFilename)
			},
			check: func(t *testing.T, v []byte) {
				defer os.Remove(addrsFilename)
				defer os.Remove(seedsFilename)
				testutil.RequireFileNotExists(t, addrsFilename)
				require.Equal(t, "Error: -seeds-file \"seeds.csv\" already exists. Use -overwrite to force writing\n", string(v))
			},
			err: errors.New("exit status 1"),
		},
		{
			name: "fiberAddressGen --overwrite",
			args: []string{"fiberAddressGen", "--overwrite"},
			setup: func(t *testing.T) {
				testutil.RequireFileNotExists(t, addrsFilename)
				testutil.RequireFileNotExists(t, seedsFilename)
				touch(t, addrsFilename)
				touch(t, seedsFilename)
				testutil.RequireFileExists(t, addrsFilename)
				testutil.RequireFileExists(t, seedsFilename)
			},
			check: func(t *testing.T, v []byte) {
				defer os.Remove(addrsFilename)
				defer os.Remove(seedsFilename)
				testutil.RequireFileExists(t, addrsFilename)
				testutil.RequireFileExists(t, seedsFilename)
				addrs := checkAddrsFile(t, addrsFilename, 100)
				checkSeedsFile(t, seedsFilename, 128, addrs)
			},
		},
		{
			name: "fiberAddressGen --addrs-file=fooaddrs.txt --seeds-file=fooseeds.csv",
			args: []string{"fiberAddressGen", "--addrs-file", "fooaddrs.txt", "--seeds-file", "fooseeds.csv"},
			setup: func(t *testing.T) {
				testutil.RequireFileNotExists(t, "fooaddrs.txt")
				testutil.RequireFileNotExists(t, "fooseeds.csv")
			},
			check: func(t *testing.T, v []byte) {
				defer os.Remove("fooaddrs.txt")
				defer os.Remove("fooseeds.csv")
				testutil.RequireFileExists(t, "fooaddrs.txt")
				testutil.RequireFileExists(t, "fooseeds.csv")
				addrs := checkAddrsFile(t, "fooaddrs.txt", 100)
				checkSeedsFile(t, "fooseeds.csv", 128, addrs)
			},
		},
		{
			name: "fiberAddressGen positional-args-not-allowed",
			args: []string{"fiberAddressGen", "foo"},
			check: func(t *testing.T, v []byte) {
				testutil.RequireFileNotExists(t, addrsFilename)
				testutil.RequireFileNotExists(t, seedsFilename)
				require.Equal(t, "Error: this command does not take any positional arguments\n", string(v))
			},
			err: errors.New("exit status 1"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(t)
			}

			output, err := execCommandCombinedOutput(tc.args...)
			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			tc.check(t, output)
		})
	}
}

func TestStableListWallets(t *testing.T) {
	if !doStable(t) {
		return
	}

	walletName, clean := createUnencryptedWallet(t)
	defer clean()

	output, err := execCommandCombinedOutput("listWallets", filepath.Dir(walletName))
	require.NoError(t, err, output)

	var wlts struct {
		Wallets []cli.WalletEntry `json:"wallets"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wlts)
	require.NoError(t, err)

	var expect struct {
		Wallets []cli.WalletEntry `json:"wallets"`
	}
	checkGoldenFile(t, "list-wallets.golden", TestData{wlts, &expect})
}

func TestLiveListWallets(t *testing.T) {
	if !doLive(t) {
		return
	}

	fn := requireWalletEnv(t)

	output, err := execCommandCombinedOutput("listWallets", filepath.Dir(fn))
	require.NoError(t, err)

	var wlts struct {
		Directory string            `json:"directory"`
		Wallets   []cli.WalletEntry `json:"wallets"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wlts)
	require.NoError(t, err)
	require.NotEmpty(t, wlts.Wallets)
	require.Equal(t, filepath.Dir(fn), wlts.Directory)

	// Defaults to $DATA_DIR/wallets when no arguments are specified
	output, err = execCommandCombinedOutput("listWallets")
	require.NoError(t, err)

	cfg := showConfig(t)
	defaultDir := filepath.Join(cfg.DataDir, "wallets")
	var wlts2 struct {
		Directory string            `json:"directory"`
		Wallets   []cli.WalletEntry `json:"wallets"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wlts2)
	require.NoError(t, err)
	require.NotEmpty(t, wlts2.Wallets)
	require.Equal(t, defaultDir, wlts2.Directory)
}

func showConfig(t *testing.T) cli.Config {
	output, err := execCommandCombinedOutput("showConfig")
	require.NoError(t, err)

	var ret cli.Config
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&ret)
	require.NoError(t, err)

	return ret
}

func TestStableListAddress(t *testing.T) {
	if !doStable(t) {
		return
	}
	seed := "radar erase claw much slush custom symbol cable poem apology genre edit"
	wlt := createTempWallet(t, "test", seed, false, nil)

	output, err := execCommandCombinedOutput("listAddresses", wlt.Meta.Filename)
	require.NoError(t, err)

	var wltAddresses struct {
		Addresses []string `json:"addresses"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltAddresses)
	require.NoError(t, err)

	var expect struct {
		Addresses []string `json:"addresses"`
	}
	checkGoldenFile(t, "list-addresses.golden", TestData{wltAddresses, &expect})
}

func TestLiveListAddresses(t *testing.T) {
	if !doLive(t) {
		return
	}

	fn := requireWalletEnv(t)

	output, err := execCommandCombinedOutput("listAddresses", fn)
	require.NoError(t, err)

	var wltAddresses struct {
		Addresses []string `json:"addresses"`
	}

	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltAddresses)
	require.NoError(t, err)

	require.NotEmpty(t, wltAddresses.Addresses)
}

func TestStableAddressBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := execCommandCombinedOutput("addressBalance", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt")
	require.NoError(t, err)

	var addrBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrBalance)
	require.NoError(t, err)

	var expect cli.BalanceResult
	checkGoldenFile(t, "address-balance.golden", TestData{addrBalance, &expect})
}

func TestLiveAddressBalance(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := execCommandCombinedOutput("addressBalance", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt")
	require.NoError(t, err)

	var addrBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrBalance)
	require.NoError(t, err)
}

func TestStableWalletBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	seed := "recall large warrior cargo harbor ask moral strong mixture small october aerobic"
	wlt := createTempWallet(t, "test-stable-wallet-balance", seed, false, nil)

	output, err := execCommandCombinedOutput("walletBalance", wlt.Meta.Filename)
	require.NoError(t, err, output)

	var wltBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltBalance)
	require.NoError(t, err)

	var expect cli.BalanceResult
	checkGoldenFile(t, "wallet-balance.golden", TestData{wltBalance, &expect})
}

func TestLiveWalletBalance(t *testing.T) {
	if !doLive(t) {
		return
	}

	fn := requireWalletEnv(t)

	output, err := execCommandCombinedOutput("walletBalance", fn)
	require.NoError(t, err)

	var wltBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltBalance)
	require.NoError(t, err)

	require.NotEmpty(t, wltBalance.Confirmed.Coins)
	require.NotEmpty(t, wltBalance.Addresses)
}

func TestStableWalletOutputs(t *testing.T) {
	if !doStable(t) {
		return
	}

	seed := "crush dice soccer what dress bread cancel predict rose relax truck side"
	wlt := createTempWallet(t, "test-stable-wallet-outputs", seed, false, nil)

	output, err := execCommandCombinedOutput("walletOutputs", wlt.Meta.Filename)
	require.NoError(t, err, output)

	var wltOutput cli.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltOutput)
	require.NoError(t, err)

	var expect cli.OutputsResult
	checkGoldenFile(t, "wallet-outputs.golden", TestData{wltOutput, &expect})
}

func TestLiveWalletOutputs(t *testing.T) {
	if !doLive(t) {
		return
	}

	fn := requireWalletEnv(t)

	output, err := execCommandCombinedOutput("walletOutputs", fn)
	require.NoError(t, err)

	var wltOutput cli.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltOutput)
	require.NoError(t, err)

	require.NotEmpty(t, wltOutput.Outputs.HeadOutputs)
}

func TestStableAddressOutputs(t *testing.T) {
	if !doStable(t) {
		return
	}

	tt := []struct {
		name       string
		args       []string
		goldenFile string
		err        error
		errMsg     string
	}{
		{
			name:       "addressOutputs one address",
			args:       []string{"addressOutputs", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
			goldenFile: "address-outputs.golden",
		},
		{
			name:       "addressOutputs two address",
			args:       []string{"addressOutputs", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt", "ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm"},
			goldenFile: "two-addresses-outputs.golden",
		},
		{
			name:   "addressOutputs two address one invalid",
			args:   []string{"addressOutputs", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt", "badaddress"},
			err:    errors.New("exit status 1"),
			errMsg: "Error: invalid address: badaddress, err: Invalid address length\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := execCommandCombinedOutput(tc.args...)

			if tc.err != nil {
				testutil.RequireError(t, err, tc.err.Error())
				require.Equal(t, tc.errMsg, string(output))
				return
			}

			require.NoError(t, err)

			var addrOutputs cli.OutputsResult
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrOutputs)
			require.NoError(t, err)

			var expect cli.OutputsResult
			checkGoldenFile(t, tc.goldenFile, TestData{addrOutputs, &expect})
		})
	}
}

func TestLiveAddressOutputs(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := execCommandCombinedOutput("addressOutputs", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt")
	require.NoError(t, err)

	var addrOutputs cli.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrOutputs)
	require.NoError(t, err)
}

func TestStableShowConfig(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := execCommandCombinedOutput("showConfig")
	require.NoError(t, err)

	var ret cli.Config
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&ret)
	require.NoError(t, err)

	// DataDir can't be checked perfectly without essentially
	// reimplementing cli.LoadConfig to compare values
	require.NotEmpty(t, ret.DataDir)
	require.True(t, strings.HasSuffix(ret.DataDir, ".skycoin"))

	ret.DataDir = "IGNORED/.skycoin"

	goldenFile := "show-config.golden"

	var expect cli.Config
	td := TestData{ret, &expect}
	loadGoldenFile(t, goldenFile, td)

	// The RPC port is not always the same between runs of the stable integration tests,
	// so use the RPC_ADDR envvar instead of the golden file value for comparison
	goldenRPCAddress := expect.RPCAddress
	expect.RPCAddress = rpcAddress()

	require.Equal(t, expect, ret)

	// Restore goldenfile's value before checking if JSON fields were added or removed
	expect.RPCAddress = goldenRPCAddress
	checkGoldenFileObjectChanges(t, goldenFile, TestData{ret, &expect})
}

func TestLiveShowConfig(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := execCommandCombinedOutput("showConfig")
	require.NoError(t, err)

	var ret cli.Config
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&ret)
	require.NoError(t, err)

	// DataDir can't be checked perfectly without essentially
	// reimplementing cli.LoadConfig to compare values
	require.NotEmpty(t, ret.DataDir)
	require.True(t, strings.HasSuffix(ret.DataDir, ".skycoin"))

	coin := os.Getenv("COIN")
	if coin == "" {
		coin = "skycoin"
	}
	require.Equal(t, coin, ret.Coin)

	require.Equal(t, rpcAddress(), ret.RPCAddress)
}

func TestStableStatus(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := execCommandCombinedOutput("status")
	require.NoError(t, err)

	var ret cli.StatusResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&ret)
	require.NoError(t, err)

	// TimeSinceLastBlock is not stable
	ret.Status.BlockchainMetadata.TimeSinceLastBlock = wh.FromDuration(time.Duration(0))
	// Version is not stable
	ret.Status.Version = readable.BuildInfo{}
	// Uptime is not stable
	ret.Status.Uptime = wh.FromDuration(time.Duration(0))
	// StartedAt is not stable
	ret.Status.StartedAt = 0
	goldenFile := "status"
	if useCSRF(t) {
		goldenFile += "-csrf-enabled"
	}
	if !doHeaderCheck(t) {
		goldenFile += "-header-check-disabled"
	}
	if dbNoUnconfirmed(t) {
		goldenFile += "-no-unconfirmed"
	}
	goldenFile += ".golden"

	var expect cli.StatusResult
	td := TestData{ret, &expect}
	loadGoldenFile(t, goldenFile, td)

	// The RPC port is not always the same between runs of the stable integration tests,
	// so use the RPC_ADDR envvar instead of the golden file value for comparison
	goldenRPCAddress := expect.Config.RPCAddress
	expect.Config.RPCAddress = rpcAddress()

	require.Equal(t, expect, ret)

	// Restore goldenfile's value before checking if JSON fields were added or removed
	expect.Config.RPCAddress = goldenRPCAddress
	checkGoldenFileObjectChanges(t, goldenFile, TestData{ret, &expect})
}

func TestLiveStatus(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := execCommandCombinedOutput("status")
	require.NoError(t, err)

	var ret cli.StatusResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&ret)
	require.NoError(t, err)
	require.Equal(t, rpcAddress(), ret.Config.RPCAddress)
}

func TestStableTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	type testCase struct {
		name       string
		args       []string
		err        error
		errMsg     string
		goldenFile string
	}

	tt := []testCase{
		{
			name:       "invalid txid",
			args:       []string{"abcd"},
			err:        errors.New("exit status 1"),
			errMsg:     "Error: invalid txid\n",
			goldenFile: "",
		},
		{
			name:       "not exist",
			args:       []string{"540582ee4128b733f810f149e908d984a5f403ad2865108e6c1c5423aeefc759"},
			err:        errors.New("exit status 1"),
			errMsg:     "Error: 404 Not Found\n",
			goldenFile: "",
		},
		{
			name:       "empty txid",
			args:       []string{""},
			err:        errors.New("exit status 1"),
			errMsg:     "Error: txid is empty\n",
			goldenFile: "",
		},
		{
			name:       "genesis transaction",
			args:       []string{"d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add"},
			err:        nil,
			errMsg:     "",
			goldenFile: "genesis-transaction-cli.golden",
		},
	}

	if !dbNoUnconfirmed(t) {
		tt = append(tt, testCase{
			name:       "unconfirmed",
			args:       []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			goldenFile: "unconfirmed-transaction-cli.golden",
		})
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"transaction"}, tc.args...)
			o, err := execCommandCombinedOutput(args...)
			if tc.err != nil {
				testutil.RequireError(t, err, tc.err.Error())
				require.Equal(t, tc.errMsg, string(o))
				return
			}

			require.NoError(t, err)

			var tx cli.TxnResult
			err = json.NewDecoder(bytes.NewReader(o)).Decode(&tx)
			require.NoError(t, err)

			var expect cli.TxnResult
			checkGoldenFile(t, tc.goldenFile, TestData{tx, &expect})
		})
	}

	scanTransactions(t, true)
}

func TestLiveTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	o, err := execCommandCombinedOutput("transaction", "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add")
	require.NoError(t, err)
	var tx cli.TxnResult
	err = json.NewDecoder(bytes.NewReader(o)).Decode(&tx)
	require.NoError(t, err)

	var expect cli.TxnResult

	loadGoldenFile(t, "genesis-transaction.golden", TestData{tx, &expect})
	require.Equal(t, expect.Transaction.Transaction, tx.Transaction.Transaction)

	scanTransactions(t, *liveTxFull)

	// scan pending transactions
	scanPendingTransactions(t)
}

func prepareCSVFile(t *testing.T, toAddrs [][]string) (csvFile string, teardown func(t *testing.T)) {
	fn := "create_txn_test.csv"
	tmpDir, err := ioutil.TempDir("", "create_raw_transaction")
	require.NoError(t, err)
	csvFile = filepath.Join(tmpDir, fn)

	f, err := os.Create(csvFile)
	require.NoError(t, err)
	defer f.Close()
	w := csv.NewWriter(f)

	for _, to := range toAddrs {
		err := w.Write(to)
		require.NoError(t, err)
	}
	w.Flush()
	require.NoError(t, w.Error())

	return csvFile, func(t *testing.T) {
		require.NoError(t, os.Remove(csvFile))
	}
}

func TestLiveCreateRawTransactionV2(t *testing.T) {
	if !doLive(t) {
		return
	}

	walletFile := requireWalletEnv(t)
	w, err := wallet.Load(walletFile)
	require.NoError(t, err)
	addrs, err := w.GetAddresses()
	require.NoError(t, err)
	require.Truef(t, len(addrs) >= 2, "wallet must have at least 2 addresses")

	// prepare csv file for testing
	toAddrs := [][]string{
		{addrs[0].String(), "0.001"},
		{addrs[1].String(), "0.001"},
	}
	csvFile, teardown := prepareCSVFile(t, toAddrs)
	defer teardown(t)

	var testCases = []struct {
		name   string
		args   func(t *testing.T) []string
		verify func(t *testing.T, data []byte)
	}{
		{
			name: "unsigned=true json=false",
			args: func(t *testing.T) []string {
				return []string{
					walletFile,
					addrs[0].String(), // to address
					"1.01",
					"--unsign",
				}
			},
			verify: func(t *testing.T, data []byte) {
				s := strings.TrimSuffix(string(data), "\n")
				txn, err := coin.DeserializeTransactionHex(string(s))
				require.NoError(t, err)
				require.Equal(t, 1, len(txn.Sigs))
				require.Equal(t, cipher.Sig{}, txn.Sigs[0])
			},
		},
		{
			name: "unsigned=true json=true",
			args: func(t *testing.T) []string {
				return []string{
					walletFile,
					addrs[0].String(), // to address
					"1.01",
					"--unsign",
					"--json",
				}
			},
			verify: func(t *testing.T, data []byte) {
				var rsp api.CreateTransactionResponse
				err := json.NewDecoder(bytes.NewReader(data)).Decode(&rsp)
				require.NoError(t, err)
			},
		},
		{
			name: "unsigned=false json=false",
			args: func(t *testing.T) []string {
				args := []string{
					walletFile,
					addrs[0].String(), // to address
					"1.01",
				}

				// Require password if the wallet is encrypted
				if w.IsEncrypted() {
					password := os.Getenv("WALLET_PASSWORD")
					if len(password) == 0 {
						t.Fatal("missing WALLET_PASSWORD environment variable")
						return nil
					}
					args = append(args, "-p", password)
				}
				return args
			},
			verify: func(t *testing.T, data []byte) {
				s := strings.TrimSuffix(string(data), "\n")
				txn, err := coin.DeserializeTransactionHex(string(s))
				require.NoError(t, err)
				require.Equal(t, 1, len(txn.Sigs))
				require.NotEqual(t, cipher.Sig{}, txn.Sigs[0])
			},
		},
		{
			name: "unsigned=true json=false change-address",
			args: func(t *testing.T) []string {
				return []string{
					walletFile,
					addrs[0].String(), // to address
					"0.001",
					"--unsign",
					"--change-address",
					addrs[1].String(),
				}
			},
			verify: func(t *testing.T, data []byte) {
				s := strings.TrimSuffix(string(data), "\n")
				txn, err := coin.DeserializeTransactionHex(string(s))
				require.NoError(t, err)
				require.Equal(t, 1, len(txn.Sigs))
				require.Equal(t, cipher.Sig{}, txn.Sigs[0])
				require.Equal(t, 2, len(txn.Out))
				addrOutMap := make(map[string]struct{})
				for _, o := range txn.Out {
					addrOutMap[o.Address.String()] = struct{}{}
				}

				// Confirms that the toAddr exists in txn.Out
				toAddr := addrs[1].String()
				_, ok := addrOutMap[toAddr]
				require.True(t, ok)
				// Confirms that the change address exists in txn.Out
				_, ok = addrOutMap[addrs[1].String()]
				require.True(t, ok)
			},
		},
		{
			name: "unsigned=true json=false from-addrss",
			args: func(t *testing.T) []string {
				return []string{
					walletFile,
					addrs[1].String(), // to address
					"0.001",
					"--unsign",
					"--from-address",
					addrs[0].String(),
				}
			},
			verify: func(t *testing.T, data []byte) {
				s := strings.TrimSuffix(string(data), "\n")
				txn, err := coin.DeserializeTransactionHex(string(s))
				require.NoError(t, err)
				require.Equal(t, 1, len(txn.Sigs))
				require.Equal(t, cipher.Sig{}, txn.Sigs[0])
				// Get the uxouts of from address
				uxoutsMap := getAddressOutputs(t, addrs[0].String())
				for _, in := range txn.In {
					_, ok := uxoutsMap[in.String()]
					require.True(t, ok)
				}
			},
		},
		{
			name: "unsigned=true json=false -csv",
			args: func(t *testing.T) []string {
				return []string{
					walletFile,
					"--unsign",
					"--csv",
					csvFile,
				}
			},
			verify: func(t *testing.T, data []byte) {
				s := strings.TrimSuffix(string(data), "\n")
				txn, err := coin.DeserializeTransactionHex(string(s))
				require.NoError(t, err)
				require.Equal(t, 1, len(txn.Sigs))
				require.Equal(t, cipher.Sig{}, txn.Sigs[0])
				require.True(t, len(txn.Out) >= 2)

				// Confirms that the txn.Out contains the receiver address
				addrOutMap := make(map[string]coin.TransactionOutput)
				for i, o := range txn.Out {
					addrOutMap[o.Address.String()] = txn.Out[i]
				}

				for _, to := range toAddrs {
					out, ok := addrOutMap[to[0]]
					require.True(t, ok)
					coins, err := droplet.FromString(to[1])
					require.NoError(t, err)
					require.True(t, out.Coins >= coins)
				}
			},
		},
	}

	for _, tc := range testCases {
		args := append([]string{"createRawTransactionV2"}, tc.args(t)...)
		o, err := execCommandCombinedOutput(args...)
		require.NoError(t, err)
		tc.verify(t, o)
	}
}

func getAddressOutputs(t *testing.T, address string) map[string]struct{} {
	output, err := execCommandCombinedOutput("addressOutputs", address)
	require.NoError(t, err)

	var addrOutputs cli.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrOutputs)
	require.NoError(t, err)
	uxoutsMap := make(map[string]struct{})
	for _, o := range addrOutputs.Outputs.HeadOutputs {
		uxoutsMap[o.Hash] = struct{}{}
	}
	return uxoutsMap
}

// TODO cli doesn't have command to querying pending transactions yet.
func scanPendingTransactions(t *testing.T) {
}

// scanTransactions scans transactions against blockchain.
// If fullTest is true, scan the whole blockchain, and test every transactions,
// otherwise just test random transactions.
func scanTransactions(t *testing.T, fullTest bool) {
	// Gets blockchain height through "status" command
	output, err := execCommandCombinedOutput("status")
	require.NoError(t, err)

	d := json.NewDecoder(bytes.NewReader(output))
	d.DisallowUnknownFields()

	var status cli.StatusResult
	err = d.Decode(&status)
	require.NoError(t, err)

	txids := getTxids(t, status.Status.BlockchainMetadata.Head.BkSeq)

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
	nReq := 4
	if useCSRF(t) {
		nReq = 1
	}
	for i := 0; i < nReq; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for txid := range txC {
				t.Run(fmt.Sprintf("%v", txid), func(t *testing.T) {
					o, err := execCommandCombinedOutput("transaction", txid)
					require.NoError(t, err)
					var txRlt cli.TxnResult
					err = json.NewDecoder(bytes.NewReader(o)).Decode(&txRlt)
					require.NoError(t, err)
					require.Equal(t, txid, txRlt.Transaction.Transaction.Hash)
					require.True(t, txRlt.Transaction.Status.Confirmed)
				})
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
	for i := 0; i < n; i++ {
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
	o, err := execCommandCombinedOutput("blocks", s, e)
	require.NoError(t, err)
	var blocks readable.Blocks
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
	output, err := execCommandCombinedOutput("blocks", "180", "181")
	require.NoError(t, err)

	var blocks readable.Blocks
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&blocks)
	require.NoError(t, err)

	var expect readable.Blocks
	checkGoldenFile(t, "blocks180.golden", TestData{blocks, &expect})
}

func TestLiveBlocks(t *testing.T) {
	if !doLive(t) {
		return
	}

	testKnownBlocks(t)

	// These blocks were affected by the coinhour overflow issue, make sure that they can be queried
	blockSeqs := []int{11685, 11707, 11710, 11709, 11705, 11708, 11711, 11706, 11699}

	for _, seq := range blockSeqs {
		output, err := execCommandCombinedOutput("blocks", strconv.Itoa(seq))
		require.NoError(t, err)
		var blocks readable.Blocks
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
			output, err := execCommandCombinedOutput(tc.args...)
			require.NoError(t, err)

			var blocks readable.Blocks
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&blocks)
			require.NoError(t, err)

			var expect readable.Blocks
			checkGoldenFile(t, tc.goldenFile, TestData{blocks, &expect})
		})
	}

	scanBlocks(t, "0", "180")
}

func scanBlocks(t *testing.T, start, end string) { //nolint:unparam
	outputs, err := execCommandCombinedOutput("blocks", start, end)
	require.NoError(t, err)

	var blocks readable.Blocks
	err = json.NewDecoder(bytes.NewReader(outputs)).Decode(&blocks)
	require.NoError(t, err)

	var preBlocks readable.Block
	preBlocks.Head.Hash = "0000000000000000000000000000000000000000000000000000000000000000"
	for _, b := range blocks.Blocks {
		require.Equal(t, b.Head.PreviousHash, preBlocks.Head.Hash)
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
			output, err := execCommandCombinedOutput(tc.args...)

			if bytes.Contains(output, []byte("Error: ")) {
				require.Equal(t, string(tc.errMsg), string(output))
				return
			}

			require.NoError(t, err)

			var blocks readable.Blocks
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&blocks)
			require.NoError(t, err)

			var expect readable.Blocks
			checkGoldenFile(t, tc.goldenFile, TestData{blocks, &expect})
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
			output, err := execCommandCombinedOutput(tc.args...)
			require.NoError(t, err)

			var blocks readable.Blocks
			err = json.NewDecoder(bytes.NewReader(output)).Decode(&blocks)
			require.NoError(t, err)
		})
	}
}

// TestLiveSend sends coin from specific wallet file, user should manually specify the
// wallet file by setting the CLI_WALLET_FILE environment variable
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

	fn := requireWalletEnv(t)

	// prepares wallet and confirms the wallet has at least 2 coins and 16 coin hours.
	w, totalCoins, _ := prepareAndCheckWallet(t, 2e6, 16)
	entries, err := w.GetEntries()
	require.NoError(t, err)

	if w.IsEncrypted() {
		t.Skip("CLI wallet integration tests do not support encrypted wallets yet")
		return
	}

	tt := []struct {
		name     string
		args     func() []string
		errMsg   []byte
		checkTxn func(t *testing.T, txid string)
	}{
		{
			// Send all coins to the first address to one output.
			name: "send all coins to the first address",
			args: func() []string {
				coins, err := droplet.ToString(totalCoins)
				require.NoError(t, err)
				return []string{"send", fn, entries[0].Address.String(), coins}
			},
			checkTxn: func(t *testing.T, txid string) {
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
						entries[1].Address.String(),
						"0.5",
					},
					{
						entries[2].Address.String(),
						"0.5",
					},
				}

				v, err := json.Marshal(addrCoins)
				require.NoError(t, err)

				return []string{"send", fn, "-m", string(v)}
			},
			checkTxn: func(t *testing.T, txid string) {
				tx := getTransaction(t, txid)
				// Confirms the second address receives 0.5 coin and 1 coinhour in this transaction
				checkCoinsAndCoinhours(t, tx, entries[1].Address.String(), 5e5, 1)
				// Confirms the third address receives 0.5 coin and 1 coinhour in this transaction
				checkCoinsAndCoinhours(t, tx, entries[2].Address.String(), 5e5, 1)
				// Confirms the first address has at least 1 coin left.
				coins, _ := getAddressBalance(t, entries[0].Address.String())
				require.True(t, coins >= 1e6)
			},
		},
		{
			// Send 0.001 coin from the third address to the second address.
			// Set the second as change address, so the 0.499 change coin will also be sent to the second address.
			// After sending, the second address should have 1 coin and 1 coin hour.
			name: "send with -c(change address) -a(from address) options",
			args: func() []string {
				return []string{"send", fn, "-c", entries[1].Address.String(),
					"-a", entries[2].Address.String(), entries[1].Address.String(), "0.001"}
			},
			checkTxn: func(t *testing.T, txid string) {
				tx := getTransaction(t, txid)
				// Confirms the second address receives 0.5 coin and 0 coinhour in this transaction
				checkCoinsAndCoinhours(t, tx, entries[1].Address.String(), 5e5, 0)
				// Confirms the second address have 1 coin and 1 coin hour
				coins, hours := getAddressBalance(t, entries[1].Address.String())
				require.Equal(t, uint64(1e6), coins)
				require.Equal(t, uint64(1), hours)
			},
		},
		{
			// Send 1 coin from second to the the third address, this will spend three outputs(0.2, 0.3. 0.5 coin),
			// and burn out the remaining 1 coin hour.
			name: "send to burn all coin hour",
			args: func() []string {
				return []string{"send", fn, "-a", entries[1].Address.String(),
					entries[2].Address.String(), "1"}
			},
			checkTxn: func(t *testing.T, txid string) {
				// Confirms that the third address has 1 coin and 0 coin hour
				coins, hours := getAddressBalance(t, entries[2].Address.String())
				require.Equal(t, uint64(1e6), coins)
				require.Equal(t, uint64(0), hours)
			},
		},
		{
			// Send with 0 coin hour, this test should fail.
			name: "send 0 coin hour",
			args: func() []string {
				return []string{"send", fn, "-a", entries[2].Address.String(),
					entries[1].Address.String(), "1"}
			},
			errMsg:   []byte("See 'skycoin-cli send --help'\nError: Transaction has zero coinhour fee"),
			checkTxn: func(t *testing.T, txid string) {},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := execCommandCombinedOutput(tc.args()...)

			output = bytes.TrimRight(output, "\n")

			if tc.errMsg != nil {
				require.Equal(t, err.Error(), "exit status 1")
				require.Equal(t, tc.errMsg, output)
				return
			}

			require.NoError(t, err)

			// output: "txid:$txid_string"
			// split the output to get txid value
			v := bytes.Split(output, []byte(":"))
			require.Len(t, v, 2)
			txid := string(v[1])
			fmt.Println("txid:", txid)
			_, err = cipher.SHA256FromHex(txid)
			require.NoError(t, err)

			// Wait until transaction is confirmed.
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

			tc.checkTxn(t, txid)
		})
	}
}

func TestLiveSendNotEnoughDecimals(t *testing.T) {
	if !doLive(t) {
		return
	}

	fn := requireWalletEnv(t)

	// prepares wallet and confirms the wallet has at least 2 coins and 16 coin hours.
	w, _, _ := prepareAndCheckWallet(t, 2e6, 16)

	if w.IsEncrypted() {
		t.Skip("CLI wallet integration tests do not support encrypted wallets yet")
		return
	}

	entries, err := w.GetEntries()
	require.NoError(t, err)

	// Send with too small decimal value
	// CLI send is a litte bit slow, almost 300ms each. so we only test 20 invalid decimal coin.
	errMsg := []byte("See 'skycoin-cli send --help'\nError: Transaction violates soft constraint: invalid amount, too many decimal places")
	for i := uint64(1); i < uint64(20); i++ {
		v, err := droplet.ToString(i)
		require.NoError(t, err)
		name := fmt.Sprintf("send %v", v)
		t.Run(name, func(t *testing.T) {
			output, err := execCommandCombinedOutput("send", fn, entries[0].Address.String(), v)
			require.Error(t, err)
			require.Equal(t, err.Error(), "exit status 1")
			output = bytes.TrimRight(output, "\n")
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

	fn := requireWalletEnv(t)

	// prepares wallet and confirms the wallet has at least 2 coins and 2 coin hours.
	w, totalCoins, _ := prepareAndCheckWallet(t, 2e6, 2)

	if w.IsEncrypted() {
		t.Skip("CLI wallet integration tests do not support encrypted wallets yet")
		return
	}

	var tmpCSVFile string

	defer func() {
		if tmpCSVFile != "" {
			err := os.Remove(tmpCSVFile)
			require.NoError(t, err)
		}
	}()

	entries, err := w.GetEntries()
	require.NoError(t, err)

	tt := []struct {
		name     string
		args     func() []string
		errMsg   string
		checkTxn func(t *testing.T, txid string)
	}{
		{
			// Send all coins to the first address to one output.
			name: "send all coins to the first address",
			args: func() []string {
				coins, err := droplet.ToString(totalCoins)
				require.NoError(t, err)
				return []string{"createRawTransaction", fn, entries[0].Address.String(), coins}
			},
			checkTxn: func(t *testing.T, txid string) {
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
						entries[1].Address.String(),
						"0.5",
					},
					{
						entries[2].Address.String(),
						"0.5",
					},
				}

				v, err := json.Marshal(addrCoins)
				require.NoError(t, err)

				return []string{"createRawTransaction", fn, "-m", string(v)}
			},
			checkTxn: func(t *testing.T, txid string) {
				// Confirms the first address has at least 1 coin left.
				coins, _ := getAddressBalance(t, entries[0].Address.String())
				require.True(t, coins >= 1e6)
			},
		},
		{
			// Send 0.5 coin to the second address.
			// Send 0.5 coin to the third address.
			// After sending, the first address should have at least 1 coin left.
			name: "send to multiple address with --csv option",
			args: func() []string {
				fields := [][]string{
					{entries[1].Address.String(), "0.5"},
					{entries[2].Address.String(), "0.5"},
				}

				f, err := ioutil.TempFile("", "createrawtxn")
				require.NoError(t, err)
				defer f.Close()

				w := csv.NewWriter(f)

				err = w.WriteAll(fields)
				require.NoError(t, err)

				w.Flush()
				err = w.Error()
				require.NoError(t, err)

				tmpCSVFile = f.Name()

				return []string{"createRawTransaction", fn, "--csv", tmpCSVFile}
			},
			checkTxn: func(t *testing.T, txid string) {
				// Confirms the first address has at least 1 coin left.
				coins, _ := getAddressBalance(t, entries[0].Address.String())
				require.True(t, coins >= 1e6)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Create raw transaction first
			output, err := execCommandCombinedOutput(tc.args()...)
			if err != nil {
				t.Fatalf("err: %v, output: %v", err, string(output))
				return
			}
			require.NoError(t, err)
			output = bytes.TrimRight(output, "\n")
			if bytes.Contains(output, []byte("Error:")) {
				require.Equal(t, tc.errMsg, string(output))
				return
			}

			// Broadcast transaction
			output, err = execCommandCombinedOutput("broadcastTransaction", string(output))
			require.NoError(t, err, string(output))

			txid := string(bytes.TrimRight(output, "\n"))
			fmt.Println("txid:", txid)
			_, err = cipher.SHA256FromHex(txid)
			require.NoError(t, err)

			// Wait until transaction is confirmed.
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

			tc.checkTxn(t, txid)
		})
	}

	// Send with too small decimal value
	errMsg := "Error: Transaction violates soft constraint: invalid amount, too many decimal places"
	for i := uint64(1); i < uint64(20); i++ {
		v, err := droplet.ToString(i)
		require.NoError(t, err)
		name := fmt.Sprintf("send %v", v)
		t.Run(name, func(t *testing.T) {
			output, err := execCommandCombinedOutput("createRawTransaction", fn, entries[0].Address.String(), v)
			require.Error(t, err)
			output = bytes.Trim(output, "\n")
			require.Equal(t, errMsg, string(output))
		})
	}
}

func getTransaction(t *testing.T, txid string) *cli.TxnResult {
	output, err := execCommandCombinedOutput("transaction", txid)
	if err != nil {
		t.Log(string(output))
		return nil
	}

	require.NoError(t, err)

	var tx cli.TxnResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&tx)
	require.NoError(t, err)

	return &tx
}

func isTxConfirmed(t *testing.T, txid string) bool {
	tx := getTransaction(t, txid)
	require.NotNil(t, tx)
	require.NotNil(t, tx.Transaction)
	require.NotNil(t, tx.Transaction.Status)
	return tx.Transaction.Status.Confirmed
}

// checkCoinhours checks if the address coinhours in transaction are correct
func checkCoinsAndCoinhours(t *testing.T, tx *cli.TxnResult, addr string, coins, coinhours uint64) { //nolint:unparam
	addrCoinhoursMap := make(map[string][]readable.TransactionOutput)
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
// Returns wallet.Wallet, total coin, total hours.
// Confirms that the wallet meets the minimal requirements of coins and coinhours.
func prepareAndCheckWallet(t *testing.T, miniCoins, miniCoinHours uint64) (wallet.Wallet, uint64, uint64) { //nolint:unparam
	walletFile := requireWalletEnv(t)
	// Checks if the wallet does exist
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		t.Fatalf("Wallet file: %v does not exist", walletFile)
	}

	// Loads the wallet
	w, err := wallet.Load(walletFile)
	if err != nil {
		t.Fatalf("Load wallet failed: %v", err)
	}

	el, err := w.EntriesLen()
	require.NoError(t, err)

	if el < 3 {
		// Generates addresses
		_, err = w.GenerateAddresses(uint64(3 - el))
		if err != nil {
			t.Fatalf("Wallet generateAddress failed: %v", err)
		}
	}

	outputs := getWalletOutputs(t, walletFile)
	// Confirms the wallet is not empty.
	if len(outputs) == 0 {
		t.Fatalf("Wallet %v has no coin", walletFile)
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

	if err := wallet.Save(w, filepath.Dir(walletFile)); err != nil {
		t.Fatalf("%v", err)
	}
	return w, totalCoins, totalCoinhours
}

func getAddressBalance(t *testing.T, addr string) (uint64, uint64) {
	output, err := execCommandCombinedOutput("addressBalance", addr)
	require.NoError(t, err, string(output))

	var addrBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrBalance)
	require.NoError(t, err)
	coins, err := droplet.FromString(addrBalance.Confirmed.Coins)
	require.NoError(t, err)

	hours, err := strconv.ParseUint(addrBalance.Confirmed.Hours, 10, 64)
	require.NoError(t, err)
	return coins, hours
}

func getWalletOutputs(t *testing.T, walletFile string) readable.UnspentOutputs {
	output, err := execCommandCombinedOutput("walletOutputs", walletFile)
	require.NoError(t, err, string(output))

	var wltOutput cli.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltOutput)
	require.NoError(t, err)

	return wltOutput.Outputs.HeadOutputs
}

func TestStableWalletHistory(t *testing.T) {
	if !doStable(t) {
		return
	}

	seed := "visit harbor excite frown flat nothing reduce price wrist label destroy citizen"
	wlt := createTempWallet(t, "test-stable-wallet-history", seed, false, nil)

	output, err := execCommandCombinedOutput("walletHistory", wlt.Meta.Filename)
	require.NoError(t, err, output)

	var history []cli.AddrHistory
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&history)
	require.NoError(t, err)

	var expect []cli.AddrHistory
	checkGoldenFile(t, "wallet-history.golden", TestData{history, &expect})
}

func TestLiveWalletHistory(t *testing.T) {
	if !doLive(t) {
		return
	}

	fn := requireWalletEnv(t)

	output, err := execCommandCombinedOutput("walletHistory", fn)
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
		result string
		errMsg string
	}{
		{
			name:   "no signature",
			dbPath: "../../visor/testdata/data.db.nosig",
			errMsg: "Error: checkdb failed: Signature not found for block seq=1000 hash=71852c1a8ab5e470bd14e5fce8e1116697151181a188d4262b545542fb3d526c\n",
		},
		{
			name:   "invalid database",
			dbPath: "../../visor/testdata/data.db.garbage",
			errMsg: "Error: open db failed: invalid database\n",
		},
		{
			name:   "valid database",
			dbPath: "../../api/integration/testdata/blockchain-180.db",
			result: "check db success\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := execCommandCombinedOutput("checkdb", tc.dbPath)
			if err != nil {
				t.Log(string(output))
				require.Equal(t, tc.errMsg, string(output))
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.result, string(output))
		})
	}
}

func TestVersion(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	// Gets version in json format.
	output, err := execCommandCombinedOutput("version", "-j")
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
	output, err = execCommandCombinedOutput("version")
	require.NoError(t, err)

	// Confirms the result contains 4 version componments
	output = bytes.TrimRight(output, "\n")
	vers := bytes.Split(output, []byte("\n"))
	require.Len(t, vers, 4)
}

func TestStableWalletCreateXPubFlow(t *testing.T) {
	if !doStable(t) {
		return
	}

	// The flow to create an xpub wallet is:
	// - Create a bip44 wallet
	// - Export an xpub key from the bip44 wallet
	// - Create an xpub wallet from the xpub key

	dir, clean := createTempWalletDir(t)
	defer clean()

	bip44Filename := filepath.Join(dir, "bip44.wlt")

	// Create a bip44 wallet
	args := []string{"walletCreate", bip44Filename, "-t", "bip44", "-n", "10"}
	_, err := execCommandCombinedOutput(args...)
	require.NoError(t, err)

	// Export the xpub key from the bip44 wallet subpath 0'/0
	args = []string{"walletKeyExport", bip44Filename, "-k", "xpub", "--path", "0/0"}
	output, err := execCommandCombinedOutput(args...)
	require.NoError(t, err)

	xpub := strings.TrimSpace(string(output))

	xpubFilename := filepath.Join(dir, "xpub.wlt")

	// Create an xpub wallet
	args = []string{"walletCreate", xpubFilename, "-t", "xpub", "--xpub", xpub, "-n", "10"}
	_, err = execCommandCombinedOutput(args...)
	require.NoError(t, err)

	// Compare the entries of both wallets: they should match
	w, err := wallet.Load(bip44Filename)
	require.NoError(t, err)

	w2, err := wallet.Load(xpubFilename)
	require.NoError(t, err)

	entries, err := w.GetEntries()
	require.NoError(t, err)

	for i, e := range entries {
		e2, err := w2.GetEntryAt(i)
		require.NoError(t, err)
		require.Equal(t, e.Public, e2.Public)
		require.Equal(t, e.Address, e2.Address)
		require.False(t, e.Secret.Null())
		require.True(t, e2.Secret.Null())
		require.Equal(t, e.ChildNumber, uint32(i))
		require.Equal(t, e.ChildNumber, e2.ChildNumber)
		require.Equal(t, e.Change, e2.Change)
	}
}

func TestStableWalletCreate(t *testing.T) {
	if !doStable(t) {
		return
	}

	tt := []struct {
		name        string
		filename    string
		args        []string
		setup       func(t *testing.T)
		errMsg      string
		errMsgFunc  func(filename string) string
		checkWallet func(t *testing.T, w wallet.Wallet)
	}{
		{
			name: "generate wallet with -r option",
			args: []string{"-r", "-e=false", "-t", wallet.WalletTypeDeterministic},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				// Confirms the seed is a valid hex string
				_, err := hex.DecodeString(w.Seed())
				require.NoError(t, err)

				// Confirms the label is empty
				require.Equal(t, "test-stable-wallet-create", w.Label())
			},
		},
		{
			name:   "generate wallet with -r option bip44",
			args:   []string{"-r", "-e=false", "-t", wallet.WalletTypeBip44},
			errMsg: "Error: -r can't be used for bip44 wallets\n",
		},
		{
			name: "generate wallet with -m option",
			args: []string{"-m", "-e=false"},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				// Confirms the seed has 12 words
				seed := w.Seed()
				words := strings.Split(seed, " ")
				require.Len(t, words, 12)

				err := bip39.ValidateMnemonic(seed)
				require.NoError(t, err)

				// Confirms the label is empty
				require.Equal(t, "test-stable-wallet-create", w.Label())
			},
		},
		{
			name: "generate wallet with -s option",
			args: []string{"-e=false", "-t", wallet.WalletTypeDeterministic, "-s", "great duck trophy inhale dad pluck include maze smart mechanic ring merge"},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "test-stable-wallet-create", w.Label())

				require.Equal(t, "great duck trophy inhale dad pluck include maze smart mechanic ring merge", w.Seed())
				e, err := w.GetEntryAt(0)
				require.NoError(t, err)
				require.Equal(t, "2amA8sxKJhNRp3wfWrE5JfTEUjr9S3C2BaU", e.Address.String())
				require.Equal(t, "02b4a4b63f2f8ba56f9508712815eca3c088693333715eaf7a73275d8928e1be5a", e.Public.Hex())
				require.Equal(t, "f4a281d094a6e9e95a84c23701a7d01a0e413c838758e94ad86a10b9b83e0434", e.Secret.Hex())
			},
		},
		{
			name: "generate wallet with -s option bip44",
			args: []string{"-e=false", "-t", wallet.WalletTypeBip44, "-s", "great duck trophy inhale dad pluck include maze smart mechanic ring merge"},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "test-stable-wallet-create", w.Label())

				require.Equal(t, "great duck trophy inhale dad pluck include maze smart mechanic ring merge", w.Seed())
				e, err := w.GetEntryAt(0)
				require.NoError(t, err)
				require.Equal(t, "skPFvAokn63RTHh8MR8cZuHUpzZiutFeks", e.Address.String())
				require.Equal(t, "0381a3a0ed879eae12a612d24c73c39ff1d9f3c238e10ecf29c318db11e84e1143", e.Public.Hex())
				require.Equal(t, "7519c63c890d593a11eeebfe4b4552f3d8a01094c086262e7c68a3e7adc61677", e.Secret.Hex())
			},
		},
		{
			name: "generate wallet with -s option and seed-passphrase bip44",
			args: []string{"-e=false", "-t", wallet.WalletTypeBip44, "--seed-passphrase", "foobar", "-s", "great duck trophy inhale dad pluck include maze smart mechanic ring merge"},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "test-stable-wallet-create", w.Label())

				e, err := w.GetEntryAt(0)
				require.NoError(t, err)
				require.Equal(t, "foobar", w.SeedPassphrase())
				require.Equal(t, "great duck trophy inhale dad pluck include maze smart mechanic ring merge", w.Seed())
				require.Equal(t, "29bFTouyAtgRhdhHqiC6kccom3avjBpfnxS", e.Address.String())
				require.Equal(t, "02dbaefd35105e27ac7428505186844dcb9472021a5fdfb057c106722f3396bdf6", e.Public.Hex())
				require.Equal(t, "c2edb45b4c5c7828d6b879a163cb2bf016746e9c03efe5df24e9ded878b0b772", e.Secret.Hex())
			},
		},
		{
			name:   "generate wallet with -s option and seed-passphrase deterministic",
			args:   []string{"-e=false", "-t", wallet.WalletTypeDeterministic, "--seed-passphrase", "foobar", "-s", "great duck trophy inhale dad pluck include maze smart mechanic ring merge"},
			errMsg: "Error: 400 Bad Request - seedPassphrase is only used for \"bip44\" wallets\n",
		},
		{
			name: "generate wallet with -n option, bip44",
			args: []string{"-e=false", "-n", "5", "-t", wallet.WalletTypeBip44},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, wallet.WalletTypeBip44, w.Type())
				require.Equal(t, "test-stable-wallet-create", w.Label())

				// get external entries length
				extLen, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, extLen, 5)

				// get change entries length
				chgLen, err := w.EntriesLen(wallet.OptionChange(true))
				require.NoError(t, err)
				require.Equal(t, chgLen, 1)
			},
		},
		{
			name: "generate wallet with -n option",
			args: []string{"-e=false", "-n", "5"},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, wallet.WalletTypeDeterministic, w.Type())
				require.Equal(t, "test-stable-wallet-create", w.Label())
				// Confirms wallet has 5 address entries
				l, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, l, 5)
			},
		},
		{
			name: "generate wallet with duplicate seed for deterministic wallet",
			args: []string{"-s", "session eyebrow giant vote volcano eight code ahead return yard essay copy", "-e=false"},
			setup: func(t *testing.T) {
				seed := "session eyebrow giant vote volcano eight code ahead return yard essay copy"
				createTempWallet(t, "test-stable-wallet-create", seed, false, nil)
			},
			errMsg: "Error: 400 Bad Request - fingerprint conflict for \"deterministic\" wallet\n",
		},
		{
			name: "encrypt=true",
			args: []string{"-e", "-p", "pwd"},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				// Confirms the wallet is encrypted
				require.True(t, w.IsEncrypted())
				require.Empty(t, w.Seed())
				require.Empty(t, w.LastSeed())

				// Confirms the secrets in address entries are empty
				entries, err := w.GetEntries()
				require.NoError(t, err)
				for _, e := range entries {
					require.Equal(t, cipher.SecKey{}, e.Secret)
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(t)
			}

			// Run command with arguments
			args := append([]string{"walletCreate", "test-stable-wallet-create"}, tc.args...)
			output, err := execCommandCombinedOutput(args...)
			if err != nil {
				require.EqualError(t, err, "exit status 1")
				require.Equal(t, tc.errMsg, string(output))
				return
			}

			require.Empty(t, tc.errMsg, "the command did not return an error, but we expected one")

			var w api.WalletResponse
			require.NoError(t, json.Unmarshal(output, &w))

			c := newClient()
			walletDir := getWalletDir(t, c)

			// Confirms the wallet does exist
			walletPath := filepath.Join(walletDir, w.Meta.Filename)
			_, err = os.Stat(walletPath)
			require.NoError(t, err)

			// Loads the wallet and confirms that the wallet has the same seed
			wl, err := wallet.Load(walletPath)
			require.NoError(t, err)

			if !wl.IsEncrypted() {
				// Confirms all entries and lastSeed are derived from seed.
				checkWalletEntriesAndLastSeed(t, wl)
			}

			// Checks the wallet with provided checking method.
			require.NotNil(t, wl)
			require.NotNil(t, tc.checkWallet)
			tc.checkWallet(t, wl)
		})
	}
}

// checkWalletEntriesAndLastSeed confirms the wallet entries and lastSeed are derivied
// from the seed.
func checkWalletEntriesAndLastSeed(t *testing.T, w wallet.Wallet) {
	switch w.Type() {
	case wallet.WalletTypeDeterministic:
		checkWalletEntriesAndLastSeedDeterministic(t, w)
	case wallet.WalletTypeBip44:
		checkWalletEntriesAndLastSeedBip44(t, w)
	case wallet.WalletTypeCollection:
		checkWalletEntriesAndLastSeedCollection(t, w)
	default:
		t.Fatalf("unknown wallet type %q", w.Type())
	}
}

func checkWalletEntriesAndLastSeedDeterministic(t *testing.T, w wallet.Wallet) {
	seed := w.Seed()
	require.NotEmpty(t, seed)

	entries, err := w.GetEntries()
	require.NoError(t, err)
	newSeed, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), len(entries))
	require.Len(t, seckeys, len(entries))
	for i, sk := range seckeys {
		require.Equal(t, entries[i].Secret, sk)
		pk := cipher.MustPubKeyFromSecKey(sk)
		require.Equal(t, entries[i].Public, pk)
	}

	lastSeed := w.LastSeed()
	require.NotEmpty(t, lastSeed)
	require.Equal(t, lastSeed, hex.EncodeToString(newSeed))

	require.Empty(t, w.SeedPassphrase())
}

func checkWalletEntriesAndLastSeedBip44(t *testing.T, w wallet.Wallet) {
	seed := w.Seed()
	require.NotEmpty(t, seed)
	// bip44 wallet seed must be a valid bip39 mnemonic
	err := bip39.ValidateMnemonic(seed)
	require.NoError(t, err)

	entries, err := w.GetEntries()
	require.NoError(t, err)

	for _, e := range entries {
		require.False(t, e.Secret.Null())
		require.False(t, e.Public.Null())
	}

	// lastSeed is only for "deterministic" type wallet
	require.Empty(t, w.LastSeed())
}

func checkWalletEntriesAndLastSeedCollection(t *testing.T, w wallet.Wallet) {
	require.Empty(t, w.Seed())
	require.Empty(t, w.SeedPassphrase())
	require.Empty(t, w.LastSeed())
}

// TestLiveGUIInjectTransaction does almost the same procedure as TestCreateAndBroadcastRawTransaction.
// The only difference is we broadcast the raw transaction through the gui /injectTransaction api.
func TestLiveGUIInjectTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	fn := requireWalletEnv(t)
	require.NotEmpty(t, fn)

	c := newClient()
	// prepares wallet and confirms the wallet has at least 2 coins and 2 coin hours.
	w, totalCoins, _ := prepareAndCheckWallet(t, 2e6, 2)

	if w.IsEncrypted() {
		t.Skip("CLI wallet integration tests do not support encrypted wallets yet")
		return
	}

	entries, err := w.GetEntries()
	require.NoError(t, err)

	tt := []struct {
		name     string
		args     func() []string
		errMsg   []byte
		checkTxn func(t *testing.T, txid string)
	}{
		{
			// Send all coins to the first address to one output.
			name: "send all coins to the first address",
			args: func() []string {
				coins, err := droplet.ToString(totalCoins)
				require.NoError(t, err)
				return []string{"createRawTransaction", fn, entries[0].Address.String(), coins}
			},
			checkTxn: func(t *testing.T, txid string) {
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
						entries[1].Address.String(),
						"0.5",
					},
					{
						entries[2].Address.String(),
						"0.5",
					},
				}

				v, err := json.Marshal(addrCoins)
				require.NoError(t, err)

				return []string{"createRawTransaction", fn, "-m", string(v)}
			},
			checkTxn: func(t *testing.T, txid string) {
				// Confirms the first address has at least 1 coin left.
				coins, _ := getAddressBalance(t, entries[0].Address.String())
				require.True(t, coins >= 1e6)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Create raw transaction first
			output, err := execCommandCombinedOutput(tc.args()...)
			if err != nil {
				t.Fatalf("err: %v, output: %v", err, string(output))
				return
			}

			require.NoError(t, err)
			output = bytes.TrimRight(output, "\n")
			if bytes.Contains(output, []byte("Error:")) {
				require.Equal(t, tc.errMsg, output)
				return
			}

			// Broadcast raw transaction with gui /injectTransaction
			txid, err := c.InjectEncodedTransaction(string(output))
			require.NoError(t, err)

			txid = strings.TrimRight(txid, "\n")
			t.Logf("txid: %s", txid)
			_, err = cipher.SHA256FromHex(txid)
			require.NoError(t, err)

			// Wait until transaction is confirmed.
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

			tc.checkTxn(t, txid)
		})
	}
}

func TestEncryptWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	tt := []struct {
		name        string
		args        []string
		setup       func(t *testing.T) (string, func())
		errMsg      []byte
		errWithHelp bool
		checkWallet func(t *testing.T, w wallet.Wallet)
	}{
		{
			name:  "wallet is not encrypted",
			args:  []string{"-p", "pwd"},
			setup: createUnencryptedWallet,
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.True(t, w.IsEncrypted())
				require.Empty(t, w.Seed())
				require.Empty(t, w.LastSeed())

				// Confirms that secrets in address entries are empty
				entries, err := w.GetEntries()
				require.NoError(t, err)
				for _, e := range entries {
					require.Equal(t, cipher.SecKey{}, e.Secret)
				}
			},
		},
		{
			name:   "wallet is encrypted",
			args:   []string{"-p", "pwd"},
			setup:  createEncryptedWallet,
			errMsg: []byte("Error: wallet is encrypted\n"),
		},
		{
			name: "wallet doesn't exist",
			args: []string{"-p", "pwd"},
			setup: func(t *testing.T) (string, func()) {
				return "not-exist.wlt", func() {}
			},
			errWithHelp: true,
			errMsg:      []byte("not-exist.wlt\" doesn't exist"),
		},
	}

	for _, tc := range tt {
		for _, ct := range cryptoTypes {
			name := fmt.Sprintf("name=%v crypto type=%v", tc.name, ct)
			t.Run(name, func(t *testing.T) {
				walletFile, clean := tc.setup(t)
				defer clean()
				args := append([]string{"encryptWallet", walletFile, "-x", string(ct)}, tc.args[:]...)
				output, err := execCommandCombinedOutput(args...)
				if err != nil {
					require.EqualError(t, err, "exit status 1")
					if tc.errWithHelp {
						require.True(t, bytes.Contains(output, tc.errMsg), string(output))
					} else {
						require.Equal(t, tc.errMsg, output)
					}
					return
				}

				w, err := wallet.Load(walletFile)
				require.NoError(t, err)
				tc.checkWallet(t, w)
			})
		}
	}
}

func TestDecryptWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	tt := []struct {
		name        string
		args        []string
		setup       func(t *testing.T) (string, func())
		errMsg      []byte
		errWithHelp bool
		checkWallet func(t *testing.T, w wallet.Wallet)
	}{
		{
			name:  "wallet is encrypted",
			args:  []string{"-p", "pwd"},
			setup: createEncryptedWallet,
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.False(t, w.IsEncrypted())
				require.Empty(t, w.Secrets())
				require.NotEmpty(t, w.Seed())
				require.NotEmpty(t, w.LastSeed())

				entries, err := w.GetEntries()
				require.NoError(t, err)
				for _, e := range entries {
					require.False(t, e.Secret.Null())
				}
			},
		},
		{
			name:   "wallet is not encrypted",
			args:   []string{"-p", "pwd"},
			setup:  createUnencryptedWallet,
			errMsg: []byte("Error: wallet is not encrypted\n"),
		},
		{
			name:   "invalid password",
			args:   []string{"-p", "wrong password"},
			setup:  createEncryptedWallet,
			errMsg: []byte("Error: invalid password\n"),
		},
		{
			name: "wallet doesn't exist",
			args: []string{"-p", "pwd"},
			setup: func(t *testing.T) (string, func()) {
				return "not-exist.wlt", func() {}
			},
			errWithHelp: true,
			errMsg:      []byte("not-exist.wlt\" doesn't exist"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			walletFile, clean := tc.setup(t)
			defer clean()
			args := append([]string{"decryptWallet", walletFile}, tc.args...)
			output, err := execCommandCombinedOutput(args...)
			if err != nil {
				require.EqualError(t, err, "exit status 1")
				if tc.errWithHelp {
					require.True(t, bytes.Contains(output, tc.errMsg), string(output))
				} else {
					require.Equal(t, tc.errMsg, output)
				}
				return
			}

			w, err := wallet.Load(walletFile)
			require.NoError(t, err)
			tc.checkWallet(t, w)
		})
	}
}

func TestWalletShowSeed(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	tt := []struct {
		name         string
		args         []string
		setup        func(t *testing.T) (string, func())
		errWithHelp  bool
		errMsg       []byte
		expectOutput []byte
	}{
		{
			name:         "unencrypted wallet",
			setup:        createUnencryptedWallet,
			expectOutput: []byte("exchange stage green marine palm tobacco decline shadow cereal chapter lamp copy\n"),
		},
		{
			name:         "unencrypted wallet with -j option",
			args:         []string{"-j"},
			setup:        createUnencryptedWallet,
			expectOutput: []byte("{\n    \"seed\": \"exchange stage green marine palm tobacco decline shadow cereal chapter lamp copy\"\n}\n"),
		},
		{
			name:         "encrypted wallet",
			setup:        createEncryptedWallet,
			args:         []string{"-p", "pwd"},
			expectOutput: []byte("exchange stage green marine palm tobacco decline shadow cereal chapter lamp copy\n"),
		},
		{
			name:         "encrypted wallet with -j option",
			setup:        createEncryptedWallet,
			args:         []string{"-p", "pwd", "-j"},
			expectOutput: []byte("{\n    \"seed\": \"exchange stage green marine palm tobacco decline shadow cereal chapter lamp copy\"\n}\n"),
		},
		{
			name:   "encrypted wallet with invalid password",
			setup:  createEncryptedWallet,
			args:   []string{"-p", "wrong password"},
			errMsg: []byte("Error: invalid password\n"),
		},
		{
			name: "wallet doesn't exist",
			setup: func(t *testing.T) (string, func()) {
				return "not-exist.wlt", func() {}
			},
			errWithHelp: true,
			errMsg:      []byte("not-exist.wlt\" doesn't exist"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			wltFilename, clean := tc.setup(t)
			defer clean()
			args := append([]string{"showSeed", wltFilename}, tc.args...)
			output, err := execCommandCombinedOutput(args...)
			if err != nil {
				require.EqualError(t, err, "exit status 1")
				if tc.errWithHelp {
					require.True(t, bytes.Contains(output, tc.errMsg), string(output))
				} else {
					require.Equal(t, tc.errMsg, output)
				}
				return
			}

			require.Equal(t, tc.expectOutput, output)
		})
	}
}

func loadDeterministicWalletFromBytes(t *testing.T, data []byte) *deterministic.Wallet {
	var w deterministic.Wallet
	err := w.Deserialize(data)
	require.NoError(t, err)
	return &w
}
func getWalletDir(t *testing.T, c *api.Client) string {
	wf, err := c.WalletFolderName()
	if err != nil {
		t.Fatalf("%v", err)
	}
	return wf.Address
}
