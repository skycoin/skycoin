package benchmark

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tsenart/vegeta/lib"
)

func TestGetWalletDirWrk(t *testing.T) {
	url := "http://localhost:6420/api/v1/wallets/folderName"
	cmd := exec.Command("wrk", "-d", "10", "-c", "10", url)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("got err. stdout: \n%s", out))
	log.Printf("%s", out)
}

func TestGetWalletDirVegeta(t *testing.T) {
	rate := uint64(8000) // per second
	duration := 10 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:6420/api/v1/wallets/folderName",
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
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

func TestGetWalletDirHey(t *testing.T) {
	url := "http://localhost:6420/api/v1/wallets/folderName"
	cmd := exec.Command("hey", "-q", "15000", "-z", "10s", url)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("got err. stdout: \n%s", out))
	log.Printf("%s", out)
}

func TestGetWalletDirBombardier(t *testing.T) {
	url := "http://localhost:6420/api/v1/wallets/folderName"
	cmd := exec.Command("bombardier", "-n", "100000", url)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("got err. stdout: \n%s", out))
	log.Printf("%s", out)
}

func TestGetWalletDirSniper(t *testing.T) {
	url := "http://localhost:6420/api/v1/wallets/folderName"
	cmd := exec.Command("sniper", "-t", "10", "-n", "100", url)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("got err. stdout: \n%s", out))
	log.Printf("%s", out)
}

func TestGetWalletDirGobench(t *testing.T) {
	url := "http://localhost:6420/api/v1/wallets/folderName"
	cmd := exec.Command("gobench", "-t", "10", "-u", url)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("got err. stdout: \n%s", out))
	log.Printf("%s", out)
}
