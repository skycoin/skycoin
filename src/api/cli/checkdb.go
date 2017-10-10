package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor"
	gcli "github.com/urfave/cli"
)

const (
	genesisPubkey = "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a"
)

func checkdbCmd() gcli.Command {
	name := "checkdb"
	return gcli.Command{
		Name:         name,
		Usage:        "Verify the database",
		ArgsUsage:    "[db path]",
		Description:  "If no argument is specificed, the default data.db in $HOME/.$COIN/ will be checked.",
		OnUsageError: onCommandUsageError(name),
		Action:       checkdb,
	}
}

func checkdb(c *gcli.Context) error {
	cfg := ConfigFromContext(c)

	// get db path
	dbpath, err := resolveDBPath(cfg, c.Args().First())
	if err != nil {
		return err
	}

	// check if this file is exist
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		return fmt.Errorf("db file: %v does not exist", dbpath)
	}

	db, err := bolt.Open(dbpath, 0600, &bolt.Options{
		Timeout: 5 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("open db failed: %v", err)
	}
	pubkey, err := cipher.PubKeyFromHex(genesisPubkey)
	if err != nil {
		return fmt.Errorf("decode genesis pubkey failed: %v", err)
	}

	if err := IntegrityCheck(db, pubkey); err != nil {
		return fmt.Errorf("checkdb failed: %v", err)
	}

	fmt.Println("check db success")
	return nil
}

func IntegrityCheck(db *bolt.DB, genesisPubkey cipher.PubKey) error {
	_, err := visor.NewBlockchain(db, genesisPubkey, visor.Arbitrating(true))
	return err
}
