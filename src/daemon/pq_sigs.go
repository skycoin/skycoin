package daemon

import (
	"errors"
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
)

// pqSigsProtocolVersion is the first protocol version that understands
// GivePQSigsMessage. Older peers would disconnect on it, so they never get one.
const pqSigsProtocolVersion int32 = 3

const (
	// maxPQSigsPerMessage keeps a GivePQSigsMessage near 160 KB, under the
	// 256 KB default MaxOutgoingMessageLength.
	maxPQSigsPerMessage = 64
	// pqSigCacheSize bounds the signatures held for blocks not yet executed.
	pqSigCacheSize = 4096
)

// BlockPQSig is the ML-DSA-44 signature of one block, by its header hash.
type BlockPQSig struct {
	Hash cipher.SHA256
	Sig  []byte
}

// GivePQSigsMessage carries the post-quantum signatures of blocks. A peer sends it
// right before the GiveBlocksMessage holding those blocks, on the same connection.
type GivePQSigsMessage struct {
	Sigs []BlockPQSig
	c    *gnet.MessageContext `enc:"-"`
}

// EncodeSize implements gnet.Serializer
func (m *GivePQSigsMessage) EncodeSize() uint64 {
	n := uint64(4)
	for _, s := range m.Sigs {
		n += 32 + 4 + uint64(len(s.Sig))
	}
	return n
}

// Encode implements gnet.Serializer
func (m *GivePQSigsMessage) Encode(buf []byte) error {
	if uint64(len(buf)) < m.EncodeSize() {
		return encoder.ErrBufferUnderflow
	}
	if len(m.Sigs) > maxPQSigsPerMessage {
		return encoder.ErrMaxLenExceeded
	}
	e := &encoder.Encoder{Buffer: buf}
	e.Uint32(uint32(len(m.Sigs))) //nolint:gosec // bounded by maxPQSigsPerMessage
	for _, s := range m.Sigs {
		if len(s.Sig) > pqsig.SigSize {
			return encoder.ErrMaxLenExceeded
		}
		e.CopyBytes(s.Hash[:])
		e.Uint32(uint32(len(s.Sig))) //nolint:gosec // bounded by pqsig.SigSize
		e.CopyBytes(s.Sig)
	}
	return nil
}

// Decode implements gnet.Serializer
func (m *GivePQSigsMessage) Decode(buf []byte) (uint64, error) {
	d := &encoder.Decoder{Buffer: buf}
	n, err := d.Uint32()
	if err != nil {
		return 0, err
	}
	if n > maxPQSigsPerMessage {
		return 0, encoder.ErrMaxLenExceeded
	}
	m.Sigs = make([]BlockPQSig, n)
	for i := range m.Sigs {
		if len(d.Buffer) < 32 {
			return 0, encoder.ErrBufferUnderflow
		}
		copy(m.Sigs[i].Hash[:], d.Buffer[:32])
		d.Buffer = d.Buffer[32:]

		l, err := d.Uint32()
		if err != nil {
			return 0, err
		}
		if l > pqsig.SigSize {
			return 0, encoder.ErrMaxLenExceeded
		}
		if uint64(len(d.Buffer)) < uint64(l) {
			return 0, encoder.ErrBufferUnderflow
		}
		m.Sigs[i].Sig = append([]byte(nil), d.Buffer[:l]...)
		d.Buffer = d.Buffer[l:]
	}
	return uint64(len(buf) - len(d.Buffer)), nil //nolint:gosec // d.Buffer is a suffix of buf
}

// Handle handle message
func (m *GivePQSigsMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	m.c = mc
	return daemon.(daemoner).recordMessageEvent(m, mc)
}

// process stores the signatures until their blocks arrive
func (m *GivePQSigsMessage) process(d daemoner) {
	if d.DaemonConfig().DisableNetworking {
		return
	}
	d.addPQSigs(m.Sigs)
}

// newGivePQSigsMessages splits the post-quantum signatures of blocks into messages.
func newGivePQSigsMessages(blocks []coin.SignedBlock) []*GivePQSigsMessage {
	var msgs []*GivePQSigsMessage
	var cur *GivePQSigsMessage
	for _, b := range blocks {
		if len(b.PQSig) == 0 {
			continue
		}
		if cur == nil || len(cur.Sigs) == maxPQSigsPerMessage {
			cur = &GivePQSigsMessage{}
			msgs = append(msgs, cur)
		}
		cur.Sigs = append(cur.Sigs, BlockPQSig{Hash: b.HashHeader(), Sig: b.PQSig})
	}
	return msgs
}

// pqSigCache holds received post-quantum signatures until their blocks execute.
// The oldest entries go first once it is full.
type pqSigCache struct {
	mu    sync.Mutex
	sigs  map[cipher.SHA256][]byte
	order []cipher.SHA256
}

func newPQSigCache() *pqSigCache {
	return &pqSigCache{sigs: make(map[cipher.SHA256][]byte)}
}

func (c *pqSigCache) add(sigs []BlockPQSig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, s := range sigs {
		if _, ok := c.sigs[s.Hash]; ok {
			continue
		}
		for len(c.order) >= pqSigCacheSize {
			delete(c.sigs, c.order[0])
			c.order = c.order[1:]
		}
		c.sigs[s.Hash] = s.Sig
		c.order = append(c.order, s.Hash)
	}
}

func (c *pqSigCache) get(hash cipher.SHA256) []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sigs[hash]
}

// addPQSigs keeps signatures from a peer for blocks that will follow
func (dm *Daemon) addPQSigs(sigs []BlockPQSig) {
	dm.pqSigs.add(sigs)
}

// pqSigFor returns a received post-quantum signature for a block, if any
func (dm *Daemon) pqSigFor(hash cipher.SHA256) []byte {
	return dm.pqSigs.get(hash)
}

// sendPQSigs sends the post-quantum signatures of blocks to addr ahead of the
// blocks themselves, if the peer understands them.
func (dm *Daemon) sendPQSigs(addr string, blocks []coin.SignedBlock) error {
	c := dm.connections.get(addr)
	if c == nil {
		return errors.New("sendPQSigs: no connection")
	}
	if c.ProtocolVersion < pqSigsProtocolVersion {
		return nil
	}
	for _, m := range newGivePQSigsMessages(blocks) {
		if err := dm.sendMessage(addr, m); err != nil {
			return err
		}
	}
	return nil
}

// broadcastPQSigs sends the post-quantum signatures of blocks to every
// introduced peer that understands them.
func (dm *Daemon) broadcastPQSigs(blocks []coin.SignedBlock) error {
	msgs := newGivePQSigsMessages(blocks)
	if len(msgs) == 0 {
		return nil
	}
	var addrs []string
	for _, c := range dm.connections.all() {
		if c.HasIntroduced() && c.ProtocolVersion >= pqSigsProtocolVersion {
			addrs = append(addrs, c.Addr)
		}
	}
	if len(addrs) == 0 {
		return nil
	}
	for _, m := range msgs {
		if _, err := dm.pool.Pool.BroadcastMessage(m, addrs); err != nil {
			return err
		}
	}
	return nil
}
