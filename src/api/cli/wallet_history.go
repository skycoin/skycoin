package cli

import (
	"fmt"

	"path/filepath"

	"strings"

	"bytes"
	"encoding/json"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

type addrHistory struct {
	Txid      string `json:"txid"`
	Address   string `json:"address"`
	Amount    int64  `json:"amount"`
	Timestamp int    `json:"timestamp"`
	Status    int    `json:"status"`
}

func walletHisCMD() gcli.Command {
	name := "walletHistory"
	return gcli.Command{
		Name:         name,
		Usage:        "Display the transaction history of specific wallet",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "f",
				Usage: "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.",
			},
		},
		Action: func(c *gcli.Context) error {
			f := c.String("f")
			if f == "" {
				f = filepath.Join(cfg.WalletDir, cfg.DefaultWalletName)
			}

			// check the file extension.
			if !strings.HasSuffix(f, walletExt) {
				return errWalletName
			}

			// check if file name contains path.
			if filepath.Base(f) != f {
				af, err := filepath.Abs(f)
				if err != nil {
					return fmt.Errorf("invalid wallet file:%v, err:%v", f, err)
				}
				f = af
			} else {
				f = filepath.Join(cfg.WalletDir, f)
			}

			addrs, err := getAddresses(f)
			if err != nil {
				return err
			}

			uxouts, err := getAddrUxOuts(addrs)
			if err != nil {
				return err
			}

			addrHis := []addrHistory{}

			// Uxid          string `json:"uxid"`
			// Time          uint64 `json:"time"`
			// SrcBkSeq      uint64 `json:"src_block_seq"`
			// SrcTx         string `json:"src_tx"`
			// OwnerAddress  string `json:"owner_address"`
			// Coins         uint64 `json:"coins"`
			// Hours         uint64 `json:"hours"`
			// SpentBlockSeq uint64 `json:"spent_block_seq"` // block seq that spent the output.
			// SpentTxID     string `json:"spent_tx"`        // id of tx which spent this output.
			for _, ux := range uxouts {
				for _, u := range ux.UxOuts {
					addrHis = append(addrHis, addrHistory{
						Address:   ux.Address,
						Txid:      u.SrcTx,
						Amount:    int64(u.Coins),
						Timestamp: int(u.Time),
						Status:    1,
					})

					if u.SpentBlockSeq != 0 {
						// get spent transaction timestamp.
					}
				}

			}

			return nil
		},
	}
}

func getAddrUxOuts(addrs []string) ([]webrpc.AddrUxoutResult, error) {
	req, err := webrpc.NewRequest("get_address_uxouts", addrs, "1")
	if err != nil {
		return nil, fmt.Errorf("create rpc request failed:%v", err)
	}

	rsp, err := webrpc.Do(req, cfg.RPCAddress)
	if err != nil {
		return nil, fmt.Errorf("do rpc request failed:%v", err)
	}

	if rsp.Error != nil {
		return nil, fmt.Errorf("do rpc request failed:%+v", *rsp.Error)
	}

	fmt.Println(string(rsp.Result))

	uxouts := []webrpc.AddrUxoutResult{}
	if err := json.NewDecoder(bytes.NewReader(rsp.Result)).Decode(&uxouts); err != nil {
		return nil, fmt.Errorf("decode result failed, err:%v", err)
	}

	return uxouts, nil
}

func getAddresses(f string) ([]string, error) {
	wlt, err := wallet.Load(f)
	if err != nil {
		return nil, err
	}

	addrs := make([]string, len(wlt.Entries))
	for i, entry := range wlt.Entries {
		addrs[i] = entry.Address.String()
	}
	return addrs, nil
}
