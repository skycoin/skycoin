package btc

import (
	"errors"
	"fmt"
	"sync"

	"github.com/skycoin/skycoin/src/electrum"
)

// FailoverElectrum is a Backend over a list of Electrum servers. It connects on
// first use, so a server that is unreachable at startup is tried again later,
// and it keeps one server until a call through it fails. Then it moves to the
// next and retries the call once. Staying on one server limits how many
// operators learn the wallet's addresses.
type FailoverElectrum struct {
	servers []string
	dial    electrum.DialFunc

	mu   sync.Mutex
	next int
	cur  *ElectrumBackend
}

// NewFailoverElectrum returns a backend for servers, which must not be empty.
// dial may be nil for the clearnet dialer.
func NewFailoverElectrum(servers []string, dial electrum.DialFunc) (*FailoverElectrum, error) {
	if len(servers) == 0 {
		return nil, errors.New("no electrum servers")
	}
	return &FailoverElectrum{servers: servers, dial: dial}, nil
}

// backend returns the connected server, connecting to the first one that
// answers, starting where the last one failed.
func (f *FailoverElectrum) backend() (*ElectrumBackend, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cur != nil {
		return f.cur, nil
	}
	var errs []error
	for range f.servers {
		s := f.servers[f.next]
		b, err := NewElectrumBackendWithDialer(s, f.dial)
		if err == nil {
			f.cur = b
			return b, nil
		}
		errs = append(errs, err)
		f.next = (f.next + 1) % len(f.servers)
	}
	return nil, fmt.Errorf("no electrum server answered: %w", errors.Join(errs...))
}

// drop forgets b after a failed call, so the next call uses the next server.
func (f *FailoverElectrum) drop(b *ElectrumBackend) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cur != b {
		return
	}
	b.Close() //nolint:errcheck,gosec
	f.cur = nil
	f.next = (f.next + 1) % len(f.servers)
}

func do[T any](f *FailoverElectrum, call func(*ElectrumBackend) (T, error)) (T, error) {
	var zero T
	for attempt := 0; attempt < 2; attempt++ {
		b, err := f.backend()
		if err != nil {
			return zero, err
		}
		v, err := call(b)
		if err == nil {
			return v, nil
		}
		f.drop(b)
		if attempt == 1 {
			return zero, err
		}
	}
	return zero, nil
}

// GetBalance implements Backend.
func (f *FailoverElectrum) GetBalance(addresses []string) (map[string]AddressBalance, error) {
	return do(f, func(b *ElectrumBackend) (map[string]AddressBalance, error) { return b.GetBalance(addresses) })
}

// ListUnspent implements Backend.
func (f *FailoverElectrum) ListUnspent(addresses []string) ([]UTXO, error) {
	return do(f, func(b *ElectrumBackend) ([]UTXO, error) { return b.ListUnspent(addresses) })
}

// GetHistory implements Backend.
func (f *FailoverElectrum) GetHistory(addresses []string) ([]Transaction, error) {
	return do(f, func(b *ElectrumBackend) ([]Transaction, error) { return b.GetHistory(addresses) })
}

// BroadcastTransaction implements Backend. A retry on the next server sends
// the same transaction, which is harmless.
func (f *FailoverElectrum) BroadcastTransaction(rawTx string) (string, error) {
	return do(f, func(b *ElectrumBackend) (string, error) { return b.BroadcastTransaction(rawTx) })
}

// EstimateFee implements Backend.
func (f *FailoverElectrum) EstimateFee(confirmBlocks int) (int64, error) {
	return do(f, func(b *ElectrumBackend) (int64, error) { return b.EstimateFee(confirmBlocks) })
}

// Health implements Backend.
func (f *FailoverElectrum) Health() (int, error) {
	return do(f, func(b *ElectrumBackend) (int, error) { return b.Health() })
}

// Close implements Backend.
func (f *FailoverElectrum) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cur == nil {
		return nil
	}
	err := f.cur.Close()
	f.cur = nil
	return err
}
