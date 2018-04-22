package cli

import (
	"fmt"

	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

func showSeedCmd(cfg Config) gcli.Command {
	name := "showSeed"
	return gcli.Command{
		Name:  name,
		Usage: "Show wallet seed",
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
				Usage: "[password] Wallet password, if encrypted",
			},
			gcli.BoolFlag{
				Name:  "j,json",
				Usage: "Returns the results in JSON format",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			cfg := ConfigFromContext(c)

			w, err := resolveWalletPath(cfg, "")
			if err != nil {
				return err
			}

			seed, err := getSeed(w, []byte(c.String("p")))
			switch err.(type) {
			case nil:
			case WalletLoadError:
				errorWithHelp(c, err)
				return nil
			default:
				return err
			}

			if c.Bool("j") {
				v := struct {
					Seed string `json:"seed"`
				}{
					Seed: seed,
				}

				printJSON(v)
				return nil
			}

			fmt.Println(seed)
			return nil
		},
	}
}

func getSeed(walletFile string, password []byte) (string, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return "", WalletLoadError{err}
	}

	if !wlt.IsEncrypted() {
		return wlt.Meta["seed"], nil
	}

	if len(password) == 0 {
		var err error
		password, err = readPasswordFromTerminal()
		if err != nil {
			return "", err
		}
	}

	var seed string
	if err := wlt.GuardView(password, func(w *wallet.Wallet) error {
		seed = w.Meta["seed"]
		return nil
	}); err != nil {
		return "", err
	}

	return seed, nil
}
