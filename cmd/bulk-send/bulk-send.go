package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cli"
	"github.com/skycoin/skycoin/src/wallet"
)

func run() error {
	csvFile := flag.String("csv", "", "csv file to load (format: skyaddress,coins). coins are in whole numbers")
	walletFile := flag.String("wallet", "", "wallet file")
	rpcAddr := flag.String("rpc-addr", "http://127.0.0.1:6420", "rpc interface address")

	flag.Parse()

	if *csvFile == "" {
		return errors.New("csv required")
	}
	if *walletFile == "" {
		return errors.New("wallet required")
	}

	wlt, err := wallet.Load(*walletFile)
	if err != nil {
		return err
	}

	if len(wlt.Entries) == 0 {
		return errors.New("Wallet is empty")
	}

	changeAddr := wlt.Entries[0].Address.String()

	f, err := os.Open(*csvFile)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	fields, err := r.ReadAll()
	if err != nil {
		return err
	}

	var sends []cli.SendAmount
	var errs []error
	for _, f := range fields {
		addr := f[0]

		addr = strings.TrimSpace(addr)

		if _, err := cipher.DecodeBase58Address(addr); err != nil {
			err = fmt.Errorf("Invalid address %s: %v", addr, err)
			errs = append(errs, err)
			continue
		}

		amt, err := strconv.ParseInt(f[1], 10, 64)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if amt <= 0 {
			err := fmt.Errorf("Invalid amount %s", f[1])
			errs = append(errs, err)
			continue
		}

		sends = append(sends, cli.SendAmount{
			Addr:  addr,
			Coins: uint64(amt * 1e6),
		})
	}

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println("ERROR:", err)
		}
		return errs[0]
	}

	c, err := webrpc.NewClient(*rpcAddr)
	if err != nil {
		return err
	}
	c.UseCSRF = true

	tx, err := cli.CreateRawTxFromWallet(c, *walletFile, changeAddr, sends, nil)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", tx)

	// txid, err := c.InjectTransaction(tx)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("txid:", txid)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
