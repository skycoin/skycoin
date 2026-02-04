package cli

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/gogo/protobuf/proto"

	"github.com/spf13/cobra"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	skyWallet "github.com/skycoin/skycoin/src/hardware-wallet/skywallet"
)

func init() {
	transactionSignCmd.Flags().StringSliceVar(&inputHash, "inputHash", []string{}, "Hash of the Input of the transaction we expect the device to sign")
	transactionSignCmd.Flags().StringSliceVar(&prevHash, "prevHash", []string{}, "Hash of the previous transaction we expect the device to sign")
	transactionSignCmd.Flags().IntSliceVar(&inputIndex, "inputIndex", []int{}, "Index of the input in the wallet")
	transactionSignCmd.Flags().StringSliceVar(&outputAddress, "outputAddress", []string{}, "Addresses of the output for the transaction")
	transactionSignCmd.Flags().Int64SliceVar(&coins, "coins", []int64{}, "Amount of coins")
	transactionSignCmd.Flags().Int64SliceVar(&hours, "hours", []int64{}, "Number of hours")
	transactionSignCmd.Flags().IntSliceVar(&addressIndex, "addressIndex", []int{}, "If the address is a return address tell its index in the wallet")
	transactionSignCmd.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
	transactionSignCmd.Flags().StringVar(&coinTypeStr, "coinTypeStr", "SKY", "Coin type to use on hardware-wallet.")
}

var transactionSignCmd = &cobra.Command{
	Use:   "transactionSign",
	Short: "Ask the device to sign a transaction using the provided information.",
	RunE: func(_ *cobra.Command, _ []string) error {
		coinType, err := skyWallet.CoinTypeFromString(coinTypeStr)
		if err != nil {
			return err
		}
		if coinType != skyWallet.SkycoinCoinType && len(inputHash) > 0 {
			return fmt.Errorf("coin type %v doesn't need input hash", coinType)
		}

		if coinType != skyWallet.BitcoinCoinType && len(prevHash) > 0 {
			return fmt.Errorf("coin type %v doesn't need previous hash", coinType)
		}

		device := skyWallet.NewDevice(skyWallet.DeviceTypeFromString(deviceType))
		if device == nil {
			return fmt.Errorf("failed to create device")
		}
		defer device.Close()

		if os.Getenv("AUTO_PRESS_BUTTONS") == "1" && device.Driver.DeviceType() == skyWallet.DeviceTypeEmulator && runtime.GOOS == "linux" {
			err := device.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				return err
			}
		}

		if len(outputAddress) != len(coins) {
			return fmt.Errorf("every given output should have a coin value")
		}

		switch coinType {
		case skyWallet.SkycoinCoinType:
			err = transactionSkycoinSign(device, inputHash, outputAddress, coins, hours, inputIndex, addressIndex)
			if err != nil {
				return err
			}
		case skyWallet.BitcoinCoinType:
			err = transactionBitcoinSign(device, prevHash, outputAddress, coins, inputIndex, addressIndex)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("TransactionSign is not implemented for %v yet", coinType)
		}
		return nil
	},
}

func transactionSkycoinSign(device *skyWallet.Device, inputs, outputs []string, coins, hours []int64, inputIndex, addressIndex []int) error {
	if len(inputs) != len(inputIndex) {
		return fmt.Errorf("every given input hash should have the an inputIndex")
	}
	if len(outputs) != len(hours) {
		return fmt.Errorf("every given output should have a coin value")
	}
	var transactionInputs []*messages.TxAck_TransactionType_TxInputType
	var transactionOutputs []*messages.TxAck_TransactionType_TxOutputType

	for i, input := range inputs {
		transactionInputs = append(transactionInputs, &messages.TxAck_TransactionType_TxInputType{
			AddressN: []uint32{*proto.Uint32(uint32((inputIndex)[i]))}, //nolint:gosec // index fits in uint32
			HashIn:   proto.String(input),
		})
	}
	for i, output := range outputs {
		transactionOutputs = append(transactionOutputs, &messages.TxAck_TransactionType_TxOutputType{
			Address: proto.String(output),
			Coins:   proto.Uint64(uint64((coins)[i])), //nolint:gosec // coin values from CLI args are validated
			Hours:   proto.Uint64(uint64((hours)[i])), //nolint:gosec // hour values from CLI args are validated
		})
		if i < len(addressIndex) {
			transactionOutputs[len(transactionOutputs)-1].AddressN = []uint32{uint32((addressIndex)[i])} //nolint:gosec // index fits in uint32
		}
	}
	signer := skyWallet.SkycoinTransactionSigner{
		Inputs:   transactionInputs,
		Outputs:  transactionOutputs,
		Version:  1,
		LockTime: 0,
	}

	signatures, err := device.GeneralTransactionSign(&signer)
	if err != nil {
		return err
	}

	// Construct full transaction
	txn := coin.Transaction{
		Type: 0,
	}

	// Parse input hashes (UTXO hashes being spent)
	for _, inputHashStr := range inputs {
		hash, err := cipher.SHA256FromHex(inputHashStr)
		if err != nil {
			return fmt.Errorf("invalid input hash %s: %v", inputHashStr, err)
		}
		txn.In = append(txn.In, hash)
	}

	// Parse output addresses and create outputs
	for i, outputAddr := range outputs {
		addr, err := cipher.DecodeBase58Address(outputAddr)
		if err != nil {
			return fmt.Errorf("invalid output address %s: %v", outputAddr, err)
		}
		txn.Out = append(txn.Out, coin.TransactionOutput{
			Address: addr,
			Coins:   uint64(coins[i]), //nolint:gosec // coin values from CLI args
			Hours:   uint64(hours[i]), //nolint:gosec // hour values from CLI args
		})
	}

	// Parse signatures from device
	for _, sigHex := range signatures {
		sigBytes, err := hex.DecodeString(sigHex)
		if err != nil {
			return fmt.Errorf("invalid signature hex %s: %v", sigHex, err)
		}
		var sig cipher.Sig
		copy(sig[:], sigBytes)
		txn.Sigs = append(txn.Sigs, sig)
	}

	// Compute inner hash and update transaction
	if err := txn.UpdateHeader(); err != nil {
		return fmt.Errorf("failed to update transaction header: %v", err)
	}

	// Serialize to hex
	txnHex, err := txn.SerializeHex()
	if err != nil {
		return fmt.Errorf("failed to serialize transaction: %v", err)
	}

	fmt.Println(txnHex)
	return nil
}

func transactionBitcoinSign(device *skyWallet.Device, prevHashes, outputs []string, coins []int64, inputIndex, addressIndex []int) error {
	if len(prevHashes) != len(inputIndex) {
		return fmt.Errorf("Every given input index should have a hash of previous the tx")
	}

	var transactionInputs []*messages.BitcoinTransactionInput
	var transactionOutputs []*messages.BitcoinTransactionOutput
	for i, prevHash := range prevHashes {
		var transactionInput messages.BitcoinTransactionInput
		transactionInput.AddressN = proto.Uint32(uint32(inputIndex[i])) //nolint:gosec // index fits in uint32
		decoded, err := hex.DecodeString(prevHash)
		if err != nil {
			return err
		}
		transactionInput.PrevHash = decoded
		transactionInputs = append(transactionInputs, &transactionInput)
	}
	for i, output := range outputs {
		var transactionOutput messages.BitcoinTransactionOutput
		transactionOutput.Address = proto.String(output)
		transactionOutput.Coin = proto.Uint64(uint64(coins[i])) //nolint:gosec // coin values from CLI args
		if i < len(addressIndex) {
			transactionOutput.AddressIndex = proto.Uint32(uint32(addressIndex[i])) //nolint:gosec // index fits in uint32
		}
		transactionOutputs = append(transactionOutputs, &transactionOutput)
	}

	signer := skyWallet.BitcoinTransactionSigner{
		Inputs:   transactionInputs,
		Outputs:  transactionOutputs,
		Version:  1,
		LockTime: 0,
	}

	signatures, err := device.GeneralTransactionSign(&signer)
	if err != nil {
		return err
	}
	fmt.Println(signatures)
	return nil
}
