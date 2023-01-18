package limitorders

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

var maintenanceMargin = big.NewInt(1e5)
var spreadRatioThreshold = big.NewInt(20 * 1e4)
var BASE_PRECISION = big.NewInt(1e6)
var SIZE_BASE_PRECISION = big.NewInt(1e12)

type Liquidable struct {
	Address        common.Address
	Size           *big.Int
	MarginFraction *big.Int
	FilledSize     *big.Int
}

func (liq Liquidable) GetUnfilledSize() *big.Int {
	return big.NewInt(0).Sub(liq.Size, liq.FilledSize)
}

func (db *InMemoryDatabase) GetLiquidableTraders(market Market, oraclePrice *big.Int) (longPositions []Liquidable, shortPositions []Liquidable) {
	longPositions, shortPositions = []Liquidable{}, []Liquidable{}
	markPrice := db.lastPrice[market]

	overSpreadLimit := isOverSpreadLimit(markPrice, oraclePrice)

	for addr, trader := range db.traderMap {
		position := trader.Positions[market]
		if position != nil {
			margin := getMarginForTrader(trader, market)
			marginFraction := getMarginFraction(margin, markPrice, position)

			if overSpreadLimit {
				oracleBasedMarginFraction := getMarginFraction(margin, oraclePrice, position)
				if oracleBasedMarginFraction.Cmp(marginFraction) == 1 {
					marginFraction = oracleBasedMarginFraction
				}
			}

			if marginFraction.Cmp(maintenanceMargin) == -1 {
				liquidable := Liquidable{
					Address:        addr,
					Size:           position.Size,
					MarginFraction: marginFraction,
					FilledSize:     big.NewInt(0),
				}
				if position.Size.Sign() == -1 {
					shortPositions = append(shortPositions, liquidable)
				} else {
					longPositions = append(longPositions, liquidable)
				}
			}
		}
	}

	// lower margin fraction positions should be liquidated first
	sortLiquidableSliceByMarginFraction(longPositions)
	sortLiquidableSliceByMarginFraction(shortPositions)
	return longPositions, shortPositions
}

func sortLiquidableSliceByMarginFraction(positions []Liquidable) []Liquidable {
	sort.SliceStable(positions, func(i, j int) bool {
		return positions[i].MarginFraction.Cmp(positions[j].MarginFraction) == -1
	})
	return positions
}

func isOverSpreadLimit(markPrice *big.Int, oraclePrice *big.Int) bool {
	// diff := abs(markPrice - oraclePrice)
	diff := multiplyBasePrecision(big.NewInt(0).Abs(big.NewInt(0).Sub(markPrice, oraclePrice)))
	// spreadRatioAbs := diff * 100 / oraclePrice
	spreadRatioAbs := big.NewInt(0).Div(diff, oraclePrice)
	if spreadRatioAbs.Cmp(spreadRatioThreshold) >= 0 {
		return true
	} else {
		return false
	}
}

func getNormalisedMargin(trader *Trader) *big.Int {
	return trader.Margins[USDC]

	// this will change after multi collateral
	// var normalisedMargin *big.Int
	// for coll, margin := range trader.Margins {
	// 	normalisedMargin += margin * priceMap[coll] * collateralWeightMap[coll]
	// }

	// return normalisedMargin
}

func getMarginForTrader(trader *Trader, market Market) *big.Int {
	if position, ok := trader.Positions[market]; ok {
		if position.UnrealisedFunding != nil {
			return big.NewInt(0).Sub(getNormalisedMargin(trader), position.UnrealisedFunding)
		}
	}
	return getNormalisedMargin(trader)
}

func getNotionalPosition(price *big.Int, size *big.Int) *big.Int {
	//notional position is base precision 1e6
	return big.NewInt(0).Abs(dividePrecisionSize(big.NewInt(0).Mul(size, price)))
}

func getUnrealisedPnl(price *big.Int, position *Position) *big.Int {
	notionalPosition := getNotionalPosition(price, position.Size)
	if position.Size.Sign() == 1 {
		return big.NewInt(0).Sub(notionalPosition, position.OpenNotional)
	} else {
		return big.NewInt(0).Sub(position.OpenNotional, notionalPosition)
	}
}

func getMarginFraction(margin *big.Int, price *big.Int, position *Position) *big.Int {
	notionalPosition := getNotionalPosition(price, position.Size)
	unrealisedPnl := getUnrealisedPnl(price, position)
	effectionMargin := big.NewInt(0).Add(margin, unrealisedPnl)
	mf := big.NewInt(0).Div(multiplyBasePrecision(effectionMargin), notionalPosition)
	if mf.Sign() == -1 {
		return big.NewInt(0)
	}
	return mf
}

func multiplyBasePrecision(number *big.Int) *big.Int {
	return big.NewInt(0).Mul(number, BASE_PRECISION)
}

func multiplyPrecisionSize(number *big.Int) *big.Int {
	return big.NewInt(0).Mul(number, SIZE_BASE_PRECISION)
}

func dividePrecisionSize(number *big.Int) *big.Int {
	return big.NewInt(0).Div(number, SIZE_BASE_PRECISION)
}
