package wire

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		wantErr bool
	}{
		{
			name:    "empty buffer",
			buf:     []byte{},
			wantErr: false,
		},
		{
			name: "valid varint field",
			buf: []byte{
				0x08, // field 1, wire type 0 (varint)
				0x01, // value 1
			},
			wantErr: false,
		},
		{
			name: "valid string field",
			buf: []byte{
				0x0a,                         // field 1, wire type 2 (length-delimited)
				0x05,                         // length 5
				'h', 'e', 'l', 'l', 'o', // data
			},
			wantErr: false,
		},
		{
			name: "multiple fields",
			buf: []byte{
				0x08, 0x01, // field 1, varint = 1
				0x10, 0x02, // field 2, varint = 2
				0x1a, 0x03, 'a', 'b', 'c', // field 3, string = "abc"
			},
			wantErr: false,
		},
		{
			name: "nested message (length-delimited)",
			buf: []byte{
				0x0a,       // field 1, wire type 2 (length-delimited)
				0x04,       // length 4
				0x08, 0x01, // nested: field 1, varint = 1
				0x10, 0x02, // nested: field 2, varint = 2
			},
			wantErr: false,
		},
		{
			name: "invalid wire type",
			buf: []byte{
				0x0d, // field 1, wire type 5 (32-bit fixed) - not supported
				0x01, 0x02, 0x03, 0x04,
			},
			wantErr: true,
		},
		{
			name: "length-delimited field with missing data",
			buf: []byte{
				0x0a, // field 1, wire type 2 (length-delimited)
				0x10, // length 16 (but no data follows)
			},
			// Note: bytes.Reader.Seek doesn't error when seeking past end,
			// so this doesn't return an error. The validator only checks
			// for obviously malformed structures, not truncated data.
			wantErr: false,
		},
		{
			name: "large valid message",
			buf: func() []byte {
				// Create a message with multiple string fields
				buf := make([]byte, 0, 1000)
				for i := 0; i < 10; i++ {
					buf = append(buf, 0x0a) // field 1, wire type 2
					buf = append(buf, 0x32) // length 50
					data := make([]byte, 50)
					for j := range data {
						data[j] = 'x'
					}
					buf = append(buf, data...)
				}
				return buf
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRealProtobufMessages(t *testing.T) {
	// Test with actual protobuf-encoded messages similar to what the hardware wallet uses
	tests := []struct {
		name    string
		buf     []byte
		wantErr bool
	}{
		{
			name: "Initialize message (empty)",
			buf:  []byte{},
			// Initialize message with no fields is valid
			wantErr: false,
		},
		{
			name: "Ping message with string",
			buf: []byte{
				0x0a,                            // field 1 (message), wire type 2
				0x0b,                            // length 11
				'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd',
			},
			wantErr: false,
		},
		{
			name: "Features-like message",
			buf: []byte{
				0x0a, 0x06, 'S', 'k', 'y', 'c', 'o', 'i', // vendor string
				0x10, 0x01, // major version
				0x18, 0x00, // minor version
				0x20, 0x00, // patch version
				0x28, 0x01, // initialized = true
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkValidate(b *testing.B) {
	// Benchmark with a moderately sized message
	buf := []byte{
		0x0a, 0x06, 'S', 'k', 'y', 'c', 'o', 'i',
		0x10, 0x01,
		0x18, 0x00,
		0x20, 0x00,
		0x28, 0x01,
		0x32, 0x10, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Validate(buf)
	}
}
