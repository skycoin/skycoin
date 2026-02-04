package skywallet

import (
	"testing"
)

func TestFirmwareFeaturesMarshal(t *testing.T) {
	tests := []struct {
		name     string
		features FirmwareFeatures
		wantBits uint64
	}{
		{
			name:     "all false",
			features: FirmwareFeatures{},
			wantBits: 0,
		},
		{
			name: "RequireGetEntropyConfirm only",
			features: FirmwareFeatures{
				RequireGetEntropyConfirm: true,
			},
			wantBits: 1, // bit 0
		},
		{
			name: "IsGetEntropyEnabled only",
			features: FirmwareFeatures{
				IsGetEntropyEnabled: true,
			},
			wantBits: 2, // bit 1
		},
		{
			name: "IsEmulator only",
			features: FirmwareFeatures{
				IsEmulator: true,
			},
			wantBits: 4, // bit 2
		},
		{
			name: "RDP level 1",
			features: FirmwareFeatures{
				FirmwareFeaturesRdpLevel: 1,
			},
			wantBits: 8, // bit 3
		},
		{
			name: "RDP level 2",
			features: FirmwareFeatures{
				FirmwareFeaturesRdpLevel: 2,
			},
			wantBits: 16, // bit 4
		},
		{
			name: "RDP level 3",
			features: FirmwareFeatures{
				FirmwareFeaturesRdpLevel: 3,
			},
			wantBits: 24, // bits 3 and 4
		},
		{
			name: "all flags set with RDP level 2",
			features: FirmwareFeatures{
				RequireGetEntropyConfirm: true,
				IsGetEntropyEnabled:      true,
				IsEmulator:               true,
				FirmwareFeaturesRdpLevel: 2,
			},
			wantBits: 1 + 2 + 4 + 16, // 23
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.features.Marshal()
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if got != tt.wantBits {
				t.Errorf("Marshal() = %d (0b%b), want %d (0b%b)", got, got, tt.wantBits, tt.wantBits)
			}
		})
	}
}

func TestFirmwareFeaturesUnmarshal(t *testing.T) {
	tests := []struct {
		name               string
		flags              uint64
		wantRequireEntropy bool
		wantEntropyEnabled bool
		wantEmulator       bool
		wantRdpLevel       uint8
	}{
		{
			name:  "all zero",
			flags: 0,
		},
		{
			name:               "bit 0 set",
			flags:              1,
			wantRequireEntropy: true,
		},
		{
			name:               "bit 1 set",
			flags:              2,
			wantEntropyEnabled: true,
		},
		{
			name:         "bit 2 set",
			flags:        4,
			wantEmulator: true,
		},
		{
			name:         "bit 3 set (RDP level 1)",
			flags:        8,
			wantRdpLevel: 1,
		},
		{
			name:         "bit 4 set (RDP level 2)",
			flags:        16,
			wantRdpLevel: 2,
		},
		{
			name:         "bits 3 and 4 set (RDP level 3)",
			flags:        24,
			wantRdpLevel: 3,
		},
		{
			name:         "bit 4 only (RDP level 2) - flags 16",
			flags:        16, // 0b10000 = bit 4 set
			wantRdpLevel: 2,
		},
		{
			name:               "all flags",
			flags:              31, // 0b11111
			wantRequireEntropy: true,
			wantEntropyEnabled: true,
			wantEmulator:       true,
			wantRdpLevel:       3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := &FirmwareFeatures{flags: tt.flags}
			err := ff.Unmarshal()
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			if ff.RequireGetEntropyConfirm != tt.wantRequireEntropy {
				t.Errorf("RequireGetEntropyConfirm = %v, want %v", ff.RequireGetEntropyConfirm, tt.wantRequireEntropy)
			}
			if ff.IsGetEntropyEnabled != tt.wantEntropyEnabled {
				t.Errorf("IsGetEntropyEnabled = %v, want %v", ff.IsGetEntropyEnabled, tt.wantEntropyEnabled)
			}
			if ff.IsEmulator != tt.wantEmulator {
				t.Errorf("IsEmulator = %v, want %v", ff.IsEmulator, tt.wantEmulator)
			}
			if ff.FirmwareFeaturesRdpLevel != tt.wantRdpLevel {
				t.Errorf("FirmwareFeaturesRdpLevel = %v, want %v", ff.FirmwareFeaturesRdpLevel, tt.wantRdpLevel)
			}
		})
	}
}

func TestFirmwareFeaturesRoundTrip(t *testing.T) {
	// Test that marshaling and unmarshaling produces the same values
	tests := []FirmwareFeatures{
		{},
		{RequireGetEntropyConfirm: true},
		{IsGetEntropyEnabled: true},
		{IsEmulator: true},
		{FirmwareFeaturesRdpLevel: 1},
		{FirmwareFeaturesRdpLevel: 2},
		{FirmwareFeaturesRdpLevel: 3},
		{
			RequireGetEntropyConfirm: true,
			IsGetEntropyEnabled:      true,
			IsEmulator:               true,
			FirmwareFeaturesRdpLevel: 2,
		},
	}

	for _, original := range tests {
		// Marshal
		flags, err := original.Marshal()
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}

		// Create new instance and unmarshal
		restored := &FirmwareFeatures{flags: flags}
		err = restored.Unmarshal()
		if err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}

		// Compare
		if original.RequireGetEntropyConfirm != restored.RequireGetEntropyConfirm {
			t.Errorf("RequireGetEntropyConfirm mismatch: %v vs %v", original.RequireGetEntropyConfirm, restored.RequireGetEntropyConfirm)
		}
		if original.IsGetEntropyEnabled != restored.IsGetEntropyEnabled {
			t.Errorf("IsGetEntropyEnabled mismatch: %v vs %v", original.IsGetEntropyEnabled, restored.IsGetEntropyEnabled)
		}
		if original.IsEmulator != restored.IsEmulator {
			t.Errorf("IsEmulator mismatch: %v vs %v", original.IsEmulator, restored.IsEmulator)
		}
		if original.FirmwareFeaturesRdpLevel != restored.FirmwareFeaturesRdpLevel {
			t.Errorf("FirmwareFeaturesRdpLevel mismatch: %v vs %v", original.FirmwareFeaturesRdpLevel, restored.FirmwareFeaturesRdpLevel)
		}
	}
}

func TestHasRdpMemProtectEnabled(t *testing.T) {
	tests := []struct {
		rdpLevel uint8
		want     bool
	}{
		{0, false},
		{1, false},
		{2, true}, // Only level 2 means RDP is enabled
		{3, false},
	}

	for _, tt := range tests {
		ff := FirmwareFeatures{FirmwareFeaturesRdpLevel: tt.rdpLevel}
		if got := ff.HasRdpMemProtectEnabled(); got != tt.want {
			t.Errorf("HasRdpMemProtectEnabled() with RDP level %d = %v, want %v", tt.rdpLevel, got, tt.want)
		}
	}
}

func TestNewFirmwareFeatures(t *testing.T) {
	flags := uint64(23) // Some test value
	bf := NewFirmwareFeatures(flags)

	ff, ok := bf.(*FirmwareFeatures)
	if !ok {
		t.Fatal("NewFirmwareFeatures did not return *FirmwareFeatures")
	}

	if ff.flags != flags {
		t.Errorf("flags = %d, want %d", ff.flags, flags)
	}
}

func TestFirmwareFeaturesString(t *testing.T) {
	ff := FirmwareFeatures{
		RequireGetEntropyConfirm: true,
		IsGetEntropyEnabled:      true,
		IsEmulator:               false,
		FirmwareFeaturesRdpLevel: 2,
	}

	str := ff.String()
	if str == "" {
		t.Error("String() returned empty string")
	}

	// Should contain JSON
	if str[0] != '{' {
		t.Errorf("String() should return JSON, got: %s", str)
	}
}

func TestBitHelperFunctions(t *testing.T) {
	// Test bitStatusInByte
	tests := []struct {
		data   uint8
		bitPos uint8
		want   bool
	}{
		{0b00000001, 0, true},
		{0b00000010, 1, true},
		{0b00000100, 2, true},
		{0b00001000, 3, true},
		{0b00010000, 4, true},
		{0b00000000, 0, false},
		{0b11111110, 0, false},
		{0b11111111, 7, true},
	}

	for _, tt := range tests {
		got := bitStatusInByte(tt.data, tt.bitPos)
		if got != tt.want {
			t.Errorf("bitStatusInByte(0b%08b, %d) = %v, want %v", tt.data, tt.bitPos, got, tt.want)
		}
	}

	// Test setBitInByte
	var b uint8
	setBitInByte(&b, true, 0)
	if b != 1 {
		t.Errorf("setBitInByte set bit 0: got %d, want 1", b)
	}

	setBitInByte(&b, true, 2)
	if b != 5 { // 0b101
		t.Errorf("setBitInByte set bit 2: got %d, want 5", b)
	}

	setBitInByte(&b, false, 0)
	if b != 4 { // 0b100
		t.Errorf("setBitInByte clear bit 0: got %d, want 4", b)
	}
}
