// Code generated by github.com/skycoin/skyencoder. DO NOT EDIT.
package blockdb

import "github.com/skycoin/skycoin/src/cipher/encoder"

// EncodeSizeSig computes the size of an encoded object of type Sig
func EncodeSizeSig(obj *Sig) int {
	i0 := 0

	// obj.Sig
	i0 += 65

	return i0
}

// EncodeSig encodes an object of type Sig to the buffer in encoder.Encoder.
// The buffer must be large enough to encode the object, otherwise an error is returned.
func EncodeSig(buf []byte, obj *Sig) error {
	e := &encoder.Encoder{
		Buffer: buf[:],
	}

	// obj.Sig
	e.CopyBytes(obj.Sig[:])

	return nil
}

// DecodeSig decodes an object of type Sig from the buffer in encoder.Decoder.
// Returns the number of bytes used from the buffer to decode the object.
func DecodeSig(buf []byte, obj *Sig) (int, error) {
	d := &encoder.Decoder{
		Buffer: buf[:],
	}

	{
		// obj.Sig
		if len(d.Buffer) < len(obj.Sig) {
			return len(buf) - len(d.Buffer), encoder.ErrBufferUnderflow
		}
		copy(obj.Sig[:], d.Buffer[:len(obj.Sig)])
		d.Buffer = d.Buffer[len(obj.Sig):]
	}

	return len(buf) - len(d.Buffer), nil
}