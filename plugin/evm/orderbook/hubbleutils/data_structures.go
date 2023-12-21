package hubbleutils

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type MarginMode = uint8

const (
	Maintenance_Margin MarginMode = iota
	Min_Allowable_Margin
)

type Collateral struct {
	Price    *big.Int // scaled by 1e6
	Weight   *big.Int // scaled by 1e6
	Decimals uint8
}

type Market = int

type Position struct {
	OpenNotional *big.Int `json:"open_notional"`
	Size         *big.Int `json:"size"`
}

type Trader struct {
	Positions map[Market]*Position `json:"positions"` // position for every market
	Margin    Margin               `json:"margin"`    // available margin/balance for every market
}

type Margin struct {
	Reserved  *big.Int                `json:"reserved"`
	Deposited map[Collateral]*big.Int `json:"deposited"`
}

type Side uint8

const (
	Long Side = iota
	Short
	Liquidation
)

type OrderStatus uint8

// has to be exact same as IOrderHandler
const (
	Invalid OrderStatus = iota
	Placed
	Filled
	Cancelled
)

type OrderType uint8

const (
	Limit OrderType = iota
	IOC
	Signed
)

func (o OrderType) String() string {
	return [...]string{"limit", "ioc"}[o]
}

type BaseOrder struct {
	AmmIndex          *big.Int       `json:"ammIndex"`
	Trader            common.Address `json:"trader"`
	BaseAssetQuantity *big.Int       `json:"baseAssetQuantity"`
	Price             *big.Int       `json:"price"`
	Salt              *big.Int       `json:"salt"`
	ReduceOnly        bool           `json:"reduceOnly"`
}

// LimitOrder type is copy of Order struct defined in LimitOrderbook contract
type LimitOrder struct {
	BaseOrder
	PostOnly bool `json:"postOnly"`
}

// IOCOrder type is copy of IOCOrder struct defined in Orderbook contract
type IOCOrder struct {
	BaseOrder
	OrderType uint8    `json:"orderType"`
	ExpireAt  *big.Int `json:"expireAt"`
}

// IOCOrder type is copy of IOCOrder struct defined in Orderbook contract
type SignedOrder struct {
	LimitOrder
	OrderType uint8    `json:"orderType"`
	ExpireAt  *big.Int `json:"expireAt"`
	Sig       []byte   `json:"sig"`
}

// LimitOrder
func (order *LimitOrder) EncodeToABIWithoutType() ([]byte, error) {
	limitOrderType, err := getOrderType("limit")
	if err != nil {
		return nil, err
	}
	encodedLimitOrder, err := abi.Arguments{{Type: limitOrderType}}.Pack(order)
	if err != nil {
		return nil, err
	}
	return encodedLimitOrder, nil
}

func (order *LimitOrder) EncodeToABI() ([]byte, error) {
	encodedLimitOrder, err := order.EncodeToABIWithoutType()
	if err != nil {
		return nil, fmt.Errorf("limit order packing failed: %w", err)
	}
	orderType, _ := abi.NewType("uint8", "uint8", nil)
	orderBytesType, _ := abi.NewType("bytes", "bytes", nil)
	// 0 means ordertype = limit order
	encodedOrder, err := abi.Arguments{{Type: orderType}, {Type: orderBytesType}}.Pack(uint8(0) /* Limit Order */, encodedLimitOrder)
	if err != nil {
		return nil, fmt.Errorf("order encoding failed: %w", err)
	}
	return encodedOrder, nil
}

func (order *LimitOrder) DecodeFromRawOrder(rawOrder interface{}) {
	marshalledOrder, _ := json.Marshal(rawOrder)
	json.Unmarshal(marshalledOrder, &order)
}

func (order *LimitOrder) Map() map[string]interface{} {
	return map[string]interface{}{
		"ammIndex":          order.AmmIndex,
		"trader":            order.Trader,
		"baseAssetQuantity": utils.BigIntToFloat(order.BaseAssetQuantity, 18),
		"price":             utils.BigIntToFloat(order.Price, 6),
		"reduceOnly":        order.ReduceOnly,
		"postOnly":          order.PostOnly,
		"salt":              order.Salt,
	}
}

func DecodeLimitOrder(encodedOrder []byte) (*LimitOrder, error) {
	limitOrderType, err := getOrderType("limit")
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}
	order, err := abi.Arguments{{Type: limitOrderType}}.Unpack(encodedOrder)
	if err != nil {
		return nil, err
	}
	limitOrder := &LimitOrder{}
	limitOrder.DecodeFromRawOrder(order[0])
	return limitOrder, nil
}

func (order *LimitOrder) Hash() (common.Hash, error) {
	data, err := order.EncodeToABIWithoutType()
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(crypto.Keccak256(data)), nil
}

// ----------------------------------------------------------------------------
// IOCOrder

func (order *IOCOrder) EncodeToABIWithoutType() ([]byte, error) {
	iocOrderType, err := getOrderType("ioc")
	if err != nil {
		return nil, err
	}
	encodedOrder, err := abi.Arguments{{Type: iocOrderType}}.Pack(order)
	if err != nil {
		return nil, err
	}
	return encodedOrder, nil
}

func (order *IOCOrder) EncodeToABI() ([]byte, error) {
	encodedIOCOrder, err := order.EncodeToABIWithoutType()
	if err != nil {
		return nil, fmt.Errorf("limit order packing failed: %w", err)
	}

	orderType, _ := abi.NewType("uint8", "uint8", nil)
	orderBytesType, _ := abi.NewType("bytes", "bytes", nil)
	// 1 means ordertype = IOC/market order
	encodedOrder, err := abi.Arguments{{Type: orderType}, {Type: orderBytesType}}.Pack(uint8(IOC), encodedIOCOrder)
	if err != nil {
		return nil, fmt.Errorf("order encoding failed: %w", err)
	}
	return encodedOrder, nil
}

func (order *IOCOrder) DecodeFromRawOrder(rawOrder interface{}) {
	marshalledOrder, _ := json.Marshal(rawOrder)
	json.Unmarshal(marshalledOrder, &order)
}

func (order *IOCOrder) Map() map[string]interface{} {
	return map[string]interface{}{
		"ammIndex":          order.AmmIndex,
		"trader":            order.Trader,
		"baseAssetQuantity": utils.BigIntToFloat(order.BaseAssetQuantity, 18),
		"price":             utils.BigIntToFloat(order.Price, 6),
		"reduceOnly":        order.ReduceOnly,
		"salt":              order.Salt,
		"orderType":         order.OrderType,
		"expireAt":          order.ExpireAt,
	}
}

func DecodeIOCOrder(encodedOrder []byte) (*IOCOrder, error) {
	iocOrderType, err := getOrderType("ioc")
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}
	order, err := abi.Arguments{{Type: iocOrderType}}.Unpack(encodedOrder)
	if err != nil {
		return nil, err
	}
	iocOrder := &IOCOrder{}
	iocOrder.DecodeFromRawOrder(order[0])
	return iocOrder, nil
}

func (order *IOCOrder) Hash() (hash common.Hash, err error) {
	data, err := order.EncodeToABIWithoutType()
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(crypto.Keccak256(data)), nil
}

// ----------------------------------------------------------------------------
// SignedOrder

func (order *SignedOrder) EncodeToABIWithoutType() ([]byte, error) {
	signedOrderType, err := getOrderType("signed")
	if err != nil {
		return nil, err
	}
	encodedOrder, err := abi.Arguments{{Type: signedOrderType}}.Pack(order)
	if err != nil {
		return nil, err
	}
	return encodedOrder, nil
}

func (order *SignedOrder) EncodeToABI() ([]byte, error) {
	encodedSignedOrder, err := order.EncodeToABIWithoutType()
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}

	orderType, _ := abi.NewType("uint8", "uint8", nil)
	orderBytesType, _ := abi.NewType("bytes", "bytes", nil)
	// 1 means ordertype = IOC/market order
	encodedOrder, err := abi.Arguments{{Type: orderType}, {Type: orderBytesType}}.Pack(uint8(Signed), encodedSignedOrder)
	if err != nil {
		return nil, fmt.Errorf("order encoding failed: %w", err)
	}

	return encodedOrder, nil
}

func (o *SignedOrder) UnmarshalJSON(data []byte) error {
	// Alias types to avoid recursive call to UnmarshalJSON
	// type AliasSignedOrder SignedOrder
	// type AliasBaseOrder BaseOrder
	// type AliasLimitOrder LimitOrder

	// Redefine the structs with simple types for JSON unmarshalling
	aux := &struct {
		AmmIndex          uint64         `json:"ammIndex"`
		Trader            common.Address `json:"trader"`
		BaseAssetQuantity string         `json:"baseAssetQuantity"`
		Price             string         `json:"price"`
		Salt              string         `json:"salt"`
		ReduceOnly        bool           `json:"reduceOnly"`
		PostOnly          bool           `json:"postOnly"`
		OrderType         uint8          `json:"orderType"`
		ExpireAt          uint64         `json:"expireAt"`
		Sig               string         `json:"sig"`
	}{}

	// Perform the unmarshalling
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Convert and assign the values to the original struct
	o.AmmIndex = new(big.Int).SetUint64(aux.AmmIndex)

	o.Trader = aux.Trader

	o.BaseAssetQuantity = new(big.Int)
	o.BaseAssetQuantity.SetString(aux.BaseAssetQuantity, 10)

	o.Price = new(big.Int)
	o.Price.SetString(aux.Price, 10)

	o.Salt = new(big.Int)
	o.Salt.SetBytes(common.FromHex(aux.Salt))

	o.ReduceOnly = aux.ReduceOnly
	o.PostOnly = aux.PostOnly
	o.OrderType = aux.OrderType

	o.ExpireAt = new(big.Int).SetUint64(aux.ExpireAt)
	o.Sig = common.FromHex(aux.Sig)
	return nil
}

func (order *SignedOrder) DecodeFromRawOrder(rawOrder interface{}) {
	// marshalledOrder, err := json.Marshal(rawOrder)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	order_, ok := rawOrder.(string)
	fmt.Println("ok", ok)

	orderJson := []byte(order_)
	err := json.Unmarshal(orderJson, &order)
	if err != nil {
		fmt.Println("err in DecodeFromRawOrder")
		fmt.Println(err)
	}
	// json.Unmarshal(marshalledOrder, &order)
}

func (order *SignedOrder) Map() map[string]interface{} {
	return map[string]interface{}{
		"ammIndex":          order.AmmIndex,
		"trader":            order.Trader,
		"baseAssetQuantity": utils.BigIntToFloat(order.BaseAssetQuantity, 18),
		"price":             utils.BigIntToFloat(order.Price, 6),
		"reduceOnly":        order.ReduceOnly,
		"postOnly":          order.PostOnly,
		"salt":              order.Salt,
		"orderType":         order.OrderType,
		"expireAt":          order.ExpireAt,
		"sig":               order.Sig,
	}
}

func DecodeSignedOrder(encodedOrder []byte) (*SignedOrder, error) {
	orderType, err := getOrderType("signed")
	if err != nil {
		return nil, fmt.Errorf("failed getting abi type: %w", err)
	}
	order, err := abi.Arguments{{Type: orderType}}.Unpack(encodedOrder)
	if err != nil {
		return nil, err
	}
	signedOrder := &SignedOrder{}
	signedOrder.DecodeFromRawOrder(order[0])
	return signedOrder, nil
}

// ----------------------------------------------------------------------------
// Helper functions
type DecodeStep struct {
	OrderType    OrderType
	EncodedOrder []byte
}

func DecodeTypeAndEncodedOrder(data []byte) (*DecodeStep, error) {
	orderType, _ := abi.NewType("uint8", "uint8", nil)
	orderBytesType, _ := abi.NewType("bytes", "bytes", nil)
	decodedValues, err := abi.Arguments{{Type: orderType}, {Type: orderBytesType}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return &DecodeStep{
		OrderType:    OrderType(decodedValues[0].(uint8)),
		EncodedOrder: decodedValues[1].([]byte),
	}, nil
}

func getOrderType(orderType string) (abi.Type, error) {
	if orderType == "limit" {
		return abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "ammIndex", Type: "uint256"},
			{Name: "trader", Type: "address"},
			{Name: "baseAssetQuantity", Type: "int256"},
			{Name: "price", Type: "uint256"},
			{Name: "salt", Type: "uint256"},
			{Name: "reduceOnly", Type: "bool"},
			{Name: "postOnly", Type: "bool"},
		})
	}
	if orderType == "ioc" {
		return abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "orderType", Type: "uint8"},
			{Name: "expireAt", Type: "uint256"},
			{Name: "ammIndex", Type: "uint256"},
			{Name: "trader", Type: "address"},
			{Name: "baseAssetQuantity", Type: "int256"},
			{Name: "price", Type: "uint256"},
			{Name: "salt", Type: "uint256"},
			{Name: "reduceOnly", Type: "bool"},
		})
	}
	if orderType == "signed" {
		return abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "orderType", Type: "uint8"},
			{Name: "expireAt", Type: "uint256"},
			{Name: "ammIndex", Type: "uint256"},
			{Name: "trader", Type: "address"},
			{Name: "baseAssetQuantity", Type: "int256"},
			{Name: "price", Type: "uint256"},
			{Name: "salt", Type: "uint256"},
			{Name: "reduceOnly", Type: "bool"},
			{Name: "postOnly", Type: "bool"},
		})
	}
	return abi.Type{}, fmt.Errorf("invalid order type")
}
