package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/util/droplet"
)

func distributeGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "distributeGenesis [genesis address secret key]",
		Short: "Distributes the genesis block coins into the configured distribution addresses",
		Long: `Distributes the genesis block coins into the configured distribution addresses.

    The genesis block contains a single transaction with a single output that creates all coins
    in existence. Skycoin expects the second block to be a "distribution" transaction, where
    the genesis coins are split into N distribution addresses, each holding an equal amount.

    RPC_ADDR must be set, to communicate with a running Skycoin node.`,
		RunE: distributeGenesisHandler,
	}

	cmd.Flags().StringP("genesis-seckey", "s", "", "Genesis address secret key")

	return cmd
}

func distributeGenesisHandler(c *cobra.Command, args []string) error {
	sk, err := cipher.SecKeyFromHex(args[0])
	if err != nil {
		return errors.New("invalid genesis secret key")
	}

	// Obtain the genesis uxid from the node
	uxID, err := getGenesisUxID()
	if err != nil {
		return err
	}

	txn, err := createDistributionTransaction(uxID, sk, params.MainNetDistribution)
	if err != nil {
		return err
	}

	health, err := apiClient.Health()
	if err != nil {
		return err
	}

	// If the node is a block publisher node, we can skip the broadcast and
	// avoid any broadcast related errors. Otherwise, InjectTransaction would
	// require that the block publisher node be connected to another peer.
	// Otherwise, inject the transaction normally, which requires the node
	// to be connected to peers, or else an error is returned.
	// This allows the user to use this command against a block publisher node
	// during setup, or against a local node, after setting up a block publisher
	// node elsewhere.
	if health.BlockPublisher {
		if _, err := apiClient.InjectTransactionNoBroadcast(txn); err != nil {
			return err
		}
	} else {
		if _, err := apiClient.InjectTransaction(txn); err != nil {
			return err
		}
	}

	return nil
}

func getGenesisUxID() (string, error) {
	// Check that the only block is the genesis block
	bm, err := apiClient.BlockchainMetadata()
	if err != nil {
		return "", err
	}
	if bm.Head.BkSeq != 0 {
		return "", errors.New("genesis output has already been distributed")
	}

	b, err := apiClient.BlockBySeq(0)
	if err != nil {
		return "", err
	}

	// Sanity checks
	if len(b.Body.Transactions) != 1 {
		return "", errors.New("genesis block has multiple transactions")
	}
	if len(b.Body.Transactions[0].Out) != 1 {
		return "", errors.New("genesis block has multiple outputs")
	}

	return b.Body.Transactions[0].Out[0].Hash, nil
}

func createDistributionTransaction(uxID string, genesisSecKey cipher.SecKey, p params.Distribution) (*coin.Transaction, error) {
	// Sanity check
	addrs := p.AddressesDecoded()
	if p.MaxCoinSupply%uint64(len(addrs)) != 0 {
		return nil, errors.New("the number of distribution addresses must divide MaxCoinSupply exactly")
	}

	var txn coin.Transaction

	output, err := cipher.SHA256FromHex(uxID)
	if err != nil {
		return nil, fmt.Errorf("invalid uxid")
	}

	if err := txn.PushInput(output); err != nil {
		return nil, err
	}

	coins := (p.MaxCoinSupply / uint64(len(addrs))) * droplet.Multiplier
	hours := uint64(1000)

	for _, addr := range addrs {
		if err := txn.PushOutput(addr, coins, hours); err != nil {
			return nil, err
		}
	}

	seckeys := make([]cipher.SecKey, 1)
	seckey := genesisSecKey.Hex()
	seckeys[0] = cipher.MustSecKeyFromHex(seckey)
	txn.SignInputs([]cipher.SecKey{
		genesisSecKey,
	})

	if err := txn.UpdateHeader(); err != nil {
		return nil, err
	}

	if err := txn.Verify(); err != nil {
		return nil, err
	}

	return &txn, nil
}
