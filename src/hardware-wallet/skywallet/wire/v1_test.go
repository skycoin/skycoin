package wire

import (
	"bytes"
	"testing"
)

func TestMessageWriteTo(t *testing.T) {
	tests := []struct {
		name     string
		msg      Message
		wantData []byte
	}{
		{
			name: "empty message",
			msg:  Message{Kind: 0, Data: []byte{}},
			wantData: func() []byte {
				// Header: ? # # + kind(2) + size(4) = 9 bytes
				// Then padded to 64 bytes
				buf := make([]byte, 64)
				buf[0] = '?'
				buf[1] = '#'
				buf[2] = '#'
				// kind = 0 (big endian)
				buf[3] = 0
				buf[4] = 0
				// size = 0 (big endian)
				buf[5] = 0
				buf[6] = 0
				buf[7] = 0
				buf[8] = 0
				return buf
			}(),
		},
		{
			name: "small message",
			msg:  Message{Kind: 1, Data: []byte{0x01, 0x02, 0x03}},
			wantData: func() []byte {
				buf := make([]byte, 64)
				buf[0] = '?'
				buf[1] = '#'
				buf[2] = '#'
				// kind = 1 (big endian)
				buf[3] = 0
				buf[4] = 1
				// size = 3 (big endian)
				buf[5] = 0
				buf[6] = 0
				buf[7] = 0
				buf[8] = 3
				// data
				buf[9] = 0x01
				buf[10] = 0x02
				buf[11] = 0x03
				return buf
			}(),
		},
		{
			name: "message with kind 17 (Features)",
			msg:  Message{Kind: 17, Data: []byte{0x0a, 0x07, 'S', 'k', 'y', 'c', 'o', 'i', 'n'}},
			wantData: func() []byte {
				buf := make([]byte, 64)
				buf[0] = '?'
				buf[1] = '#'
				buf[2] = '#'
				// kind = 17 (big endian)
				buf[3] = 0
				buf[4] = 17
				// size = 9 (big endian)
				buf[5] = 0
				buf[6] = 0
				buf[7] = 0
				buf[8] = 9
				// data
				buf[9] = 0x0a
				buf[10] = 0x07
				buf[11] = 'S'
				buf[12] = 'k'
				buf[13] = 'y'
				buf[14] = 'c'
				buf[15] = 'o'
				buf[16] = 'i'
				buf[17] = 'n'
				return buf
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			_, err := tt.msg.WriteTo(&buf)
			if err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			got := buf.Bytes()
			if !bytes.Equal(got, tt.wantData) {
				t.Errorf("WriteTo() got = %v, want %v", got, tt.wantData)
			}
		})
	}
}

func TestMessageWriteToMultiPacket(t *testing.T) {
	// Create a message that spans multiple packets
	// Single packet can hold 64-9=55 bytes of data in first packet
	// and 64-1=63 bytes in subsequent packets
	data := make([]byte, 100) // Will need 2 packets
	for i := range data {
		data[i] = byte(i)
	}

	msg := Message{Kind: 55, Data: data}
	var buf bytes.Buffer
	_, err := msg.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}

	// Should have 2 packets (128 bytes total)
	got := buf.Bytes()
	if len(got) != 128 {
		t.Errorf("WriteTo() wrote %d bytes, want 128", len(got))
	}

	// Check first packet header
	if got[0] != '?' || got[1] != '#' || got[2] != '#' {
		t.Error("First packet missing header magic")
	}

	// Check second packet marker
	if got[64] != '?' {
		t.Errorf("Second packet missing marker, got %c", got[64])
	}
}

func TestReadFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantKind uint16
		wantData []byte
		wantErr  bool
	}{
		{
			name: "empty message",
			input: func() []byte {
				buf := make([]byte, 64)
				buf[0] = '?'
				buf[1] = '#'
				buf[2] = '#'
				return buf
			}(),
			wantKind: 0,
			wantData: []byte{},
		},
		{
			name: "small message",
			input: func() []byte {
				buf := make([]byte, 64)
				buf[0] = '?'
				buf[1] = '#'
				buf[2] = '#'
				buf[3] = 0
				buf[4] = 1
				buf[5] = 0
				buf[6] = 0
				buf[7] = 0
				buf[8] = 3
				buf[9] = 0x01
				buf[10] = 0x02
				buf[11] = 0x03
				return buf
			}(),
			wantKind: 1,
			wantData: []byte{0x01, 0x02, 0x03},
		},
		{
			name: "message with kind 17",
			input: func() []byte {
				buf := make([]byte, 64)
				buf[0] = '?'
				buf[1] = '#'
				buf[2] = '#'
				buf[3] = 0
				buf[4] = 17
				buf[5] = 0
				buf[6] = 0
				buf[7] = 0
				buf[8] = 5
				buf[9] = 'h'
				buf[10] = 'e'
				buf[11] = 'l'
				buf[12] = 'l'
				buf[13] = 'o'
				return buf
			}(),
			wantKind: 17,
			wantData: []byte("hello"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.input)
			msg, err := ReadFrom(r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if msg.Kind != tt.wantKind {
				t.Errorf("ReadFrom() kind = %v, want %v", msg.Kind, tt.wantKind)
			}
			if !bytes.Equal(msg.Data, tt.wantData) {
				t.Errorf("ReadFrom() data = %v, want %v", msg.Data, tt.wantData)
			}
		})
	}
}

func TestReadFromMultiPacket(t *testing.T) {
	// Create a multi-packet message
	data := make([]byte, 100)
	for i := range data {
		data[i] = byte(i)
	}

	// First packet: header + first 55 bytes of data
	buf := make([]byte, 128)
	buf[0] = '?'
	buf[1] = '#'
	buf[2] = '#'
	buf[3] = 0
	buf[4] = 55 // kind
	buf[5] = 0
	buf[6] = 0
	buf[7] = 0
	buf[8] = 100 // size
	copy(buf[9:64], data[:55])

	// Second packet: marker + remaining 45 bytes (padded)
	buf[64] = '?'
	copy(buf[65:], data[55:])

	r := bytes.NewReader(buf)
	msg, err := ReadFrom(r)
	if err != nil {
		t.Fatalf("ReadFrom() error = %v", err)
	}

	if msg.Kind != 55 {
		t.Errorf("ReadFrom() kind = %v, want 55", msg.Kind)
	}
	if !bytes.Equal(msg.Data, data) {
		t.Errorf("ReadFrom() data mismatch")
	}
}

func TestRoundTrip(t *testing.T) {
	// Test that writing and reading produces the same message
	tests := []struct {
		name string
		msg  Message
	}{
		{
			name: "empty",
			msg:  Message{Kind: 0, Data: []byte{}},
		},
		{
			name: "small",
			msg:  Message{Kind: 1, Data: []byte{1, 2, 3}},
		},
		{
			name: "medium",
			msg:  Message{Kind: 17, Data: make([]byte, 50)},
		},
		{
			name: "large",
			msg:  Message{Kind: 100, Data: make([]byte, 200)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill data with recognizable pattern
			for i := range tt.msg.Data {
				tt.msg.Data[i] = byte(i % 256)
			}

			var buf bytes.Buffer
			_, err := tt.msg.WriteTo(&buf)
			if err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			r := bytes.NewReader(buf.Bytes())
			got, err := ReadFrom(r)
			if err != nil {
				t.Fatalf("ReadFrom() error = %v", err)
			}

			if got.Kind != tt.msg.Kind {
				t.Errorf("Round trip kind = %v, want %v", got.Kind, tt.msg.Kind)
			}
			if !bytes.Equal(got.Data, tt.msg.Data) {
				t.Errorf("Round trip data mismatch: got %d bytes, want %d bytes", len(got.Data), len(tt.msg.Data))
			}
		})
	}
}
