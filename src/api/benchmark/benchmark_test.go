package benchmark

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	vegeta "github.com/tsenart/vegeta/lib"
)

const (
	testModeStable           = "stable"
	testModeLive             = "live"
	testModeDisableWalletAPI = "disable-wallet-api"
	testModeDisableSeedAPI   = "disable-seed-api"
)

func nodeAddress() string {
	addr := os.Getenv("SKYCOIN_NODE_HOST")
	if addr == "" {
		return "http://127.0.0.1:6420"
	}
	return addr
}

func doLiveOrStable(t *testing.B) bool {
	if enabled() {
		switch mode(t) {
		case testModeStable, testModeLive:
			return true
		}
	}

	t.Skip("Live and stable tests disabled")
	return false
}

func enabled() bool {
	return os.Getenv("SKYCOIN_BENCHMARK_TESTS") == "1"
}

func mode(t *testing.B) string {
	mode := os.Getenv("SKYCOIN_BENCHMARK_TEST_MODE")
	switch mode {
	case "":
		mode = testModeStable
	case testModeLive,
		testModeStable,
		testModeDisableWalletAPI,
		testModeDisableSeedAPI:
	default:
		t.Fatal("Invalid test mode, must be stable, live or disable-wallet-api")
	}
	return mode
}

func BenchmarkGetWalletDirWrk(t *testing.B) {
	if !doLiveOrStable(t) {
		return
	}
	url := nodeAddress() + "/api/v1/wallets/folderName"
	cmd := exec.Command("wrk", "-d", "10", "-c", "10", url)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("got err. stdout: \n%s", out))
	log.Printf("%s", out)
}

func BenchmarkGetWalletDirVegeta(t *testing.B) {
	if !doLiveOrStable(t) {
		return
	}
	url := nodeAddress() + "/api/v1/wallets/folderName"
	rate := uint64(8000) // per second
	duration := 10 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    url,
	})
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "folderName") {
		if res.Error != "" {
			log.Printf("getWalletDir. err: %v", res.Error)
		}
		metrics.Add(res)
	}
	metrics.Close()
	reporter := vegeta.NewTextReporter(&metrics)
	reporter.Report(os.Stdout)
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func BenchmarkGetWalletDirHey(t *testing.B) {
	if !doLiveOrStable(t) {
		return
	}
	url := nodeAddress() + "/api/v1/wallets/folderName"
	cmd := exec.Command("hey", "-q", "15000", "-z", "10s", url)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("got err. stdout: \n%s", out))
	log.Printf("%s", out)
}

func BenchmarkGetWalletDirBombardier(t *testing.B) {
	if !doLiveOrStable(t) {
		return
	}
	url := nodeAddress() + "/api/v1/wallets/folderName"
	cmd := exec.Command("bombardier", "-n", "10000", url)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("got err. stdout: \n%s", out))
	log.Printf("%s", out)
}

func BenchmarkGetWalletDirSniper(t *testing.B) {
	if !doLiveOrStable(t) {
		return
	}
	url := nodeAddress() + "/api/v1/wallets/folderName"
	cmd := exec.Command("sniper", "-t", "10", "-n", "100", url)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("got err. stdout: \n%s", out))
	log.Printf("%s", out)
}

func BenchmarkGetWalletDirGobench(t *testing.B) {
	if !doLiveOrStable(t) {
		return
	}
	url := nodeAddress() + "/api/v1/wallets/folderName"
	cmd := exec.Command("gobench", "-t", "10", "-u", url)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("got err. stdout: \n%s", out))
	log.Printf("%s", out)
}
