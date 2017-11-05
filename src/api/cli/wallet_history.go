package cli

import (
	"errors"
	"fmt"

	"time"

	"sort"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

type addrHistory struct {
	BlockSeq  uint64    `json:"-"`
	Txid      string    `json:"txid"`
	Address   string    `json:"address"`
	Amount    string    `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`

	coins uint64 `json:"-"`
}

type byTime []addrHistory

func (obt byTime) Less(i, j int) bool {
	return obt[i].Timestamp.Unix() < obt[j].Timestamp.Unix()
}

func (obt byTime) Swap(i, j int) {
	obt[i], obt[j] = obt[j], obt[i]
}

func (obt byTime) Len() int {
	return len(obt)
}

func walletHisCmd() gcli.Command {
	name := "walletHistory"
	return gcli.Command{
		Name:         name,
		Usage:        "Display the transaction history of specific wallet. Requires skycoin node rpc.",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "f",
				Usage: "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.",
			},
		},
		Action: walletHistoryAction,
	}
}

func walletHistoryAction(c *gcli.Context) error {
	cfg := ConfigFromContext(c)
	rpcClient := RpcClientFromContext(c)

	if c.NArg() > 0 {
		fmt.Printf("Error: invalid argument\n\n")
		gcli.ShowSubcommandHelp(c)
		return nil
	}

	w, err := resolveWalletPath(cfg, c.String("f"))
	if err != nil {
		return err
	}

	// get all addresses in the wallet.
	addrs, err := getAddresses(w)
	if err != nil {
		return err
	}

	if len(addrs) == 0 {
		return errors.New("Wallet is empty")
	}

	// get all the addresses affected uxouts
	uxouts, err := rpcClient.GetAddressUxOuts(addrs)
	if err != nil {
		return err
	}

	// transmute the uxout to addrHistory, and sort the items by time in ascend order.
	totalAddrHis := []addrHistory{}
	for _, ux := range uxouts {
		addrHis, err := makeAddrHisArray(rpcClient, ux)
		if err != nil {
			return err
		}
		totalAddrHis = append(totalAddrHis, addrHis...)
	}

	sort.Sort(byTime(totalAddrHis))

	// print the addr history
	return printJson(totalAddrHis)
}

func makeAddrHisArray(c *webrpc.Client, ux webrpc.AddrUxoutResult) ([]addrHistory, error) {
	if len(ux.UxOuts) == 0 {
		return nil, nil
	}

	var addrHis, spentHis, realHis []addrHistory
	var spentBlkSeqMap = map[uint64]bool{}

	for _, u := range ux.UxOuts {
		amount, err := droplet.ToString(u.Coins)
		if err != nil {
			return nil, err
		}

		addrHis = append(addrHis, addrHistory{
			BlockSeq:  u.SrcBkSeq,
			Txid:      u.SrcTx,
			Address:   ux.Address,
			Amount:    amount,
			Timestamp: time.Unix(int64(u.Time), 0).UTC(),
			Status:    1,
			coins:     u.Coins,
		})

		// the SpentBlockSeq will be 0 if the uxout has not been spent yet.
		if u.SpentBlockSeq != 0 {
			spentBlkSeqMap[u.SpentBlockSeq] = true
			spentHis = append(spentHis, addrHistory{
				BlockSeq: u.SpentBlockSeq,
				Address:  ux.Address,
				Txid:     u.SpentTxID,
				Amount:   "-" + amount,
				Status:   1,
				coins:    u.Coins,
			})
		}
	}

	if len(spentBlkSeqMap) > 0 {
		spentBlkSeq := make([]uint64, 0, len(spentBlkSeqMap))
		for seq := range spentBlkSeqMap {
			spentBlkSeq = append(spentBlkSeq, seq)
		}

		getBlkTime, err := createBlkTimeFinder(c, spentBlkSeq)
		if err != nil {
			return nil, err
		}

		for i, his := range spentHis {
			spentHis[i].Timestamp = time.Unix(getBlkTime(his.BlockSeq), 0).UTC()
		}
	}

	type historyRecord struct {
		received []addrHistory
		spent    []addrHistory
	}

	// merge history in the same transaction.
	hisMap := map[string]historyRecord{}
	for _, his := range addrHis {
		hr := hisMap[his.Txid]
		hr.received = append(hr.received, his)
		hisMap[his.Txid] = hr
	}
	for _, his := range spentHis {
		hr := hisMap[his.Txid]
		hr.spent = append(hr.spent, his)
		hisMap[his.Txid] = hr
	}

	for txid, hs := range hisMap {
		var receivedCoins, spentCoins, coins uint64
		for _, h := range hs.received {
			receivedCoins += h.coins
		}
		for _, h := range hs.spent {
			spentCoins += h.coins
		}

		isNegative := spentCoins > receivedCoins

		if spentCoins > receivedCoins {
			coins = spentCoins - receivedCoins
		} else {
			coins = receivedCoins - spentCoins
		}

		amount, err := droplet.ToString(coins)
		if err != nil {
			return nil, err
		}

		if isNegative {
			amount = "-" + amount
		}

		var his addrHistory
		if len(hs.received) > 0 {
			his = hs.received[0]
		} else {
			his = hs.spent[0]
		}

		realHis = append(realHis, addrHistory{
			BlockSeq:  his.BlockSeq,
			Txid:      txid,
			Address:   ux.Address,
			Amount:    amount,
			Timestamp: his.Timestamp,
			Status:    1,
		})
	}

	return realHis, nil
}

func createBlkTimeFinder(c *webrpc.Client, ss []uint64) (func(uint64) int64, error) {
	// get spent blocks
	blks, err := c.GetBlocksBySeq(ss)
	if err != nil {
		return nil, err
	}

	if len(blks.Blocks) == 0 {
		return nil, fmt.Errorf("found no block")
	}

	return func(seq uint64) int64 {
		for _, b := range blks.Blocks {
			if seq == b.Head.BkSeq {
				return int64(b.Head.Time)
			}
		}
		panic("block not found")
	}, nil
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
