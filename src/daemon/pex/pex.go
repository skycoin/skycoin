// Package pex is a toolkit for implementing a peer exchange system
package pex

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"

	"github.com/skycoin/skycoin/src/util/logging"
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
	// PeerCacheFilename filename for disk-cached peers
	PeerCacheFilename = "peers.json"
	// oldPeerCacheFilename previous filename for disk-cached peers. The cache loader will fall back onto this filename if it can't load peers.json
	oldPeerCacheFilename = "peers.txt"
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
	HasIncomingPort bool   // Whether this peer has accessible public port
	RetryTimes      int    `json:"-"` // records the retry times
}

// NewPeer returns a *Peer initialized by an address string of the form ip:port
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
	peer.LastSeen = time.Now().UTC().Unix()
}

// IncreaseRetryTimes adds the retry times
func (peer *Peer) IncreaseRetryTimes() {
	peer.RetryTimes++
	logger.Debugf("Increase retry times of %v: %v", peer.Addr, peer.RetryTimes)
}

// ResetRetryTimes resets the retry time
func (peer *Peer) ResetRetryTimes() {
	peer.RetryTimes = 0
	logger.Debugf("Reset retry times of %v", peer.Addr)
}

// CanTry returns whether this peer is tryable base on the exponential backoff algorithm
func (peer *Peer) CanTry() bool {
	// Exponential backoff
	mod := (math.Exp2(float64(peer.RetryTimes)) - 1) * 5
	if mod == 0 {
		return true
	}

	// Random time elapsed
	now := time.Now().UTC().Unix()
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
	// Set all peers as untrusted (even if loaded from DefaultConnections)
	DisableTrustedPeers bool
	// Load peers from this file on disk. NOTE: this is different from the peers file cache in the data directory
	CustomPeersFile string
	// Default "trusted" connections
	DefaultConnections []string
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
		DisableTrustedPeers: false,
		CustomPeersFile:     "",
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
func New(cfg Config) (*Pex, error) {
	pex := &Pex{
		Config:   cfg,
		peerlist: newPeerlist(),
		quit:     make(chan struct{}),
		done:     make(chan struct{}),
	}

	// Load peers from disk
	if err := pex.loadCache(); err != nil {
		logger.Critical().WithError(err).Error("pex.loadCache failed")
		return nil, err
	}

	// Load default hardcoded peers
	for _, addr := range cfg.DefaultConnections {
		// Default peers will mark as trusted peers.
		if err := pex.AddPeer(addr); err != nil {
			logger.Critical().WithError(err).Error("Add default peer failed")
			return nil, err
		}
		if err := pex.SetTrusted(addr); err != nil {
			logger.Critical().WithError(err).Error("pex.SetTrusted for default peer failed")
			return nil, err
		}
	}

	if cfg.DisableTrustedPeers {
		// Unset trusted status from any existing peers
		pex.setAllUntrusted()
	}

	// Add custom peers
	if cfg.CustomPeersFile != "" {
		if err := pex.loadCustom(cfg.CustomPeersFile); err != nil {
			logger.Critical().WithError(err).Errorf("Failed to load custom peers from %s", cfg.CustomPeersFile)
			return nil, err
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
				logger.Errorf("Failed to download peers list: %v", err)
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
			logger.Errorf("Save peers failed: %v", err)
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
		logger.Errorf("Failed to download peers from %s. err: %s", px.Config.PeerListURL, err.Error())
		return err
	}

	peers := parseRemotePeerList(body)
	logger.Infof("Downloaded peers list from %s, got %d peers", px.Config.PeerListURL, len(peers))

	n := px.AddPeers(peers)
	logger.Infof("Added %d/%d peers from downloaded peers list", n, len(peers))

	return nil
}

func (px *Pex) loadCache() error {
	px.Lock()
	defer px.Unlock()

	fp := filepath.Join(px.Config.DataDirectory, PeerCacheFilename)
	peers, err := loadCachedPeersFile(fp)

	if err != nil {
		return err
	}

	// If the PeerCacheFilename peers.json file does not exist, try to load the old peers.txt file
	if peers == nil {
		logger.Infof("Peer cache %s not found, falling back on %s", PeerCacheFilename, oldPeerCacheFilename)

		fp := filepath.Join(px.Config.DataDirectory, oldPeerCacheFilename)
		peers, err = loadCachedPeersFile(fp)
		if err != nil {
			return err
		}

		if peers == nil {
			logger.Infof("Fallback peer cache %s not found", oldPeerCacheFilename)
			return nil
		}
	}

	// remove invalid peers and limit the max number of peers to pex.Config.Max
	var validPeers []Peer
	for addr, p := range peers {
		if _, err := validateAddress(addr, px.Config.AllowLocalhost); err != nil {
			logger.Errorf("Invalid peer address: %v", err)
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

func (px *Pex) loadCustom(fn string) error {
	px.Lock()
	defer px.Unlock()

	f, err := os.Open(fn)
	if err != nil {
		return err
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	peers, err := parseLocalPeerList(string(data), px.Config.AllowLocalhost)
	if err != nil {
		return err
	}

	logger.Infof("Loaded %d peers from %s", len(peers), fn)

	px.peerlist.addPeers(peers)
	return nil
}

// SavePeers persists the peerlist
func (px *Pex) save() error {
	px.Lock()
	defer px.Unlock()

	fn := filepath.Join(px.Config.DataDirectory, PeerCacheFilename)
	return px.peerlist.save(fn)
}

// AddPeer adds a peer to the peer list, given an address. If the peer list is
// full, PeerlistFullError is returned
func (px *Pex) AddPeer(addr string) error {
	px.Lock()
	defer px.Unlock()

	cleanAddr, err := validateAddress(addr, px.Config.AllowLocalhost)
	if err != nil {
		logger.Errorf("Invalid address %s: %v", addr, err)
		return ErrInvalidAddress
	}

	if !px.peerlist.hasPeer(cleanAddr) {
		if px.Config.Max > 0 && px.peerlist.len() >= px.Config.Max {
			return ErrPeerlistFull
		}

		px.peerlist.addPeer(cleanAddr)
	}

	return nil
}

// AddPeers add multiple peers at once. Any errors will be logged, but not returned
// Returns the number of peers that were added without error. Note that
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
			logger.Infof("Add peers sees an invalid address %s: %v", addr, err)
			continue
		}
		validAddrs = append(validAddrs, a)
	}
	addrs = validAddrs

	// Shuffle the addresses before capping them
	rand.Shuffle(len(addrs), func(i, j int) {
		addrs[i], addrs[j] = addrs[j], addrs[i]
	})

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
		logger.Errorf("Invalid address %s: %v", addr, err)
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
		logger.Errorf("Invalid address %s: %v", addr, err)
		return ErrInvalidAddress
	}

	return px.peerlist.setTrusted(cleanAddr, true)
}

// setAllUntrusted unsets the trusted field on all peers
func (px *Pex) setAllUntrusted() {
	px.Lock()
	defer px.Unlock()

	px.peerlist.setAllUntrusted()
}

// SetHasIncomingPort sets if the peer has public port
func (px *Pex) SetHasIncomingPort(addr string, hasPublicPort bool) error {
	px.Lock()
	defer px.Unlock()

	cleanAddr, err := validateAddress(addr, px.Config.AllowLocalhost)
	if err != nil {
		logger.Errorf("Invalid address %s: %v", addr, err)
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
	return px.peerlist.getPeers([]Filter{isTrusted})
}

// Private returns private peers
func (px *Pex) Private() Peers {
	px.RLock()
	defer px.RUnlock()
	return px.peerlist.getCanTryPeers([]Filter{isPrivate})
}

// TrustedPublic returns trusted public peers
func (px *Pex) TrustedPublic() Peers {
	px.RLock()
	defer px.RUnlock()
	return px.peerlist.getCanTryPeers([]Filter{isPublic, isTrusted})
}

// RandomPublic returns N random public peers
func (px *Pex) RandomPublic(n int) Peers {
	px.RLock()
	defer px.RUnlock()
	return px.peerlist.random(n, []Filter{isPublic})
}

// RandomExchangeable returns N random exchangeable peers
func (px *Pex) RandomExchangeable(n int) Peers {
	px.RLock()
	defer px.RUnlock()
	return px.peerlist.random(n, isExchangeable)
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
		logger.Errorf("waiting %v to retry downloadText, error: %v", wait, err)
	}

	operation := func() error {
		logger.Infof("Trying to download peers list from %s", url)
		var err error
		body, err = downloadText(url)
		return err
	}

	if err := backoff.RetryNotify(operation, b, notify); err != nil {
		logger.Infof("Gave up dowloading peers list from %s: %v", url, err)
		return "", err
	}

	logger.Infof("Peers list downloaded from %s", url)

	return body, nil
}

// parseRemotePeerList parses a remote peers.txt file
// The peers list format is newline separated list of ip:port strings
// Any lines that don't parse to an ip:port are skipped, otherwise they return an error
// Localhost ip:port addresses are ignored
// NOTE: this does not parse the cached peers.json file in the data directory, which is a JSON file
// and is loaded by loadCachedPeersFile
func parseRemotePeerList(body string) []string {
	var peers []string
	for _, addr := range strings.Split(body, "\n") {
		addr = whitespaceFilter.ReplaceAllString(addr, "")
		if addr == "" {
			continue
		}

		// Never allow localhost addresses from the remote peers list
		a, err := validateAddress(addr, false)
		if err != nil {
			err = fmt.Errorf("Peers list has invalid address %s: %v", addr, err)
			logger.WithError(err).Error()
			continue
		}

		peers = append(peers, a)
	}

	return peers
}

// parseLocalPeerList parses a local peers.txt file
// The peers list format is newline separated list of ip:port strings
// Empty lines and lines that begin with # are treated as comment lines
// Otherwise, the line is parsed as an ip:port
// If the line fails to parse, an error is returned
// Localhost addresses are allowed if allowLocalhost is true
// NOTE: this does not parse the cached peers.json file in the data directory, which is a JSON file
// and is loaded by loadCachedPeersFile
func parseLocalPeerList(body string, allowLocalhost bool) ([]string, error) {
	var peers []string
	for _, addr := range strings.Split(body, "\n") {
		addr = whitespaceFilter.ReplaceAllString(addr, "")
		if addr == "" {
			continue
		}

		if strings.HasPrefix(addr, "#") {
			continue
		}

		a, err := validateAddress(addr, allowLocalhost)
		if err != nil {
			err = fmt.Errorf("Peers list has invalid address %s: %v", addr, err)
			logger.WithError(err).Error()
			return nil, err
		}

		peers = append(peers, a)
	}

	return peers, nil
}
