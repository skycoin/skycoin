package cli

import (
	"errors"
	"fmt"

	"time"

	"sort"

	cobra "github.com/spf13/cobra"

	"github.com/SkycoinProject/skycoin/src/api"
	"github.com/SkycoinProject/skycoin/src/readable"
	"github.com/SkycoinProject/skycoin/src/util/droplet"
	"github.com/SkycoinProject/skycoin/src/wallet"
)

// AddrHistory represents a transactional event for an address
type AddrHistory struct {
	BlockSeq  uint64    `json:"-"`
	Txid      string    `json:"txid"`
	Address   string    `json:"address"`
	Amount    string    `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`

	coins uint64
}

type byTime []AddrHistory

func (obt byTime) Less(i, j int) bool {
	return obt[i].Timestamp.Unix() < obt[j].Timestamp.Unix()
}

func (obt byTime) Swap(i, j int) {
	obt[i], obt[j] = obt[j], obt[i]
}

func (obt byTime) Len() int {
	return len(obt)
}

func walletHisCmd() *cobra.Command {
	walletHisCmd := &cobra.Command{
		Short:        "Display the transaction history of specific wallet. Requires skycoin node rpc.",
		Use:          "walletHistory [wallet]",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE:         walletHistoryAction,
	}

	return walletHisCmd
}

func walletHistoryAction(c *cobra.Command, args []string) error {
	w := args[0]

	// Get all addresses in the wallet
	addrs, err := getAddresses(w)
	if err != nil {
		return err
	}

	if len(addrs) == 0 {
		return errors.New("Wallet is empty")
	}

	// Get all the addresses' historical uxouts
	totalAddrHis := []AddrHistory{}
	for _, addr := range addrs {
		uxouts, err := apiClient.AddressUxOuts(addr)
		if err != nil {
			return err
		}

		addrHis, err := makeAddrHisArray(apiClient, addr, uxouts)
		if err != nil {
			return err
		}
		totalAddrHis = append(totalAddrHis, addrHis...)
	}

	// Sort the uxouts by time ascending
	sort.Sort(byTime(totalAddrHis))

	return printJSON(totalAddrHis)
}

func makeAddrHisArray(c *api.Client, addr string, uxOuts []readable.SpentOutput) ([]AddrHistory, error) {
	if len(uxOuts) == 0 {
		return nil, nil
	}

	var addrHis, spentHis, realHis []AddrHistory
	var spentBlkSeqMap = map[uint64]bool{}

	for _, u := range uxOuts {
		amount, err := droplet.ToString(u.Coins)
		if err != nil {
			return nil, err
		}

		addrHis = append(addrHis, AddrHistory{
			BlockSeq:  u.SrcBkSeq,
			Txid:      u.SrcTx,
			Address:   addr,
			Amount:    amount,
			Timestamp: time.Unix(int64(u.Time), 0).UTC(),
			Status:    1,
			coins:     u.Coins,
		})

		// the SpentBlockSeq will be 0 if the uxout has not been spent yet.
		if u.SpentBlockSeq != 0 {
			spentBlkSeqMap[u.SpentBlockSeq] = true
			spentHis = append(spentHis, AddrHistory{
				BlockSeq: u.SpentBlockSeq,
				Address:  addr,
				Txid:     u.SpentTxnID,
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
		received []AddrHistory
		spent    []AddrHistory
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

		var his AddrHistory
		if len(hs.received) > 0 {
			his = hs.received[0]
		} else {
			his = hs.spent[0]
		}

		realHis = append(realHis, AddrHistory{
			BlockSeq:  his.BlockSeq,
			Txid:      txid,
			Address:   addr,
			Amount:    amount,
			Timestamp: his.Timestamp,
			Status:    1,
		})
	}

	return realHis, nil
}

func createBlkTimeFinder(c *api.Client, ss []uint64) (func(uint64) int64, error) {
	// get spent blocks
	blocks := make([]*readable.Block, 0, len(ss))
	for _, s := range ss {
		block, err := c.BlockBySeq(s)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, block)
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("found no block")
	}

	return func(seq uint64) int64 {
		for _, b := range blocks {
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

	addrs := wlt.GetAddresses()

	strAddrs := make([]string, len(addrs))
	for i, a := range addrs {
		strAddrs[i] = a.String()
	}

	return strAddrs, nil
}
