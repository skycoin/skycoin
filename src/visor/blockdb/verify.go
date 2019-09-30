package blockdb

import (
	"errors"
	"reflect"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
)

var (
	// ErrVerifyStopped is returned when database verification is interrupted
	ErrVerifyStopped = errors.New("database verification stopped")
)

// VerifyDBSkyencoderSafe verifies that the skyencoder generated code has the same result as the encoder
// for all data in the blockchain
func VerifyDBSkyencoderSafe(tx *dbutil.Tx, quit <-chan struct{}) error {
	if quit == nil {
		quit = make(chan struct{})
	}

	if err := dbutil.ForEach(tx, BlockSigsBkt, func(_, v []byte) error {
		select {
		case <-quit:
			return ErrVerifyStopped
		default:
		}

		var sig1 sigWrapper
		if err := decodeSigWrapperExact(v, &sig1); err != nil {
			return err
		}

		var sig2 cipher.Sig
		if err := encoder.DeserializeRawExact(v, &sig2); err != nil {
			return err
		}

		if sig1.Sig != sig2 {
			return errors.New("BlockSigsBkt sig decode mismatch")
		}

		return nil
	}); err != nil {
		return err
	}

	if err := dbutil.ForEach(tx, BlocksBkt, func(_, v []byte) error {
		select {
		case <-quit:
			return ErrVerifyStopped
		default:
		}

		var b1 coin.Block
		if err := decodeBlockExact(v, &b1); err != nil {
			return err
		}

		var b2 coin.Block
		if err := encoder.DeserializeRawExact(v, &b2); err != nil {
			return err
		}

		if !reflect.DeepEqual(b1, b2) {
			return errors.New("BlocksBkt block mismatch")
		}

		return nil
	}); err != nil {
		return err
	}

	if err := dbutil.ForEach(tx, TreeBkt, func(_, v []byte) error {
		select {
		case <-quit:
			return ErrVerifyStopped
		default:
		}

		var b1 hashPairsWrapper
		if err := decodeHashPairsWrapperExact(v, &b1); err != nil {
			return err
		}

		var b2 []coin.HashPair
		if err := encoder.DeserializeRawExact(v, &b2); err != nil {
			return err
		}

		if !reflect.DeepEqual(b1.HashPairs, b2) {
			return errors.New("TreeBkt hash pairs mismatch")
		}

		return nil
	}); err != nil {
		return err
	}

	if err := dbutil.ForEach(tx, UnspentPoolBkt, func(_, v []byte) error {
		select {
		case <-quit:
			return ErrVerifyStopped
		default:
		}

		var b1 coin.UxOut
		if err := decodeUxOutExact(v, &b1); err != nil {
			return err
		}

		var b2 coin.UxOut
		if err := encoder.DeserializeRawExact(v, &b2); err != nil {
			return err
		}

		if !reflect.DeepEqual(b1, b2) {
			return errors.New("UnspentPoolBkt ux out mismatch")
		}

		return nil
	}); err != nil {
		return err
	}

	if err := dbutil.ForEach(tx, UnspentPoolAddrIndexBkt, func(_, v []byte) error {
		select {
		case <-quit:
			return ErrVerifyStopped
		default:
		}

		var b1 hashesWrapper
		if err := decodeHashesWrapperExact(v, &b1); err != nil {
			return err
		}

		var b2 []cipher.SHA256
		if err := encoder.DeserializeRawExact(v, &b2); err != nil {
			return err
		}

		if !reflect.DeepEqual(b1.Hashes, b2) {
			return errors.New("UnspentPoolAddrIndexBkt sha256 hashes mismatch")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
