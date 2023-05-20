package hubblebibliophile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	CLEARING_HOUSE_GENESIS_ADDRESS       = "0x0300000000000000000000000000000000000071"
	AMMS_SLOT                      int64 = 12
)

type MarginMode uint8

const (
	Maintenance_Margin MarginMode = iota
	Min_Allowable_Margin
)

func GetMarginMode(marginMode uint8) MarginMode {
	if marginMode == 0 {
		return Maintenance_Margin
	}
	return Min_Allowable_Margin
}

func marketsStorageSlot() *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256(common.LeftPadBytes(big.NewInt(AMMS_SLOT).Bytes(), 32)))
}

func GetActiveMarketsCount(stateDB contract.StateDB) int64 {
	rawVal := stateDB.GetState(common.HexToAddress(CLEARING_HOUSE_GENESIS_ADDRESS), common.BytesToHash(common.LeftPadBytes(big.NewInt(AMMS_SLOT).Bytes(), 32)))
	return new(big.Int).SetBytes(rawVal.Bytes()).Int64()
}

func GetMarkets(stateDB contract.StateDB) []common.Address {
	numMarkets := GetActiveMarketsCount(stateDB)
	markets := make([]common.Address, numMarkets)
	baseStorageSlot := marketsStorageSlot()
	for i := int64(0); i < numMarkets; i++ {
		amm := stateDB.GetState(common.HexToAddress(CLEARING_HOUSE_GENESIS_ADDRESS), common.BigToHash(new(big.Int).Add(baseStorageSlot, big.NewInt(i))))
		markets[i] = common.BytesToAddress(amm.Bytes())

	}
	return markets
}

func GetNotionalPositionAndMargin(stateDB contract.StateDB, input *GetNotionalPositionAndMarginInput) GetNotionalPositionAndMarginOutput {
	margin := GetNormalizedMargin(stateDB, input.Trader)
	if input.IncludeFundingPayments {
		margin.Sub(margin, GetTotalFunding(stateDB, &input.Trader))
	}
	notionalPosition, unrealizedPnl := GetTotalNotionalPositionAndUnrealizedPnl(stateDB, &input.Trader, margin, GetMarginMode(input.Mode))
	return GetNotionalPositionAndMarginOutput{
		NotionalPosition: notionalPosition,
		Margin:           new(big.Int).Add(margin, unrealizedPnl),
	}
}

func GetTotalNotionalPositionAndUnrealizedPnl(stateDB contract.StateDB, trader *common.Address, margin *big.Int, marginMode MarginMode) (*big.Int, *big.Int) {
	notionalPosition := big.NewInt(0)
	unrealizedPnl := big.NewInt(0)
	for _, market := range GetMarkets(stateDB) {
		lastPrice := getLastPrice(stateDB, market)
		// oraclePrice := getUnderlyingPrice(stateDB, market) // TODO
		oraclePrice := multiply1e6(big.NewInt(1800))
		_notionalPosition, _unrealizedPnl := getOptimalPnl(stateDB, market, oraclePrice, lastPrice, trader, margin, marginMode)
		notionalPosition.Add(notionalPosition, _notionalPosition)
		unrealizedPnl.Add(unrealizedPnl, _unrealizedPnl)
	}
	return notionalPosition, unrealizedPnl
}

func GetTotalFunding(stateDB contract.StateDB, trader *common.Address) *big.Int {
	totalFunding := big.NewInt(0)
	for _, market := range GetMarkets(stateDB) {
		totalFunding.Add(totalFunding, getPendingFundingPayment(stateDB, market, trader))
	}
	return totalFunding
}
