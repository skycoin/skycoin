package skywallet

import (
	"encoding/binary"
	"encoding/json"
)

// BitEncodedFlags allow you to work with bit field encoded in integer in a high level
type BitEncodedFlags interface {
	// Marshal encode the fields to a uint64
	Marshal() (uint64, error)

	// decode fields
	Unmarshal() error

	// return tru if the RDP is enabled
	HasRdpMemProtectEnabled() bool
}

// FirmwareFeatures handle the features in firmware as a BitEncodedFlags implementation
type FirmwareFeatures struct {
	flags                    uint64
	RequireGetEntropyConfirm bool
	IsGetEntropyEnabled      bool
	IsEmulator               bool
	FirmwareFeaturesRdpLevel uint8
}

// NewFirmwareFeatures return a new BitEncodedFlags initialized with the internal fields
func NewFirmwareFeatures(flags uint64) BitEncodedFlags {
	return &FirmwareFeatures{flags: flags}
}

// Marshal encode the FirmwareFeatures internal field in a uint64 value
// return this number and keeps an internal copy
func (ff *FirmwareFeatures) Marshal() (uint64, error) {
	ff.flags = 0
	bs := make([]byte, 8)
	setBitInByte(&bs[7], ff.RequireGetEntropyConfirm, 0)
	setBitInByte(&bs[7], ff.IsGetEntropyEnabled, 1)
	setBitInByte(&bs[7], ff.IsEmulator, 2)
	setBitInByte(&bs[7], ff.FirmwareFeaturesRdpLevel == 1 || ff.FirmwareFeaturesRdpLevel == 3, 3)
	setBitInByte(&bs[7], ff.FirmwareFeaturesRdpLevel == 2 || ff.FirmwareFeaturesRdpLevel == 3, 4)
	ff.flags = binary.BigEndian.Uint64(bs)
	return ff.flags, nil
}

// Unmarshal fill all the struct fields based on the encoded info in the flags field
func (ff *FirmwareFeatures) Unmarshal() error {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, ff.flags)
	ff.RequireGetEntropyConfirm = bitStatusInByte(bs[7], 0)
	ff.IsGetEntropyEnabled = bitStatusInByte(bs[7], 1)
	ff.IsEmulator = bitStatusInByte(bs[7], 2)
	setBitInByte(&ff.FirmwareFeaturesRdpLevel, bitStatusInByte(bs[7], 3), 0)
	setBitInByte(&ff.FirmwareFeaturesRdpLevel, bitStatusInByte(bs[7], 4), 1)
	return nil
}

// HasRdpMemProtectEnabled return true if rdp == true
func (ff FirmwareFeatures) HasRdpMemProtectEnabled() bool {
	return ff.FirmwareFeaturesRdpLevel == 2
}

// String allow pretty print in cli applications
func (ff FirmwareFeatures) String() string {
	b, err := json.Marshal(ff)
	if err != nil {
		return "error rendering FirmwareFeatures " + err.Error()
	}
	return string(b)
}

func bitStatusInByte(data, bitPos uint8) bool {
	return (data & (uint8)(1<<bitPos)) != 0
}

func setBitInByte(data *uint8, val bool, bitPos uint8) {
	mask := (uint8)(1 << bitPos)
	if val {
		*data |= mask
	} else {
		*data &= mask ^ 255
	}
}
