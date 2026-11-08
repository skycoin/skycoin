package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/skycoin/skycoin/src/util/droplet"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"

	"github.com/spf13/cobra"
)

// TxnResult wraps readable.TransactionWithStatus
type TxnResult struct {
	Transaction *readable.TransactionWithStatus `json:"transaction"`
}

func transactionCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Show detail info of specific transaction",
		Use:                   "transaction [transaction id]",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			txid := args[0]
			if txid == "" {
				return errors.New("txid is empty")
			}

			// validate the txid
			_, err := cipher.SHA256FromHex(txid)
			if err != nil {
				return errors.New("invalid txid")
			}

			txn, err := apiClient.Transaction(txid)
			if err != nil {
				return err
			}

			return printJSON(TxnResult{
				Transaction: txn,
			})
		},
	}
}

func decodeRawTxnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Decode raw transaction",
		Use:   "decodeRawTransaction [raw transaction]",
		Long: `Decode a raw transaction to JSON.

    With --review it prints each input with its owner, coins, hours and whether it
    is signed, each output, the fee, and whether the transaction is fully signed.
    The inputs are looked up on the node. Add --wallet to mark the wallet's own
    inputs and outputs.`,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			txn, err := coin.DeserializeTransactionHexWithPQ(args[0])
			if err != nil {
				return fmt.Errorf("invalid raw transaction: %v", err)
			}

			review, err := c.Flags().GetBool("review")
			if err != nil {
				return err
			}
			wltID, err := c.Flags().GetString("wallet")
			if err != nil {
				return err
			}
			if review || wltID != "" {
				var mine map[string]struct{}
				if wltID != "" {
					w, err := apiClient.Wallet(wltID)
					if err != nil {
						return err
					}
					mine = walletAddressSet(w)
				}
				inputs, err := lookupInputs(apiClient, txnInputHashes(&txn))
				if err != nil {
					return err
				}
				fmt.Print(txnReview(&txn, inputs, mine))
				return nil
			}

			// Assume the transaction is not malformed and if it has no inputs
			// that it is the genesis block's transaction
			isGenesis := len(txn.In) == 0
			rTxn, err := readable.NewTransaction(txn, isGenesis)
			if err != nil {
				return err
			}

			return printJSON(rTxn)
		},
	}
	cmd.Flags().Bool("review", false, "Print a readable review, looking up the inputs on the node")
	cmd.Flags().StringP("wallet", "w", "", "Wallet whose inputs and outputs to mark in the review. Implies --review")
	return cmd
}

func encodeJSONTxnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short:                 "Encode JSON transaction",
		Use:                   "encodeJsonTransaction [file path or -]",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			jsonOutput, err := c.Flags().GetBool("json")
			if err != nil {
				return err
			}

			jsonFilePath := args[0]
			var jsonFile *os.File
			if jsonFilePath == "-" {
				jsonFile = os.Stdin
				err = nil
				jsonFilePath = "<stdin>"
			} else {
				jsonFile, err = os.Open(jsonFilePath) //nolint:gosec
			}
			if err != nil {
				return fmt.Errorf("open file failed %s: %v", jsonFilePath, err)
			}
			var rTxn readable.Transaction
			err = json.NewDecoder(jsonFile).Decode(&rTxn)
			if err != nil {
				return fmt.Errorf("invalid JSON transaction: %v", err)
			}

			// fixHashes recomputes the txn hashes
			// the recomputation of hash is needed when txn info for e.g. the sig is changed
			fixHashes, err := c.Flags().GetBool("fix")
			if err != nil {
				return err
			}

			if fixHashes {
				err = recomputeHashes(&rTxn)
				if err != nil {
					return err
				}
			}

			txn, err := readableToCreatedTransaction(&rTxn).ToTransaction()
			if err != nil {
				return err
			}
			rawTxn, err := txn.SerializeHexWithPQ()
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(struct {
					RawTx string `json:"rawtx"`
				}{
					RawTx: rawTxn,
				})
			}
			fmt.Println(rawTxn)
			return nil
		},
	}
	cmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format.")
	cmd.Flags().BoolP("fix", "f", false, "Recompute transaction inner and outer hashes")
	return cmd
}

func readableToCreatedTransaction(rTxn *readable.Transaction) *api.CreatedTransaction {
	inputs := make([]api.CreatedTransactionInput, len(rTxn.In))
	outputs := make([]api.CreatedTransactionOutput, len(rTxn.Out))
	for i, rIn := range rTxn.In {
		inputs[i] = api.CreatedTransactionInput{
			UxID: rIn,
		}
	}
	for i, rOut := range rTxn.Out {
		outputs[i] = api.CreatedTransactionOutput{
			UxID:    rOut.Hash,
			Address: rOut.Address,
			Coins:   rOut.Coins,
			Hours:   strconv.FormatUint(rOut.Hours, 10),
		}
	}
	cTxn := api.CreatedTransaction{
		Length:    rTxn.Length,
		Type:      rTxn.Type,
		TxID:      rTxn.Hash,
		InnerHash: rTxn.InnerHash,
		Fee:       "",
		Sigs:      rTxn.Sigs[:],
		In:        inputs,
		Out:       outputs,
		PQ:        rTxn.PQ,
	}
	return &cTxn
}

func recomputeHashes(rTxn *readable.Transaction) error {
	t := coin.Transaction{}

	t.Length = rTxn.Length
	t.Type = rTxn.Type

	var err error

	sigs := make([]cipher.Sig, len(rTxn.Sigs))
	for i, s := range rTxn.Sigs {
		sigs[i], err = cipher.SigFromHex(s)
		if err != nil {
			return err
		}
	}

	t.Sigs = sigs

	in := make([]cipher.SHA256, len(rTxn.In))
	for i, UxID := range rTxn.In {
		in[i], err = cipher.SHA256FromHex(UxID)
		if err != nil {
			return err
		}
	}

	t.In = in

	out := make([]coin.TransactionOutput, len(rTxn.Out))
	for i, o := range rTxn.Out {
		addr, err := cipher.DecodeBase58Address(o.Address)
		if err != nil {
			return err
		}

		coins, err := droplet.FromString(o.Coins)
		if err != nil {
			return err
		}

		out[i] = coin.TransactionOutput{
			Address: addr,
			Coins:   coins,
			Hours:   o.Hours,
		}
	}

	t.Out = out

	// recompute inner hash
	rTxn.InnerHash = t.HashInner().Hex()
	t.InnerHash, err = cipher.SHA256FromHex(rTxn.InnerHash)
	if err != nil {
		return err
	}

	// recompute txid
	rTxn.Hash = t.Hash().Hex()

	return nil
}

func addressTransactionsCmd() *cobra.Command {
	return &cobra.Command{
		Short: "Show detail for transaction associated with one or more specified addresses",
		Use:   "addressTransactions [address list]",
		Long: `Display transactions for specific addresses, separate multiple addresses with a space,
        example: addressTransactions addr1 addr2 addr3`,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getAddressTransactionsCmd,
	}
}

func getAddressTransactionsCmd(c *cobra.Command, args []string) error {
	// Build the list of addresses from the command line arguments
	addrs := make([]string, len(args))
	var err error
	for i := 0; i < len(args); i++ {
		addrs[i] = args[i]
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
	}

	// If one or more addresses have been provided, request their transactions - otherwise report an error
	if len(addrs) > 0 {
		outputs, err := apiClient.TransactionsVerbose(addrs)
		if err != nil {
			return err
		}

		return printJSON(outputs)
	}

	return fmt.Errorf("at least one address must be specified. Example: %s addr1 addr2 addr3", c.Name())
}

func verifyTransactionCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Verify if the specific transaction is spendable",
		Use:                   "verifyTransaction [encoded transaction]",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			encodedTxn := args[0]
			if encodedTxn == "" {
				return errors.New("transaction is empty")
			}

			_, err := apiClient.VerifyTransaction(api.VerifyTransactionRequest{
				EncodedTransaction: encodedTxn,
			})
			if err != nil {
				return err
			}

			fmt.Println("transaction is spendable")

			return nil
		},
	}
}

func pendingTransactionsCmd() *cobra.Command {
	pendingTxnsCmd := &cobra.Command{
		Short:                 "Get all unconfirmed transactions",
		Use:                   "pendingTransactions",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			isVerbose, err := c.Flags().GetBool("verbose")
			if err != nil {
				return err
			}

			if isVerbose {
				pendingTxns, err := apiClient.PendingTransactionsVerbose()
				if err != nil {
					return err
				}

				return printJSON(pendingTxns)
			}

			pendingTxns, err := apiClient.PendingTransactions()
			if err != nil {
				return err
			}

			return printJSON(pendingTxns)
		},
	}

	pendingTxnsCmd.Flags().BoolP("verbose", "v", false,
		`Require the transaction inputs to include the owner address, coins, hours and calculated hours.
	The hours are the original hours the output was created with.
	The calculated hours are calculated based upon the current system time, and provide an approximate
	coin hour value of the output if it were to be confirmed at that instant.`)

	return pendingTxnsCmd
}

func signTxnCmd() *cobra.Command {
	signTxnCmd := &cobra.Command{
		Short: "Sign an unsigned transaction with specific wallet",
		Use:   "signTransaction [wallet] [raw transaction]",
		Long: `Sign the inputs of a transaction that the wallet owns.

    A joint transaction spends inputs of several wallets. Each wallet signs only
    its own inputs, found by looking up each input's owner on the node, and the
    result goes to the next signer. Use --sign-indexes to choose the inputs by hand.`,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			rawTxn := args[1]

			signIndexes, err := c.Flags().GetIntSlice("sign-indexes")
			if err != nil {
				return err
			}

			txn, err := coin.DeserializeTransactionHexWithPQ(rawTxn)
			if err != nil {
				return err
			}
			if txn.IsFullySigned() {
				return fmt.Errorf("Transaction already signed")
			}

			// id is the node's wallet ID, which is what the signing request sends too
			w, err := apiClient.Wallet(id)
			if err != nil {
				return err
			}

			if len(signIndexes) == 0 {
				signIndexes, err = walletSignIndexes(apiClient, &txn, w)
				if err != nil {
					return err
				}
			}

			req := api.WalletSignTransactionRequest{
				WalletID:           id,
				EncodedTransaction: rawTxn,
				SignIndexes:        signIndexes,
			}

			if w.Meta.Encrypted {
				v, err := readPasswordFromTerminal()
				if err != nil {
					return err
				}
				req.Password = string(v)
				defer func() {
					// Wipe out the password from memory
					v = []byte{}
					req.Password = ""
				}()
			}

			signedTxn, err := apiClient.WalletSignTransaction(req)
			if err != nil {
				return err
			}

			signed, err := coin.DeserializeTransactionHexWithPQ(signedTxn.EncodedTransaction)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "signed %d of %d inputs\n", countSigned(&signed), len(signed.In))

			return printJSON(signedTxn)
		},
	}

	signTxnCmd.Flags().IntSlice("sign-indexes", nil, "Input indexes to sign. Defaults to the unsigned inputs the wallet owns")
	return signTxnCmd
}

// walletSignIndexes returns the unsigned inputs the wallet owns, or nil when
// those are all the unsigned inputs, so the node signs them all as before
func walletSignIndexes(c outputsLookup, txn *coin.Transaction, w *api.WalletResponse) ([]int, error) {
	inputs, err := lookupInputs(c, txnInputHashes(txn))
	if err != nil {
		return nil, err
	}
	idx := ownedUnsignedIndexes(txn, inputs, walletAddressSet(w))
	if len(idx) == 0 {
		return nil, errors.New("the wallet owns none of the unsigned inputs")
	}
	if len(idx) == len(txn.In)-countSigned(txn) {
		return nil, nil
	}
	return idx, nil
}
