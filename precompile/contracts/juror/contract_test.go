// Code generated
// This file is a generated precompile contract test with the skeleton of test functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package juror

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// These tests are run against the precompile contract directly with
// the given input and expected output. They're just a guide to
// help you write your own tests. These tests are for general cases like
// allowlist, readOnly behaviour, and gas cost. You should write your own
// tests for specific cases.
var (
	tests = map[string]testutils.PrecompileTest{
		"insufficient gas for getNotionalPositionAndMargin should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := GetNotionalPositionAndMarginInput{
					Trader: common.Address{1},
				}
				input, err := PackGetNotionalPositionAndMargin(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: GetNotionalPositionAndMarginGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for validateCancelLimitOrder should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := ValidateCancelLimitOrderInput{
					Order: ILimitOrderBookOrder{
						AmmIndex:          big.NewInt(0),
						Trader:            common.Address{1},
						BaseAssetQuantity: big.NewInt(0),
						Price:             big.NewInt(0),
						Salt:              big.NewInt(0),
						ReduceOnly:        false,
						PostOnly:          false,
					},
					Sender: common.Address{1},
				}
				input, err := PackValidateCancelLimitOrder(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidateCancelLimitOrderGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for validateLiquidationOrderAndDetermineFillPrice should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := ValidateLiquidationOrderAndDetermineFillPriceInput{
					LiquidationAmount: big.NewInt(0),
				}
				input, err := PackValidateLiquidationOrderAndDetermineFillPrice(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidateLiquidationOrderAndDetermineFillPriceGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for validateOrdersAndDetermineFillPrice should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := ValidateOrdersAndDetermineFillPriceInput{FillAmount: big.NewInt(0)}
				input, err := PackValidateOrdersAndDetermineFillPrice(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidateOrdersAndDetermineFillPriceGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for validatePlaceIOCOrder should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := ValidatePlaceIOCOrderInput{
					Order: IImmediateOrCancelOrdersOrder{
						OrderType:         uint8(IOC),
						ExpireAt:          big.NewInt(0),
						AmmIndex:          big.NewInt(0),
						Trader:            common.Address{1},
						BaseAssetQuantity: big.NewInt(0),
						Price:             big.NewInt(0),
						Salt:              big.NewInt(0),
						ReduceOnly:        false,
					},
					Sender: common.Address{1},
				}
				input, err := PackValidatePlaceIOCOrder(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidatePlaceIOCOrderGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for validatePlaceLimitOrder should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := ValidatePlaceLimitOrderInput{
					Order: ILimitOrderBookOrder{
						AmmIndex:          big.NewInt(0),
						Trader:            common.Address{1},
						BaseAssetQuantity: big.NewInt(0),
						Price:             big.NewInt(0),
						Salt:              big.NewInt(0),
						ReduceOnly:        false,
						PostOnly:          false,
					},
					Sender: common.Address{1},
				}
				input, err := PackValidatePlaceLimitOrder(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidatePlaceLimitOrderGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	}
)

// TestJurorRun tests the Run function of the precompile contract.
func TestJurorRun(t *testing.T) {
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func BenchmarkJuror(b *testing.B) {
	// Benchmark tests.
	for name, test := range tests {
		b.Run(name, func(b *testing.B) {
			test.Bench(b, Module, state.NewTestStateDB(b))
		})
	}
}
