package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/apputil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

const (
	blockchainPubkey = "02933015bd2fa1e0a885c05fb08eb7c647bf8c3188ed5120b51d0d09ccaf525036"
)

// wrapDB calls dbutil.WrapDB and disables all logging
func wrapDB(db *bolt.DB) *dbutil.DB {
	wdb := dbutil.WrapDB(db)
	wdb.ViewLog = false
	wdb.ViewTrace = false
	wdb.UpdateLog = false
	wdb.UpdateTrace = false
	wdb.DurationLog = false
	return wdb
}

func checkDBCmd() *cobra.Command {
	return &cobra.Command{
		Short: "Verify the database",
		Use:   "checkdb [db path]",
		Long: `Checks if the given database file contains valid skycoin blockchain data.
    If no argument is specificed, the default data.db in $HOME/.$COIN/ will be checked.`,
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  checkDB,
	}
}

func checkDB(_ *cobra.Command, args []string) error {
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

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{
		Timeout:  5 * time.Second,
		ReadOnly: true,
	})

	if err != nil {
		return fmt.Errorf("open db failed: %v", err)
	}
	pubkey, err := cipher.PubKeyFromHex(blockchainPubkey)
	if err != nil {
		return fmt.Errorf("decode blockchain pubkey failed: %v", err)
	}

	go func() {
		apputil.CatchInterrupt(quitChan)
	}()

	if err := visor.CheckDatabase(wrapDB(db), pubkey, quitChan); err != nil {
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

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{
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
