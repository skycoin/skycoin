package visor

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/params"
)

// Config configuration parameters for the Visor
type Config struct {
	// Is this a block publishing node
	IsBlockPublisher bool

	// Public key of the blockchain
	BlockchainPubkey cipher.PubKey

	// Secret key of the blockchain (required if block publisher)
	BlockchainSeckey cipher.SecKey

	// Transaction verification parameters used for unconfirmed transactions
	UnconfirmedVerifyTxn params.VerifyTxn
	// Transaction verification parameters used when creating a block
	CreateBlockVerifyTxn params.VerifyTxn
	// Maximum size of a block, in bytes for creating blocks
	MaxBlockTransactionsSize uint32

	// Coin distribution parameters (necessary for txn verification)
	Distribution params.Distribution

	// Where the blockchain is saved
	BlockchainFile string
	// Where the block signatures are saved
	BlockSigsFile string

	//address for genesis
	GenesisAddress cipher.Address
	// Genesis block sig
	GenesisSignature cipher.Sig
	// Genesis block timestamp
	GenesisTimestamp uint64
	// Number of coins in genesis block
	GenesisCoinVolume uint64
	// enable arbitrating mode
	Arbitrating bool
}

// NewConfig creates Config
func NewConfig() Config {
	c := Config{
		IsBlockPublisher: false,

		BlockchainPubkey: cipher.PubKey{},
		BlockchainSeckey: cipher.SecKey{},

		UnconfirmedVerifyTxn:     params.UserVerifyTxn,
		CreateBlockVerifyTxn:     params.UserVerifyTxn,
		MaxBlockTransactionsSize: params.UserVerifyTxn.MaxTransactionSize,

		GenesisAddress:    cipher.Address{},
		GenesisSignature:  cipher.Sig{},
		GenesisTimestamp:  0,
		GenesisCoinVolume: 0, //100e12, 100e6 * 10e6
	}

	return c
}

// Verify verifies the configuration
func (c Config) Verify() error {
	if c.IsBlockPublisher {
		if c.BlockchainPubkey != cipher.MustPubKeyFromSecKey(c.BlockchainSeckey) {
			return errors.New("Cannot run as block publisher: invalid seckey for pubkey")
		}
	}

	if err := c.UnconfirmedVerifyTxn.Validate(); err != nil {
		return err
	}

	if err := c.CreateBlockVerifyTxn.Validate(); err != nil {
		return err
	}

	if c.UnconfirmedVerifyTxn.BurnFactor < params.UserVerifyTxn.BurnFactor {
		return fmt.Errorf("UnconfirmedVerifyTxn.BurnFactor must be >= params.UserVerifyTxn.BurnFactor (%d)", params.UserVerifyTxn.BurnFactor)
	}

	if c.CreateBlockVerifyTxn.BurnFactor < params.UserVerifyTxn.BurnFactor {
		return fmt.Errorf("CreateBlockVerifyTxn.BurnFactor must be >= params.UserVerifyTxn.BurnFactor (%d)", params.UserVerifyTxn.BurnFactor)
	}

	if c.UnconfirmedVerifyTxn.MaxTransactionSize < params.UserVerifyTxn.MaxTransactionSize {
		return fmt.Errorf("UnconfirmedVerifyTxn.MaxTransactionSize must be >= params.UserVerifyTxn.MaxTransactionSize (%d)", params.UserVerifyTxn.MaxTransactionSize)
	}

	if c.CreateBlockVerifyTxn.MaxTransactionSize < params.UserVerifyTxn.MaxTransactionSize {
		return fmt.Errorf("CreateBlockVerifyTxn.MaxTransactionSize must be >= params.UserVerifyTxn.MaxTransactionSize (%d)", params.UserVerifyTxn.MaxTransactionSize)
	}

	if c.UnconfirmedVerifyTxn.MaxDropletPrecision < params.UserVerifyTxn.MaxDropletPrecision {
		return fmt.Errorf("UnconfirmedVerifyTxn.MaxDropletPrecision must be >= params.UserVerifyTxn.MaxDropletPrecision (%d)", params.UserVerifyTxn.MaxDropletPrecision)
	}

	if c.CreateBlockVerifyTxn.MaxDropletPrecision < params.UserVerifyTxn.MaxDropletPrecision {
		return fmt.Errorf("CreateBlockVerifyTxn.MaxDropletPrecision must be >= params.UserVerifyTxn.MaxDropletPrecision (%d)", params.UserVerifyTxn.MaxDropletPrecision)
	}

	if c.MaxBlockTransactionsSize < c.CreateBlockVerifyTxn.MaxTransactionSize {
		return errors.New("MaxBlockTransactionsSize must be >= CreateBlockVerifyTxn.MaxTransactionSize")
	}

	if err := c.Distribution.Validate(); err != nil {
		return err
	}

	return nil
}
