package wire

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	repMarker = '?'
	repMagic  = '#'
	packetLen = 64
)

type Message struct {
	Kind uint16
	Data []byte
}

func (m *Message) WriteTo(w io.Writer) (int64, error) {
	var (
		rep  [packetLen]byte
		kind = m.Kind
		size = uint32(len(m.Data))
	)
	// pack header
	rep[0] = repMarker
	rep[1] = repMagic
	rep[2] = repMagic
	binary.BigEndian.PutUint16(rep[3:], kind)
	binary.BigEndian.PutUint32(rep[5:], size)

	var (
		written = 0 // number of written bytes
		offset  = 9 // just after the header
	)
	for written < len(m.Data) {
		n := copy(rep[offset:], m.Data[written:])
		written += n
		offset += n
		if offset >= len(rep) {
			_, err := w.Write(rep[:])
			if err != nil {
				return int64(written), err
			}
			offset = 1 // just after the marker
		}
	}
	if offset != 1 {
		for offset < len(rep) {
			rep[offset] = 0x00
			offset++
		}
		_, err := w.Write(rep[:])
		if err != nil {
			return int64(written), err
		}
	}

	return int64(written), nil
}

var (
	ErrMalformedMessage = errors.New("malformed wire format")
)

func (m *Message) ReadFrom(r io.Reader) (int64, error) {
	var (
		rep  [packetLen]byte
		read = 0 // number of read bytes
	)
	n, err := r.Read(rep[:])
	if err != nil {
		return int64(read), err
	}
	read += n
	if rep[0] != repMarker || rep[1] != repMagic || rep[2] != repMagic {
		return int64(read), ErrMalformedMessage
	}

	// parse header
	var (
		kind = binary.BigEndian.Uint16(rep[3:])
		size = binary.BigEndian.Uint32(rep[5:])
		data = make([]byte, 0, size)
	)
	data = append(data, rep[9:]...) // read data after header

	for uint32(len(data)) < size {
		n, err := r.Read(rep[:])
		if err != nil {
			return int64(read), err
		}
		if rep[0] != repMarker {
			return int64(read), ErrMalformedMessage
		}
		read += n
		data = append(data, rep[1:]...) // read data after marker
	}
	data = data[:size]

	m.Kind = kind
	m.Data = data

	return int64(read), nil
}
