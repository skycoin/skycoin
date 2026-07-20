// Package message pkg/skychat/message/message.go c4-app-chat
// format: length-prefixed frames and the chat-msg/chat-ack envelope.
//
// Before this package the same framing (a 4-byte big-endian length followed by
// that many payload bytes, capped at 64 KiB) was implemented three times — the
// native app (cmd/apps/skychat framedConn), the browser-tab wasm visor
// (cmd/wasm-visor chatConn), and the group relay (group/session.go) — each with
// its own constant and write mutex, free to drift. They are now one codec.
//
// The package is pure Go (only io/net/encoding/sync/time), so it compiles under
// GOOS=js GOARCH=wasm as well as native — the wasm visor and the native app share it.
//
// Wire history: pre-2026-05-12 the protocol was "one Write = one Read" (a raw
// byte slice per message). That held on dmsg (noise-framed) but broke on skynet
// routes, where a VStream can split one Write across several Reads at arbitrary
// boundaries, so a long message arrived split into two chat entries. The
// length-prefixed frame fixed that; old unframed binaries can't talk to framed
// ones. The chat-msg/chat-ack envelope (opt-in, JSON) rides inside a frame; a
// default send writes the plain-text body with no envelope, keeping the wire
// byte-identical for ordinary chat.
package message

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// MaxFrameSize bounds a single framed payload. A larger frame on read is rejected
// as an out-of-sync or hostile peer (or an old, unframed binary).
const MaxFrameSize = 64 * 1024

var errEmptyPayload = errors.New("skychat: empty payload")

// WriteFrame writes payload as one length-prefixed frame: a 4-byte big-endian
// length followed by exactly that many bytes. Rejects an empty or oversize payload.
func WriteFrame(w io.Writer, payload []byte) error {
	if len(payload) == 0 {
		return errEmptyPayload
	}
	if len(payload) > MaxFrameSize {
		return fmt.Errorf("skychat: payload %d > max %d", len(payload), MaxFrameSize)
	}
	var hdr [4]byte
	binary.BigEndian.PutUint32(hdr[:], uint32(len(payload))) //nolint:gosec // len bounded by MaxFrameSize above
	if _, err := w.Write(hdr[:]); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

// ReadFrame reads exactly one length-prefixed frame. Rejects a zero-length or
// oversize frame — the latter usually means the peer is running the pre-framing,
// unbounded protocol.
func ReadFrame(r io.Reader) ([]byte, error) {
	var hdr [4]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(hdr[:])
	if length == 0 {
		return nil, errors.New("skychat: zero-length frame")
	}
	if length > MaxFrameSize {
		return nil, fmt.Errorf("skychat: frame %d > max %d (peer running old unframed protocol?)", length, MaxFrameSize)
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// Conn wraps a net.Conn with length-prefixed framing and a write mutex. The write
// mutex is required because two callers can race to write to the same underlying
// conn (e.g. an outbound message and the ack reply for an inbound one);
// interleaving one frame's length prefix with another's payload would desync the
// receiver permanently. The read path is expected to have a single owner, so
// there is no read mutex.
type Conn struct {
	net.Conn
	writeMu sync.Mutex
}

// NewConn wraps c with framing.
func NewConn(c net.Conn) *Conn { return &Conn{Conn: c} }

// WriteFrame writes a length-prefixed frame, serialized against other writers.
func (c *Conn) WriteFrame(payload []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return WriteFrame(c.Conn, payload)
}

// WriteFrameDeadline is WriteFrame bounded by a per-call write deadline, so a
// stalled peer (transport half-dead, peer not draining) can't pin the caller's
// goroutine indefinitely while holding the write mutex. timeout <= 0 falls back
// to an unbounded WriteFrame. The deadline is acquired AFTER the write mutex so
// its scope is the on-wire write rather than time spent queued behind another
// writer, and reset to "no deadline" (time.Time{}) on the way out so a later call
// isn't pre-expired.
func (c *Conn) WriteFrameDeadline(payload []byte, timeout time.Duration) error {
	if timeout <= 0 {
		return c.WriteFrame(payload)
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	// SetWriteDeadline failure is rare on a healthy net.Conn; fall through to a
	// best-effort unbounded write rather than failing the message on a metadata
	// error. The reset defer is likewise best-effort.
	_ = c.Conn.SetWriteDeadline(time.Now().Add(timeout))        //nolint:errcheck,gosec
	defer func() { _ = c.Conn.SetWriteDeadline(time.Time{}) }() //nolint:errcheck,gosec
	return WriteFrame(c.Conn, payload)
}

// ReadFrame reads one length-prefixed frame from the underlying conn.
func (c *Conn) ReadFrame() ([]byte, error) { return ReadFrame(c.Conn) }

// Envelope is the JSON payload shape for the opt-in peer-receipt protocol:
//
//	chat-msg: {"type":"chat-msg","id":"<hex>","body":"...","ack":true}
//	chat-ack: {"type":"chat-ack","id":"<hex>"}
//
// A default (no --wait) send writes the plain-text body with no envelope, so the
// wire stays byte-identical for ordinary chat; the envelope is opt-in at the sender.
type Envelope struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Body string `json:"body,omitempty"`
	Ack  bool   `json:"ack,omitempty"` // chat-msg only: request ack-on-receipt
}

// Envelope type values.
const (
	TypeMsg = "chat-msg"
	TypeAck = "chat-ack"
)

// Marshal encodes the envelope as its JSON wire bytes.
func (e Envelope) Marshal() ([]byte, error) { return json.Marshal(e) }

// ParseEnvelope conservatively recognizes a chat-msg/chat-ack envelope in a frame
// payload: after trimming surrounding whitespace the payload must begin with '{',
// be valid JSON, and carry a known chat-* type. Anything else returns ok=false so
// a plain-text message that happens to look like JSON still reaches its peer as
// text (only known types are consumed as envelopes).
func ParseEnvelope(payload []byte) (env Envelope, ok bool) {
	t := bytes.TrimSpace(payload)
	if len(t) == 0 || t[0] != '{' {
		return Envelope{}, false
	}
	if err := json.Unmarshal(t, &env); err != nil {
		return Envelope{}, false
	}
	switch env.Type {
	case TypeMsg, TypeAck:
		return env, true
	}
	return Envelope{}, false
}
