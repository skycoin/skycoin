package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	bbolt "github.com/0magnet/bbolt"
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/fiber"
	"github.com/skycoin/skycoin/src/util/apputil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

const (
	blockchainPubkey = "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a"
)

// wrapDB calls dbutil.WrapDB and disables all logging
func wrapDB(db *bbolt.DB) *dbutil.DB {
	wdb := dbutil.WrapDB(db)
	wdb.ViewLog = false
	wdb.ViewTrace = false
	wdb.UpdateLog = false
	wdb.UpdateTrace = false
	wdb.DurationLog = false
	return wdb
}

func checkDBCmd() *cobra.Command {
	c := &cobra.Command{
		Short: "Verify the database",
		Use:   "checkdb [db path]",
		Long: `Checks if the given database file contains valid skycoin blockchain data.
    If no argument is specificed, the default data.db in $HOME/.$COIN/ will be checked.
    With FIBER_TOML set, the block publisher keys come from that coin's config.`,
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  checkDB,
	}
	c.Flags().String("pq-pubkey", "", "post-quantum (ML-DSA-44) block publisher public key, in hex")
	c.Flags().Uint64("pq-fork-seq", 0, "first block that needs a post-quantum signature")
	return c
}

func checkDB(c *cobra.Command, args []string) error {
	// get db path
	dbPath := ""
	if len(args) > 0 {
		dbPath = args[0]
	}
	dbPath, err := resolveDBPath(cliConfig, dbPath)
	if err != nil {
		return err
	}

	// check if this file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("db file: %v does not exist", dbPath)
	}

	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{
		Timeout:  5 * time.Second,
		ReadOnly: true,
	})

	if err != nil {
		return fmt.Errorf("open db failed: %v", err)
	}
	node, err := checkDBNodeConfig()
	if err != nil {
		return err
	}
	pubkey, err := cipher.PubKeyFromHex(node.BlockchainPubkeyStr)
	if err != nil {
		return fmt.Errorf("decode blockchain pubkey failed: %v", err)
	}

	var pq coin.PQBlockRule
	pqPubkey, err := c.Flags().GetString("pq-pubkey")
	if err != nil {
		return err
	}
	pqForkSeq, err := c.Flags().GetUint64("pq-fork-seq")
	if err != nil {
		return err
	}
	if pqPubkey == "" {
		pqPubkey = node.BlockchainPQPubkeyStr
	}
	if !c.Flags().Changed("pq-fork-seq") {
		pqForkSeq = node.PQForkSeq
	}
	if pqPubkey != "" {
		if pq.Pubkey, err = pqsig.PubKeyFromHex(pqPubkey); err != nil {
			return fmt.Errorf("decode post-quantum pubkey failed: %v", err)
		}
		pq.ForkSeq = pqForkSeq
	}

	go func() {
		apputil.CatchInterrupt(quitChan)
	}()

	if err := visor.CheckDatabase(wrapDB(db), pubkey, pq, quitChan); err != nil {
		if err == visor.ErrVerifyStopped {
			return nil
		}
		return fmt.Errorf("checkdb failed: %v", err)
	}

	fmt.Println("check db success")
	return nil
}

func checkDBEncodingCmd() *cobra.Command {
	return &cobra.Command{
		Short: "Verify the database data encoding",
		Use:   "checkDBDecoding [db path]",
		Long: `Verify the generated binary encoders match the dynamic encoders for database data.
    If no argument is specificed, the default data.db in $HOME/.$COIN/ will be checked.`,
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  checkDBDecoding,
	}
}

func checkDBDecoding(_ *cobra.Command, args []string) error {
	// get db path
	dbPath := ""
	if len(args) > 0 {
		dbPath = args[0]
	}
	dbPath, err := resolveDBPath(cliConfig, dbPath)
	if err != nil {
		return err
	}

	// check if this file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("db file: %v does not exist", dbPath)
	}

	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{
		Timeout:  5 * time.Second,
		ReadOnly: true,
	})

	if err != nil {
		return fmt.Errorf("open db failed: %v", err)
	}

	go func() {
		apputil.CatchInterrupt(quitChan)
	}()

	if err := visor.VerifyDBSkyencoderSafe(wrapDB(db), quitChan); err != nil {
		if err == visor.ErrVerifyStopped {
			return nil
		}
		return fmt.Errorf("checkDBDecoding failed: %v", err)
	}

	fmt.Println("check db decoding success")
	return nil

}

// checkDBNodeConfig returns the publisher keys of the coin FIBER_TOML names,
// or skycoin's when it is not set.
func checkDBNodeConfig() (fiber.NodeConfig, error) {
	path := os.Getenv("FIBER_TOML")
	if path == "" {
		return fiber.NodeConfig{BlockchainPubkeyStr: blockchainPubkey}, nil
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return fiber.NodeConfig{}, err
	}
	cfg, err := fiber.NewConfig(filepath.Base(path), filepath.Dir(path))
	if err != nil {
		return fiber.NodeConfig{}, fmt.Errorf("load %s failed: %v", path, err)
	}
	return cfg.Node, nil
}
