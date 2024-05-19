// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package juror

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"

	_ "embed"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// Gas costs for each function. These are set to 1 by default.
	// You should set a gas cost for each function in your contract.
	// Generally, you should not set gas costs very low as this may cause your network to be vulnerable to DoS attacks.
	// There are some predefined gas costs in contract/utils.go that you can use.
	GetNotionalPositionAndMarginGasCost                  uint64 = 69
	ValidateCancelLimitOrderGasCost                      uint64 = 69
	ValidateLiquidationOrderAndDetermineFillPriceGasCost uint64 = 69
	ValidateOrdersAndDetermineFillPriceGasCost           uint64 = 69
	ValidatePlaceIOCOrderGasCost                         uint64 = 69
	ValidatePlaceLimitOrderGasCost                       uint64 = 69
)

// CUSTOM CODE STARTS HERE
// Reference imports to suppress errors from unused imports. This code and any unnecessary imports can be removed.
var (
	_ = abi.JSON
	_ = errors.New
	_ = big.NewInt
)

// Singleton StatefulPrecompiledContract and signatures.
var (

	// JurorRawABI contains the raw ABI of Juror contract.
	//go:embed contract.abi
	JurorRawABI string

	JurorABI = contract.ParseABI(JurorRawABI)

	JurorPrecompile = createJurorPrecompile()
)

// IClearingHouseInstruction is an auto generated low-level Go binding around an user-defined struct.
type IClearingHouseInstruction struct {
	AmmIndex  *big.Int
	Trader    common.Address
	OrderHash [32]byte
	Mode      uint8
}

// IImmediateOrCancelOrdersOrder is an auto generated low-level Go binding around an user-defined struct.
type IImmediateOrCancelOrdersOrder struct {
	OrderType         uint8
	ExpireAt          *big.Int
	AmmIndex          *big.Int
	Trader            common.Address
	BaseAssetQuantity *big.Int
	Price             *big.Int
	Salt              *big.Int
	ReduceOnly        bool
}

// ILimitOrderBookOrder is an auto generated low-level Go binding around an user-defined struct.
type ILimitOrderBookOrder struct {
	AmmIndex          *big.Int
	Trader            common.Address
	BaseAssetQuantity *big.Int
	Price             *big.Int
	Salt              *big.Int
	ReduceOnly        bool
	PostOnly          bool
}

// IOrderHandlerCancelOrderRes is an auto generated low-level Go binding around an user-defined struct.
type IOrderHandlerCancelOrderRes struct {
	UnfilledAmount *big.Int
	Amm            common.Address
}

// IOrderHandlerLiquidationMatchingValidationRes is an auto generated low-level Go binding around an user-defined struct.
type IOrderHandlerLiquidationMatchingValidationRes struct {
	Instruction  IClearingHouseInstruction
	OrderType    uint8
	EncodedOrder []byte
	FillPrice    *big.Int
	FillAmount   *big.Int
}

// IOrderHandlerMatchingValidationRes is an auto generated low-level Go binding around an user-defined struct.
type IOrderHandlerMatchingValidationRes struct {
	Instructions  [2]IClearingHouseInstruction
	OrderTypes    [2]uint8
	EncodedOrders [2][]byte
	FillPrice     *big.Int
}

// IOrderHandlerPlaceOrderRes is an auto generated low-level Go binding around an user-defined struct.
type IOrderHandlerPlaceOrderRes struct {
	ReserveAmount *big.Int
	Amm           common.Address
}

type GetNotionalPositionAndMarginInput struct {
	Trader                 common.Address
	IncludeFundingPayments bool
	Mode                   uint8
}

type GetNotionalPositionAndMarginOutput struct {
	NotionalPosition *big.Int
	Margin           *big.Int
}

type ValidateCancelLimitOrderInput struct {
	Order           ILimitOrderBookOrder
	Sender          common.Address
	AssertLowMargin bool
}

type ValidateCancelLimitOrderOutput struct {
	Err       string
	OrderHash [32]byte
	Res       IOrderHandlerCancelOrderRes
}

type ValidateLiquidationOrderAndDetermineFillPriceInput struct {
	Data              []byte
	LiquidationAmount *big.Int
}

type ValidateLiquidationOrderAndDetermineFillPriceOutput struct {
	Err     string
	Element uint8
	Res     IOrderHandlerLiquidationMatchingValidationRes
}

type ValidateOrdersAndDetermineFillPriceInput struct {
	Data       [2][]byte
	FillAmount *big.Int
}

type ValidateOrdersAndDetermineFillPriceOutput struct {
	Err     string
	Element uint8
	Res     IOrderHandlerMatchingValidationRes
}

type ValidatePlaceIOCOrderInput struct {
	Order  IImmediateOrCancelOrdersOrder
	Sender common.Address
}

type ValidatePlaceIOCOrderOutput struct {
	Err       string
	OrderHash [32]byte
}

type ValidatePlaceLimitOrderInput struct {
	Order  ILimitOrderBookOrder
	Sender common.Address
}

type ValidatePlaceLimitOrderOutput struct {
	Err       string
	Orderhash [32]byte
	Res       IOrderHandlerPlaceOrderRes
}

// UnpackGetNotionalPositionAndMarginInput attempts to unpack [input] as GetNotionalPositionAndMarginInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackGetNotionalPositionAndMarginInput(input []byte) (GetNotionalPositionAndMarginInput, error) {
	inputStruct := GetNotionalPositionAndMarginInput{}
	err := JurorABI.UnpackInputIntoInterface(&inputStruct, "getNotionalPositionAndMargin", input, true)

	return inputStruct, err
}

// PackGetNotionalPositionAndMargin packs [inputStruct] of type GetNotionalPositionAndMarginInput into the appropriate arguments for getNotionalPositionAndMargin.
func PackGetNotionalPositionAndMargin(inputStruct GetNotionalPositionAndMarginInput) ([]byte, error) {
	return JurorABI.Pack("getNotionalPositionAndMargin", inputStruct.Trader, inputStruct.IncludeFundingPayments, inputStruct.Mode)
}

// PackGetNotionalPositionAndMarginOutput attempts to pack given [outputStruct] of type GetNotionalPositionAndMarginOutput
// to conform the ABI outputs.
func PackGetNotionalPositionAndMarginOutput(outputStruct GetNotionalPositionAndMarginOutput) ([]byte, error) {
	return JurorABI.PackOutput("getNotionalPositionAndMargin",
		outputStruct.NotionalPosition,
		outputStruct.Margin,
	)
}

func getNotionalPositionAndMargin(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, GetNotionalPositionAndMarginGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the GetNotionalPositionAndMarginInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackGetNotionalPositionAndMarginInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	bibliophile := bibliophile.NewBibliophileClient(accessibleState)
	output := GetNotionalPositionAndMargin(bibliophile, &inputStruct)
	packedOutput, err := PackGetNotionalPositionAndMarginOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackValidateCancelLimitOrderInput attempts to unpack [input] as ValidateCancelLimitOrderInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackValidateCancelLimitOrderInput(input []byte) (ValidateCancelLimitOrderInput, error) {
	inputStruct := ValidateCancelLimitOrderInput{}
	err := JurorABI.UnpackInputIntoInterface(&inputStruct, "validateCancelLimitOrder", input, true)

	return inputStruct, err
}

// PackValidateCancelLimitOrder packs [inputStruct] of type ValidateCancelLimitOrderInput into the appropriate arguments for validateCancelLimitOrder.
func PackValidateCancelLimitOrder(inputStruct ValidateCancelLimitOrderInput) ([]byte, error) {
	return JurorABI.Pack("validateCancelLimitOrder", inputStruct.Order, inputStruct.Sender, inputStruct.AssertLowMargin)
}

// PackValidateCancelLimitOrderOutput attempts to pack given [outputStruct] of type ValidateCancelLimitOrderOutput
// to conform the ABI outputs.
func PackValidateCancelLimitOrderOutput(outputStruct ValidateCancelLimitOrderOutput) ([]byte, error) {
	return JurorABI.PackOutput("validateCancelLimitOrder",
		outputStruct.Err,
		outputStruct.OrderHash,
		outputStruct.Res,
	)
}

func validateCancelLimitOrder(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ValidateCancelLimitOrderGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the ValidateCancelLimitOrderInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackValidateCancelLimitOrderInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	bibliophile := bibliophile.NewBibliophileClient(accessibleState)
	output := ValidateCancelLimitOrder(bibliophile, &inputStruct)
	packedOutput, err := PackValidateCancelLimitOrderOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackValidateLiquidationOrderAndDetermineFillPriceInput attempts to unpack [input] as ValidateLiquidationOrderAndDetermineFillPriceInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackValidateLiquidationOrderAndDetermineFillPriceInput(input []byte) (ValidateLiquidationOrderAndDetermineFillPriceInput, error) {
	inputStruct := ValidateLiquidationOrderAndDetermineFillPriceInput{}
	err := JurorABI.UnpackInputIntoInterface(&inputStruct, "validateLiquidationOrderAndDetermineFillPrice", input, true)

	return inputStruct, err
}

// PackValidateLiquidationOrderAndDetermineFillPrice packs [inputStruct] of type ValidateLiquidationOrderAndDetermineFillPriceInput into the appropriate arguments for validateLiquidationOrderAndDetermineFillPrice.
func PackValidateLiquidationOrderAndDetermineFillPrice(inputStruct ValidateLiquidationOrderAndDetermineFillPriceInput) ([]byte, error) {
	return JurorABI.Pack("validateLiquidationOrderAndDetermineFillPrice", inputStruct.Data, inputStruct.LiquidationAmount)
}

// PackValidateLiquidationOrderAndDetermineFillPriceOutput attempts to pack given [outputStruct] of type ValidateLiquidationOrderAndDetermineFillPriceOutput
// to conform the ABI outputs.
func PackValidateLiquidationOrderAndDetermineFillPriceOutput(outputStruct ValidateLiquidationOrderAndDetermineFillPriceOutput) ([]byte, error) {
	return JurorABI.PackOutput("validateLiquidationOrderAndDetermineFillPrice",
		outputStruct.Err,
		outputStruct.Element,
		outputStruct.Res,
	)
}

func validateLiquidationOrderAndDetermineFillPrice(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ValidateLiquidationOrderAndDetermineFillPriceGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the ValidateLiquidationOrderAndDetermineFillPriceInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackValidateLiquidationOrderAndDetermineFillPriceInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	bibliophile := bibliophile.NewBibliophileClient(accessibleState)
	output := ValidateLiquidationOrderAndDetermineFillPrice(bibliophile, &inputStruct)
	packedOutput, err := PackValidateLiquidationOrderAndDetermineFillPriceOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackValidateOrdersAndDetermineFillPriceInput attempts to unpack [input] as ValidateOrdersAndDetermineFillPriceInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackValidateOrdersAndDetermineFillPriceInput(input []byte) (ValidateOrdersAndDetermineFillPriceInput, error) {
	inputStruct := ValidateOrdersAndDetermineFillPriceInput{}
	err := JurorABI.UnpackInputIntoInterface(&inputStruct, "validateOrdersAndDetermineFillPrice", input, true)

	return inputStruct, err
}

// PackValidateOrdersAndDetermineFillPrice packs [inputStruct] of type ValidateOrdersAndDetermineFillPriceInput into the appropriate arguments for validateOrdersAndDetermineFillPrice.
func PackValidateOrdersAndDetermineFillPrice(inputStruct ValidateOrdersAndDetermineFillPriceInput) ([]byte, error) {
	return JurorABI.Pack("validateOrdersAndDetermineFillPrice", inputStruct.Data, inputStruct.FillAmount)
}

// PackValidateOrdersAndDetermineFillPriceOutput attempts to pack given [outputStruct] of type ValidateOrdersAndDetermineFillPriceOutput
// to conform the ABI outputs.
func PackValidateOrdersAndDetermineFillPriceOutput(outputStruct ValidateOrdersAndDetermineFillPriceOutput) ([]byte, error) {
	return JurorABI.PackOutput("validateOrdersAndDetermineFillPrice",
		outputStruct.Err,
		outputStruct.Element,
		outputStruct.Res,
	)
}

func validateOrdersAndDetermineFillPrice(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ValidateOrdersAndDetermineFillPriceGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the ValidateOrdersAndDetermineFillPriceInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackValidateOrdersAndDetermineFillPriceInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	bibliophile := bibliophile.NewBibliophileClient(accessibleState)
	output := ValidateOrdersAndDetermineFillPrice(bibliophile, &inputStruct)
	packedOutput, err := PackValidateOrdersAndDetermineFillPriceOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackValidatePlaceIOCOrderInput attempts to unpack [input] as ValidatePlaceIOCOrderInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackValidatePlaceIOCOrderInput(input []byte) (ValidatePlaceIOCOrderInput, error) {
	inputStruct := ValidatePlaceIOCOrderInput{}
	err := JurorABI.UnpackInputIntoInterface(&inputStruct, "validatePlaceIOCOrder", input, true)

	return inputStruct, err
}

// PackValidatePlaceIOCOrder packs [inputStruct] of type ValidatePlaceIOCOrderInput into the appropriate arguments for validatePlaceIOCOrder.
func PackValidatePlaceIOCOrder(inputStruct ValidatePlaceIOCOrderInput) ([]byte, error) {
	return JurorABI.Pack("validatePlaceIOCOrder", inputStruct.Order, inputStruct.Sender)
}

// PackValidatePlaceIOCOrderOutput attempts to pack given [outputStruct] of type ValidatePlaceIOCOrderOutput
// to conform the ABI outputs.
func PackValidatePlaceIOCOrderOutput(outputStruct ValidatePlaceIOCOrderOutput) ([]byte, error) {
	return JurorABI.PackOutput("validatePlaceIOCOrder",
		outputStruct.Err,
		outputStruct.OrderHash,
	)
}

func validatePlaceIOCOrder(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ValidatePlaceIOCOrderGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the ValidatePlaceIOCOrderInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackValidatePlaceIOCOrderInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	bibliophile := bibliophile.NewBibliophileClient(accessibleState)
	output := ValidatePlaceIOCorder(bibliophile, &inputStruct)
	packedOutput, err := PackValidatePlaceIOCOrderOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackValidatePlaceLimitOrderInput attempts to unpack [input] as ValidatePlaceLimitOrderInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackValidatePlaceLimitOrderInput(input []byte) (ValidatePlaceLimitOrderInput, error) {
	inputStruct := ValidatePlaceLimitOrderInput{}
	err := JurorABI.UnpackInputIntoInterface(&inputStruct, "validatePlaceLimitOrder", input, true)

	return inputStruct, err
}

// PackValidatePlaceLimitOrder packs [inputStruct] of type ValidatePlaceLimitOrderInput into the appropriate arguments for validatePlaceLimitOrder.
func PackValidatePlaceLimitOrder(inputStruct ValidatePlaceLimitOrderInput) ([]byte, error) {
	return JurorABI.Pack("validatePlaceLimitOrder", inputStruct.Order, inputStruct.Sender)
}

// PackValidatePlaceLimitOrderOutput attempts to pack given [outputStruct] of type ValidatePlaceLimitOrderOutput
// to conform the ABI outputs.
func PackValidatePlaceLimitOrderOutput(outputStruct ValidatePlaceLimitOrderOutput) ([]byte, error) {
	return JurorABI.PackOutput("validatePlaceLimitOrder",
		outputStruct.Err,
		outputStruct.Orderhash,
		outputStruct.Res,
	)
}

func validatePlaceLimitOrder(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ValidatePlaceLimitOrderGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the ValidatePlaceLimitOrderInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackValidatePlaceLimitOrderInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	bibliophile := bibliophile.NewBibliophileClient(accessibleState)
	output := ValidatePlaceLimitOrder(bibliophile, &inputStruct)
	packedOutput, err := PackValidatePlaceLimitOrderOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// createJurorPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.

func createJurorPrecompile() contract.StatefulPrecompiledContract {
	var functions []*contract.StatefulPrecompileFunction

	abiFunctionMap := map[string]contract.RunStatefulPrecompileFunc{
		"getNotionalPositionAndMargin":                  getNotionalPositionAndMargin,
		"validateCancelLimitOrder":                      validateCancelLimitOrder,
		"validateLiquidationOrderAndDetermineFillPrice": validateLiquidationOrderAndDetermineFillPrice,
		"validateOrdersAndDetermineFillPrice":           validateOrdersAndDetermineFillPrice,
		"validatePlaceIOCOrder":                         validatePlaceIOCOrder,
		"validatePlaceLimitOrder":                       validatePlaceLimitOrder,
	}

	for name, function := range abiFunctionMap {
		method, ok := JurorABI.Methods[name]
		if !ok {
			panic(fmt.Errorf("given method (%s) does not exist in the ABI", name))
		}
		functions = append(functions, contract.NewStatefulPrecompileFunction(method.ID, function))
	}
	// Construct the contract with no fallback function.
	statefulContract, err := contract.NewStatefulPrecompileContract(nil, functions)
	if err != nil {
		panic(err)
	}
	return statefulContract
}
