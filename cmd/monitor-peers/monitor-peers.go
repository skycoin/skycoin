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
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/cmd/monitor-peers/connection"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/util/logging"
)

// PeerState is a current state of the peer
type PeerState string

const (
	// StateUnreachable is set when a peer couldn't be reached
	StateUnreachable = "unreachable"
	// StateReachable is set when a connection to the peer was successful
	StateReachable = "reachable"
	// StateSentIntroduction is set when an introduction message was received from the peer
	// and successfully parsed
	StateSentIntroduction = "introduced"
)

// Report contains remote `peers.txt` report data.
type Report []ReportEntry

// ReportEntry contains report data of a peer.
type ReportEntry struct {
	Address            string
	State              PeerState
	Introduction       *daemon.IntroductionMessage
	IntroValidationErr error
}

func (re ReportEntry) String() string {
	uaCoin := "-"
	uaVersion := "-"
	uaRemark := "-"
	verifyTxBurnFactor := "-"
	verifyTxMaxTxSize := "-"
	verifyTxMaxDropletPrecision := "-"
	introValidationErr := "-"

	if re.Introduction != nil {
		if re.Introduction.UserAgent.Coin != "" {
			uaCoin = re.Introduction.UserAgent.Coin
		}
		if re.Introduction.UserAgent.Version != "" {
			uaVersion = re.Introduction.UserAgent.Version
		}
		if re.Introduction.UserAgent.Remark != "" {
			uaRemark = re.Introduction.UserAgent.Remark
		}

		verifyTxBurnFactor = strconv.FormatUint(uint64(re.Introduction.UnconfirmedVerifyTxn.BurnFactor), 10)
		verifyTxMaxTxSize = strconv.FormatUint(uint64(re.Introduction.UnconfirmedVerifyTxn.MaxTransactionSize), 10)
		verifyTxMaxDropletPrecision = strconv.
			FormatUint(uint64(re.Introduction.UnconfirmedVerifyTxn.MaxDropletPrecision), 10)
	}

	if re.IntroValidationErr != nil {
		introValidationErr = re.IntroValidationErr.Error()
	}

	return fmt.Sprintf(reportFormat, re.Address, re.State, uaCoin, uaVersion, uaRemark,
		verifyTxBurnFactor, verifyTxMaxTxSize, verifyTxMaxDropletPrecision, introValidationErr)
}

// Append constructs the new report entry and appends it to the report
func (r Report) Append(addr string, state PeerState, introduction *daemon.IntroductionMessage,
	introValidationErr error) Report {
	entry := ReportEntry{
		Address:            addr,
		State:              state,
		IntroValidationErr: introValidationErr,
		Introduction:       introduction,
	}

	return append(r, entry)
}

const (
	blockchainPubKey                 = "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a"
	defaultConnectTimeout            = "1s"
	defaultReadTimeout               = "1s"
	defaultPeersFile                 = "peers.txt"
	addrWidth                        = "25"
	stateWidth                       = "15"
	uaCoinWidth                      = "10"
	uaVersionWidth                   = "10"
	uaRemarkWidth                    = "10"
	verifyTxBurnFactorWidth          = "10"
	verifyTxMaxTxSizeWidth           = "10"
	verifyTxMaxDropletPrecisionWidth = "20"
	reportFormat                     = "%-" + addrWidth + "s\t%-" + stateWidth + "s\t%-" + uaCoinWidth + "s\t%-" +
		uaVersionWidth + "s\t%-" + uaRemarkWidth + "s\t%-" + verifyTxBurnFactorWidth + "s\t%-" +
		verifyTxMaxTxSizeWidth + "s\t%-" + verifyTxMaxDropletPrecisionWidth + "s\t%v\n"
)

var (
	logger = logging.MustGetLogger("main")
	// For removing inadvertent whitespace from addresses
	whitespaceFilter = regexp.MustCompile(`\s`)
	help             = fmt.Sprintf(`monitor-peers checks the status of peers.

By default it gets peers list from %s. May be overridden with -f flag.

The default connect timeout is %s. May be overridden with -ctimeout flag. The timeout is parsed by time.ParseDuration.

The default read timeout is %s. May be overridden with -rtimeout flag. The timeout is parsed by time.ParseDuration.

It generates a report for each peer which contains the peer address and status. Status may be one of the following:

- unreachable
No connection made.

- reachable
Connection made, no introduction message received.

- introduced
Connection made, introduction message received.
`, defaultPeersFile, defaultConnectTimeout, defaultReadTimeout)
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n\nUsage of %s:\n", help, os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	peersFile := flag.String("f", defaultPeersFile, "file containing peers")
	connectTimeoutStr := flag.String("ctimeout", defaultConnectTimeout, "connect timeout for each peer")
	readTimeoutStr := flag.String("rtimeout", defaultReadTimeout, "read timeout for each peer")

	flag.Parse()

	connectTimeout, err := time.ParseDuration(*connectTimeoutStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Bad connect timeout: ", *connectTimeoutStr)
		os.Exit(1)
	}

	logger.Infof("Peer connection threshold is %v", connectTimeout)

	readTimeout, err := time.ParseDuration(*readTimeoutStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Bad read timeout: ", *readTimeoutStr)
		os.Exit(1)
	}

	logger.Infof("Peer read threshold is %v", readTimeout)

	peers, err := getPeersListFromFile(*peersFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	report := getPeersReport(peers, connectTimeout, readTimeout)
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

// getPeersReport loops through `peers`, connects to each and tries to read the introduction
// message. Builds and returns the report
func getPeersReport(peers []string, connectTimeout, readTimeout time.Duration) Report {
	dc := daemon.NewDaemonConfig()
	dc.BlockchainPubkey = cipher.MustPubKeyFromHex(blockchainPubKey)

	report := make(Report, 0, len(peers))

	for _, addr := range peers {
		conn, err := connection.NewConnection(addr, connectTimeout, readTimeout)
		if err != nil {
			logger.WithError(err).Error()
			continue
		}

		if err := conn.Connect(); err != nil {
			report = report.Append(addr, StateUnreachable, nil, nil)
			continue
		}

		introduction, err := conn.TryReadIntroductionMessage()
		if err != nil {
			report = report.Append(addr, StateReachable, nil, nil)
			continue
		}

		if err := introduction.Verify(dc, logrus.Fields{
			"addr": addr,
		}); err != nil {
			report = report.Append(addr, StateSentIntroduction, introduction, err)
			continue
		}

		report = report.Append(addr, StateSentIntroduction, introduction, nil)

		if err := conn.Disconnect(); err != nil {
			logger.WithError(err).Error()
		}
	}

	return report
}

// buildReport builds a report to a string output
func buildReport(report Report) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(reportFormat, "Address", "Status", "Coin", "Version", "Remark",
		"Burn factor", "Max tx size", "Max droplet precision", "Intro validation error"))
	for _, entry := range report {
		sb.WriteString(entry.String())
	}

	return sb.String()
}

// validateAddress returns a sanitized address if valid, otherwise an error
func validateAddress(ipPort string, allowLocalhost bool) (string, error) {
	ipPort = whitespaceFilter.ReplaceAllString(ipPort, "")
	pts := strings.Split(ipPort, ":")
	if len(pts) != 2 {
		return "", pex.ErrInvalidAddress
	}

	ip := net.ParseIP(pts[0])
	if ip == nil {
		return "", pex.ErrInvalidAddress
	} else if ip.IsLoopback() {
		if !allowLocalhost {
			return "", pex.ErrNoLocalhost
		}
	} else if !ip.IsGlobalUnicast() {
		return "", pex.ErrNotExternalIP
	}

	port, err := strconv.ParseUint(pts[1], 10, 16)
	if err != nil {
		return "", pex.ErrInvalidAddress
	}

	if port < 1024 {
		return "", pex.ErrPortTooLow
	}

	return ipPort, nil
}
