package historydb

import (
	"errors"
	"reflect"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
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

	if err := dbutil.ForEach(tx, AddressTxnsBkt, func(_, v []byte) error {
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
			return errors.New("AddressTxnsBkt sha256 hashes mismatch")
		}

		return nil
	}); err != nil {
		return err
	}

	if err := dbutil.ForEach(tx, AddressUxBkt, func(_, v []byte) error {
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
			return errors.New("AddressUxBkt sha256 hashes mismatch")
		}

		return nil
	}); err != nil {
		return err
	}

	if err := dbutil.ForEach(tx, UxOutsBkt, func(_, v []byte) error {
		select {
		case <-quit:
			return ErrVerifyStopped
		default:
		}

		var b1 UxOut
		if err := decodeUxOutExact(v, &b1); err != nil {
			return err
		}

		var b2 UxOut
		if err := encoder.DeserializeRawExact(v, &b2); err != nil {
			return err
		}

		if !reflect.DeepEqual(b1, b2) {
			return errors.New("UxOutsBkt ux out mismatch")
		}

		return nil
	}); err != nil {
		return err
	}

	if err := dbutil.ForEach(tx, TransactionsBkt, func(_, v []byte) error {
		select {
		case <-quit:
			return ErrVerifyStopped
		default:
		}

		var b1 Transaction
		if err := decodeTransactionExact(v, &b1); err != nil {
			return err
		}

		var b2 Transaction
		if err := encoder.DeserializeRawExact(v, &b2); err != nil {
			return err
		}

		if !reflect.DeepEqual(b1, b2) {
			return errors.New("TransactionsBkt ux out mismatch")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
