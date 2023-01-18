package evm

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/plugin/evm/limitorders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetOrderBookContractFileLocation(t *testing.T) {
	newFileLocation := "new/location"
	SetOrderBookContractFileLocation(newFileLocation)
	assert.Equal(t, newFileLocation, orderBookContractFileLocation)
}

func newVM(t *testing.T) *VM {
	txFeeCap := float64(11)
	enabledEthAPIs := []string{"debug"}
	configJSON := fmt.Sprintf("{\"rpc-tx-fee-cap\": %g,\"eth-apis\": %s}", txFeeCap, fmt.Sprintf("[%q]", enabledEthAPIs[0]))
	_, vm, _, _ := GenesisVM(t, false, "", configJSON, "")
	return vm
}

func newLimitOrderProcesser(t *testing.T, db limitorders.LimitOrderDatabase, lotp limitorders.LimitOrderTxProcessor) LimitOrderProcesser {
	vm := newVM(t)
	lop := NewLimitOrderProcesser(
		vm.ctx,
		vm.txPool,
		vm.shutdownChan,
		&vm.shutdownWg,
		vm.eth.APIBackend,
		vm.eth.BlockChain(),
		db,
		lotp,
	)
	return lop
}
func TestNewLimitOrderProcesser(t *testing.T) {
	_, _, lop := setupDependencies(t)
	assert.NotNil(t, lop)
}

func setupDependencies(t *testing.T) (*MockLimitOrderDatabase, *MockLimitOrderTxProcessor, LimitOrderProcesser) {
	db := NewMockLimitOrderDatabase()
	lotp := NewMockLimitOrderTxProcessor()
	lop := newLimitOrderProcesser(t, db, lotp)
	return db, lotp, lop
}

func TestRunMatchingEngine(t *testing.T) {
	t.Run("when no long orders are present in memorydb", func(t *testing.T) {
		t.Run("when no short orders are present, matching engine does not call ExecuteMatchedOrders", func(t *testing.T) {
			db, lotp, lop := setupDependencies(t)
			longOrders := make([]limitorders.LimitOrder, 0)
			shortOrders := make([]limitorders.LimitOrder, 0)
			db.On("GetLongOrders").Return(longOrders)
			db.On("GetShortOrders").Return(shortOrders)
			lotp.On("PurgeLocalTx").Return(nil)
			lop.RunMatchingEngine()
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
		t.Run("when short orders are present, matching engine does not call ExecuteMatchedOrders", func(t *testing.T) {
			db, lotp, lop := setupDependencies(t)
			longOrders := make([]limitorders.LimitOrder, 0)
			shortOrders := make([]limitorders.LimitOrder, 0)
			shortOrders = append(shortOrders, getShortOrder())
			db.On("GetLongOrders").Return(longOrders)
			db.On("GetShortOrders").Return(shortOrders)
			lotp.On("PurgeLocalTx").Return(nil)
			lop.RunMatchingEngine()
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
	})
	t.Run("when long orders are present in memorydb", func(t *testing.T) {
		t.Run("when no short orders are present in memorydb, matching engine does not call ExecuteMatchedOrders", func(t *testing.T) {
			db, lotp, lop := setupDependencies(t)
			longOrders := make([]limitorders.LimitOrder, 0)
			shortOrders := make([]limitorders.LimitOrder, 0)
			longOrder := getLongOrder()
			longOrders = append(longOrders, longOrder)
			db.On("GetLongOrders").Return(longOrders)
			db.On("GetShortOrders").Return(shortOrders)
			lotp.On("PurgeLocalTx").Return(nil)
			lop.RunMatchingEngine()
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
		t.Run("When short orders are present in db", func(t *testing.T) {
			t.Run("when longOrder.price < shortOrder.price, matching engine does not call ExecuteMatchedOrders", func(t *testing.T) {
				db, lotp, lop := setupDependencies(t)
				longOrders := make([]limitorders.LimitOrder, 0)
				shortOrders := make([]limitorders.LimitOrder, 0)
				longOrder := getLongOrder()
				longOrders = append(longOrders, longOrder)
				shortOrder := getShortOrder()
				shortOrder.Price = shortOrder.Price + 2
				shortOrders = append(shortOrders, shortOrder)
				db.On("GetLongOrders").Return(longOrders)
				db.On("GetShortOrders").Return(shortOrders)
				lotp.On("PurgeLocalTx").Return(nil)
				lop.RunMatchingEngine()
				lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
			})
			t.Run("when longOrder.price >= shortOrder.price", func(t *testing.T) {
				t.Run("When long order and short order's unfulfilled quantity is same", func(t *testing.T) {
					t.Run("When long order and short order's base asset quantity is same, matching engine calls ExecuteMatchedOrders", func(t *testing.T) {
						//Add 2 long orders
						db, lotp, lop := setupDependencies(t)
						longOrder1 := getLongOrder()
						longOrder2 := getLongOrder()
						longOrder2.Price = longOrder1.Price + 1
						longOrder2.Signature = []byte("Here is a 2nd long order")
						//slice sorted by higher price
						longOrders := []limitorders.LimitOrder{longOrder2, longOrder1}

						// Add 2 short orders
						shortOrder1 := getShortOrder()
						shortOrder2 := getShortOrder()
						shortOrder2.Price = shortOrder1.Price - 1
						shortOrder2.Signature = []byte("Here is a 2nd short order")
						//slice sorted by lower price
						shortOrders := []limitorders.LimitOrder{shortOrder2, shortOrder1}

						db.On("GetLongOrders").Return(longOrders)
						db.On("GetShortOrders").Return(shortOrders)
						lotp.On("PurgeLocalTx").Return(nil)
						fillAmount1 := uint(longOrder1.BaseAssetQuantity)
						fillAmount2 := uint(longOrder2.BaseAssetQuantity)
						lotp.On("ExecuteMatchedOrdersTx", longOrder1, shortOrder1, fillAmount1).Return(nil)
						lotp.On("ExecuteMatchedOrdersTx", longOrder2, shortOrder2, fillAmount2).Return(nil)
						lop.RunMatchingEngine()
						lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder1, shortOrder1, fillAmount1)
						lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder2, shortOrder2, fillAmount2)
					})
					t.Run("When long order and short order's base asset quantity is different, matching engine calls ExecuteMatchedOrders", func(t *testing.T) {
						db, lotp, lop := setupDependencies(t)

						longOrder := getLongOrder()
						longOrder.BaseAssetQuantity = 20
						longOrder.FilledBaseAssetQuantity = 5
						longOrders := []limitorders.LimitOrder{longOrder}

						shortOrder := getShortOrder()
						shortOrder.BaseAssetQuantity = -30
						shortOrder.FilledBaseAssetQuantity = -15
						shortOrders := []limitorders.LimitOrder{shortOrder}

						fillAmount := uint(longOrder.BaseAssetQuantity - longOrder.FilledBaseAssetQuantity)
						db.On("GetLongOrders").Return(longOrders)
						db.On("GetShortOrders").Return(shortOrders)
						lotp.On("PurgeLocalTx").Return(nil)
						lotp.On("ExecuteMatchedOrdersTx", longOrder, shortOrder, fillAmount).Return(nil)
						lop.RunMatchingEngine()
						lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder, shortOrder, fillAmount)
					})
				})
				t.Run("When long order and short order's unfulfilled quantity is not same, matching engine calls ExecuteMatchedOrders", func(t *testing.T) {
					db, lotp, lop := setupDependencies(t)
					longOrder1 := getLongOrder()
					longOrder1.BaseAssetQuantity = 20
					longOrder1.FilledBaseAssetQuantity = 5
					longOrder2 := getLongOrder()
					longOrder2.BaseAssetQuantity = 40
					longOrder2.FilledBaseAssetQuantity = 0
					longOrder2.Price = longOrder1.Price + 1
					longOrder2.Signature = []byte("Here is a 2nd long order")
					longOrder3 := getLongOrder()
					longOrder3.BaseAssetQuantity = 10
					longOrder3.FilledBaseAssetQuantity = 3
					longOrder3.Signature = []byte("Here is a 3rd long order")
					longOrder3.Price = longOrder2.Price + 1
					//slice sorted by higher price
					longOrders := []limitorders.LimitOrder{longOrder3, longOrder2, longOrder1}

					// Add 2 short orders
					shortOrder1 := getShortOrder()
					shortOrder1.BaseAssetQuantity = -30
					shortOrder1.FilledBaseAssetQuantity = -2
					shortOrder2 := getShortOrder()
					shortOrder2.BaseAssetQuantity = -50
					shortOrder2.FilledBaseAssetQuantity = -20
					shortOrder2.Price = shortOrder1.Price - 1
					shortOrder2.Signature = []byte("Here is a 2nd short order")
					shortOrder3 := getShortOrder()
					shortOrder3.BaseAssetQuantity = -20
					shortOrder3.FilledBaseAssetQuantity = -10
					shortOrder3.Price = shortOrder2.Price - 1
					shortOrder3.Signature = []byte("Here is a 3rd short order")
					//slice sorted by lower price
					shortOrders := []limitorders.LimitOrder{shortOrder3, shortOrder2, shortOrder1}

					lotp.On("ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(5)

					db.On("GetLongOrders").Return(longOrders)
					db.On("GetShortOrders").Return(shortOrders)
					lotp.On("PurgeLocalTx").Return(nil)
					lop.RunMatchingEngine()

					//During 1st  matching iteration
					fillAmount := uint(math.Min(float64(getUnFilledBaseAssetQuantity(longOrder3)), float64(-(getUnFilledBaseAssetQuantity(shortOrder3)))))
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder3, shortOrder3, fillAmount)
					//After 1st matching iteration longOrder3 has been matched fully but shortOrder3 has not
					longOrder3.FilledBaseAssetQuantity = longOrder3.FilledBaseAssetQuantity + int(fillAmount)
					shortOrder3.FilledBaseAssetQuantity = shortOrder3.FilledBaseAssetQuantity - int(fillAmount)

					//During 2nd iteration longOrder2 with be matched with shortOrder3
					fillAmount = uint(math.Min(float64(getUnFilledBaseAssetQuantity(longOrder2)), float64(-(getUnFilledBaseAssetQuantity(shortOrder3)))))
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder2, shortOrder3, fillAmount)
					//After 2nd matching iteration shortOrder3 has been matched fully but longOrder2 has not
					longOrder2.FilledBaseAssetQuantity = longOrder2.FilledBaseAssetQuantity + int(fillAmount)
					shortOrder3.FilledBaseAssetQuantity = shortOrder3.FilledBaseAssetQuantity - int(fillAmount)

					//During 3rd iteration longOrder2 with be matched with shortOrder2
					fillAmount = uint(math.Min(float64(getUnFilledBaseAssetQuantity(longOrder2)), float64(-(getUnFilledBaseAssetQuantity(shortOrder2)))))
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder2, shortOrder2, fillAmount)
					//After 3rd matching iteration shortOrder2 has been matched fully but longOrder2 has not
					longOrder2.FilledBaseAssetQuantity = longOrder2.FilledBaseAssetQuantity + int(fillAmount)
					shortOrder2.FilledBaseAssetQuantity = shortOrder2.FilledBaseAssetQuantity - int(fillAmount)

					//During 4th iteration longOrder2 with be matched with shortOrder1
					fillAmount = uint(math.Min(float64(getUnFilledBaseAssetQuantity(longOrder2)), float64(-(getUnFilledBaseAssetQuantity(shortOrder1)))))
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder2, shortOrder1, fillAmount)
					//After 4rd matching iteration shortOrder2 has been matched fully but longOrder3 has not
					longOrder2.FilledBaseAssetQuantity = longOrder2.FilledBaseAssetQuantity + int(fillAmount)
					shortOrder1.FilledBaseAssetQuantity = shortOrder1.FilledBaseAssetQuantity - int(fillAmount)

					//During 5th iteration longOrder1 with be matched with shortOrder1
					fillAmount = uint(math.Min(float64(getUnFilledBaseAssetQuantity(longOrder1)), float64(-(getUnFilledBaseAssetQuantity(shortOrder1)))))
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder1, shortOrder1, fillAmount)
				})
			})
		})
	})
}

func getShortOrder() limitorders.LimitOrder {
	signature := []byte("Here is a short order")
	salt := time.Now().Unix()
	shortOrder := createLimitOrder("short", "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", -10, 20.01, "unfulfilled", salt, signature, 2)
	return shortOrder
}

func getLongOrder() limitorders.LimitOrder {
	signature := []byte("Here is a long order")
	salt := time.Now().Unix()
	longOrder := createLimitOrder("long", "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", 10, 20.01, "unfulfilled", salt, signature, 2)
	return longOrder
}

func createLimitOrder(positionType string, userAddress string, baseAssetQuantity int, price float64, status string, salt int64, signature []byte, blockNumber uint64) limitorders.LimitOrder {
	return limitorders.LimitOrder{
		PositionType:            positionType,
		UserAddress:             userAddress,
		BaseAssetQuantity:       baseAssetQuantity,
		FilledBaseAssetQuantity: 0,
		Price:                   price,
		Status:                  status,
		Salt:                    salt,
		Signature:               signature,
		BlockNumber:             blockNumber,
	}
}
