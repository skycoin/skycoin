package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
	gcli "github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/apputil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

const (
	blockchainPubkey = "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a"
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

func checkdbCmd() *gcli.Command {
	return &gcli.Command{
		Short: "Verify the database",
		Use:   "checkdb [db path]",
		Long:  `Checks if the given database file contains valid skycoin blockchain data.
    If no argument is specificed, the default data.db in $HOME/.$COIN/ will be checked.`,
		Args:  gcli.MaximumNArgs(1),
        DisableFlagsInUseLine: true,
        SilenceUsage: true,
		RunE:  checkdb,
	}
}

func checkdb(c *gcli.Command, args []string) error {
	// get db path
	dbpath, err := resolveDBPath(cliConfig, args[0])
	if err != nil {
		return err
	}

	// check if this file is exist
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		return fmt.Errorf("db file: %v does not exist", dbpath)
	}

	db, err := bolt.Open(dbpath, 0600, &bolt.Options{
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
