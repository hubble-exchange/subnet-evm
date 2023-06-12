// Code generated
// This file is a generated precompile contract test with the skeleton of test functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package hubblebibliophile

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRun tests the Run function of the precompile contract.
// These tests are run against the precompile contract directly with
// the given input and expected output. They're just a guide to
// help you write your own tests. These tests are for general cases like
// allowlist, readOnly behaviour, and gas cost. You should write your own
// tests for specific cases.
func TestRun(t *testing.T) {
	trader := common.HexToAddress("0x6900000000000000000000000000000000000069")
	tests := map[string]testutils.PrecompileTest{
		"insufficient gas for getNotionalPositionAndMargin should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				testInput := GetNotionalPositionAndMarginInput{
					Trader: trader,
				}
				input, err := PackGetNotionalPositionAndMargin(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: GetNotionalPositionAndMarginGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for getPositionSizes should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				testInput := common.HexToAddress("0x0300000000000000000000000000000000000000")
				input, err := PackGetPositionSizes(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: GetPositionSizesGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas for validateLiquidationOrderAndDetermineFillPrice should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := ValidateLiquidationOrderAndDetermineFillPriceInput{
					Order: IHubbleBibliophileOrder{
						AmmIndex:          big.NewInt(0),
						Trader:            trader,
						BaseAssetQuantity: big.NewInt(0),
						Price:             big.NewInt(10),
						Salt:              big.NewInt(0),
						ReduceOnly:        false,
					},
					FillAmount: big.NewInt(0),
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
				testInput := ValidateOrdersAndDetermineFillPriceInput{
					Orders: [2]IHubbleBibliophileOrder{
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(0),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(0),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
					},
					FillAmount: big.NewInt(0),
				}
				input, err := PackValidateOrdersAndDetermineFillPrice(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidateOrdersAndDetermineFillPriceGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"ErrNotLongOrder_0": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				testInput := ValidateOrdersAndDetermineFillPriceInput{
					Orders: [2]IHubbleBibliophileOrder{
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(0),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(0),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
					},
					FillAmount: big.NewInt(0),
				}
				input, err := PackValidateOrdersAndDetermineFillPrice(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidateOrdersAndDetermineFillPriceGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrNotLongOrder.Error(),
		},
		"ErrNotLongOrder": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				testInput := ValidateOrdersAndDetermineFillPriceInput{
					Orders: [2]IHubbleBibliophileOrder{
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(-1),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(0),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
					},
					FillAmount: big.NewInt(0),
				}
				input, err := PackValidateOrdersAndDetermineFillPrice(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidateOrdersAndDetermineFillPriceGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrNotLongOrder.Error(),
		},
		"ErrNotShortOrder_0": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				testInput := ValidateOrdersAndDetermineFillPriceInput{
					Orders: [2]IHubbleBibliophileOrder{
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(1),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(0),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
					},
					FillAmount: big.NewInt(0),
				}
				input, err := PackValidateOrdersAndDetermineFillPrice(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidateOrdersAndDetermineFillPriceGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrNotShortOrder.Error(),
		},
		"ErrNotShortOrder": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				testInput := ValidateOrdersAndDetermineFillPriceInput{
					Orders: [2]IHubbleBibliophileOrder{
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(1),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(1),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
					},
					FillAmount: big.NewInt(0),
				}
				input, err := PackValidateOrdersAndDetermineFillPrice(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidateOrdersAndDetermineFillPriceGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrNotShortOrder.Error(),
		},
		"ErrNotSameAMM": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				testInput := ValidateOrdersAndDetermineFillPriceInput{
					Orders: [2]IHubbleBibliophileOrder{
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(1),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(1),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(-1),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
					},
					FillAmount: big.NewInt(0),
				}
				input, err := PackValidateOrdersAndDetermineFillPrice(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidateOrdersAndDetermineFillPriceGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrNotSameAMM.Error(),
		},
		"ErrNoMatch": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				testInput := ValidateOrdersAndDetermineFillPriceInput{
					Orders: [2]IHubbleBibliophileOrder{
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(1),
							Price:             big.NewInt(10),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
						IHubbleBibliophileOrder{
							AmmIndex:          big.NewInt(0),
							Trader:            trader,
							BaseAssetQuantity: big.NewInt(-1),
							Price:             big.NewInt(11),
							Salt:              big.NewInt(0),
							ReduceOnly:        false,
						},
					},
					FillAmount: big.NewInt(0),
				}
				input, err := PackValidateOrdersAndDetermineFillPrice(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ValidateOrdersAndDetermineFillPriceGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrNoMatch.Error(),
		},
	}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
}

func TestValidateOrdersAndDetermineFillPrice(t *testing.T) {
	oraclePrice := multiply1e6(big.NewInt(20))                                                  // $10
	spreadLimit := new(big.Int).Mul(big.NewInt(50), big.NewInt(1e4))                            // 50%
	upperbound := divide1e6(new(big.Int).Mul(oraclePrice, new(big.Int).Add(_1e6, spreadLimit))) // $10
	lowerbound := divide1e6(new(big.Int).Mul(oraclePrice, new(big.Int).Sub(_1e6, spreadLimit))) // $30
	Taker := uint8(0)
	Maker := uint8(1)

	t.Run("long order came first", func(t *testing.T) {
		blockPlaced0 := big.NewInt(69)
		blockPlaced1 := big.NewInt(70)
		t.Run("long price < lower bound", func(t *testing.T) {
			t.Run("short price < long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, multiply1e6(big.NewInt(9)), multiply1e6(big.NewInt(8)), blockPlaced0, blockPlaced1)
				assert.Nil(t, output)
				assert.Equal(t, ErrTooLow, err)
			})

			t.Run("short price == long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, multiply1e6(big.NewInt(7)), multiply1e6(big.NewInt(7)), blockPlaced0, blockPlaced1)
				assert.Nil(t, output)
				assert.Equal(t, ErrTooLow, err)
			})
		})

		t.Run("long price == lower bound", func(t *testing.T) {
			longPrice := lowerbound
			t.Run("short price < long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, new(big.Int).Sub(longPrice, big.NewInt(1)), blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{longPrice, Maker, Taker}, *output)
			})

			t.Run("short price == long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, longPrice, blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{longPrice, Maker, Taker}, *output)
			})
		})

		t.Run("lowerbound < long price < oracle", func(t *testing.T) {
			longPrice := multiply1e6(big.NewInt(15))
			t.Run("short price < long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, new(big.Int).Sub(longPrice, big.NewInt(1)), blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{longPrice, Maker, Taker}, *output)
			})

			t.Run("short price == long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, longPrice, blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{longPrice, Maker, Taker}, *output)
			})
		})

		t.Run("long price == oracle", func(t *testing.T) {
			longPrice := oraclePrice
			t.Run("short price < long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, new(big.Int).Sub(longPrice, big.NewInt(1)), blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{longPrice, Maker, Taker}, *output)
			})

			t.Run("short price == long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, longPrice, blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{longPrice, Maker, Taker}, *output)
			})
		})

		t.Run("oracle < long price < upper bound", func(t *testing.T) {
			longPrice := multiply1e6(big.NewInt(25))
			t.Run("short price < long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, new(big.Int).Sub(longPrice, big.NewInt(1)), blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{longPrice, Maker, Taker}, *output)
			})

			t.Run("short price == long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, longPrice, blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{longPrice, Maker, Taker}, *output)
			})
		})

		t.Run("long price == upper bound", func(t *testing.T) {
			longPrice := upperbound
			t.Run("short price < long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, new(big.Int).Sub(longPrice, big.NewInt(1)), blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{longPrice, Maker, Taker}, *output)
			})

			t.Run("short price == long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, longPrice, blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{longPrice, Maker, Taker}, *output)
			})
		})

		t.Run("upper bound < long price", func(t *testing.T) {
			longPrice := new(big.Int).Add(upperbound, big.NewInt(42))
			t.Run("upper < short price < long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, new(big.Int).Add(upperbound, big.NewInt(1)), blockPlaced0, blockPlaced1)
				assert.Nil(t, output)
				assert.Equal(t, ErrTooHigh, err)
			})

			t.Run("upper == short price < long price", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, upperbound, blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{upperbound, Maker, Taker}, *output)
			})

			t.Run("short price < upper", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, new(big.Int).Sub(upperbound, big.NewInt(1)), blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{upperbound, Maker, Taker}, *output)
			})

			t.Run("short price < lower", func(t *testing.T) {
				output, err := determineFillPrice(oraclePrice, spreadLimit, longPrice, new(big.Int).Sub(lowerbound, big.NewInt(1)), blockPlaced0, blockPlaced1)
				assert.Nil(t, err)
				assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{upperbound, Maker, Taker}, *output)
			})
		})
	})

	t.Run("both orders came in same block", func(t *testing.T) {
		blockPlaced0 := big.NewInt(70)
		blockPlaced1 := big.NewInt(69) // short order came first
		for i := 0; i < 2; i++ {
			if i == 1 {
				blockPlaced0 = blockPlaced1 // both orders came in same block
			}
			t.Run("short price < lower bound", func(t *testing.T) {
				t.Run("short price < long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, multiply1e6(big.NewInt(9)), multiply1e6(big.NewInt(8)), blockPlaced0, blockPlaced1)
					assert.Nil(t, output)
					assert.Equal(t, ErrTooLow, err)
				})

				t.Run("short price == long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, multiply1e6(big.NewInt(7)), multiply1e6(big.NewInt(7)), blockPlaced0, blockPlaced1)
					assert.Nil(t, output)
					assert.Equal(t, ErrTooLow, err)
				})
			})

			t.Run("short price == lower bound", func(t *testing.T) {
				shortPrice := lowerbound
				t.Run("short price < long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, new(big.Int).Add(shortPrice, big.NewInt(67)), shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, err)
					assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{shortPrice, Taker, Maker}, *output)
				})

				t.Run("short price == long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, shortPrice, shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, err)
					assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{shortPrice, Taker, Maker}, *output)
				})
			})

			t.Run("lowerbound < short price < oracle", func(t *testing.T) {
				shortPrice := multiply1e6(big.NewInt(15))
				t.Run("short price < long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, new(big.Int).Add(shortPrice, big.NewInt(58)), shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, err)
					assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{shortPrice, Taker, Maker}, *output)
				})

				t.Run("short price == long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, shortPrice, shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, err)
					assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{shortPrice, Taker, Maker}, *output)
				})
			})

			t.Run("short price == oracle", func(t *testing.T) {
				shortPrice := oraclePrice
				t.Run("short price < long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, new(big.Int).Add(shortPrice, big.NewInt(99)), shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, err)
					assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{shortPrice, Taker, Maker}, *output)
				})

				t.Run("short price == long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, shortPrice, shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, err)
					assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{shortPrice, Taker, Maker}, *output)
				})
			})

			t.Run("oracle < short price < upper bound", func(t *testing.T) {
				shortPrice := multiply1e6(big.NewInt(25))
				t.Run("short price < long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, new(big.Int).Add(shortPrice, big.NewInt(453)), shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, err)
					assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{shortPrice, Taker, Maker}, *output)
				})

				t.Run("short price == long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, shortPrice, shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, err)
					assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{shortPrice, Taker, Maker}, *output)
				})
			})

			t.Run("short price == upper bound", func(t *testing.T) {
				shortPrice := upperbound
				t.Run("short price < long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, new(big.Int).Add(shortPrice, big.NewInt(896)), shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, err)
					assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{shortPrice, Taker, Maker}, *output)
				})

				t.Run("short price == long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, shortPrice, shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, err)
					assert.Equal(t, ValidateOrdersAndDetermineFillPriceOutput{shortPrice, Taker, Maker}, *output)
				})
			})

			t.Run("upper bound < short price", func(t *testing.T) {
				shortPrice := new(big.Int).Add(upperbound, big.NewInt(42))
				t.Run("short price < long price", func(t *testing.T) {
					output, err := determineFillPrice(oraclePrice, spreadLimit, new(big.Int).Add(shortPrice, big.NewInt(896)), shortPrice, blockPlaced0, blockPlaced1)
					assert.Nil(t, output)
					assert.Equal(t, ErrTooHigh, err)
				})
			})
		}
	})
}
