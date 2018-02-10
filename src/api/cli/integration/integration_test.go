package cli_integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/require"
)

const (
	binaryName = "skycoin-cli"
	walletName = "integration_test.wlt"

	testModeStable = "stable"
	testModeLive   = "live"
)

var (
	binaryPath string
	walletDir  string
)

func TestGenerateAddresses(t *testing.T) {
	output, err := exec.Command(binaryPath, "generateAddresses").CombinedOutput()
	require.NoError(t, err)
	o := strings.Trim(string(output), "\n")
	require.Equal(t, "7g3M372kxwNwwQEAmrronu4anXTW8aD1XC", o)

	wltPath := filepath.Join(walletDir, walletName)
	var w wallet.ReadableWallet
	loadJSON(t, wltPath, &w)

	var expect wallet.ReadableWallet
	loadJSON(t, "generateAddresses.golden", &expect)
	require.Equal(t, expect, w)
}

func TestStabeStatus(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "status").CombinedOutput()
	require.NoError(t, err)
	var ret struct {
		webrpc.StatusResult
		RPCAddress string `json:"webrpc_address"`
	}

	if err := json.NewDecoder(bytes.NewReader(output)).Decode(&ret); err != nil {
		fmt.Fprintf(os.Stderr, "Decode result failed: %v", err)
		os.Exit(1)
	}

	// TimeSinceLastBlock is not stable
	ret.TimeSinceLastBlock = ""

	var expect struct {
		webrpc.StatusResult
		RPCAddress string `json:"webrpc_address"`
	}
	loadJSON(t, "status.golden", &expect)

	require.Equal(t, expect, ret)
}

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

	dir, clean, err := createTempWalletFile("integration_test.wlt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer clean()

	walletDir = dir

	os.Setenv("WALLET_DIR", dir)
	os.Setenv("WALLET_NAME", walletName)

	ret := m.Run()
	// Remove the geneate cli binary file.
	if err := os.Remove(binaryPath); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Delete %v failed: %v", binaryName, err))
		os.Exit(1)
	}

	os.Exit(ret)
}

// createTempWalletFile creates a temporary dir, and copy the ./integration_test.wlt file to dir.
// returns the temporary dir path, cleanup callback function, and error if any.
func createTempWalletFile(fromPath string) (string, func(), error) {
	dir, err := ioutil.TempDir("", "integration_test")
	if err != nil {
		return "", nil, fmt.Errorf("Get temporary dir failed: %v", err)
	}

	// Copy the ./integration_test.wlt to the temporary dir.
	wltPath := filepath.Join(dir, walletName)
	f, err := os.Create(wltPath)
	if err != nil {
		return "", nil, fmt.Errorf("Create temporary file: %v failed: %v", wltPath, err)
	}

	defer f.Close()

	rf, err := os.Open(fromPath)
	if err != nil {
		return "", nil, fmt.Errorf("Open ./integration_test.wlt failed: %v", err)
	}

	defer rf.Close()
	io.Copy(f, rf)

	fun := func() {
		// Delete the temporary dir
		os.RemoveAll(dir)
	}

	return dir, fun, nil
}

func loadJSON(t *testing.T, filename string, obj interface{}) {
	f, err := os.Open(filename)
	require.NoError(t, err)
	defer f.Close()

	err = json.NewDecoder(f).Decode(obj)
	require.NoError(t, err)
}

func nodeAddress() string {
	addr := os.Getenv("SKYCOIN_NODE_HOST")
	if addr == "" {
		return "http://127.0.0.1:6420"
	}
	return addr
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
