package message

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"strings"
	"testing"
)

// TestFrameRoundTrip verifies the exact on-wire layout (4-byte big-endian length
// + payload) and that ReadFrame recovers what WriteFrame wrote, back to back.
func TestFrameRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	msgs := [][]byte{[]byte("hi"), []byte("a longer message with spaces"), {0x00, 0x01, 0x02}}
	for _, m := range msgs {
		if err := WriteFrame(&buf, m); err != nil {
			t.Fatalf("WriteFrame(%q): %v", m, err)
		}
	}
	// Exact layout of the first frame.
	if got := binary.BigEndian.Uint32(buf.Bytes()[:4]); int(got) != len(msgs[0]) {
		t.Fatalf("first frame length prefix = %d, want %d", got, len(msgs[0]))
	}
	for i, want := range msgs {
		got, err := ReadFrame(&buf)
		if err != nil {
			t.Fatalf("ReadFrame #%d: %v", i, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("ReadFrame #%d = %q, want %q", i, got, want)
		}
	}
}

func TestWriteFrameRejects(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteFrame(&buf, nil); err == nil {
		t.Fatal("empty payload: want error, got nil")
	}
	if err := WriteFrame(&buf, make([]byte, MaxFrameSize+1)); err == nil {
		t.Fatal("oversize payload: want error, got nil")
	}
	if buf.Len() != 0 {
		t.Fatalf("rejected writes should emit no bytes, got %d", buf.Len())
	}
}

func TestReadFrameRejects(t *testing.T) {
	// zero-length frame
	zero := []byte{0, 0, 0, 0}
	if _, err := ReadFrame(bytes.NewReader(zero)); err == nil {
		t.Fatal("zero-length frame: want error, got nil")
	}
	// oversize length prefix
	var big [4]byte
	binary.BigEndian.PutUint32(big[:], MaxFrameSize+1)
	if _, err := ReadFrame(bytes.NewReader(big[:])); err == nil {
		t.Fatal("oversize frame: want error, got nil")
	}
	// truncated payload
	trunc := []byte{0, 0, 0, 8, 'a', 'b'}
	if _, err := ReadFrame(bytes.NewReader(trunc)); err != io.ErrUnexpectedEOF {
		t.Fatalf("truncated payload: want ErrUnexpectedEOF, got %v", err)
	}
}

// TestConnFramingCompat proves the *Conn wrapper writes the same bytes the raw
// WriteFrame does (so it interoperates with any peer using either API).
func TestConnFramingCompat(t *testing.T) {
	a, b := net.Pipe()
	fc := NewConn(a)
	// net.Pipe write blocks until the other side reads, so write from a goroutine
	// and collect its error afterward.
	errc := make(chan error, 1)
	go func() { errc <- fc.WriteFrame([]byte("ping")) }()
	got, err := ReadFrame(b)
	if err != nil {
		t.Fatalf("ReadFrame from raw side: %v", err)
	}
	if werr := <-errc; werr != nil {
		t.Fatalf("Conn.WriteFrame: %v", werr)
	}
	if string(got) != "ping" {
		t.Fatalf("got %q, want ping", got)
	}
}

func TestParseEnvelope(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantOK  bool
		wantTyp string
		wantID  string
		wantAck bool
	}{
		{"chat-msg with ack", `{"type":"chat-msg","id":"deadbeef","body":"hi","ack":true}`, true, TypeMsg, "deadbeef", true},
		{"chat-ack", `{"type":"chat-ack","id":"deadbeef"}`, true, TypeAck, "deadbeef", false},
		{"leading whitespace", "  {\"type\":\"chat-ack\",\"id\":\"x\"}", true, TypeAck, "x", false},
		{"plain text", "hello world", false, "", "", false},
		{"json-looking but unknown type", `{"type":"something-else"}`, false, "", "", false},
		{"literal brace text", "{not json}", false, "", "", false},
		{"empty", "", false, "", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			env, ok := ParseEnvelope([]byte(c.in))
			if ok != c.wantOK {
				t.Fatalf("ok = %v, want %v", ok, c.wantOK)
			}
			if !ok {
				return
			}
			if env.Type != c.wantTyp || env.ID != c.wantID || env.Ack != c.wantAck {
				t.Fatalf("env = %+v, want type=%s id=%s ack=%v", env, c.wantTyp, c.wantID, c.wantAck)
			}
		})
	}
}

// TestEnvelopeMarshalDefaultBytes locks the wire bytes: a chat-ack omits body+ack.
func TestEnvelopeMarshalDefaultBytes(t *testing.T) {
	b, err := Envelope{Type: TypeAck, ID: "abc"}.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	// Check for the quoted JSON keys, not bare substrings — the type value
	// "chat-ack" itself contains "ack".
	if got := string(b); strings.Contains(got, `"body"`) || strings.Contains(got, `"ack"`) {
		t.Fatalf("chat-ack should omit body/ack keys, got %s", got)
	}
	env, ok := ParseEnvelope(b)
	if !ok || env.Type != TypeAck || env.ID != "abc" {
		t.Fatalf("round-trip failed: ok=%v env=%+v", ok, env)
	}
}
