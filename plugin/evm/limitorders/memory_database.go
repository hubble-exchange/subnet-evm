package limitorders

import (
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Market int

const (
	AvaxPerp Market = iota
	EthPerp
)

func GetActiveMarkets() []Market {
	return []Market{AvaxPerp}
}

type Collateral int

const (
	USDC Collateral = iota
	Avax
	Eth
)

var collateralWeightMap map[Collateral]float64 = map[Collateral]float64{USDC: 1, Avax: 0.8, Eth: 0.8}

type Status string

const (
	Unfulfilled = "unfulfilled"
	Fulfilled   = "fulfilled"
	Cancelled   = "cancelled"
)

type LimitOrder struct {
	id                      uint64
	Market                  Market
	PositionType            string
	UserAddress             string
	BaseAssetQuantity       *big.Int
	FilledBaseAssetQuantity *big.Int
	Price                   *big.Int
	Status                  Status
	Salt                    string
	Signature               []byte
	RawOrder                interface{}
	BlockNumber             *big.Int // block number order was placed on
	// RawSignature interface{}
}

func (order LimitOrder) GetUnFilledBaseAssetQuantity() *big.Int {
	return big.NewInt(0).Sub(order.BaseAssetQuantity, order.FilledBaseAssetQuantity)
}

type Position struct {
	OpenNotional        *big.Int
	Size                *big.Int
	UnrealisedFunding   *big.Int
	LastPremiumFraction *big.Int
}

type Trader struct {
	Positions   map[Market]*Position    // position for every market
	Margins     map[Collateral]*big.Int // available margin/balance for every market
	BlockNumber *big.Int
}

type LimitOrderDatabase interface {
	GetAllOrders() []LimitOrder
	Add(order *LimitOrder)
	UpdateFilledBaseAssetQuantity(quantity *big.Int, signature []byte)
	GetLongOrders(market Market) []LimitOrder
	GetShortOrders(market Market) []LimitOrder
	UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int)
	UpdateMargin(trader common.Address, collateral Collateral, addAmount *big.Int)
	UpdateUnrealisedFunding(market Market, cumulativePremiumFraction *big.Int)
	ResetUnrealisedFunding(market Market, trader common.Address, cumulativePremiumFraction *big.Int)
	UpdateNextFundingTime(nextFundingTime uint64)
	GetNextFundingTime() uint64
	GetLiquidableTraders(market Market, oraclePrice *big.Int) (longPositions []Liquidable, shortPositions []Liquidable)
	UpdateLastPrice(market Market, lastPrice *big.Int)
	GetLastPrice(market Market) *big.Int
}

type InMemoryDatabase struct {
	orderMap        map[string]*LimitOrder     // signature => order
	traderMap       map[common.Address]*Trader // address => trader info
	nextFundingTime uint64
	lastPrice       map[Market]*big.Int
}

func NewInMemoryDatabase() *InMemoryDatabase {
	orderMap := map[string]*LimitOrder{}
	lastPrice := map[Market]*big.Int{AvaxPerp: big.NewInt(0)}
	traderMap := map[common.Address]*Trader{}
	nextFundingTime := uint64(getNextHour().Unix())

	return &InMemoryDatabase{
		orderMap:        orderMap,
		traderMap:       traderMap,
		nextFundingTime: nextFundingTime,
		lastPrice:       lastPrice,
	}
}

func (db *InMemoryDatabase) GetAllOrders() []LimitOrder {
	allOrders := []LimitOrder{}
	for _, order := range db.orderMap {
		allOrders = append(allOrders, *order)
	}
	return allOrders
}

func (db *InMemoryDatabase) Add(order *LimitOrder) {
	db.orderMap[string(order.Signature)] = order
}

func (db *InMemoryDatabase) UpdateFilledBaseAssetQuantity(quantity *big.Int, signature []byte) {
	limitOrder := db.orderMap[string(signature)]
	if limitOrder.BaseAssetQuantity.Cmp(quantity) == 0 {
		deleteOrder(db, signature)
		return
	} else {
		if limitOrder.PositionType == "long" {
			limitOrder.FilledBaseAssetQuantity.Add(limitOrder.FilledBaseAssetQuantity, quantity) // filled = filled + quantity
		}
		if limitOrder.PositionType == "short" {
			limitOrder.FilledBaseAssetQuantity.Sub(limitOrder.FilledBaseAssetQuantity, quantity) // filled = filled - quantity
		}
	}
}

func (db *InMemoryDatabase) GetNextFundingTime() uint64 {
	return db.nextFundingTime
}

func (db *InMemoryDatabase) UpdateNextFundingTime(nextFundingTime uint64) {
	db.nextFundingTime = nextFundingTime
}

func (db *InMemoryDatabase) GetLongOrders(market Market) []LimitOrder {
	var longOrders []LimitOrder
	for _, order := range db.orderMap {
		if order.PositionType == "long" && order.Market == market {
			longOrders = append(longOrders, *order)
		}
	}
	sortLongOrders(longOrders)
	return longOrders
}

func (db *InMemoryDatabase) GetShortOrders(market Market) []LimitOrder {
	var shortOrders []LimitOrder
	for _, order := range db.orderMap {
		if order.PositionType == "short" && order.Market == market {
			shortOrders = append(shortOrders, *order)
		}
	}
	sortShortOrders(shortOrders)
	return shortOrders
}

func (db *InMemoryDatabase) UpdateMargin(trader common.Address, collateral Collateral, addAmount *big.Int) {
	if _, ok := db.traderMap[trader]; !ok {
		db.traderMap[trader] = &Trader{
			Positions: map[Market]*Position{},
			Margins:   map[Collateral]*big.Int{},
		}
	}

	if _, ok := db.traderMap[trader].Margins[collateral]; !ok {
		db.traderMap[trader].Margins[collateral] = big.NewInt(0)
	}

	db.traderMap[trader].Margins[collateral].Add(db.traderMap[trader].Margins[collateral], addAmount)
}

func (db *InMemoryDatabase) UpdatePosition(trader common.Address, market Market, size *big.Int, openNotional *big.Int) {
	if _, ok := db.traderMap[trader]; !ok {
		db.traderMap[trader] = &Trader{
			Positions: map[Market]*Position{},
			Margins:   map[Collateral]*big.Int{},
		}
	}

	if _, ok := db.traderMap[trader].Positions[market]; !ok {
		db.traderMap[trader].Positions[market] = &Position{}
	}

	db.traderMap[trader].Positions[market].Size = size
	db.traderMap[trader].Positions[market].OpenNotional = openNotional
}

func (db *InMemoryDatabase) UpdateUnrealisedFunding(market Market, cumulativePremiumFraction *big.Int) {
	for _, trader := range db.traderMap {
		position := trader.Positions[market]
		if position != nil {
			position.UnrealisedFunding = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(0).Sub(cumulativePremiumFraction, position.LastPremiumFraction), position.Size), BASE_PRECISION)
		}
	}
}

func (db *InMemoryDatabase) ResetUnrealisedFunding(market Market, trader common.Address, cumulativePremiumFraction *big.Int) {
	if db.traderMap[trader] != nil {
		if _, ok := db.traderMap[trader].Positions[market]; ok {
			db.traderMap[trader].Positions[market].UnrealisedFunding = big.NewInt(0)
			db.traderMap[trader].Positions[market].LastPremiumFraction = cumulativePremiumFraction
		}
	}
}

func (db *InMemoryDatabase) GetAllTraders() map[common.Address]*Trader {
	return db.traderMap
}

func (db *InMemoryDatabase) UpdateLastPrice(market Market, lastPrice *big.Int) {
	db.lastPrice[market] = lastPrice
}

func (db *InMemoryDatabase) GetLastPrice(market Market) *big.Int {
	return db.lastPrice[market]
}

func (trader *Trader) GetNormalisedMargin() *big.Int {
	return trader.Margins[USDC]

	// this will change after multi collateral
	// var normalisedMargin float64
	// for coll, margin := range trader.Margins {
	// 	normalisedMargin += margin * priceMap[coll] * collateralWeightMap[coll]
	// }

	// return normalisedMargin
}

func sortLongOrders(orders []LimitOrder) []LimitOrder {
	sort.SliceStable(orders, func(i, j int) bool {
		if orders[i].Price.Cmp(orders[j].Price) == 1 {
			return true
		}
		if orders[i].Price.Cmp(orders[j].Price) == 0 {
			if orders[i].BlockNumber.Cmp(orders[j].BlockNumber) == -1 {
				return true
			}
		}
		return false
	})
	return orders
}

func sortShortOrders(orders []LimitOrder) []LimitOrder {
	sort.SliceStable(orders, func(i, j int) bool {
		if orders[i].Price.Cmp(orders[j].Price) == -1 {
			return true
		}
		if orders[i].Price.Cmp(orders[j].Price) == 0 {
			if orders[i].BlockNumber.Cmp(orders[j].BlockNumber) == -1 {
				return true
			}
		}
		return false
	})
	return orders
}

func getNextHour() time.Time {
	now := time.Now().UTC()
	nextHour := now.Round(time.Hour)
	if time.Since(nextHour) >= 0 {
		nextHour = nextHour.Add(time.Hour)
	}
	return nextHour
}

func deleteOrder(db *InMemoryDatabase, signature []byte) {
	delete(db.orderMap, string(signature))
}
