// Code generated by github.com/skycoin/skyencoder. DO NOT EDIT.
package historydb

import "github.com/skycoin/skycoin/src/cipher/encoder"

// EncodeSizeUxOut computes the size of an encoded object of type UxOut
func EncodeSizeUxOut(obj *UxOut) int {
	i0 := 0

	// obj.Out.Head.Time
	i0 += 8

	// obj.Out.Head.BkSeq
	i0 += 8

	// obj.Out.Body.SrcTransaction
	i0 += 32

	// obj.Out.Body.Address.Version
	i0++

	// obj.Out.Body.Address.Key
	i0 += 20

	// obj.Out.Body.Coins
	i0 += 8

	// obj.Out.Body.Hours
	i0 += 8

	// obj.SpentTxnID
	i0 += 32

	// obj.SpentBlockSeq
	i0 += 8

	return i0
}

// EncodeUxOut encodes an object of type UxOut to the buffer in encoder.Encoder.
// The buffer must be large enough to encode the object, otherwise an error is returned.
func EncodeUxOut(buf []byte, obj *UxOut) error {
	e := &encoder.Encoder{
		Buffer: buf[:],
	}

	// obj.Out.Head.Time
	e.Uint64(obj.Out.Head.Time)

	// obj.Out.Head.BkSeq
	e.Uint64(obj.Out.Head.BkSeq)

	// obj.Out.Body.SrcTransaction
	e.CopyBytes(obj.Out.Body.SrcTransaction[:])

	// obj.Out.Body.Address.Version
	e.Uint8(obj.Out.Body.Address.Version)

	// obj.Out.Body.Address.Key
	e.CopyBytes(obj.Out.Body.Address.Key[:])

	// obj.Out.Body.Coins
	e.Uint64(obj.Out.Body.Coins)

	// obj.Out.Body.Hours
	e.Uint64(obj.Out.Body.Hours)

	// obj.SpentTxnID
	e.CopyBytes(obj.SpentTxnID[:])

	// obj.SpentBlockSeq
	e.Uint64(obj.SpentBlockSeq)

	return nil
}

// DecodeUxOut decodes an object of type UxOut from the buffer in encoder.Decoder.
// Returns the number of bytes used from the buffer to decode the object.
func DecodeUxOut(buf []byte, obj *UxOut) (int, error) {
	d := &encoder.Decoder{
		Buffer: buf[:],
	}

	{
		// obj.Out.Head.Time
		i, err := d.Uint64()
		if err != nil {
			return len(buf) - len(d.Buffer), err
		}
		obj.Out.Head.Time = i
	}

	{
		// obj.Out.Head.BkSeq
		i, err := d.Uint64()
		if err != nil {
			return len(buf) - len(d.Buffer), err
		}
		obj.Out.Head.BkSeq = i
	}

	{
		// obj.Out.Body.SrcTransaction
		if len(d.Buffer) < len(obj.Out.Body.SrcTransaction) {
			return len(buf) - len(d.Buffer), encoder.ErrBufferUnderflow
		}
		copy(obj.Out.Body.SrcTransaction[:], d.Buffer[:len(obj.Out.Body.SrcTransaction)])
		d.Buffer = d.Buffer[len(obj.Out.Body.SrcTransaction):]
	}

	{
		// obj.Out.Body.Address.Version
		i, err := d.Uint8()
		if err != nil {
			return len(buf) - len(d.Buffer), err
		}
		obj.Out.Body.Address.Version = i
	}

	{
		// obj.Out.Body.Address.Key
		if len(d.Buffer) < len(obj.Out.Body.Address.Key) {
			return len(buf) - len(d.Buffer), encoder.ErrBufferUnderflow
		}
		copy(obj.Out.Body.Address.Key[:], d.Buffer[:len(obj.Out.Body.Address.Key)])
		d.Buffer = d.Buffer[len(obj.Out.Body.Address.Key):]
	}

	{
		// obj.Out.Body.Coins
		i, err := d.Uint64()
		if err != nil {
			return len(buf) - len(d.Buffer), err
		}
		obj.Out.Body.Coins = i
	}

	{
		// obj.Out.Body.Hours
		i, err := d.Uint64()
		if err != nil {
			return len(buf) - len(d.Buffer), err
		}
		obj.Out.Body.Hours = i
	}

	{
		// obj.SpentTxnID
		if len(d.Buffer) < len(obj.SpentTxnID) {
			return len(buf) - len(d.Buffer), encoder.ErrBufferUnderflow
		}
		copy(obj.SpentTxnID[:], d.Buffer[:len(obj.SpentTxnID)])
		d.Buffer = d.Buffer[len(obj.SpentTxnID):]
	}

	{
		// obj.SpentBlockSeq
		i, err := d.Uint64()
		if err != nil {
			return len(buf) - len(d.Buffer), err
		}
		obj.SpentBlockSeq = i
	}

	return len(buf) - len(d.Buffer), nil
}