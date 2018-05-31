package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	gcli "github.com/urfave/cli"

	"github.com/skycoin/skycoin/src/wallet"
)

func decryptWalletCmd(cfg Config) gcli.Command {
	name := "decryptWallet"
	return gcli.Command{
		Name:  name,
		Usage: "Decrypt wallet",
		Description: fmt.Sprintf(`
		The default wallet (%s) will be
		used if no wallet was specified.
		
		Use caution when using the "-p" command. If you have command history enabled 
		your wallet encryption password can be recovered from the history log. If you
		do not include the "-p" option you will be prompted to enter your password
		after you enter your command.`, cfg.FullWalletPath()),
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "p",
				Usage: "[password] Wallet password",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			cfg := ConfigFromContext(c)

			w, err := resolveWalletPath(cfg, "")
			if err != nil {
				return err
			}

			pr := NewPasswordReader([]byte(c.String("p")))

			wlt, err := decryptWallet(w, pr)
			switch err.(type) {
			case nil:
			case WalletLoadError:
				errorWithHelp(c, err)
				return nil
			case WalletSaveError:
				return errors.New("save wallet failed")
			default:
				return err
			}

			printJSON(wallet.NewReadableWallet(wlt))
			return nil
		},
	}
}

func decryptWallet(walletFile string, pr PasswordReader) (*wallet.Wallet, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, WalletLoadError{err}
	}

	if !wlt.IsEncrypted() {
		return nil, wallet.ErrWalletNotEncrypted
	}

	if pr == nil {
		return nil, wallet.ErrMissingPassword
	}

	password, err := pr.Password()
	if err != nil {
		return nil, err
	}

	unlockedWlt, err := wlt.Unlock(password)
	if err != nil {
		return nil, err
	}

	dir, err := filepath.Abs(filepath.Dir(walletFile))
	if err != nil {
		return nil, err
	}

	// save the wallet
	if err := unlockedWlt.Save(dir); err != nil {
		return nil, WalletLoadError{err}
	}

	return unlockedWlt, nil
}
