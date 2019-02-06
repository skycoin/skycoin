// Code generated by github.com/skycoin/skyencoder. DO NOT EDIT.
package blockdb

import (
	"errors"
	"math"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
)

// EncodeSizeHashPairs computes the size of an encoded object of type HashPairs
func EncodeSizeHashPairs(obj *HashPairs) int {
	i0 := 0

	// obj.HashPairs
	i0 += 4
	{
		i1 := 0

		// x.Hash
		i1 += 32

		// x.PreHash
		i1 += 32

		i0 += len(obj.HashPairs) * i1
	}

	return i0
}

// EncodeHashPairs encodes an object of type HashPairs to the buffer in encoder.Encoder.
// The buffer must be large enough to encode the object, otherwise an error is returned.
func EncodeHashPairs(buf []byte, obj *HashPairs) error {
	e := &encoder.Encoder{
		Buffer: buf[:],
	}

	// obj.HashPairs length check
	if len(obj.HashPairs) > math.MaxUint32 {
		return errors.New("obj.HashPairs length exceeds math.MaxUint32")
	}

	// obj.HashPairs length
	e.Uint32(uint32(len(obj.HashPairs)))

	// obj.HashPairs
	for _, x := range obj.HashPairs {

		// x.Hash
		e.CopyBytes(x.Hash[:])

		// x.PreHash
		e.CopyBytes(x.PreHash[:])

	}

	return nil
}

// DecodeHashPairs decodes an object of type HashPairs from the buffer in encoder.Decoder.
// Returns the number of bytes used from the buffer to decode the object.
func DecodeHashPairs(buf []byte, obj *HashPairs) (int, error) {
	d := &encoder.Decoder{
		Buffer: buf[:],
	}

	{
		// obj.HashPairs

		ul, err := d.Uint32()
		if err != nil {
			return len(buf) - len(d.Buffer), err
		}

		length := int(ul)
		if length < 0 || length > len(d.Buffer) {
			return len(buf) - len(d.Buffer), encoder.ErrBufferUnderflow
		}

		if length != 0 {
			obj.HashPairs = make([]coin.HashPair, length)

			for z1 := range obj.HashPairs {
				{
					// obj.HashPairs[z1].Hash
					if len(d.Buffer) < len(obj.HashPairs[z1].Hash) {
						return len(buf) - len(d.Buffer), encoder.ErrBufferUnderflow
					}
					copy(obj.HashPairs[z1].Hash[:], d.Buffer[:len(obj.HashPairs[z1].Hash)])
					d.Buffer = d.Buffer[len(obj.HashPairs[z1].Hash):]
				}

				{
					// obj.HashPairs[z1].PreHash
					if len(d.Buffer) < len(obj.HashPairs[z1].PreHash) {
						return len(buf) - len(d.Buffer), encoder.ErrBufferUnderflow
					}
					copy(obj.HashPairs[z1].PreHash[:], d.Buffer[:len(obj.HashPairs[z1].PreHash)])
					d.Buffer = d.Buffer[len(obj.HashPairs[z1].PreHash):]
				}

			}
		}
	}

	return len(buf) - len(d.Buffer), nil
}