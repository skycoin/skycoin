// Package cli src/cli/sweep_rule.go c3-coin-cli
package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/wallet/sweep"
)

// sweepFlags are the flags of sweepRule.
type sweepFlags struct {
	shared       []uint
	addresses    []string
	toAccount    uint
	toWallet     string
	toAddress    string
	maps         []string
	includeSpent bool
	gap          uint
	enable       bool
	disable      bool
}

// parseSource reads the source of a --map: an account number or "account:N",
// else an address, optionally written "address:X".
func parseSource(s string) (sweep.Source, error) {
	if v, ok := strings.CutPrefix(s, "account:"); ok {
		s = v
	} else if v, ok := strings.CutPrefix(s, "address:"); ok {
		return sweep.Source{Address: v}, nil
	}
	if n, err := strconv.ParseUint(s, 10, 32); err == nil {
		a := uint32(n)
		return sweep.Source{Account: &a}, nil
	}
	if s == "" {
		return sweep.Source{}, errors.New("the source is empty")
	}
	return sweep.Source{Address: s}, nil
}

// parseDestination reads the destination of a --map: "account:N",
// "next-account", "wallet" for the wallet itself, "wallet:ID", or an address,
// optionally written "address:X".
func parseDestination(s string) (sweep.Destination, error) {
	switch {
	case s == "next-account":
		return sweep.Destination{Kind: sweep.DestAccount}, nil
	case s == "wallet":
		return sweep.Destination{Kind: sweep.DestWallet}, nil
	case strings.HasPrefix(s, "account:"):
		n, err := strconv.ParseUint(strings.TrimPrefix(s, "account:"), 10, 32)
		if err != nil {
			return sweep.Destination{}, fmt.Errorf("bad account in %q", s)
		}
		a := uint32(n)
		return sweep.Destination{Kind: sweep.DestAccount, Account: &a}, nil
	case strings.HasPrefix(s, "wallet:"):
		return sweep.Destination{Kind: sweep.DestWallet, Wallet: strings.TrimPrefix(s, "wallet:")}, nil
	case s == "":
		return sweep.Destination{}, errors.New("the destination is empty")
	}
	return sweep.Destination{Kind: sweep.DestAddress, Address: strings.TrimPrefix(s, "address:")}, nil
}

func parseMapping(s string) (sweep.Mapping, error) {
	from, to, ok := strings.Cut(s, "=")
	if !ok {
		return sweep.Mapping{}, fmt.Errorf("--map %q: want <source>=<destination>", s)
	}
	src, err := parseSource(from)
	if err != nil {
		return sweep.Mapping{}, fmt.Errorf("--map %q: %w", s, err)
	}
	dst, err := parseDestination(to)
	if err != nil {
		return sweep.Mapping{}, fmt.Errorf("--map %q: %w", s, err)
	}
	return sweep.Mapping{From: src, To: dst}, nil
}

// applySweepFlags changes base by the flags that were set.
func applySweepFlags(base sweep.Rule, f sweepFlags, changed func(string) bool) (sweep.Rule, error) {
	r := base
	if changed("shared-account") {
		r.SharedAccounts = make([]uint32, len(f.shared))
		for i, a := range f.shared {
			r.SharedAccounts[i] = uint32(a) //nolint:gosec // account index
		}
	}
	if changed("address") {
		r.Addresses = f.addresses
	}
	n := 0
	for _, name := range []string{"to-account", "to-next-account", "to-wallet", "to-address"} {
		if changed(name) {
			n++
		}
	}
	if n > 1 {
		return sweep.Rule{}, errors.New("give one of --to-account, --to-next-account, --to-wallet and --to-address")
	}
	switch {
	case changed("to-account"):
		a := uint32(f.toAccount) //nolint:gosec // account index
		r.Destination = sweep.Destination{Kind: sweep.DestAccount, Account: &a}
	case changed("to-next-account"):
		r.Destination = sweep.Destination{Kind: sweep.DestAccount}
	case changed("to-wallet"):
		r.Destination = sweep.Destination{Kind: sweep.DestWallet, Wallet: f.toWallet}
	case changed("to-address"):
		r.Destination = sweep.Destination{Kind: sweep.DestAddress, Address: f.toAddress}
	}
	if changed("clear-maps") && changed("map") {
		return sweep.Rule{}, errors.New("give --map or --clear-maps, not both")
	}
	if changed("clear-maps") {
		r.Mappings = nil
	}
	if changed("map") {
		r.Mappings = make([]sweep.Mapping, 0, len(f.maps))
		for _, s := range f.maps {
			m, err := parseMapping(s)
			if err != nil {
				return sweep.Rule{}, err
			}
			r.Mappings = append(r.Mappings, m)
		}
	}
	if changed("include-spent") {
		r.IncludeSpent = nil
		if !f.includeSpent {
			off := false
			r.IncludeSpent = &off
		}
	}
	if changed("gap") {
		r.Gap = uint32(f.gap) //nolint:gosec // address count
	}
	if f.enable && f.disable {
		return sweep.Rule{}, errors.New("give --enable or --disable, not both")
	}
	if f.enable {
		r.Enabled = true
	}
	if f.disable {
		r.Enabled = false
	}
	return r, nil
}

func sweepRuleCmd() *cobra.Command {
	var f sweepFlags
	var toNext, clearMaps bool
	cmd := &cobra.Command{
		Use:   "sweepRule [wallet]",
		Short: "Show or set how the node moves coins off exposed addresses of a wallet",
		Long: `Show or set the rule by which the running node moves coins off exposed
    addresses of a wallet it holds, in the block after they arrive. With no flags
    it prints the rule and its status. Flags change the rule, and flags that are
    not given keep their current value.

    An address is exposed when its private key could be computed from public
    data: it has sent, so its public key is on the blockchain, or it belongs to
    a bip44 account whose xpub was shared, such as a skywire reward xpub. A
    single address that has never sent is never swept.

    Level 1 is --enable. It sends every exposed address to one destination. For
    a bip44 wallet that is by default the next unused address of the first
    account after the shared and source accounts, which the node creates when
    needed. --to-account picks an account, --to-wallet the next unused address
    of another wallet, and --to-address a fixed address.

    Level 2 is --map <source>=<destination>, repeatable. The source is an
    address, or an account number for the exposed addresses of that account.
    The destination is account:N, next-account, wallet (this wallet),
    wallet:ID, or an address. Mappings sweep whether level 1 is on or not, and
    exposed addresses no mapping covers go to the level 1 destination when it
    is on, else they stay. --map replaces every mapping and --clear-maps
    removes them.

    --shared-account names a bip44 account whose xpub was shared, so every
    address of it is swept, including those that never sent. Addresses that
    have sent are swept too unless --include-spent=false, which sweeps only
    shared accounts, --address and mapping sources.

    Every destination that is exposed is refused. Each sweep burns every coin
    hour of its inputs, the highest fee the rules allow, so a competing spend
    of the same outputs can at best tie. An encrypted wallet is swept only while
    the node holds its password, given with sweepUnlock.

    Wallets are named by their id in the node, the file name.`,
		Example: `  skycoin-cli sweepRule reward.wlt --shared-account 0 --enable
  skycoin-cli sweepRule reward.wlt --map 2=wallet:cold.wlt --map 2YyBJp...=account:5
  skycoin-cli sweepRule reward.wlt`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			resp, err := apiClient.WalletSweep(args[0])
			if err != nil {
				return err
			}
			if c.Flags().NFlag() == 0 {
				return printJSON(resp)
			}
			base := sweep.Rule{Wallet: args[0]}
			if resp.Rule != nil {
				base = *resp.Rule
			}
			r, err := applySweepFlags(base, f, c.Flags().Changed)
			if err != nil {
				return err
			}
			resp, err = apiClient.SetWalletSweep(r)
			if err != nil {
				return err
			}
			return printJSON(resp)
		},
	}
	cmd.Flags().BoolVar(&f.enable, "enable", false, "turn level 1 on: sweep every exposed address no mapping covers")
	cmd.Flags().BoolVar(&f.disable, "disable", false, "turn level 1 off")
	cmd.Flags().UintSliceVar(&f.shared, "shared-account", nil, "bip44 account whose xpub was shared, repeatable")
	cmd.Flags().StringSliceVar(&f.addresses, "address", nil, "address to sweep, which must have sent, repeatable")
	cmd.Flags().UintVar(&f.toAccount, "to-account", 0, "level 1 sends to the next unused address of this account")
	cmd.Flags().BoolVar(&toNext, "to-next-account", false, "level 1 sends to the first account after the shared and source ones (the default)")
	cmd.Flags().StringVar(&f.toWallet, "to-wallet", "", "level 1 sends to the next unused address of this wallet")
	cmd.Flags().StringVar(&f.toAddress, "to-address", "", "level 1 sends to this address")
	cmd.Flags().StringArrayVar(&f.maps, "map", nil, "<source>=<destination>, repeatable, replaces every mapping")
	cmd.Flags().BoolVar(&clearMaps, "clear-maps", false, "remove every mapping")
	cmd.Flags().BoolVar(&f.includeSpent, "include-spent", true, "also sweep addresses that have sent without being listed")
	cmd.Flags().UintVar(&f.gap, "gap", 0, "unused addresses to keep generated on a shared account (default 20)")
	return cmd
}

func sweepUnlockCmd() *cobra.Command {
	var password string
	var ttl time.Duration
	cmd := &cobra.Command{
		Use:   "sweepUnlock [wallet]",
		Short: "Let the node sign sweeps of an encrypted wallet",
		Long: `Give the node the password of an encrypted wallet, so it can sign sweeps.
    The node keeps it in memory for --ttl, or until it stops or sweepLock when
    --ttl is 0, and never writes it to disk.

    Use caution when using the "-p" option. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			if ttl < 0 {
				return errors.New("--ttl cannot be negative")
			}
			pw, err := NewPasswordReader([]byte(password)).Password()
			if err != nil {
				return err
			}
			resp, err := apiClient.WalletSweepUnlock(api.SweepUnlockRequest{
				ID:         args[0],
				Password:   string(pw),
				TTLSeconds: uint64(ttl / time.Second), //nolint:gosec // non-negative
			})
			for i := range pw {
				pw[i] = 0
			}
			if err != nil {
				return err
			}
			return printJSON(resp)
		},
	}
	cmd.Flags().StringVarP(&password, "password", "p", "", "Wallet password")
	cmd.Flags().DurationVar(&ttl, "ttl", 0, "how long the node keeps the password, 0 until it stops")
	return cmd
}

func sweepLockCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "sweepLock [wallet]",
		Short:        "Make the node forget the sweep password of a wallet",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			resp, err := apiClient.WalletSweepLock(args[0])
			if err != nil {
				return err
			}
			return printJSON(resp)
		},
	}
}
