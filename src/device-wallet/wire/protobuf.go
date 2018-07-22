package wire

// this file was initially copy pasted from trezor project https://github.com/trezor/trezord-go/blob/master/wire/protobuf.go
// commit 4527402f7004dfe677225315a0dd4d4b1b74be49

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var (
	// ErrMalformedProtobuf malformed protobuf
	ErrMalformedProtobuf = errors.New("malformed protobuf")
)

// Validate checks buffer validity
func Validate(buf []byte) error {
	const (
		wireVarint   = 0               // int32, int64, uint32, uint64, sint32, sint64, bool, enum
		wireData     = 2               // string, bytes, embedded messages, packed repeated fields
		maxFieldSize = 1024 * 1024 * 4 // 4mb field size
	)

	r := bytes.NewReader(buf)

	for r.Len() > 0 {
		// read the field key (combination of tag and type)
		key, err := binary.ReadUvarint(r)
		if err != nil {
			return err
		}

		// validate the field type
		typ := key & 7
		if typ != wireVarint && typ != wireData {
			return ErrMalformedProtobuf
		}

		// read the field value
		val, err := binary.ReadUvarint(r)
		if err != nil {
			return err
		}
		if typ == wireData {
			// field is length-delimited data, skip the data
			if val > maxFieldSize || int64(val) < 0 {
				return ErrMalformedProtobuf
			}
			_, err = r.Seek(int64(val), io.SeekCurrent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
