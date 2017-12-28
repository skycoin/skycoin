// Package pex is a toolkit for implementing a peer exchange system
package pex

import (
	"errors"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"

	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/utc"
)

//TODO:
// - keep track of last time the peer was connected to
// - last time peer was connected to is more important than "seen"
// - peer "seen" means something else than use here
// - save last time connected to, use 0 for never
// - only transmit peers that have active or recent connections

const (
	// DefaultPeerListURL is the default URL to download remote peers list from, if enabled
	DefaultPeerListURL = "https://downloads.skycoin.net/blockchain/peers.txt"
	// PeerDatabaseFilename filename for disk-cached peers
	PeerDatabaseFilename = "peers.txt"
	// MaxPeerRetryTimes is the maximum number of times to retry a peer
	MaxPeerRetryTimes = 10
)

var (
	// ErrPeerlistFull is returned when the Pex is at a maximum
	ErrPeerlistFull = errors.New("Peer list full")
	// ErrInvalidAddress is returned when an address appears malformed
	ErrInvalidAddress = errors.New("Invalid address")
	// ErrNoLocalhost is returned if a localhost addresses are not allowed
	ErrNoLocalhost = errors.New("Localhost address is not allowed")
	// ErrNotExternalIP is returned if an IP address is not a global unicast address
	ErrNotExternalIP = errors.New("IP is not a valid external IP")
	// ErrPortTooLow is returned if a port is less than 1024
	ErrPortTooLow = errors.New("Port must be >= 1024")
	// ErrBlacklistedAddress returned when attempting to add a blacklisted peer
	ErrBlacklistedAddress = errors.New("Blacklisted address")

	// Logging. See http://godoc.org/github.com/op/go-logging for
	// instructions on how to include this log's output
	logger = logging.MustGetLogger("pex")
	// Default rng
	rnum = rand.New(rand.NewSource(time.Now().Unix()))
	// For removing inadvertent whitespace from addresses
	whitespaceFilter = regexp.MustCompile(`\s`)
)

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

// Peer represents a known peer
type Peer struct {
	Addr            string // An address of the form ip:port
	LastSeen        int64  // Unix timestamp when this peer was last seen
	Private         bool   // Whether it should omitted from public requests
	Trusted         bool   // Whether this peer is trusted
	HasIncomingPort bool   // Whether this peer has accessable public port
	RetryTimes      int    `json:"-"` // records the retry times
}

// NewPeer returns a *Peer initialised by an address string of the form ip:port
func NewPeer(address string) *Peer {
	p := &Peer{
		Addr:    address,
		Private: false,
		Trusted: false,
	}
	p.Seen()
	return p
}

// Seen marks the peer as seen
func (peer *Peer) Seen() {
	peer.LastSeen = utc.UnixNow()
}

// IncreaseRetryTimes adds the retry times
func (peer *Peer) IncreaseRetryTimes() {
	peer.RetryTimes++
	logger.Debug("Increase retry times of %v: %v", peer.Addr, peer.RetryTimes)
}

// ResetRetryTimes resets the retry time
func (peer *Peer) ResetRetryTimes() {
	peer.RetryTimes = 0
	logger.Debug("Reset retry times of %v", peer.Addr)
}

// CanTry returns whether this peer is tryable base on the exponential backoff algorithm
func (peer *Peer) CanTry() bool {
	// Exponential backoff
	mod := (math.Exp2(float64(peer.RetryTimes)) - 1) * 5
	if mod == 0 {
		return true
	}

	// Random time elapsed
	now := utc.UnixNow()
	t := rnum.Int63n(int64(mod))
	return now-peer.LastSeen > t
}

// String returns the peer address
func (peer *Peer) String() string {
	return peer.Addr
}

// Config pex config
type Config struct {
	// Folder where peers database should be saved
	DataDirectory string
	// Maximum number of peers to keep account of in the PeerList
	Max int
	// Cull peers after they havent been seen in this much time
	Expiration time.Duration
	// Cull expired peers on this interval
	CullRate time.Duration
	// clear old peers on this interval
	ClearOldRate time.Duration
	// How often to clear expired blacklist entries
	UpdateBlacklistRate time.Duration
	// How often to request peers via PEX
	RequestRate time.Duration
	// How many peers to send back in response to a peers request
	ReplyCount int
	// Localhost peers are allowed in the peerlist
	AllowLocalhost bool
	// Disable exchanging of peers.  Peers are still loaded from disk
	Disabled bool
	// Whether the network is disabled
	NetworkDisabled bool
	// Download peers list from remote host
	DownloadPeerList bool
	// Download peers list from this URL
	PeerListURL string
}

// NewConfig creates default pex config.
func NewConfig() Config {
	return Config{
		DataDirectory:       "./",
		Max:                 1000,
		Expiration:          time.Hour * 24 * 7,
		CullRate:            time.Minute * 10,
		ClearOldRate:        time.Minute * 10,
		UpdateBlacklistRate: time.Minute,
		RequestRate:         time.Minute,
		ReplyCount:          30,
		AllowLocalhost:      false,
		Disabled:            false,
		NetworkDisabled:     false,
		DownloadPeerList:    false,
		PeerListURL:         DefaultPeerListURL,
	}
}

// Pex manages a set of known peers and controls peer acquisition
type Pex struct {
	sync.RWMutex
	// All known peers
	peerlist peerlist
	Config   Config
	quit     chan struct{}
	done     chan struct{}
}

// New creates pex
func New(cfg Config, defaultConns []string) (*Pex, error) {
	pex := &Pex{
		Config:   cfg,
		peerlist: newPeerlist(),
		quit:     make(chan struct{}),
		done:     make(chan struct{}),
	}

	// Load peers from disk
	if err := pex.load(); err != nil {
		return nil, err
	}

	// Load default hardcoded peers
	for _, addr := range defaultConns {
		// Default peers will mark as trusted peers.
		if err := pex.AddPeer(addr); err != nil {
			logger.Critical("add peer failed:%v", err)
			continue
		}
		if err := pex.SetTrusted(addr); err != nil {
			logger.Critical("pex.SetTrust failed: %v", err)
		}
	}

	// Save peers to disk
	if err := pex.save(); err != nil {
		return nil, err
	}

	// Download peers from remote peers list
	if pex.Config.DownloadPeerList {
		go func() {
			if err := pex.downloadPeers(); err != nil {
				logger.Error("Failed to download peers list: %v", err)
			}
		}()
	}

	return pex, nil
}

// Run starts the pex service
func (px *Pex) Run() error {
	logger.Info("Pex.Run started")
	defer logger.Info("Pex.Run stopped")
	defer close(px.done)

	defer func() {
		// Save the peerlist
		logger.Info("Save peerlist")
		if err := px.save(); err != nil {
			logger.Error("Save peers failed: %v", err)
		}
	}()

	clearOldTicker := time.NewTicker(px.Config.ClearOldRate)

	for {
		select {
		case <-clearOldTicker.C:
			// Remove peers we haven't seen in a while
			if !px.Config.Disabled && !px.Config.NetworkDisabled {
				px.Lock()
				px.peerlist.clearOld(px.Config.Expiration)
				px.Unlock()
			}
		case <-px.quit:
			return nil
		}
	}
}

// Shutdown notifies the pex service to exist
func (px *Pex) Shutdown() {
	logger.Info("Shutting down pex")
	defer logger.Info("Pex shutdown")
	close(px.quit)
	<-px.done
}

func (px *Pex) downloadPeers() error {
	body, err := backoffDownloadText(px.Config.PeerListURL)
	if err != nil {
		logger.Error("Failed to download peers from %s. err: %s", px.Config.PeerListURL, err.Error())
		return err
	}

	peers := parseRemotePeerList(body)
	logger.Info("Downloaded peers list from %s, got %d peers", px.Config.PeerListURL, len(peers))

	n := px.AddPeers(peers)
	logger.Info("Added %d/%d peers from downloaded peers list", n, len(peers))

	return nil
}

func (px *Pex) load() error {
	px.Lock()
	defer px.Unlock()

	fp := filepath.Join(px.Config.DataDirectory, PeerDatabaseFilename)
	peers, err := loadPeersFromFile(fp)
	if err != nil {
		return err
	}

	// file does not exist
	if peers == nil {
		return nil
	}

	// remove invalid peers and limit the max number of peers to pex.Config.Max
	var validPeers []Peer
	for addr, p := range peers {
		if _, err := validateAddress(addr, px.Config.AllowLocalhost); err != nil {
			logger.Error("Invalid peer address: %v", err)
			continue
		}

		validPeers = append(validPeers, *p)
		if px.Config.Max > 0 && len(validPeers) >= px.Config.Max {
			break
		}
	}

	px.peerlist.setPeers(validPeers)
	return nil
}

// SavePeers persists the peerlist
func (px *Pex) save() error {
	px.Lock()
	defer px.Unlock()

	fn := filepath.Join(px.Config.DataDirectory, PeerDatabaseFilename)
	return px.peerlist.save(fn)
}

// AddPeer adds a peer to the peer list, given an address. If the peer list is
// full, PeerlistFullError is returned */
func (px *Pex) AddPeer(addr string) error {
	px.Lock()
	defer px.Unlock()

	cleanAddr, err := validateAddress(addr, px.Config.AllowLocalhost)
	if err != nil {
		logger.Error("Invalid address %s: %v", addr, err)
		return ErrInvalidAddress
	}

	if px.Config.Max > 0 && px.peerlist.len() >= px.Config.Max {
		return ErrPeerlistFull
	}

	px.peerlist.addPeer(cleanAddr)
	return nil
}

// AddPeers add multiple peers at once. Any errors will be logged, but not returned
// Returns the number of peers that were added without error.  Note that
// adding a duplicate peer will not cause an error.
func (px *Pex) AddPeers(addrs []string) int {
	px.Lock()
	defer px.Unlock()

	if px.Config.Max > 0 && px.peerlist.len() >= px.Config.Max {
		logger.Warning("Add peers failed, peer list is full")
		return 0
	}

	// validate the addresses
	var validAddrs []string
	for _, addr := range addrs {
		a, err := validateAddress(addr, px.Config.AllowLocalhost)
		if err != nil {
			logger.Info("Add peers sees an invalid address %s: %v", addr, err)
			continue
		}
		validAddrs = append(validAddrs, a)
	}
	addrs = validAddrs

	// Shuffle the addresses before capping them
	for i := len(addrs) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		addrs[i], addrs[j] = addrs[j], addrs[i]
	}

	if px.Config.Max > 0 {
		rcap := px.Config.Max - px.peerlist.len()
		if len(addrs) > rcap {
			addrs = addrs[:rcap]
		}
	}

	px.peerlist.addPeers(addrs)
	return len(addrs)
}

// SetPrivate updates peer's private value
func (px *Pex) SetPrivate(addr string, private bool) error {
	px.Lock()
	defer px.Unlock()

	cleanAddr, err := validateAddress(addr, px.Config.AllowLocalhost)
	if err != nil {
		logger.Error("Invalid address %s: %v", addr, err)
		return ErrInvalidAddress
	}

	return px.peerlist.setPrivate(cleanAddr, private)
}

// SetTrusted updates peer's trusted value
func (px *Pex) SetTrusted(addr string) error {
	px.Lock()
	defer px.Unlock()

	cleanAddr, err := validateAddress(addr, px.Config.AllowLocalhost)
	if err != nil {
		logger.Error("Invalid address %s: %v", addr, err)
		return ErrInvalidAddress
	}

	return px.peerlist.setTrusted(cleanAddr, true)
}

// SetHasIncomingPort sets if the peer has public port
func (px *Pex) SetHasIncomingPort(addr string, hasPublicPort bool) error {
	px.Lock()
	defer px.Unlock()

	cleanAddr, err := validateAddress(addr, px.Config.AllowLocalhost)
	if err != nil {
		logger.Error("Invalid address %s: %v", addr, err)
		return ErrInvalidAddress
	}

	return px.peerlist.setHasIncomingPort(cleanAddr, hasPublicPort)
}

// RemovePeer removes peer
func (px *Pex) RemovePeer(addr string) {
	px.Lock()
	defer px.Unlock()
	px.peerlist.removePeer(addr)
}

// GetPeerByAddr returns peer of given address
func (px *Pex) GetPeerByAddr(addr string) (Peer, bool) {
	px.RLock()
	defer px.RUnlock()
	return px.peerlist.getPeerByAddr(addr)
}

// Trusted returns trusted peers
func (px *Pex) Trusted() Peers {
	px.RLock()
	defer px.RUnlock()
	return px.peerlist.getPeers(isTrusted)
}

// Private returns private peers
func (px *Pex) Private() Peers {
	px.RLock()
	defer px.RUnlock()
	return px.peerlist.getPeers(isPrivate)
}

// TrustedPublic returns trusted public peers
func (px *Pex) TrustedPublic() Peers {
	px.RLock()
	defer px.RUnlock()
	return px.peerlist.getPeers(isPublic, isTrusted)
}

// RandomPublic returns N random public peers
func (px *Pex) RandomPublic(n int) Peers {
	px.RLock()
	defer px.RUnlock()
	return px.peerlist.random(n, isPublic)
}

// RandomExchangeable returns N random exchangeable peers
func (px *Pex) RandomExchangeable(n int) Peers {
	px.RLock()
	defer px.RUnlock()
	return px.peerlist.random(n, isExchangeable...)
}

// IncreaseRetryTimes increases retry times
func (px *Pex) IncreaseRetryTimes(addr string) {
	px.Lock()
	defer px.Unlock()
	px.peerlist.increaseRetryTimes(addr)
}

// ResetRetryTimes reset retry times
func (px *Pex) ResetRetryTimes(addr string) {
	px.Lock()
	defer px.Unlock()
	px.peerlist.resetRetryTimes(addr)
}

// ResetAllRetryTimes reset all peers' retry times
func (px *Pex) ResetAllRetryTimes() {
	px.Lock()
	defer px.Unlock()
	px.peerlist.resetAllRetryTimes()
}

// IsFull returns whether the peer list is full
func (px *Pex) IsFull() bool {
	px.RLock()
	defer px.RUnlock()
	return px.Config.Max > 0 && px.peerlist.len() >= px.Config.Max
}

// downloadText downloads a text format file from url.
// Returns the raw response body as a string.
// TODO -- move to util, add backoff options
func downloadText(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func backoffDownloadText(url string) (string, error) {
	var body string

	b := backoff.NewExponentialBackOff()

	notify := func(err error, wait time.Duration) {
		logger.Error("waiting %v to retry downloadText, error: %v", wait, err)
	}

	operation := func() error {
		logger.Info("Trying to download peers list from %s", url)
		var err error
		body, err = downloadText(url)
		return err
	}

	if err := backoff.RetryNotify(operation, b, notify); err != nil {
		logger.Info("Gave up dowloading peers list from %s: %v", url, err)
		return "", err
	}

	logger.Info("Peers list downloaded from %s", url)

	return body, nil
}

// parseRemotePeerList parses a remote peers.txt file
// The peers list format is newline separated ip:port
// Any lines that don't parse to an ip:port are skipped
// Localhost ip:port addresses are ignored
func parseRemotePeerList(body string) []string {
	var peers []string
	for _, addr := range strings.Split(string(body), "\n") {
		addr = whitespaceFilter.ReplaceAllString(addr, "")
		if addr == "" {
			continue
		}

		// Never allow localhost addresses from the remote peers list
		a, err := validateAddress(addr, false)
		if err != nil {
			logger.Error("Remote peers list has invalid address %s: %v", addr, err)
			continue
		}

		peers = append(peers, a)
	}

	return peers
}
