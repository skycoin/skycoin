package transaction

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
)

func TestCreateWalletParamsVerify(t *testing.T) {
	changeAddress := testutil.MakeAddress()

	toManual := []coin.TransactionOutput{
		{
			Address: testutil.MakeAddress(),
			Coins:   1e6,
			Hours:   1,
		},
		{
			Address: testutil.MakeAddress(),
			Coins:   5e6,
			Hours:   0,
		},
	}

	toAuto := []coin.TransactionOutput{
		{
			Address: testutil.MakeAddress(),
			Coins:   1e6,
		},
		{
			Address: testutil.MakeAddress(),
			Coins:   5e6,
		},
	}

	one := decimal.New(1, 0)
	negativeOne := decimal.New(-1, 0)
	onePointOne := decimal.New(11, -1)
	pointOneOne := decimal.New(11, -2)

	cases := []struct {
		name   string
		params Params
		err    string
	}{
		{
			name: "null change address",
			params: Params{
				ChangeAddress: &cipher.Address{},
			},
			err: "ChangeAddress must not be the null address",
		},

		{
			name: "no to destinations",
			params: Params{
				ChangeAddress: &changeAddress,
			},
			err: "To is required",
		},

		{
			name: "missing to coins",
			params: Params{
				ChangeAddress: &changeAddress,
				To: []coin.TransactionOutput{
					{
						Address: testutil.MakeAddress(),
						Hours:   1,
					},
				},
			},
			err: "To.Coins must not be zero",
		},

		{
			name: "missing to address",
			params: Params{
				ChangeAddress: &changeAddress,
				To: []coin.TransactionOutput{
					{
						Coins: 5,
						Hours: 1,
					},
				},
			},
			err: "To.Address must not be the null address",
		},

		{
			name: "nonzero to hours for auto selection",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toManual,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeAuto,
				},
			},
			err: "To.Hours must be zero for auto type hours selection",
		},

		{
			name: "mode missing for auto selection",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toAuto,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeAuto,
				},
			},
			err: "HoursSelection.Mode is required for auto type hours selection",
		},

		{
			name: "mode set for manual selection",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toManual,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
					Mode: HoursSelectionModeShare,
				},
			},
			err: "HoursSelection.Mode cannot be used for manual type hours selection",
		},

		{
			name: "missing hours selection type",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toAuto,
				HoursSelection: HoursSelection{
					Type: "",
				},
			},
			err: "Invalid HoursSelection.Type",
		},

		{
			name: "invalid hours selection type",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toAuto,
				HoursSelection: HoursSelection{
					Type: "invalid",
				},
			},
			err: "Invalid HoursSelection.Type",
		},

		{
			name: "invalid hours selection mode",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toAuto,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeAuto,
					Mode: "invalid",
				},
			},
			err: "Invalid HoursSelection.Mode",
		},

		{
			name: "share factor not set for split even mode",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toAuto,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeAuto,
					Mode: HoursSelectionModeShare,
				},
			},
			err: "HoursSelection.ShareFactor must be set for share mode",
		},

		{
			name: "share factor set but not split even mode",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toManual,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeManual,
					ShareFactor: &one,
				},
			},
			err: "HoursSelection.ShareFactor can only be used for share mode",
		},

		{
			name: "share factor less than 0",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toAuto,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: &negativeOne,
				},
			},
			err: "HoursSelection.ShareFactor must be >= 0 and <= 1",
		},

		{
			name: "share factor greater than 1",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toAuto,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: &onePointOne,
				},
			},
			err: "HoursSelection.ShareFactor must be >= 0 and <= 1",
		},

		{
			name: "duplicate output when manual",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            []coin.TransactionOutput{toManual[0], toManual[0]},
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
			},
			err: "To contains duplicate values",
		},

		{
			name: "duplicate output when auto",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            []coin.TransactionOutput{toAuto[0], toAuto[0]},
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: &pointOneOne,
				},
			},
			err: "To contains duplicate values",
		},

		{
			name: "valid auto split even share factor",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toAuto,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: &pointOneOne,
				},
			},
		},

		{
			name: "valid manual",
			params: Params{
				ChangeAddress: &changeAddress,
				To:            toManual,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.err != "" {
				require.Equal(t, NewError(errors.New(tc.err)), err, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
