// Package cli src/cli/auto_sweep.go c3-coin-cli
package cli

import (
	"errors"
	"time"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/wallet/sweep"
)

func autoSweepCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autoSweep",
		Short: "Have the node sweep exposed addresses of a wallet on every block",
		Long: `Have the running node move coins off exposed addresses of a wallet it holds,
    in the block after they arrive.

    An address is exposed when its private key could be computed from public
    data: it has sent, so its public key is on the blockchain, or it belongs to
    a bip44 account whose xpub was shared, such as a skywire reward xpub. A
    single address that has never sent is never swept.

    Each sweep burns every coin hour of its inputs, the highest fee the rules
    allow, so a competing spend of the same outputs can at best tie.

    An encrypted wallet is swept only while the node holds its password, set
    with "autoSweep unlock". The node keeps it in memory only.

    Wallets are named by their id in the node, the file name.`,
	}
	cmd.AddCommand(autoSweepShowCmd(), autoSweepSetCmd(), autoSweepUnlockCmd(), autoSweepLockCmd())
	return cmd
}

func autoSweepShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "show [wallet]",
		Short:        "Show the sweep rule of a wallet and its status",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			resp, err := apiClient.WalletSweep(args[0])
			if err != nil {
				return err
			}
			return printJSON(resp)
		},
	}
}

// sweepFlags are the flags of autoSweep set.
type sweepFlags struct {
	shared    []uint
	addresses []string
	toAccount uint
	toWallet  string
	toAddress string
	gap       uint
	enable    bool
	disable   bool
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

func autoSweepSetCmd() *cobra.Command {
	var f sweepFlags
	var toNext bool
	cmd := &cobra.Command{
		Use:   "set [wallet]",
		Short: "Create or change the sweep rule of a wallet",
		Long: `Create or change the sweep rule of a wallet. Flags that are not given keep
    their current value. A new rule is off until --enable.

    --shared-account names a bip44 account whose xpub was shared. Every
    address of it is swept, including those that never sent. Repeat it for
    several accounts. Addresses that have sent are swept in any case.

    The destination of a bip44 wallet is by default the next unused address of
    the first account after the shared ones, which the node creates when
    needed. --to-account picks a later account, --to-wallet the next unused
    address of another wallet, and --to-address a fixed address. A destination
    that is exposed is refused.`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			base := sweep.Rule{Wallet: args[0]}
			resp, err := apiClient.WalletSweep(args[0])
			if err != nil {
				return err
			}
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
	cmd.Flags().UintSliceVar(&f.shared, "shared-account", nil, "bip44 account whose xpub was shared, repeatable")
	cmd.Flags().StringSliceVar(&f.addresses, "address", nil, "address to sweep, which must have sent, repeatable")
	cmd.Flags().UintVar(&f.toAccount, "to-account", 0, "send to the next unused address of this account")
	cmd.Flags().BoolVar(&toNext, "to-next-account", false, "send to the first account after the shared ones (the default)")
	cmd.Flags().StringVar(&f.toWallet, "to-wallet", "", "send to the next unused address of this wallet")
	cmd.Flags().StringVar(&f.toAddress, "to-address", "", "send to this address")
	cmd.Flags().UintVar(&f.gap, "gap", 0, "unused addresses to keep generated on a shared account (default 20)")
	cmd.Flags().BoolVar(&f.enable, "enable", false, "turn the rule on")
	cmd.Flags().BoolVar(&f.disable, "disable", false, "turn the rule off")
	return cmd
}

func autoSweepUnlockCmd() *cobra.Command {
	var password string
	var ttl time.Duration
	cmd := &cobra.Command{
		Use:   "unlock [wallet]",
		Short: "Let the node sign sweeps of an encrypted wallet",
		Long: `Give the node the password of an encrypted wallet, so it can sign sweeps.
    The node keeps it in memory for --ttl, or until it stops or "autoSweep lock"
    when --ttl is 0, and never writes it to disk.

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

func autoSweepLockCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "lock [wallet]",
		Short:        "Make the node forget the password of a wallet",
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
