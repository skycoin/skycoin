package cli

import (
	"fmt"

	"path/filepath"

	"strings"

	"bytes"
	"encoding/json"

	"time"

	"sort"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

type addrHistory struct {
	BlockSeq  uint64    `json:"-"`
	Txid      string    `json:"txid"`
	Address   string    `json:"address"`
	Amount    int64     `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
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
		Action: walletHistoryAction,
	}
}

func walletHistoryAction(c *gcli.Context) error {
	if c.NArg() > 0 {
		fmt.Printf("Error: invalid argument\n\n")
		gcli.ShowSubcommandHelp(c)
		return nil
	}
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

	// get all addresses in the wallet.
	addrs, err := getAddresses(f)
	if err != nil {
		return err
	}

	// get all the addresses affected uxouts
	uxouts, err := getAddrUxOuts(addrs)
	if err != nil {
		return err
	}

	// transmute the uxout to addrHistory, and sort the items by time in ascend order.
	totalAddrHis := []addrHistory{}
	for _, ux := range uxouts {
		addrHis, err := makeAddrHisArray(ux)
		if err != nil {
			return err
		}
		totalAddrHis = append(totalAddrHis, addrHis...)
	}

	sort.Sort(byTime(totalAddrHis))

	// print the addr history
	v, err := json.MarshalIndent(totalAddrHis, "", "    ")
	if err != nil {
		return errJSONMarshal
	}
	fmt.Println(string(v))
	return nil
}

func makeAddrHisArray(ux webrpc.AddrUxoutResult) ([]addrHistory, error) {
	if len(ux.UxOuts) == 0 {
		return []addrHistory{}, nil
	}

	var (
		addrHis        = []addrHistory{}
		spentHis       = []addrHistory{}
		spentBlkSeqMap = map[uint64]bool{}
	)

	for _, u := range ux.UxOuts {
		addrHis = append(addrHis, addrHistory{
			BlockSeq:  u.SrcBkSeq,
			Txid:      u.SrcTx,
			Address:   ux.Address,
			Amount:    int64(u.Coins) / 1e6,
			Timestamp: time.Unix(int64(u.Time), 0).UTC(),
			Status:    1,
		})

		// the SpentBlockSeq will be 0 if the uxout has not been spent yet.
		if u.SpentBlockSeq != 0 {
			spentBlkSeqMap[u.SpentBlockSeq] = true
			spentHis = append(spentHis, addrHistory{
				BlockSeq: u.SpentBlockSeq,
				Address:  ux.Address,
				Txid:     u.SpentTxID,
				Amount:   (int64(u.Coins) * -1) / 1e6,
				Status:   1,
			})
		}
	}

	spentBlkSeq := make([]uint64, 0, len(spentBlkSeqMap))
	for seq := range spentBlkSeqMap {
		spentBlkSeq = append(spentBlkSeq, seq)
	}

	if len(spentBlkSeq) > 0 {
		getBlkTime, err := createBlkTimeFinder(spentBlkSeq)
		if err != nil {
			return []addrHistory{}, err
		}

		for i, his := range spentHis {
			spentHis[i].Timestamp = time.Unix(getBlkTime(his.BlockSeq), 0).UTC()
		}
		addrHis = append(addrHis, spentHis...)
	}

	// merge history in the same transaction.
	hisMap := map[string][]addrHistory{}
	for _, his := range addrHis {
		hisMap[his.Txid] = append(hisMap[his.Txid], his)
	}

	realHis := []addrHistory{}
	for txid, hs := range hisMap {
		var amt int64
		for _, h := range hs {
			amt += h.Amount
		}
		realHis = append(realHis, addrHistory{
			BlockSeq:  hs[0].BlockSeq,
			Txid:      txid,
			Address:   ux.Address,
			Amount:    amt,
			Timestamp: hs[0].Timestamp,
			Status:    1,
		})
	}

	return realHis, nil
}

func createBlkTimeFinder(ss []uint64) (func(uint64) int64, error) {
	// get spent blocks
	blks, err := getBlocksBySeq(ss)
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
