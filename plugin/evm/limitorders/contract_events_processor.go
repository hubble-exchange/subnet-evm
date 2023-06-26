package limitorders

import (
	"encoding/json"
	"math/big"
	"sort"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type ContractEventsProcessor struct {
	orderBookABI     abi.ABI
	marginAccountABI abi.ABI
	clearingHouseABI abi.ABI
	database         LimitOrderDatabase
}

func NewContractEventsProcessor(database LimitOrderDatabase) *ContractEventsProcessor {
	orderBookABI, err := abi.FromSolidityJson(string(orderBookAbi))
	if err != nil {
		panic(err)
	}

	marginAccountABI, err := abi.FromSolidityJson(string(marginAccountAbi))
	if err != nil {
		panic(err)
	}

	clearingHouseABI, err := abi.FromSolidityJson(string(clearingHouseAbi))
	if err != nil {
		panic(err)
	}
	return &ContractEventsProcessor{
		orderBookABI:     orderBookABI,
		marginAccountABI: marginAccountABI,
		clearingHouseABI: clearingHouseABI,
		database:         database,
	}
}

func (cep *ContractEventsProcessor) ProcessEvents(logs []*types.Log) {
	var (
		deletedLogs []*types.Log
		rebirthLogs []*types.Log
	)
	for i := 0; i < len(logs); i++ {
		log := logs[i]
		if log.Removed {
			deletedLogs = append(deletedLogs, log)
		} else {
			rebirthLogs = append(rebirthLogs, log)
		}
	}

	// deletedLogs are in descending order by (blockNumber, LogIndex)
	// rebirthLogs should be in ascending order by (blockNumber, LogIndex)
	sort.Slice(deletedLogs, func(i, j int) bool {
		if deletedLogs[i].BlockNumber == deletedLogs[j].BlockNumber {
			return deletedLogs[i].Index > deletedLogs[j].Index
		}
		return deletedLogs[i].BlockNumber > deletedLogs[j].BlockNumber
	})

	sort.Slice(rebirthLogs, func(i, j int) bool {
		if rebirthLogs[i].BlockNumber == rebirthLogs[j].BlockNumber {
			return rebirthLogs[i].Index < rebirthLogs[j].Index
		}
		return rebirthLogs[i].BlockNumber < rebirthLogs[j].BlockNumber
	})

	logs = append(deletedLogs, rebirthLogs...)
	for _, event := range logs {
		switch event.Address {
		case OrderBookContractAddress:
			cep.handleOrderBookEvent(event)
		case MarginAccountContractAddress:
			cep.handleMarginAccountEvent(event)
		}
	}
}

func (cep *ContractEventsProcessor) ProcessAcceptedEvents(logs []*types.Log) {
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber == logs[j].BlockNumber {
			return logs[i].Index < logs[j].Index
		}
		return logs[i].BlockNumber < logs[j].BlockNumber
	})

	for _, event := range logs {
		switch event.Address {
		case ClearingHouseContractAddress:
			cep.handleClearingHouseEvent(event)
		}
	}
}

// handles OrderBook's placeOrders and cancelOrders transactions
func (cep *ContractEventsProcessor) ProcessHeadBlock(block *types.Block) {
	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}
		switch *tx.To() {
		case OrderBookContractAddress:
			input := tx.Data()
			if len(input) < 4 {
				continue
			}
			method_ := input[:4]
			method, _ := cep.orderBookABI.MethodById(method_)
			if method == nil {
				continue
			}
			switch method.Name {
			case "placeOrders":
				inputMap := map[string]interface{}{}
				err := method.Inputs.UnpackIntoMap(inputMap, input[4:])
				if err != nil {
					log.Error("ProcessHeadBlock - error in orderBookAbi.UnpackIntoMap", "method", "placeOrders", "err", err)
					continue
				}
				orders := getOrdersFromInterface(inputMap["orders"])
				signatures, ok := inputMap["signatures"].([][]byte)
				if !ok {
					log.Error("ProcessHeadBlock - improper type", "method", "placeOrders", "err", "signatures is not [][]byte")
					continue
				}
				limitOrders := make([]*LimitOrder, len(orders))
				for i, order := range orders {
					signature := signatures[i]
					limitOrders[i] = &LimitOrder{
						Id:                      order.Hash(),
						Market:                  Market(order.AmmIndex.Int64()),
						PositionType:            getPositionTypeBasedOnBaseAssetQuantity(order.BaseAssetQuantity),
						Trader:                  order.Trader,
						BaseAssetQuantity:       order.BaseAssetQuantity,
						FilledBaseAssetQuantity: big.NewInt(0),
						Price:                   order.Price,
						ExpireAt:                order.ValidUntil.Uint64(),
						RawOrder:                order,
						Salt:                    order.Salt,
						ReduceOnly:              order.ReduceOnly,
						BlockNumber:             big.NewInt(0).Set(block.Number()),
						LifecycleList:           []Lifecycle{},
						Signature:               signature,
					}
					log.Info("OrderPlaced", "orderId", limitOrders[i].Id.String(), "order", limitOrders[i])
				}

				cep.database.CreateOrders(limitOrders)

			case "cancelOrders":
				inputMap := map[string]interface{}{}
				err := method.Inputs.UnpackIntoMap(inputMap, input[4:])
				if err != nil {
					log.Error("ProcessHeadBlock - error in orderBookAbi.UnpackIntoMap", "method", "cancelOrders", "err", err)
					continue
				}
				orders := getOrdersFromInterface(inputMap["orders"])

				for _, order := range orders {
					orderId := order.Hash()
					cep.database.SetOrderStatus(orderId, Cancelled, "", block.Number().Uint64())
					log.Info("OrderCancelled", "orderId", orderId.String())
				}
			}

		default:
			continue
		}
	}
}

func (cep *ContractEventsProcessor) handleOrderBookEvent(event *types.Log) {
	removed := event.Removed
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case cep.orderBookABI.Events["OrdersMatched"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "OrdersMatched", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrdersMatched", "err", err)
			return
		}

		order0Id := event.Topics[1]
		order1Id := event.Topics[2]
		fillAmount := args["fillAmount"].(*big.Int)
		if !removed {
			log.Info("OrdersMatched", "orderId_0", order0Id.String(), "orderId_1", order1Id.String(), "number", event.BlockNumber)
			cep.database.UpdateFilledBaseAssetQuantity(fillAmount, order0Id, event.BlockNumber)
			cep.database.UpdateFilledBaseAssetQuantity(fillAmount, order1Id, event.BlockNumber)
		} else {
			fillAmount.Neg(fillAmount)
			log.Info("OrdersMatched removed", "orderId_0", order0Id.String(), "orderId_1", order1Id.String(), "number", event.BlockNumber)
			cep.database.UpdateFilledBaseAssetQuantity(fillAmount, order0Id, event.BlockNumber)
			cep.database.UpdateFilledBaseAssetQuantity(fillAmount, order1Id, event.BlockNumber)
		}
	case cep.orderBookABI.Events["LiquidationOrderMatched"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "LiquidationOrderMatched", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "LiquidationOrderMatched", "err", err)
			return
		}
		fillAmount := args["fillAmount"].(*big.Int)

		orderId := event.Topics[2]
		// @todo update liquidable position info
		if !removed {
			log.Info("LiquidationOrderMatched", "args", args, "orderId", orderId.String())
			cep.database.UpdateFilledBaseAssetQuantity(fillAmount, orderId, event.BlockNumber)
		} else {
			log.Info("LiquidationOrderMatched removed", "args", args, "orderId", orderId.String(), "number", event.BlockNumber)
			cep.database.UpdateFilledBaseAssetQuantity(fillAmount.Neg(fillAmount), orderId, event.BlockNumber)
		}
	case cep.orderBookABI.Events["OrderMatchingError"].ID:
		err := cep.orderBookABI.UnpackIntoMap(args, "OrderMatchingError", event.Data)
		if err != nil {
			log.Error("error in orderBookAbi.UnpackIntoMap", "method", "OrderMatchingError", "err", err)
			return
		}
		orderId := event.Topics[1]
		if !removed {
			log.Info("OrderMatchingError", "args", args, "orderId", orderId.String())
			if err := cep.database.SetOrderStatus(orderId, Execution_Failed, args["err"].(string), event.BlockNumber); err != nil {
				log.Error("error in SetOrderStatus", "method", "OrderMatchingError", "err", err)
				return
			}
		} else {
			log.Info("OrderMatchingError removed", "args", args, "orderId", orderId.String(), "number", event.BlockNumber)
			if err := cep.database.RevertLastStatus(orderId); err != nil {
				log.Error("error in SetOrderStatus", "method", "OrderMatchingError", "removed", true, "err", err)
				return
			}
		}
	}
}

func (cep *ContractEventsProcessor) handleMarginAccountEvent(event *types.Log) {
	removed := event.Removed
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case cep.marginAccountABI.Events["MarginAdded"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "MarginAdded", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginAdded", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		collateral := event.Topics[2].Big().Int64()
		amount := args["amount"].(*big.Int)
		log.Info("MarginAdded", "trader", trader, "collateral", collateral, "amount", amount.Uint64(), "removed", removed)
		if !removed {
			cep.database.UpdateMargin(trader, Collateral(collateral), amount)
		} else {
			cep.database.UpdateMargin(trader, Collateral(collateral), big.NewInt(0).Neg(amount))
		}
	case cep.marginAccountABI.Events["MarginRemoved"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "MarginRemoved", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginRemoved", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		collateral := event.Topics[2].Big().Int64()
		amount := args["amount"].(*big.Int)
		log.Info("MarginRemoved", "trader", trader, "collateral", collateral, "amount", amount.Uint64(), "removed", removed)
		if !removed {
			cep.database.UpdateMargin(trader, Collateral(collateral), big.NewInt(0).Neg(amount))
		} else {
			cep.database.UpdateMargin(trader, Collateral(collateral), amount)
		}
	case cep.marginAccountABI.Events["MarginReserved"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "MarginReserved", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginReserved", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		amount := args["amount"].(*big.Int)
		log.Info("MarginReserved", "trader", trader, "amount", amount.Uint64(), "removed", removed)
		if !removed {
			cep.database.UpdateReservedMargin(trader, amount)
		} else {
			cep.database.UpdateReservedMargin(trader, big.NewInt(0).Neg(amount))
		}
	case cep.marginAccountABI.Events["MarginReleased"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "MarginReleased", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "MarginReleased", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		amount := args["amount"].(*big.Int)
		log.Info("MarginReleased", "trader", trader, "amount", amount.Uint64(), "removed", removed)
		if !removed {
			cep.database.UpdateReservedMargin(trader, big.NewInt(0).Neg(amount))
		} else {
			cep.database.UpdateReservedMargin(trader, amount)
		}
	case cep.marginAccountABI.Events["PnLRealized"].ID:
		err := cep.marginAccountABI.UnpackIntoMap(args, "PnLRealized", event.Data)
		if err != nil {
			log.Error("error in marginAccountABI.UnpackIntoMap", "method", "PnLRealized", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		realisedPnL := args["realizedPnl"].(*big.Int)
		log.Info("PnLRealized", "trader", trader, "amount", realisedPnL.Uint64(), "removed", removed)
		if !removed {
			cep.database.UpdateMargin(trader, HUSD, realisedPnL)
		} else {
			cep.database.UpdateMargin(trader, HUSD, big.NewInt(0).Neg(realisedPnL))
		}
	}
}

func (cep *ContractEventsProcessor) handleClearingHouseEvent(event *types.Log) {
	args := map[string]interface{}{}
	switch event.Topics[0] {
	case cep.clearingHouseABI.Events["FundingRateUpdated"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "FundingRateUpdated", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "FundingRateUpdated", "err", err)
			return
		}
		cumulativePremiumFraction := args["cumulativePremiumFraction"].(*big.Int)
		nextFundingTime := args["nextFundingTime"].(*big.Int)
		market := Market(int(event.Topics[1].Big().Int64()))
		log.Info("FundingRateUpdated", "args", args, "cumulativePremiumFraction", cumulativePremiumFraction, "market", market)
		cep.database.UpdateUnrealisedFunding(market, cumulativePremiumFraction)
		cep.database.UpdateNextFundingTime(nextFundingTime.Uint64())

	case cep.clearingHouseABI.Events["FundingPaid"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "FundingPaid", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "FundingPaid", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])
		market := Market(int(event.Topics[2].Big().Int64()))
		cumulativePremiumFraction := args["cumulativePremiumFraction"].(*big.Int)
		log.Info("FundingPaid", "trader", trader, "market", market, "cumulativePremiumFraction", cumulativePremiumFraction)
		cep.database.ResetUnrealisedFunding(market, trader, cumulativePremiumFraction)

	case cep.clearingHouseABI.Events["PositionModified"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "PositionModified", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PositionModified", "err", err)
			return
		}

		trader := getAddressFromTopicHash(event.Topics[1])
		market := Market(int(event.Topics[2].Big().Int64()))
		lastPrice := args["price"].(*big.Int)
		cep.database.UpdateLastPrice(market, lastPrice)

		openNotional := args["openNotional"].(*big.Int)
		size := args["size"].(*big.Int)
		log.Info("PositionModified", "trader", trader, "market", market, "args", args)
		cep.database.UpdatePosition(trader, market, size, openNotional, false)
	case cep.clearingHouseABI.Events["PositionLiquidated"].ID:
		err := cep.clearingHouseABI.UnpackIntoMap(args, "PositionLiquidated", event.Data)
		if err != nil {
			log.Error("error in clearingHouseABI.UnpackIntoMap", "method", "PositionLiquidated", "err", err)
			return
		}
		trader := getAddressFromTopicHash(event.Topics[1])

		market := Market(int(event.Topics[2].Big().Int64()))
		lastPrice := args["price"].(*big.Int)
		cep.database.UpdateLastPrice(market, lastPrice)

		openNotional := args["openNotional"].(*big.Int)
		size := args["size"].(*big.Int)
		log.Info("PositionLiquidated", "market", market, "trader", trader, "args", args)
		cep.database.UpdatePosition(trader, market, size, openNotional, true)
	}
}

func getAddressFromTopicHash(topicHash common.Hash) common.Address {
	address32 := topicHash.String() // address in 32 bytes with 0 padding
	return common.HexToAddress(address32[:2] + address32[26:])
}

func getOrderFromRawOrder(rawOrder interface{}) Order {
	order := Order{}
	marshalledOrder, _ := json.Marshal(rawOrder)
	_ = json.Unmarshal(marshalledOrder, &order)
	return order
}

func getOrdersFromInterface(rawOrders interface{}) []Order {
	orders := []Order{}
	marshalledOrders, _ := json.Marshal(rawOrders)
	_ = json.Unmarshal(marshalledOrders, &orders)
	return orders
}
