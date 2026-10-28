package daemon

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
)

// pqWitnessProtocolVersion is the first protocol version that understands
// GivePQWitnessesMessage.
const pqWitnessProtocolVersion int32 = 4

const (
	// maxPQWitnessesPerMessage bounds the witnesses in one message.
	maxPQWitnessesPerMessage = 256
	// maxPQWitnessMessageBytes keeps a GivePQWitnessesMessage under the 256 KB
	// default MaxOutgoingMessageLength, with room for the prefix.
	maxPQWitnessMessageBytes = 240 * 1024
)

// TxnPQWitness is the encoded PQ witness of one transaction (see coin.EncodePQWitness).
type TxnPQWitness struct {
	TxnHash cipher.SHA256
	Witness []byte
}

// GivePQWitnessesMessage carries the PQ witnesses of post-quantum transactions. A
// peer sends it right before the GiveBlocksMessage or GiveTxnsMessage holding them,
// on the same connection, because those encodings leave the witnesses out.
type GivePQWitnessesMessage struct {
	Witnesses []TxnPQWitness
	c         *gnet.MessageContext `enc:"-"`
}

// EncodeSize implements gnet.Serializer
func (m *GivePQWitnessesMessage) EncodeSize() uint64 {
	n := uint64(4)
	for _, w := range m.Witnesses {
		n += 32 + 4 + uint64(len(w.Witness))
	}
	return n
}

// Encode implements gnet.Serializer
func (m *GivePQWitnessesMessage) Encode(buf []byte) error {
	if uint64(len(buf)) < m.EncodeSize() {
		return encoder.ErrBufferUnderflow
	}
	if len(m.Witnesses) > maxPQWitnessesPerMessage {
		return encoder.ErrMaxLenExceeded
	}
	e := &encoder.Encoder{Buffer: buf}
	e.Uint32(uint32(len(m.Witnesses))) //nolint:gosec // bounded by maxPQWitnessesPerMessage
	for _, w := range m.Witnesses {
		if len(w.Witness) > maxPQWitnessMessageBytes {
			return encoder.ErrMaxLenExceeded
		}
		e.CopyBytes(w.TxnHash[:])
		e.Uint32(uint32(len(w.Witness))) //nolint:gosec // bounded by maxPQWitnessMessageBytes
		e.CopyBytes(w.Witness)
	}
	return nil
}

// Decode implements gnet.Serializer
func (m *GivePQWitnessesMessage) Decode(buf []byte) (uint64, error) {
	d := &encoder.Decoder{Buffer: buf}
	n, err := d.Uint32()
	if err != nil {
		return 0, err
	}
	if n > maxPQWitnessesPerMessage {
		return 0, encoder.ErrMaxLenExceeded
	}
	m.Witnesses = make([]TxnPQWitness, n)
	for i := range m.Witnesses {
		if len(d.Buffer) < 32 {
			return 0, encoder.ErrBufferUnderflow
		}
		copy(m.Witnesses[i].TxnHash[:], d.Buffer[:32])
		d.Buffer = d.Buffer[32:]

		l, err := d.Uint32()
		if err != nil {
			return 0, err
		}
		if l > maxPQWitnessMessageBytes {
			return 0, encoder.ErrMaxLenExceeded
		}
		if uint64(len(d.Buffer)) < uint64(l) {
			return 0, encoder.ErrBufferUnderflow
		}
		m.Witnesses[i].Witness = append([]byte(nil), d.Buffer[:l]...)
		d.Buffer = d.Buffer[l:]
	}
	return uint64(len(buf) - len(d.Buffer)), nil //nolint:gosec // d.Buffer is a suffix of buf
}

// Handle handle message
func (m *GivePQWitnessesMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	m.c = mc
	return daemon.(daemoner).recordMessageEvent(m, mc)
}

// process stores the witnesses until their transactions arrive
func (m *GivePQWitnessesMessage) process(d daemoner) {
	if d.DaemonConfig().DisableNetworking {
		return
	}
	d.addPQWitnesses(m.Witnesses)
}

// newGivePQWitnessesMessages collects the witnesses of the post-quantum
// transactions in txns into messages.
func newGivePQWitnessesMessages(txns []coin.Transaction) []*GivePQWitnessesMessage {
	var msgs []*GivePQWitnessesMessage
	var cur *GivePQWitnessesMessage
	var size uint64
	for i := range txns {
		txn := &txns[i]
		if txn.Type != coin.TxnTypePQ {
			continue
		}
		w := TxnPQWitness{TxnHash: txn.Hash(), Witness: coin.EncodePQWitness(txn.PQ)}
		n := uint64(32 + 4 + len(w.Witness))
		if cur == nil || len(cur.Witnesses) == maxPQWitnessesPerMessage || size+n > maxPQWitnessMessageBytes {
			cur = &GivePQWitnessesMessage{}
			msgs = append(msgs, cur)
			size = 4
		}
		cur.Witnesses = append(cur.Witnesses, w)
		size += n
	}
	return msgs
}

func blockTxns(blocks []coin.SignedBlock) []coin.Transaction {
	var txns []coin.Transaction
	for _, b := range blocks {
		txns = append(txns, b.Body.Transactions...)
	}
	return txns
}

// addPQWitnesses keeps witnesses from a peer for transactions that will follow
func (dm *Daemon) addPQWitnesses(ws []TxnPQWitness) {
	sigs := make([]BlockPQSig, len(ws))
	for i, w := range ws {
		sigs[i] = BlockPQSig{Hash: w.TxnHash, Sig: w.Witness}
	}
	dm.pqWitnesses.add(sigs)
}

// attachPQWitness gives a post-quantum transaction that arrived without its witness
// the one a peer sent ahead of it, if any.
func (dm *Daemon) attachPQWitness(txn *coin.Transaction) {
	if txn.Type != coin.TxnTypePQ || len(txn.PQ) != 0 {
		return
	}
	w := dm.pqWitnesses.get(txn.Hash())
	if w == nil {
		return
	}
	if pq, _, err := coin.DecodePQWitness(w); err == nil {
		txn.PQ = pq
	}
}

// sendPQWitnesses sends the witnesses of the post-quantum transactions in txns to
// addr ahead of the transactions themselves, if the peer understands them.
func (dm *Daemon) sendPQWitnesses(addr string, txns []coin.Transaction) error {
	c := dm.connections.get(addr)
	if c == nil {
		return errors.New("sendPQWitnesses: no connection")
	}
	if c.ProtocolVersion < pqWitnessProtocolVersion {
		return nil
	}
	for _, m := range newGivePQWitnessesMessages(txns) {
		if err := dm.sendMessage(addr, m); err != nil {
			return err
		}
	}
	return nil
}

// broadcastPQWitnesses sends the witnesses of the post-quantum transactions in txns
// to every introduced peer that understands them.
func (dm *Daemon) broadcastPQWitnesses(txns []coin.Transaction) error {
	msgs := newGivePQWitnessesMessages(txns)
	if len(msgs) == 0 {
		return nil
	}
	var addrs []string
	for _, c := range dm.connections.all() {
		if c.HasIntroduced() && c.ProtocolVersion >= pqWitnessProtocolVersion {
			addrs = append(addrs, c.Addr)
		}
	}
	if len(addrs) == 0 {
		return nil
	}
	for _, m := range msgs {
		if _, err := dm.out.BroadcastMessage(m, addrs); err != nil {
			return err
		}
	}
	return nil
}
