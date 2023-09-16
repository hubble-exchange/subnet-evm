// Code generated
// This file is a generated precompile contract test with the skeleton of test functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package ticks

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
		"insufficient gas for getBaseQuote should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := GetBaseQuoteInput{
					QuoteQuantity: big.NewInt(0),
				}
				input, err := PackGetBaseQuote(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: GetBaseQuoteGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for getPrevTick should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := GetPrevTickInput{
					Tick: big.NewInt(0),
				}
				input, err := PackGetPrevTick(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: GetPrevTickGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for getQuote should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := GetQuoteInput{
					BaseAssetQuantity: big.NewInt(0),
				}
				input, err := PackGetQuote(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: GetQuoteGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for sampleImpactAsk should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				var testInput common.Address
				input, err := PackSampleImpactAsk(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: SampleImpactAskGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for sampleImpactBid should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				var testInput common.Address
				input, err := PackSampleImpactBid(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: SampleImpactBidGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	}
)

// TestTicksRun tests the Run function of the precompile contract.
func TestTicksRun(t *testing.T) {
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func BenchmarkTicks(b *testing.B) {
	// Benchmark tests.
	for name, test := range tests {
		b.Run(name, func(b *testing.B) {
			test.Bench(b, Module, state.NewTestStateDB(b))
		})
	}
}
