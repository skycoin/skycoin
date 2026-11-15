// Package sweep src/wallet/sweep/service.go c1-coin-wallet
package sweep

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/wallet"
)

var logger = logging.MustGetLogger("sweep")

// Chain reads the blockchain for a sweep pass.
type Chain interface {
	HeadBkSeq() (uint64, bool, error)
	SweepState(addrs []cipher.Address, txids []cipher.SHA256) (ChainState, error)
}

// Injector sends a signed transaction to the network.
type Injector interface {
	InjectBroadcastTransaction(txn coin.Transaction) error
}

// Config configures a Service.
type Config struct {
	// Dir is the wallet directory, where RulesFile is kept.
	Dir string
	// Poll is how often the head block is checked. Zero means one second.
	Poll time.Duration
	// RetryAfter is how long a failed pass waits before it runs again without a
	// new block. Zero means ten seconds.
	RetryAfter time.Duration
}

// Status is the state of the sweep of one wallet.
type Status struct {
	Wallet  string `json:"wallet"`
	Enabled bool   `json:"enabled"`
	// Paused says why coins on exposed addresses are not being swept.
	Paused   string `json:"paused,omitempty"`
	Unlocked bool   `json:"unlocked"`
	// UnlockExpires is when the node forgets the password, if it ever does.
	UnlockExpires *time.Time `json:"unlock_expires,omitempty"`
	HeadSeq       uint64     `json:"head_seq"`
	Exposed       int        `json:"exposed_addresses"`
	// WaitingCoins and WaitingHours are on exposed addresses and not in a sweep yet.
	WaitingCoins uint64 `json:"waiting_coins"`
	WaitingHours uint64 `json:"waiting_hours"`
	// Contested counts waiting outputs that another unconfirmed transaction spends.
	Contested   int        `json:"contested"`
	Destination string     `json:"destination,omitempty"`
	Pending     []string   `json:"pending_txids,omitempty"`
	LastSweep   *time.Time `json:"last_sweep,omitempty"`
	LastTxids   []string   `json:"last_txids,omitempty"`
	Confirmed   int        `json:"confirmed_sweeps"`
	// Displaced lists the latest sweeps that left the pool without confirming.
	Displaced []string `json:"displaced_txids,omitempty"`
	LastError string   `json:"last_error,omitempty"`
}

type pendingSweep struct {
	txid   cipher.SHA256
	inputs []cipher.SHA256
}

type key struct {
	password []byte
	expires  time.Time
}

func (k *key) erase() {
	for i := range k.password {
		k.password[i] = 0
	}
}

// Service sweeps the exposed addresses of the wallets that have an enabled rule.
type Service struct {
	cfg     Config
	wallets *wallet.Service
	chain   Chain
	inject  Injector
	kick    chan struct{}
	passMu  sync.Mutex
	now     func() time.Time

	mu      sync.Mutex
	rules   map[string]Rule
	keys    map[string]*key
	status  map[string]*Status
	pending map[string][]pendingSweep
	retry   bool
}

// NewService loads the rules from cfg.Dir.
func NewService(cfg Config, ws *wallet.Service, chain Chain, inject Injector) (*Service, error) {
	if cfg.Poll == 0 {
		cfg.Poll = time.Second
	}
	if cfg.RetryAfter == 0 {
		cfg.RetryAfter = 10 * time.Second
	}
	rules, err := loadRules(cfg.Dir)
	if err != nil {
		return nil, err
	}
	return &Service{
		cfg:     cfg,
		wallets: ws,
		chain:   chain,
		inject:  inject,
		kick:    make(chan struct{}, 1),
		now:     time.Now,
		rules:   rules,
		keys:    make(map[string]*key),
		status:  make(map[string]*Status),
		pending: make(map[string][]pendingSweep),
	}, nil
}

// Run sweeps on every new block until quit is closed, then forgets every password.
func (s *Service) Run(quit <-chan struct{}) {
	t := time.NewTicker(s.cfg.Poll)
	defer t.Stop()
	defer s.LockAll()
	lastSeq, _, err := s.chain.HeadBkSeq()
	if err != nil {
		logger.WithError(err).Error("HeadBkSeq failed")
	}
	lastPass := s.now()
	s.Pass()
	for {
		select {
		case <-quit:
			return
		case <-s.kick:
		case <-t.C:
			seq, _, err := s.chain.HeadBkSeq()
			if err != nil {
				logger.WithError(err).Error("HeadBkSeq failed")
				continue
			}
			s.mu.Lock()
			retry := s.retry
			s.mu.Unlock()
			if seq == lastSeq && !(retry && s.now().Sub(lastPass) >= s.cfg.RetryAfter) {
				continue
			}
			lastSeq = seq
		}
		lastPass = s.now()
		s.Pass()
	}
}

func (s *Service) poke() {
	select {
	case s.kick <- struct{}{}:
	default:
	}
}

// Pass runs one sweep of every enabled rule.
func (s *Service) Pass() {
	s.passMu.Lock()
	defer s.passMu.Unlock()

	s.mu.Lock()
	rules := make([]Rule, 0, len(s.rules))
	for _, r := range s.rules {
		if r.Enabled {
			rules = append(rules, r)
		}
	}
	s.retry = false
	s.mu.Unlock()
	sort.Slice(rules, func(i, j int) bool { return rules[i].Wallet < rules[j].Wallet })

	for _, r := range rules {
		if err := s.sweep(r); err != nil {
			logger.WithError(err).WithField("wallet", r.Wallet).Error("sweep failed")
			s.update(r.Wallet, func(st *Status) { st.LastError = err.Error() })
			s.mu.Lock()
			s.retry = true
			s.mu.Unlock()
		} else {
			s.update(r.Wallet, func(st *Status) { st.LastError = "" })
		}
	}
}

func (s *Service) update(wltID string, f func(*Status)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.status[wltID]
	if !ok {
		st = &Status{Wallet: wltID}
		s.status[wltID] = st
	}
	f(st)
}

// password returns the password held for an encrypted wallet, forgetting it once
// it has expired.
func (s *Service) password(wltID string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k, ok := s.keys[wltID]
	if !ok {
		return nil, false
	}
	if !k.expires.IsZero() && !s.now().Before(k.expires) {
		k.erase()
		delete(s.keys, wltID)
		return nil, false
	}
	return append([]byte(nil), k.password...), true
}

func erase(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// listing is a wallet's addresses and the chain state that goes with them.
type listing struct {
	w     wallet.Wallet
	addrs []walletAddr
	st    ChainState
}

// list reads w and the chain, generating addresses on shared accounts until gap
// unused ones follow the last used one.
func (s *Service) list(r Rule, txids []cipher.SHA256) (listing, error) {
	for range 8 {
		w, err := s.wallets.GetWallet(r.Wallet)
		if err != nil {
			return listing{}, err
		}
		addrs, err := walletAddresses(w)
		if err != nil {
			return listing{}, err
		}
		st, err := s.chain.SweepState(addrsOf(addrs), txids)
		if err != nil {
			return listing{}, err
		}
		needs := gapNeeds(r, addrs, st.Seen)
		if len(needs) == 0 {
			return listing{w: w, addrs: addrs, st: st}, nil
		}
		for acct, n := range needs {
			if err := s.wallets.Update(r.Wallet, func(w wallet.Wallet) error {
				_, err := w.GenerateAddresses(wallet.OptionAccount(acct), wallet.OptionExternal(), wallet.OptionGenerateN(n))
				return err
			}); err != nil {
				return listing{}, err
			}
		}
	}
	return listing{}, errors.New("shared account addresses kept growing past the gap")
}

func addrsOf(addrs []walletAddr) []cipher.Address {
	out := make([]cipher.Address, len(addrs))
	for i, a := range addrs {
		out[i] = a.addr
	}
	return out
}

func (s *Service) sweep(r Rule) error {
	s.mu.Lock()
	pending := append([]pendingSweep(nil), s.pending[r.Wallet]...)
	s.mu.Unlock()
	txids := make([]cipher.SHA256, len(pending))
	for i, p := range pending {
		txids[i] = p.txid
	}

	l, err := s.list(r, txids)
	if err != nil {
		return err
	}

	// Drop the sweeps that confirmed or left the pool. Outputs of one that left
	// the pool and are still unspent are swept again below.
	var kept []pendingSweep
	var confirmed int
	var displaced []string
	covered := make(map[cipher.SHA256]bool)
	unspent := make(map[cipher.SHA256]bool, len(l.st.Unspent))
	for _, u := range l.st.Unspent {
		unspent[u.Hash()] = true
	}
	for _, p := range pending {
		state := l.st.Txns[p.txid]
		if state == TxnPending && !allIn(p.inputs, unspent) {
			// Another transaction spent an input in a block, so this one can never confirm
			state = TxnUnknown
		}
		switch state {
		case TxnPending:
			kept = append(kept, p)
			for _, h := range p.inputs {
				covered[h] = true
			}
		case TxnConfirmed:
			confirmed++
		default:
			displaced = append(displaced, p.txid.Hex())
		}
	}
	s.mu.Lock()
	s.pending[r.Wallet] = kept
	s.mu.Unlock()

	exp := exposed(r, l.addrs, l.st.Sent)
	var todo coin.UxArray
	var coins, hours uint64
	var contested int
	for _, u := range l.st.Unspent {
		if !exp[u.Body.Address] || covered[u.Hash()] {
			continue
		}
		todo = append(todo, u.UxOut)
		coins += u.Body.Coins
		if h, err := u.CoinHours(l.st.HeadTime); err == nil {
			hours += h
		}
		if len(u.SpentBy) > 0 {
			contested++
		}
	}

	pw, unlocked := s.password(r.Wallet)
	defer erase(pw)
	s.update(r.Wallet, func(st *Status) {
		st.Enabled = true
		st.HeadSeq = l.st.HeadSeq
		st.Exposed = len(exp)
		st.WaitingCoins, st.WaitingHours, st.Contested = coins, hours, contested
		st.Confirmed += confirmed
		st.Displaced = append(st.Displaced, displaced...)
		if n := len(st.Displaced); n > 10 {
			st.Displaced = st.Displaced[n-10:]
		}
		st.Pending = hexes(kept)
		st.Paused = ""
	})
	if len(todo) == 0 {
		return nil
	}

	coinsStr, err := droplet.ToString(coins)
	if err != nil {
		return err
	}
	if l.w.Encrypted() && !unlocked {
		s.update(r.Wallet, func(st *Status) {
			st.Paused = fmt.Sprintf("paused: wallet locked, %s coins waiting on exposed addresses", coinsStr)
		})
		return nil
	}

	dest, err := s.destination(r, l, pw)
	if err != nil {
		return err
	}
	s.update(r.Wallet, func(st *Status) { st.Destination = dest.String() })

	plan, err := Plan(todo, l.st.HeadTime, dest, params.UserVerifyTxn.MaxTransactionSize)
	if errors.Is(err, ErrNoHours) {
		s.update(r.Wallet, func(st *Status) {
			st.Paused = fmt.Sprintf("waiting: %s coins on exposed addresses have no coin hours to burn yet", coinsStr)
		})
		return nil
	}
	if err != nil {
		return err
	}

	if err := s.wallets.ViewSecrets(r.Wallet, nilIfEmpty(pw), func(w wallet.Wallet) error {
		return signPlan(w, plan)
	}); err != nil {
		return err
	}

	var errs []error
	var sent []pendingSweep
	for _, sw := range plan {
		if err := s.inject.InjectBroadcastTransaction(*sw.Txn); err != nil {
			errs = append(errs, err)
			continue
		}
		sent = append(sent, pendingSweep{txid: sw.Txn.Hash(), inputs: sw.Txn.In})
		logger.WithField("wallet", r.Wallet).Infof("swept %d outputs to %s in %s, burning %d coin hours", len(sw.Inputs), dest, sw.Txn.Hash().Hex(), sw.Fee)
	}

	s.mu.Lock()
	s.pending[r.Wallet] = append(s.pending[r.Wallet], sent...)
	all := append([]pendingSweep(nil), s.pending[r.Wallet]...)
	s.mu.Unlock()
	if len(sent) > 0 {
		now := s.now()
		s.update(r.Wallet, func(st *Status) {
			st.LastSweep = &now
			st.LastTxids = hexes(sent)
			st.Pending = hexes(all)
		})
	}
	return errors.Join(errs...)
}

func nilIfEmpty(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	return b
}

func hexes(ps []pendingSweep) []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = p.txid.Hex()
	}
	return out
}

// signPlan signs every input of plan with the keys of w, which must be decrypted.
func signPlan(w wallet.Wallet, plan []Sweep) error {
	entries := make(map[cipher.Address]wallet.Entry)
	add := func(es wallet.Entries) {
		for _, e := range es {
			if a, ok := e.Address.(cipher.Address); ok {
				entries[a] = e
			}
		}
	}
	if w.Type() == wallet.WalletTypeBip44 {
		for _, acct := range w.Accounts() {
			es, err := w.GetEntries(wallet.OptionAccount(acct.Index))
			if err != nil {
				return err
			}
			add(es)
		}
	} else {
		es, err := w.GetEntries()
		if err != nil {
			return err
		}
		add(es)
	}
	for _, sw := range plan {
		for i, in := range sw.Inputs {
			e, ok := entries[in.Body.Address]
			if !ok {
				return fmt.Errorf("address %s is not in the wallet", in.Body.Address)
			}
			if err := wallet.SignInput(sw.Txn, e, i); err != nil {
				return err
			}
		}
		if !sw.Txn.FullySigned() {
			return errors.New("sweep transaction is not fully signed")
		}
	}
	return nil
}

// destination resolves where r sends coins, creating an account or an address
// when there is no unused one, and refuses a destination that is exposed.
func (s *Service) destination(r Rule, l listing, pw []byte) (cipher.Address, error) {
	exp := exposed(r, l.addrs, l.st.Sent)
	d := r.Destination
	var dest cipher.Address
	var err error
	switch d.Kind {
	case DestAccount:
		dest, err = s.accountDestination(r, l, pw)
	case DestWallet:
		dest, err = s.walletDestination(r, l, pw)
	case DestAddress:
		dest, err = cipher.DecodeBase58Address(d.Address)
	default:
		err = fmt.Errorf("unknown destination kind %q", d.Kind)
	}
	if err != nil {
		return cipher.Address{}, err
	}

	sent := l.st.Sent
	if !containsAddr(l.addrs, dest) {
		st, err := s.chain.SweepState([]cipher.Address{dest}, nil)
		if err != nil {
			return cipher.Address{}, err
		}
		sent = st.Sent
	}
	if err := checkDestination(dest, exp, sent); err != nil {
		return cipher.Address{}, err
	}
	return dest, nil
}

// destinationAccount is the account r sends to: the one it names, or the first
// after the shared accounts.
func destinationAccount(r Rule) uint32 {
	if r.Destination.Account != nil {
		return *r.Destination.Account
	}
	if m, ok := r.maxShared(); ok {
		return m + 1
	}
	return 0
}

func (s *Service) accountDestination(r Rule, l listing, pw []byte) (cipher.Address, error) {
	if l.w.Type() != wallet.WalletTypeBip44 {
		return cipher.Address{}, errors.New("an account destination needs a bip44 wallet")
	}
	acct := destinationAccount(r)
	if r.shared()[acct] {
		return cipher.Address{}, fmt.Errorf("%w: account %d is shared", ErrDestinationExposed, acct)
	}
	for uint32(len(l.w.Accounts())) <= acct { //nolint:gosec // account count
		w, err := s.wallets.NewAccount(r.Wallet, nilIfEmpty(pw), "sweep")
		if err != nil {
			return cipher.Address{}, fmt.Errorf("create account %d: %w", acct, err)
		}
		l.w = w
	}
	return s.nextUnused(r.Wallet, l.w, []wallet.Option{wallet.OptionAccount(acct), wallet.OptionExternal()}, l.st.Seen, nil)
}

func (s *Service) walletDestination(r Rule, l listing, pw []byte) (cipher.Address, error) {
	id := r.Destination.Wallet
	if id == "" {
		id = r.Wallet
	}
	w := l.w
	seen := l.st.Seen
	if id != r.Wallet {
		var err error
		if w, err = s.wallets.GetWallet(id); err != nil {
			return cipher.Address{}, err
		}
		addrs, err := walletAddresses(w)
		if err != nil {
			return cipher.Address{}, err
		}
		st, err := s.chain.SweepState(addrsOf(addrs), nil)
		if err != nil {
			return cipher.Address{}, err
		}
		seen = st.Seen
		pw = nil
		if p, ok := s.password(id); ok {
			pw = p
			defer erase(p)
		}
	}
	var opts []wallet.Option
	if w.Type() == wallet.WalletTypeBip44 {
		s.mu.Lock()
		dr, ok := s.rules[id]
		s.mu.Unlock()
		if id == r.Wallet {
			dr, ok = r, true
		}
		if ok && dr.shared()[0] {
			return cipher.Address{}, fmt.Errorf("%w: account 0 of wallet %s is shared", ErrDestinationExposed, id)
		}
		opts = []wallet.Option{wallet.OptionAccount(0), wallet.OptionExternal()}
	}
	return s.nextUnused(id, w, opts, seen, pw)
}

// nextUnused returns the first address after the last one seen on chain, among
// the entries opts select, generating one when there is none.
func (s *Service) nextUnused(wltID string, w wallet.Wallet, opts []wallet.Option, seen map[cipher.Address]bool, pw []byte) (cipher.Address, error) {
	entries, err := w.GetEntries(opts...)
	if err != nil {
		return cipher.Address{}, err
	}
	last := -1
	for i, e := range entries {
		if a, ok := e.Address.(cipher.Address); ok && seen[a] {
			last = i
		}
	}
	if last+1 < len(entries) {
		if a, ok := entries[last+1].Address.(cipher.Address); ok {
			return a, nil
		}
		return cipher.Address{}, errors.New("the destination wallet does not hold skycoin addresses")
	}
	if w.Type() == wallet.WalletTypeCollection {
		return cipher.Address{}, fmt.Errorf("wallet %s has no unused address and a collection wallet cannot generate one", wltID)
	}
	gen := append(opts, wallet.OptionGenerateN(1)) //nolint:gocritic // fresh slice per call
	var addrs []cipher.Address
	if w.Type() == wallet.WalletTypeBip44 || !w.Encrypted() {
		addrs, err = s.wallets.NewAddresses(wltID, nil, gen...)
	} else if len(pw) != 0 {
		addrs, err = s.wallets.NewAddresses(wltID, pw, gen...)
	} else {
		return cipher.Address{}, fmt.Errorf("wallet %s has no unused address and is locked, so it cannot generate one", wltID)
	}
	if err != nil {
		return cipher.Address{}, err
	}
	if len(addrs) != 1 {
		return cipher.Address{}, errors.New("generated an unexpected number of addresses")
	}
	return addrs[0], nil
}

// SetRule validates r, fills in its default destination, saves it and runs a
// pass, which also sweeps what already sits on exposed addresses.
func (s *Service) SetRule(r Rule) (Rule, error) {
	w, err := s.wallets.GetWallet(r.Wallet)
	if err != nil {
		return Rule{}, err
	}
	if w.Type() == wallet.WalletTypeXPub {
		return Rule{}, errors.New("an xpub wallet holds no keys and cannot sweep")
	}
	if r.Destination.Kind == "" {
		r.Destination.Kind = DestWallet
		if w.Type() == wallet.WalletTypeBip44 {
			r.Destination.Kind = DestAccount
		}
	}
	if err := r.check(); err != nil {
		return Rule{}, err
	}
	if len(r.SharedAccounts) > 0 || r.Destination.Kind == DestAccount {
		if w.Type() != wallet.WalletTypeBip44 {
			return Rule{}, errors.New("shared accounts and account destinations need a bip44 wallet")
		}
	}
	n := uint32(len(w.Accounts())) //nolint:gosec // account count
	for _, a := range r.SharedAccounts {
		if a >= n {
			return Rule{}, fmt.Errorf("wallet %s has no account %d", r.Wallet, a)
		}
	}
	if r.Destination.Kind == DestWallet && r.Destination.Wallet != "" {
		if _, err := s.wallets.GetWallet(r.Destination.Wallet); err != nil {
			return Rule{}, err
		}
	}

	addrs, err := walletAddresses(w)
	if err != nil {
		return Rule{}, err
	}
	query := addrsOf(addrs)
	var fixed cipher.Address
	if r.Destination.Kind == DestAddress {
		fixed = cipher.MustDecodeBase58Address(r.Destination.Address)
		query = append(query, fixed)
	}
	st, err := s.chain.SweepState(query, nil)
	if err != nil {
		return Rule{}, err
	}
	if err := checkManual(r, addrs, st.Sent); err != nil {
		return Rule{}, err
	}
	if r.Destination.Kind == DestAddress {
		if err := checkDestination(fixed, exposed(r, addrs, st.Sent), st.Sent); err != nil {
			return Rule{}, err
		}
	}

	s.mu.Lock()
	rules := make(map[string]Rule, len(s.rules)+1)
	for k, v := range s.rules {
		rules[k] = v
	}
	rules[r.Wallet] = r
	if err := saveRules(s.cfg.Dir, rules); err != nil {
		s.mu.Unlock()
		return Rule{}, err
	}
	s.rules = rules
	if st, ok := s.status[r.Wallet]; ok {
		st.Enabled = r.Enabled
		if !r.Enabled {
			st.Paused = "disabled"
		}
	}
	s.mu.Unlock()
	s.poke()
	return r, nil
}

// Rule returns the rule of a wallet.
func (s *Service) Rule(wltID string) (Rule, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.rules[wltID]
	return r, ok
}

// Status returns the sweep status of a wallet.
func (s *Service) Status(wltID string) Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := Status{Wallet: wltID}
	if p, ok := s.status[wltID]; ok {
		st = *p
		st.Pending = append([]string(nil), p.Pending...)
		st.LastTxids = append([]string(nil), p.LastTxids...)
		st.Displaced = append([]string(nil), p.Displaced...)
	}
	r, ok := s.rules[wltID]
	st.Enabled = ok && r.Enabled
	if !st.Enabled {
		st.Paused = "disabled"
	}
	if k, ok := s.keys[wltID]; ok && (k.expires.IsZero() || s.now().Before(k.expires)) {
		st.Unlocked = true
		if !k.expires.IsZero() {
			e := k.expires
			st.UnlockExpires = &e
		}
	}
	return st
}

// Unlock keeps the password of an encrypted wallet in memory so sweeps can be
// signed, for ttl or, when ttl is zero, until Lock or the node stops. The
// password is never written to disk.
func (s *Service) Unlock(wltID string, password []byte, ttl time.Duration) error {
	w, err := s.wallets.GetWallet(wltID)
	if err != nil {
		return err
	}
	if !w.Encrypted() {
		return wallet.ErrWalletNotEncrypted
	}
	if len(password) == 0 {
		return wallet.ErrMissingPassword
	}
	if err := s.wallets.ViewSecrets(wltID, password, func(wallet.Wallet) error { return nil }); err != nil {
		return err
	}
	k := &key{password: append([]byte(nil), password...)}
	if ttl > 0 {
		k.expires = s.now().Add(ttl)
	}
	s.mu.Lock()
	if old, ok := s.keys[wltID]; ok {
		old.erase()
	}
	s.keys[wltID] = k
	s.mu.Unlock()
	s.poke()
	return nil
}

// Lock forgets the password of a wallet.
func (s *Service) Lock(wltID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if k, ok := s.keys[wltID]; ok {
		k.erase()
		delete(s.keys, wltID)
	}
}

// LockAll forgets every password.
func (s *Service) LockAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, k := range s.keys {
		k.erase()
		delete(s.keys, id)
	}
}

func containsAddr(addrs []walletAddr, a cipher.Address) bool {
	for _, w := range addrs {
		if w.addr == a {
			return true
		}
	}
	return false
}

func allIn(hs []cipher.SHA256, set map[cipher.SHA256]bool) bool {
	for _, h := range hs {
		if !set[h] {
			return false
		}
	}
	return true
}
