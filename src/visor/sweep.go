// Package visor src/visor/sweep.go c2-coin-visor
package visor

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/wallet/sweep"
)

// SweepState reads what a sweep pass needs about addrs and txids in one view of
// the database.
func (vs *Visor) SweepState(addrs []cipher.Address, txids []cipher.SHA256) (sweep.ChainState, error) {
	st := sweep.ChainState{
		Sent: make(map[cipher.Address]bool),
		Seen: make(map[cipher.Address]bool),
		Txns: make(map[cipher.SHA256]sweep.TxnState, len(txids)),
	}
	err := vs.db.View("SweepState", func(tx *dbutil.Tx) error {
		head, err := vs.blockchain.Head(tx)
		if err != nil {
			return err
		}
		st.HeadSeq, st.HeadTime = head.Seq(), head.Time()

		for _, a := range addrs {
			seen, err := vs.history.AddressSeen(tx, a)
			if err != nil {
				return err
			}
			if !seen {
				continue
			}
			st.Seen[a] = true
			outs, err := vs.history.GetOutputsForAddress(tx, a)
			if err != nil {
				return err
			}
			for _, o := range outs {
				if o.SpentBlockSeq != 0 {
					st.Sent[a] = true
					break
				}
			}
		}

		hashes, err := vs.blockchain.Unspent().GetUnspentHashesOfAddrs(tx, addrs)
		if err != nil {
			return err
		}
		uxs, err := vs.blockchain.Unspent().GetArray(tx, hashes.Flatten())
		if err != nil {
			return err
		}
		idx := make(map[cipher.SHA256]int, len(uxs))
		st.Unspent = make([]sweep.UxOut, len(uxs))
		for i, u := range uxs {
			st.Unspent[i].UxOut = u
			idx[u.Hash()] = i
		}

		if err := vs.unconfirmed.ForEach(tx, func(h cipher.SHA256, txn UnconfirmedTransaction) error {
			for _, in := range txn.Transaction.In {
				if i, ok := idx[in]; ok {
					st.Unspent[i].SpentBy = append(st.Unspent[i].SpentBy, h)
				}
			}
			return nil
		}); err != nil {
			return err
		}

		for _, h := range txids {
			utxn, err := vs.unconfirmed.Get(tx, h)
			if err != nil {
				return err
			}
			if utxn != nil {
				st.Txns[h] = sweep.TxnPending
				continue
			}
			txn, err := vs.history.GetTransaction(tx, h)
			if err != nil {
				return err
			}
			if txn != nil {
				st.Txns[h] = sweep.TxnConfirmed
			}
		}
		return nil
	})
	return st, err
}
