/*
monitor-peers checks the status of peers.

It takes in a list of peers (ip:ports, newline separated, skipping comments and empty lines).
The tool connects to each of the peers, waits for the introduction packet (or times out)
and produces a report with the status of the peer (unreachable, reachable, sent_introduction, introduction_parameters).
Introduction_parameters were added in v0.25.0 so will be absent for earlier peer versions.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/util/logging"

	"github.com/skycoin/skycoin/src/daemon"

	"github.com/skycoin/skycoin/src/daemon/pex"
)

// Report contains remote `peers.txt` report data.
type Report []ReportEntry

// ReportEntry contains report data of a peer.
type ReportEntry struct {
	Address string
	Status  daemon.ConnectionState
}

const (
	defaultTimeout = "1s"
	defaultURL     = pex.DefaultPeerListURL
	addrWidth      = "48"
	reportFormat   = "%-" + addrWidth + "s\t%s\n"
)

var (
	logger = logging.MustGetLogger("main")
	help   = fmt.Sprintf(`monitor-peers checks the status of peers.

By default it downloads a peer list from %s. May be overridden with -peersurl flag.

The default timeout is %s. May be overridden with -timeout flag. The timeout is parsed by time.ParseDuration.

It generates a report for each peer which contains the peer address and status. Status may be one of the following:

- pending
No connection made.

- connected
Connection made, no introduction message received.

- introduced
Connection made, introduction message received.
`, defaultURL, defaultTimeout)
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n\nUsage of %s:\n", help, os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	peersURL := flag.String("peersurl", defaultURL, "URL to fetch peers.txt")
	timeoutStr := flag.String("timeout", defaultTimeout, "timeout for each peer")

	flag.Parse()

	timeout, err := time.ParseDuration(*timeoutStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Bad timeout:", *timeoutStr)
		os.Exit(1)
	}

	logger.Infof("Peer connection threshold is %v", timeout)
	peers, err := pex.GetPeerListFromURL(*peersURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	report := getPeersReport(peers, timeout)
	logger.Infof("Report:\n%v", buildReport(report))
}

func getPeersReport(peers []string, timeout time.Duration) Report {
	c := daemon.NewConnections()

	var wg sync.WaitGroup

	var reportMu sync.Mutex
	report := make(Report, 0, len(peers))

	for _, addr := range peers {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()

			entry := ReportEntry{
				Address: addr,
				Status:  c.CheckStatus(addr, timeout),
			}

			reportMu.Lock()
			defer reportMu.Unlock()
			report = append(report, entry)
		}(addr)
	}
	wg.Wait()

	return report
}

func buildReport(report Report) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(reportFormat, "Address", "Status"))
	for _, entry := range report {
		sb.WriteString(fmt.Sprintf(reportFormat, entry.Address, entry.Status))
	}

	return sb.String()
}
