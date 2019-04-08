/*
monitor-peers checks the status of peers.

It takes in a list of peers (ip:ports, newline separated, skipping comments and empty lines).
The tool connects to each of the peers, waits for the introduction packet (or times out)
and produces a report with the status of the peer (unreachable, reachable, sent_introduction, introduction_parameters).
Introduction_parameters were added in v0.25.0 so will be absent for earlier peer versions.
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/logging"
)

// Report contains remote `peers.txt` report data.
type Report []ReportEntry

// ReportEntry contains report data of a peer.
type ReportEntry struct {
	Address string
	Status  daemon.ConnectionState
}

const (
	defaultTimeout   = "1s"
	defaultPeersFile = "peers.txt"
	addrWidth        = "48"
	reportFormat     = "%-" + addrWidth + "s\t%s\n"
)

var (
	// ErrInvalidAddress is returned when an address appears malformed
	ErrInvalidAddress = errors.New("invalid address")
	// ErrNoLocalhost is returned if a localhost addresses are not allowed
	ErrNoLocalhost = errors.New("localhost address is not allowed")
	// ErrNotExternalIP is returned if an IP address is not a global unicast address
	ErrNotExternalIP = errors.New("IP is not a valid external IP")
	// ErrPortTooLow is returned if a port is less than 1024
	ErrPortTooLow = errors.New("port must be >= 1024")
)

var (
	logger = logging.MustGetLogger("main")
	// For removing inadvertent whitespace from addresses
	whitespaceFilter = regexp.MustCompile(`\s`)
	help             = fmt.Sprintf(`monitor-peers checks the status of peers.

By default it gets peers list from %s. May be overridden with -f flag.

The default timeout is %s. May be overridden with -timeout flag. The timeout is parsed by time.ParseDuration.

It generates a report for each peer which contains the peer address and status. Status may be one of the following:

- pending
No connection made.

- connected
Connection made, no introduction message received.

- introduced
Connection made, introduction message received.
`, defaultPeersFile, defaultTimeout)
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n\nUsage of %s:\n", help, os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	peersFile := flag.String("f", defaultPeersFile, "file containing peers")
	timeoutStr := flag.String("timeout", defaultTimeout, "timeout for each peer")

	flag.Parse()

	timeout, err := time.ParseDuration(*timeoutStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Bad timeout:", *timeoutStr)
		os.Exit(1)
	}

	logger.Infof("Peer connection threshold is %v", timeout)
	peers, err := getPeersListFromFile(*peersFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	report := getPeersReport(peers, timeout)
	logger.Infof("Report:\n%v", buildReport(report))
}

// getPeersListFromFile parses a local `filePath` file
// The peers list format is newline separated list of ip:port strings
// Empty lines and lines that begin with # are treated as comment lines
// Otherwise, the line is parsed as an ip:port
// If the line fails to parse, an error is returned
// Localhost addresses are allowed if allowLocalhost is true
func getPeersListFromFile(filePath string) ([]string, error) {
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var peers []string
	for _, addr := range strings.Split(string(body), "\n") {
		addr = whitespaceFilter.ReplaceAllString(addr, "")
		if addr == "" {
			continue
		}

		if strings.HasPrefix(addr, "#") {
			continue
		}

		a, err := validateAddress(addr, true)
		if err != nil {
			err = fmt.Errorf("peers list has invalid address %s: %v", addr, err)
			logger.WithError(err).Error()
			return nil, err
		}

		peers = append(peers, a)
	}

	return peers, nil
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

// validateAddress returns a sanitized address if valid, otherwise an error
func validateAddress(ipPort string, allowLocalhost bool) (string, error) {
	ipPort = whitespaceFilter.ReplaceAllString(ipPort, "")
	pts := strings.Split(ipPort, ":")
	if len(pts) != 2 {
		return "", ErrInvalidAddress
	}

	ip := net.ParseIP(pts[0])
	if ip == nil {
		return "", ErrInvalidAddress
	} else if ip.IsLoopback() {
		if !allowLocalhost {
			return "", ErrNoLocalhost
		}
	} else if !ip.IsGlobalUnicast() {
		return "", ErrNotExternalIP
	}

	port, err := strconv.ParseUint(pts[1], 10, 16)
	if err != nil {
		return "", ErrInvalidAddress
	}

	if port < 1024 {
		return "", ErrPortTooLow
	}

	return ipPort, nil
}
