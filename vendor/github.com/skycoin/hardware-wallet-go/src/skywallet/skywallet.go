package skywallet

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/usb"

	"github.com/skycoin/skycoin/src/util/logging"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"
)

var (
	log = logging.MustGetLogger("skywallet")
)

const (
	entropyBufferSize int = 32
)

// ButtonType is emulator button press simulation type
type ButtonType int32

const (
	// ButtonLeft press left button
	ButtonLeft ButtonType = iota

	// ButtonRight press right button
	ButtonRight
	// ButtonBoth press both buttons
	ButtonBoth
)

var (
	// ErrAddressNZero is returned if addressN is 0
	ErrAddressNZero = errors.New("addresses to generate should be greater than 0")
	// ErrRemovePinNil is returned if removePin is nil
	ErrRemovePinNil = errors.New("removePin cannot be nil")
	// ErrDeviceTypeEmulator is returned if device type is emulator
	ErrDeviceTypeEmulator = errors.New("device type cannot be emulator")
	// ErrInvalidWordCount is returned if word count is not valid mnemonic word length
	ErrInvalidWordCount = errors.New("word count must be 12 or 24")
	// ErrNoDeviceConnected is returned if no device is connected to the system
	ErrNoDeviceConnected = errors.New("no device connected")
)

//go:generate mockery -name Devicer -case underscore -inpkg -testonly

// Devicer provides api for the hw wallet functions
type Devicer interface {
	AddressGen(addressN, startIndex uint32, confirmAddress bool, coinType CoinType) (wire.Message, error)
	ApplySettings(usePassphrase *bool, label string, language string) (wire.Message, error)
	Backup() (wire.Message, error)
	Cancel() (wire.Message, error)
	CheckMessageSignature(message, signature, address string) (wire.Message, error)
	ChangePin(removePin *bool) (wire.Message, error)
	Connected() bool
	Available() bool
	FirmwareUpload(payload []byte, hash [32]byte) error
	GetFeatures() (wire.Message, error)
	GenerateMnemonic(wordCount uint32, usePassphrase bool) (wire.Message, error)
	Recovery(wordCount uint32, usePassphrase *bool, dryRun bool) (wire.Message, error)
	SetMnemonic(mnemonic string) (wire.Message, error)
	TransactionSign(inputs []*messages.SkycoinTransactionInput, outputs []*messages.SkycoinTransactionOutput) (wire.Message, error)
	GeneralTransactionSign(signer TransactionSigner) ([]string, error)
	SignMessage(addressIndex int, message string) (wire.Message, error)
	Wipe() (wire.Message, error)
	PinMatrixAck(p string) (wire.Message, error)
	WordAck(word string) (wire.Message, error)
	PassphraseAck(passphrase string) (wire.Message, error)
	ButtonAck() (wire.Message, error)
	SetAutoPressButton(simulateButtonPress bool, simulateButtonType ButtonType) error
	Close()
	Connect() error
	Disconnect() error
}

// Device provides hardware wallet functions
type Device struct {
	Driver DeviceDriver

	// mutex to force connect requests to be sequential
	sync.Mutex
	dev                 usb.Device
	connected           bool
	simulateButtonPress bool
	simulateButtonType  ButtonType
}

// DeviceTypeFromString returns device type from string
func DeviceTypeFromString(deviceType string) DeviceType {
	var dtRet DeviceType
	switch deviceType {
	case DeviceTypeUSB.String():
		dtRet = DeviceTypeUSB
	case DeviceTypeEmulator.String():
		dtRet = DeviceTypeEmulator
	default:
		log.Errorf("device type not set, valid options are %s or %s",
			DeviceTypeUSB,
			DeviceTypeEmulator)
		dtRet = DeviceTypeInvalid
	}
	return dtRet
}

var devSingleCreator sync.Once
var devSingleInstance *Device

func newDevice(deviceType DeviceType) *Device {
	driver, err := NewDriver(deviceType)
	if err != nil {
		log.Fatalf("failed to create driver: %s", err)
	}

	return &Device{
		driver,
		sync.Mutex{},
		nil,
		false,
		false,
		ButtonType(-1),
	}
}

// NewDevice returns a new device instance
func NewDevice(deviceType DeviceType) *Device {
	// TODO rename NewDevice to DeviceInstance as this is a singleton
	// implementation.
	devSingleCreator.Do(func() {
		devSingleInstance = newDevice(deviceType)
	})
	return devSingleInstance
}

// Close closes the usb bus
// Device should be closed before shutdown to avoid running out of open file descriptors
func (d *Device) Close() {
	d.Driver.Close()
}

// Connect makes a connection to the connected device
func (d *Device) Connect() error {
	d.Lock()
	defer d.Unlock()
	if !d.connected {
		dev, err := d.Driver.GetDevice()
		if err == nil {
			d.dev = dev
			d.connected = true
		}
		return err
	}
	return nil
}

// Disconnect the device
func (d *Device) Disconnect() error {
	d.Lock()
	defer d.Unlock()
	if d.connected {
		err := d.dev.Close(false)
		if err == nil {
			d.dev = nil
			d.connected = false
		}
		return nil
	}
	return nil
}

// GetUsbInfo returns information from the attached usb
func (d *Device) GetUsbInfo() ([]usb.Info, error) {
	if d.Driver.DeviceType() == DeviceTypeUSB {
		if err := d.Connect(); err != nil {
			return nil, err
		}
		if err := d.Disconnect(); err != nil {
			return nil, err
		}
	}
	devInfos, err := d.Driver.GetDeviceInfos()
	if err != nil {
		return nil, err
	}
	return devInfos, nil
}

// AddressGen Ask the device to generate an address
func (d *Device) AddressGen(addressN, startIndex uint32, confirmAddress bool, coinType CoinType) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	if addressN == 0 {
		return wire.Message{}, ErrAddressNZero
	}

	addressGenChunks, err := MessageAddressGen(addressN, startIndex, confirmAddress, coinType)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, addressGenChunks)
}

// SaveDeviceEntropyInFile Ask the device to generate entropy and save it in a file
// if `outFile` is the "-" string, the output file is considered stdout
func (d *Device) SaveDeviceEntropyInFile(outFile string, entropyBytes uint32, getEntropyMsgBuilder func(entropyBytes uint32) ([][64]byte, error)) error {
	usingStdout := false
	if outFile == "-" {
		usingStdout = true
	}
	if !usingStdout {
		log.Infoln("Saving entropy to", outFile)
	}
	var receivedEntropyBytes uint32
	var processBytes func(buf []byte) error
	var err error
	var processGetEntropyResponse func(msg wire.Message) (*messages.Entropy, error)
	processGetEntropyResponse = func(msg wire.Message) (*messages.Entropy, error) {
		if err != nil || msg.Kind != uint16(messages.MessageType_MessageType_Entropy) {
			if err != nil {
				return &messages.Entropy{}, err
			}

			if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
				// Send ButtonAck
				chunks, err := MessageButtonAck()
				if err != nil {
					return &messages.Entropy{}, err
				}
				if err = sendToDeviceNoAnswer(d.dev, chunks); err != nil {
					return &messages.Entropy{}, err
				}
				// simulate button press
				if d.simulateButtonPress {
					if err := d.SimulateButtonPress(); err != nil {
						return &messages.Entropy{}, err
					}
				}

				msg, err := wire.ReadFrom(d.dev)
				if err != nil {
					return nil, err
				}
				return processGetEntropyResponse(*msg)
			}
			var msgStr string
			msgStr, err = DecodeFailMsg(msg)
			if err != nil {
				log.Errorf("Error decoding device response as fails msg %s", err)
				return &messages.Entropy{}, err
			}
			err = errors.New(msgStr)
			log.Errorf("Error getting entropy from device %s", err)
			return &messages.Entropy{}, err
		}
		entropy, err := DecodeResponseEntropyMessage(msg)
		if err != nil {
			log.Errorf("Error decoding device response %s", err)
			return &messages.Entropy{}, err
		}
		return entropy, nil
	}

	getEntropy := func(bytes uint32) (*messages.Entropy, error) {
		chunks, err := getEntropyMsgBuilder(bytes)
		if err != nil {
			return &messages.Entropy{}, err
		}
		resp, err := d.Driver.SendToDevice(d.dev, chunks)
		if err != nil {
			return &messages.Entropy{}, err
		}
		return processGetEntropyResponse(resp)
	}

	checkProducedFile := func() error {
		if !usingStdout {
			fileInfo, err := os.Stat(outFile)
			if err != nil {
				log.Error(err)
				return err
			}
			if fileInfo.Size() != int64(entropyBytes) {
				return fmt.Errorf(
					"no engout bytes saved in the file %s\n current: %d\nrequired: %d",
					outFile, fileInfo.Size(), entropyBytes)
			}
		}
		return nil
	}

	if usingStdout {
		processBytes = func(buf []byte) error {
			fmt.Print(buf)
			return nil
		}
	} else {
		pb := Progbar{total: int(entropyBytes)}
		defer func() {
			if checkProducedFile() == nil {
				pb.PrintComplete()
			}
		}()
		if _, err := os.Stat(outFile); err == nil {
			// nolint: gosec
			if err = os.Chmod(outFile, 0777); err != nil {
				log.Errorf("error with %s %s", outFile, err)
			}
		}
		file, err := os.Create(outFile)
		if err != nil {
			log.Errorf("error creating output file %s", err)
			return err
		}
		defer func() {
			if err := os.Chmod(outFile, 0444); err != nil {
				log.Error(err)
			}
		}()
		defer file.Close()

		processBytes = func(buf []byte) error {
			var wroteBytes = 0
			for wroteBytes < len(buf) {
				var res int
				if res, err = file.Write(buf[wroteBytes:]); err != nil {
					return err
				}
				wroteBytes += res
			}
			if wroteBytes != len(buf) {
				return errors.New("invalid bytes amount wrote")
			}
			pb.PrintProg(int(receivedEntropyBytes))
			return nil
		}
	}

	if err := d.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := d.Disconnect(); err != nil {
			log.Error(err)
		}
	}()

	entropy, err := getEntropy(entropyBytes)
	if err != nil {
		log.Error(err)
		return err
	}

	receivedEntropyBytes = uint32(len(entropy.GetEntropy()))
	if err := processBytes(entropy.GetEntropy()); err != nil {
		log.Errorf("error writing file %s.\n %s", outFile, err.Error())
		return err
	}

	for receivedEntropyBytes < entropyBytes {
		entropy, err := getEntropy(entropyBytes - receivedEntropyBytes)
		if err != nil {
			log.Error(err)
			return err
		}
		receivedEntropyBytes += uint32(len(entropy.GetEntropy()))
		if err := processBytes(entropy.GetEntropy()); err != nil {
			log.Errorf("error writing file %s.\n %s", outFile, err.Error())
			return err
		}
	}
	return checkProducedFile()
}

// ApplySettings send ApplySettings request to the device
func (d *Device) ApplySettings(usePassphrase *bool, label string, language string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	applySettingsChunks, err := MessageApplySettings(usePassphrase, label, language)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, applySettingsChunks)
}

// Backup ask the device to perform the seed backup
func (d *Device) Backup() (msg wire.Message, err error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()
	backupChunks, err := MessageBackup()
	if err != nil {
		return wire.Message{}, err
	}
	msg, err = d.Driver.SendToDevice(d.dev, backupChunks)
	if err != nil {
		return wire.Message{}, err
	}
	return msg, err
}

// Cancel sends a Cancel request
func (d *Device) Cancel() (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	cancelChunks, err := MessageCancel()
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, cancelChunks)
}

// CheckMessageSignature Check a message signature matches the given address.
func (d *Device) CheckMessageSignature(message, signature, address string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	// Send CheckMessageSignature
	checkMessageSignatureChunks, err := MessageCheckMessageSignature(message, signature, address)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, checkMessageSignatureChunks)
}

// ChangePin changes device's PIN code
// The message that is sent contains an encoded form of the PIN.
// The digits of the PIN are displayed in a 3x3 matrix on the Trezor,
// and the message that is sent back is a string containing the positions
// of the digits on that matrix. Below is the mapping between positions
// and characters to be sent:
// 7 8 9
// 4 5 6
// 1 2 3
// For example, if the numbers are laid out in this way on the Trezor,
// 3 1 5
// 7 8 4
// 9 6 2
// To set the PIN "12345", the positions are:
// top, bottom-right, top-left, right, top-right
// so you must send "83769".
func (d *Device) ChangePin(removePin *bool) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	if removePin == nil {
		return wire.Message{}, ErrRemovePinNil
	}

	changePinChunks, err := MessageChangePin(removePin)
	if err != nil {
		return wire.Message{}, err
	}

	msg, err := d.Driver.SendToDevice(d.dev, changePinChunks)
	if err != nil {
		return wire.Message{}, err
	}

	return msg, nil
}

// Connected checks if we can communicate with a connected skycoin wallet
func (d *Device) Connected() bool {
	var msg *wire.Message
	var err error
	if d.dev == nil {
		return false
	}

	chunks, err := MessageConnected()
	if err != nil {
		log.Error(err)
		return false
	}

	for _, element := range chunks {
		_, err = d.dev.Write(element[:])
		if err != nil {
			return false
		}
	}

	msg, err = wire.ReadFrom(d.dev)
	if err != nil {
		return false
	}

	for msg.Kind == uint16(messages.MessageType_MessageType_EntropyRequest) {
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			entropyChunks, err := MessageEntropyAck(entropyBufferSize)
			if err != nil {
				log.Errorf("failed to create entropy ack msg: %v", err)
				return
			}

			for _, element := range entropyChunks {
				_, err := d.dev.Write(element[:])
				if err != nil {
					log.Errorf("entropy ack error: %v", err)
					return
				}
			}
		}()

		msg, err = wire.ReadFrom(d.dev)
		if err != nil {
			return false
		}
		wg.Wait()
	}

	return msg.Kind == uint16(messages.MessageType_MessageType_Success)
}

// Available checks if a skycoin wallet is connected to the system
func (d *Device) Available() bool {
	infos, err := d.Driver.GetDeviceInfos()
	if err != nil {
		log.Error(err)
		return false
	}

	if len(infos) <= 0 {
		return false
	}

	return true
}

// FirmwareUpload Updates device's firmware
func (d *Device) FirmwareUpload(payload []byte, hash [32]byte) error {
	if d.Driver.DeviceType() != DeviceTypeUSB {
		return ErrDeviceTypeEmulator
	}

	if err := d.Connect(); err != nil {
		return err
	}
	defer d.Disconnect()

	if err := Initialize(d.dev); err != nil {
		return err
	}

	log.Printf("Length of firmware %d", uint32(len(payload)))

	chunks, err := MessageFirmwareErase(payload)
	if err != nil {
		return err
	}
	erasemsg, err := d.Driver.SendToDevice(d.dev, chunks)
	if err != nil {
		return err
	}

	switch erasemsg.Kind {
	case uint16(messages.MessageType_MessageType_Success):
		log.Printf("Success %d! FirmwareErase %s\n", erasemsg.Kind, erasemsg.Data)
	case uint16(messages.MessageType_MessageType_Failure):
		msg, err := DecodeFailMsg(erasemsg)
		if err != nil {
			return err
		}

		return errors.New(msg)
	default:
		return fmt.Errorf("received unexpected message type: %s", messages.MessageType(erasemsg.Kind))
	}

	log.Printf("Hash: %x\n", hash)

	chunks, err = MessageFirmwareUpload(payload, hash)
	if err != nil {
		return err
	}
	uploadmsg, err := d.Driver.SendToDevice(d.dev, chunks)
	if err != nil {
		return err
	}

	switch uploadmsg.Kind {
	case uint16(messages.MessageType_MessageType_ButtonRequest):
		log.Println("Please confirm in the device if fingerprints match")
		// Send ButtonAck
		chunks, err = MessageButtonAck()
		if err != nil {
			return err
		}
		resp, err := d.Driver.SendToDevice(d.dev, chunks)
		if err != nil {
			return err
		}
		switch resp.Kind {
		case uint16(messages.MessageType_MessageType_Success):
			return nil
		case uint16(messages.MessageType_MessageType_Failure):
			var msgStr string
			if msgStr, err = DecodeFailMsg(resp); err != nil {
				return err
			}
			return errors.New(msgStr)
		default:
			return errors.New("unknown response")
		}
	case uint16(messages.MessageType_MessageType_Failure):
		msg, err := DecodeFailMsg(erasemsg)
		if err != nil {
			return err
		}

		return errors.New(msg)
	default:
		return fmt.Errorf("received unexpected message type: %s", messages.MessageType(erasemsg.Kind))
	}
}

// GetFeatures send Features message to the device
func (d *Device) GetFeatures() (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	getFeaturesChunks, err := MessageGetFeatures()
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, getFeaturesChunks)
}

// GenerateMnemonic Ask the device to generate a mnemonic and configure itself with it.
func (d *Device) GenerateMnemonic(wordCount uint32, usePassphrase bool) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	if wordCount != 12 && wordCount != 24 {
		return wire.Message{}, ErrInvalidWordCount
	}

	generateMnemonicChunks, err := MessageGenerateMnemonic(wordCount, usePassphrase)
	if err != nil {
		return wire.Message{}, err
	}

	msg, err := d.Driver.SendToDevice(d.dev, generateMnemonicChunks)
	if err != nil {
		return msg, err
	}

	return msg, err
}

// Recovery ask the device to perform the seed backup
func (d *Device) Recovery(wordCount uint32, usePassphrase *bool, dryRun bool) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	if wordCount != 12 && wordCount != 24 {
		return wire.Message{}, ErrInvalidWordCount
	}

	log.Printf("Using passphrase %t\n", usePassphrase)
	recoveryChunks, err := MessageRecovery(wordCount, usePassphrase, dryRun)
	if err != nil {
		return wire.Message{}, err
	}

	msg, err := d.Driver.SendToDevice(d.dev, recoveryChunks)
	if err != nil {
		return wire.Message{}, err
	}
	log.Printf("Recovery device response kind is: %d\n", msg.Kind)

	return msg, nil
}

// SetMnemonic Configure the device with a mnemonic.
func (d *Device) SetMnemonic(mnemonic string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	// Send SetMnemonic
	setMnemonicChunks, err := MessageSetMnemonic(mnemonic)
	if err != nil {
		return wire.Message{}, err
	}
	msg, err := d.Driver.SendToDevice(d.dev, setMnemonicChunks)
	if err != nil {
		return wire.Message{}, err
	}

	return msg, err
}

// SignMessage Ask the device to sign a message using the secret key at given index.
func (d *Device) SignMessage(addressIndex int, message string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	signMessageChunks, err := MessageSignMessage(addressIndex, message)
	if err != nil {
		return wire.Message{}, err
	}

	msg, err := d.Driver.SendToDevice(d.dev, signMessageChunks)
	if err != nil {
		return wire.Message{}, err
	}

	return msg, err
}

// TransactionSign Ask the device to sign a transaction using the given information.
func (d *Device) TransactionSign(inputs []*messages.SkycoinTransactionInput, outputs []*messages.SkycoinTransactionOutput) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	var transactionInputs []*messages.TxAck_TransactionType_TxInputType
	var transactionOutputs []*messages.TxAck_TransactionType_TxOutputType

	for _, input := range inputs {
		transactionInputs = append(transactionInputs, &messages.TxAck_TransactionType_TxInputType{
			AddressN: make([]uint32, *input.Index),
			HashIn:   input.HashIn,
		})
	}
	for _, output := range outputs {
		transactionOutputs = append(transactionOutputs, &messages.TxAck_TransactionType_TxOutputType{
			Address: output.Address,
			Coins:   output.Coin,
			Hours:   output.Hour,
		})
		if output.AddressIndex != nil {
			transactionOutputs[len(transactionOutputs)-1].AddressN = make([]uint32, *output.AddressIndex)
		}
	}
	signer := SkycoinTransactionSigner{
		Device:   d,
		Inputs:   transactionInputs,
		Outputs:  transactionOutputs,
		Version:  1,
		LockTime: 0,
	}
	signatures, err := signer.Sign()
	if err != nil {
		return wire.Message{}, err
	}
	padding := false
	data, err := proto.Marshal(&messages.ResponseTransactionSign{
		Padding:    &padding,
		Signatures: signatures,
	})
	if err != nil {
		return wire.Message{}, err
	}
	return wire.Message{
		Kind: uint16(messages.MessageType_MessageType_ResponseTransactionSign),
		Data: data,
	}, nil
}

// GeneralTransactionSign Ask the device to sign a transaction using the given TransactionSigner
func (d *Device) GeneralTransactionSign(signer TransactionSigner) ([]string, error) {
	signer.SetDevice(d)
	return signer.Sign()
}

// SignTx Ask the device to sign a long transaction using the given information.
func (d *Device) SignTx(outputsCount int, inputsCount int, coinName string, version int, lockTime int, txHash string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	signTxChunks, err := MessageSignTx(outputsCount, inputsCount, coinName, version, lockTime, txHash)

	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, signTxChunks)
}

// TxAck Ask the device to continue a long transaction using the given information.
func (d *Device) TxAck(inputs []*messages.TxAck_TransactionType_TxInputType, outputs []*messages.TxAck_TransactionType_TxOutputType, version int, lockTime int) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()
	txAckChunks, err := MessageTxAck(inputs, outputs, version, lockTime)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, txAckChunks)
}

// BitcoinTxAck ask the device to continue a Bitcoin long transaction using the given information.
func (d *Device) BitcoinTxAck(inputs []*messages.BitcoinTransactionInput, outputs []*messages.BitcoinTransactionOutput) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()
	txAckChunks, err := BitcoinMessageTxAck(inputs, outputs)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, txAckChunks)
}

// Wipe wipes out device configuration
func (d *Device) Wipe() (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	wipeChunks, err := MessageWipe()
	if err != nil {
		return wire.Message{}, err
	}

	msg, err := d.Driver.SendToDevice(d.dev, wipeChunks)
	if err != nil {
		return wire.Message{}, err
	}

	return msg, err
}

// ButtonAck when the device is waiting for the user to press a button
// the PC need to acknowledge, showing it knows we are waiting for a user action
func (d *Device) ButtonAck() (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	// Send ButtonAck
	buttonChunks, err := MessageButtonAck()
	if err != nil {
		return wire.Message{}, err
	}

	err = sendToDeviceNoAnswer(d.dev, buttonChunks)
	if err != nil {
		return wire.Message{}, err
	}

	// simulate button press
	if d.simulateButtonPress {
		if err := d.SimulateButtonPress(); err != nil {
			return wire.Message{}, err
		}
	}

	msg, err := wire.ReadFrom(d.dev)
	if err != nil {
		return wire.Message{}, err
	}
	for msg.Kind == uint16(messages.MessageType_MessageType_EntropyRequest) {
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			entropyChunks, err := MessageEntropyAck(entropyBufferSize)
			if err != nil {
				log.Errorf("failed to create entropy ack msg: %v", err)
				return
			}

			for _, element := range entropyChunks {
				_, err := d.dev.Write(element[:])
				if err != nil {
					log.Errorf("entropy ack error: %v", err)
					return
				}
			}
		}()

		msg, err = wire.ReadFrom(d.dev)
		if err != nil {
			return wire.Message{}, err
		}
		wg.Wait()
	}

	return *msg, err
}

// PassphraseAck send this message when the device is waiting for the user to input a passphrase
func (d *Device) PassphraseAck(passphrase string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	passphraseChunks, err := MessagePassphraseAck(passphrase)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, passphraseChunks)
}

// WordAck send a word to the device during device "recovery procedure"
func (d *Device) WordAck(word string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	wordAckChunks, err := MessageWordAck(word)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, wordAckChunks)
}

// PinMatrixAck during PIN code setting use this message to send user input to device
func (d *Device) PinMatrixAck(p string) (wire.Message, error) {
	time.Sleep(1 * time.Second)
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.Disconnect()

	log.Printf("Setting pin: %s\n", p)

	pinMatrixChunks, err := MessagePinMatrixAck(p)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, pinMatrixChunks)
}

// SimulateButtonPress simulates a button press on emulator
func (d *Device) SimulateButtonPress() error {
	if d.Driver.DeviceType() != DeviceTypeEmulator {
		return fmt.Errorf("wrong device type: %s", d.Driver.DeviceType())
	}

	simulateMsg, err := MessageSimulateButtonPress(d.simulateButtonType)
	if err != nil {
		return err
	}

	_, err = d.dev.Write(simulateMsg.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// SetAutoPressButton enables and sets button press type
func (d *Device) SetAutoPressButton(simulateButtonPress bool, simulateButtonType ButtonType) error {
	if d.Driver.DeviceType() == DeviceTypeEmulator {
		d.simulateButtonPress = simulateButtonPress

		if simulateButtonPress {
			switch simulateButtonType {
			case ButtonLeft, ButtonRight, ButtonBoth:
				d.simulateButtonType = simulateButtonType
			default:
				return fmt.Errorf("invalid button type: %d", simulateButtonType)
			}
		} else {
			// set invalid button press type
			d.simulateButtonType = 3
		}
	}

	return nil
}
