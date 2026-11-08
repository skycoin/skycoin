package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/fee"
)

// jointInput is what the node knows about one input of a transaction
type jointInput struct {
	UxID    string
	Address string
	Coins   uint64
	Hours   uint64
	Known   bool
}

// outputsLookup is the part of the API client used to find transaction inputs
type outputsLookup interface {
	OutputsForHashes(hashes []string) (*readable.UnspentOutputsSummary, error)
}

func createJointTxnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Create an unsigned transaction that spends outputs of several owners",
		Use:   "createJointTransaction",
		Long: `Create an unsigned transaction from the given unspent outputs, which may belong to
    different wallets, and the given outputs with coins and coin hours set by hand.

    Each owner then signs only their own inputs with signTransaction, and anyone
    broadcasts it once every input is signed. Every given input is spent.

    The fee is whatever input hours are not given to an output. It must be at least
    one tenth of the input hours, rounded up.

    Example, a customer pays a vendor 1000 coin hours:
      createJointTransaction --in <customer uxid> --in <vendor uxid> \
        --out <customer addr>,0.001,0 --out <vendor addr>,0.001,1000`,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, _ []string) error {
			ins, err := c.Flags().GetStringSlice("in")
			if err != nil {
				return err
			}
			outs, err := c.Flags().GetStringArray("out")
			if err != nil {
				return err
			}
			csvFile, err := c.Flags().GetString("csv")
			if err != nil {
				return err
			}
			changeAddr, err := c.Flags().GetString("change-address")
			if err != nil {
				return err
			}
			jsonOutput, err := c.Flags().GetBool("json")
			if err != nil {
				return err
			}

			to, err := parseJointOutputs(outs)
			if err != nil {
				return err
			}
			if csvFile != "" {
				fields, err := openCSV(csvFile)
				if err != nil {
					return err
				}
				rows := make([]string, len(fields))
				for i, f := range fields {
					rows[i] = strings.Join(f, ",")
				}
				more, err := parseJointOutputs(rows)
				if err != nil {
					return err
				}
				to = append(to, more...)
			}

			req, err := makeJointTxnRequest(apiClient, ins, to, changeAddr)
			if err != nil {
				return err
			}

			rsp, err := apiClient.CreateTransaction(*req)
			if err != nil {
				return err
			}
			if err := checkAllInputsSpent(ins, rsp.Transaction.In); err != nil {
				return err
			}

			if jsonOutput {
				return printJSON(rsp)
			}
			fmt.Fprintf(os.Stderr, "fee: %s coin hours%s, review with: decodeRawTransaction --review\n", rsp.Transaction.Fee, feeNote(rsp.Transaction))
			fmt.Println(rsp.EncodedTransaction)
			return nil
		},
	}

	cmd.Flags().StringSlice("in", nil, "Unspent output hash to spend, yours or the counterparty's. Repeat or comma separate")
	cmd.Flags().StringArray("out", nil, `Output as "address,coins,hours". Repeat for each output`)
	cmd.Flags().String("csv", "", "CSV file of outputs, one address,coins,hours per row")
	cmd.Flags().StringP("change-address", "c", "", "Address for leftover coins. Required when the outputs do not use all input coins")
	cmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format.")
	return cmd
}

// parseJointOutputs parses "address,coins,hours" output specs
func parseJointOutputs(specs []string) ([]api.Receiver, error) {
	to := make([]api.Receiver, 0, len(specs))
	for i, s := range specs {
		f := strings.Split(s, ",")
		if len(f) != 3 {
			return nil, fmt.Errorf("output %d %q: want address,coins,hours", i, s)
		}
		addr := strings.TrimSpace(f[0])
		coins := strings.TrimSpace(f[1])
		hours := strings.TrimSpace(f[2])
		if _, err := cipher.DecodeBase58Address(addr); err != nil {
			return nil, fmt.Errorf("output %d: invalid address %s: %v", i, addr, err)
		}
		c, err := droplet.FromString(coins)
		if err != nil {
			return nil, fmt.Errorf("output %d: invalid coins %s: %v", i, coins, err)
		}
		if c == 0 {
			return nil, fmt.Errorf("output %d: coins must not be zero", i)
		}
		if c%params.UserVerifyTxn.MaxDropletDivisor() != 0 {
			return nil, fmt.Errorf("output %d: coins %s has too many decimal places", i, coins)
		}
		if _, err := strconv.ParseUint(hours, 10, 64); err != nil {
			return nil, fmt.Errorf("output %d: invalid hours %s: %v", i, hours, err)
		}
		to = append(to, api.Receiver{Address: addr, Coins: coins, Hours: hours})
	}
	return to, nil
}

func makeJointTxnRequest(c outputsLookup, ins []string, to []api.Receiver, changeAddr string) (*api.CreateTransactionRequest, error) {
	if len(ins) == 0 {
		return nil, errors.New("at least one --in is required")
	}
	if len(to) == 0 {
		return nil, errors.New("at least one --out is required")
	}
	if changeAddr != "" {
		if _, err := cipher.DecodeBase58Address(changeAddr); err != nil {
			return nil, fmt.Errorf("invalid change address %s: %v", changeAddr, err)
		}
	}

	inputs, err := lookupInputs(c, ins)
	if err != nil {
		return nil, err
	}
	var inCoins, outCoins uint64
	for _, in := range inputs {
		if !in.Known {
			return nil, fmt.Errorf("unspent output %s does not exist", in.UxID)
		}
		inCoins += in.Coins
	}
	for _, r := range to {
		coins, err := droplet.FromString(r.Coins)
		if err != nil {
			return nil, err
		}
		outCoins += coins
	}
	if outCoins > inCoins {
		return nil, fmt.Errorf("outputs spend %s coins but the inputs hold %s", dropletStr(outCoins), dropletStr(inCoins))
	}
	// The node would pick a change address among the input owners, which may be the counterparty
	if outCoins < inCoins && changeAddr == "" {
		return nil, fmt.Errorf("outputs leave %s coins unspent, set --change-address or add them to an output", dropletStr(inCoins-outCoins))
	}

	req := &api.CreateTransactionRequest{
		HoursSelection: api.HoursSelection{Type: transaction.HoursSelectionTypeManual},
		UxOuts:         ins,
		To:             to,
	}
	if changeAddr != "" {
		req.ChangeAddress = &changeAddr
	}
	return req, nil
}

// checkAllInputsSpent fails if the node left out any requested input, which it
// does when fewer outputs are enough to cover the spend
func checkAllInputsSpent(want []string, got []api.CreatedTransactionInput) error {
	have := make(map[string]struct{}, len(got))
	for _, in := range got {
		have[in.UxID] = struct{}{}
	}
	var missing []string
	for _, w := range want {
		if _, ok := have[w]; !ok {
			missing = append(missing, w)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("the node did not need these inputs and left them out: %s. Raise the output coins or hours so every input is needed", strings.Join(missing, ", "))
	}
	return nil
}

// feeNote says how far the fee is above the minimum, since extra hours are burned
func feeNote(txn api.CreatedTransaction) string {
	var in uint64
	for _, i := range txn.In {
		h, err := strconv.ParseUint(i.CalculatedHours, 10, 64)
		if err != nil {
			return ""
		}
		in += h
	}
	paid, err := strconv.ParseUint(txn.Fee, 10, 64)
	if err != nil {
		return ""
	}
	minFee := fee.RequiredFee(in, params.UserVerifyTxn.BurnFactor)
	if paid <= minFee {
		return ""
	}
	return fmt.Sprintf(" (minimum %d, give the other %d to an output or they are burned)", minFee, paid-minFee)
}

// lookupInputs returns the owner, coins and current hours of each uxid, in order
func lookupInputs(c outputsLookup, uxids []string) ([]jointInput, error) {
	for _, h := range uxids {
		if _, err := cipher.SHA256FromHex(h); err != nil {
			return nil, fmt.Errorf("invalid unspent output hash %s: %v", h, err)
		}
	}
	summary, err := c.OutputsForHashes(uxids)
	if err != nil {
		return nil, err
	}
	found := make(map[string]readable.UnspentOutput)
	for _, outs := range []readable.UnspentOutputs{summary.HeadOutputs, summary.OutgoingOutputs} {
		for _, o := range outs {
			found[o.Hash] = o
		}
	}
	inputs := make([]jointInput, len(uxids))
	for i, h := range uxids {
		inputs[i].UxID = h
		o, ok := found[h]
		if !ok {
			continue
		}
		coins, err := droplet.FromString(o.Coins)
		if err != nil {
			return nil, err
		}
		inputs[i] = jointInput{UxID: h, Address: o.Address, Coins: coins, Hours: o.CalculatedHours, Known: true}
	}
	return inputs, nil
}

func txnInputHashes(txn *coin.Transaction) []string {
	h := make([]string, len(txn.In))
	for i, in := range txn.In {
		h[i] = in.Hex()
	}
	return h
}

// walletAddressSet returns every address of a wallet, across all bip44 accounts
func walletAddressSet(w *api.WalletResponse) map[string]struct{} {
	addrs := make(map[string]struct{}, len(w.Entries))
	for _, e := range w.Entries {
		addrs[e.Address] = struct{}{}
	}
	for _, a := range w.Accounts {
		for _, e := range a.ExternalEntries {
			addrs[e.Address] = struct{}{}
		}
		for _, e := range a.ChangeEntries {
			addrs[e.Address] = struct{}{}
		}
	}
	return addrs
}

// ownedUnsignedIndexes returns the indexes of unsigned inputs owned by addrs
func ownedUnsignedIndexes(txn *coin.Transaction, inputs []jointInput, addrs map[string]struct{}) []int {
	var idx []int
	for i, in := range inputs {
		if txn.InputSigned(i) {
			continue
		}
		if _, ok := addrs[in.Address]; ok && in.Known {
			idx = append(idx, i)
		}
	}
	return idx
}

func countSigned(txn *coin.Transaction) int {
	n := 0
	for i := range txn.Sigs {
		if txn.InputSigned(i) {
			n++
		}
	}
	return n
}

func dropletStr(d uint64) string {
	s, err := droplet.ToString(d)
	if err != nil {
		return strconv.FormatUint(d, 10) + " droplets"
	}
	return s
}

// txnReview formats a transaction for a person to check before signing
func txnReview(txn *coin.Transaction, inputs []jointInput, mine map[string]struct{}) string {
	var b strings.Builder
	w := &b
	var inHours, outHours uint64
	allKnown := true
	fmt.Fprintf(w, "Transaction %s\n", txn.Hash().Hex())
	fmt.Fprintf(w, "Inputs (%d):\n", len(txn.In))
	for i, in := range inputs {
		signed := "unsigned"
		if txn.InputSigned(i) {
			signed = "signed"
		}
		if !in.Known {
			allKnown = false
			fmt.Fprintf(w, "  [%d] %s  not an unspent output on this node  %s\n", i, in.UxID, signed)
			continue
		}
		inHours += in.Hours
		owner := ""
		if mine != nil {
			owner = "  other"
			if _, ok := mine[in.Address]; ok {
				owner = "  YOURS"
			}
		}
		fmt.Fprintf(w, "  [%d] %s  %s coins  %d hours  %s%s\n", i, in.Address, dropletStr(in.Coins), in.Hours, signed, owner)
	}
	fmt.Fprintf(w, "Outputs (%d):\n", len(txn.Out))
	for i, o := range txn.Out {
		outHours += o.Hours
		owner := ""
		if mine != nil {
			owner = "  other"
			if _, ok := mine[o.Address.String()]; ok {
				owner = "  YOURS"
			}
		}
		fmt.Fprintf(w, "  [%d] %s  %s coins  %d hours%s\n", i, o.Address, dropletStr(o.Coins), o.Hours, owner)
	}
	if allKnown {
		if outHours > inHours {
			fmt.Fprintf(w, "Fee: outputs give %d hours but the inputs hold %d\n", outHours, inHours)
		} else {
			minFee := fee.RequiredFee(inHours, params.UserVerifyTxn.BurnFactor)
			fmt.Fprintf(w, "Fee: %d hours burned (minimum %d)\n", inHours-outHours, minFee)
		}
	} else {
		fmt.Fprintln(w, "Fee: unknown, some inputs are not unspent outputs on this node")
	}
	fmt.Fprintf(w, "Signed: %d of %d inputs", countSigned(txn), len(txn.In))
	if txn.IsFullySigned() {
		fmt.Fprintln(w, ", ready to broadcast")
	} else {
		fmt.Fprintln(w)
	}
	return b.String()
}
