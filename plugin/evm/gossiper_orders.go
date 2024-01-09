package evm

import (
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook"
	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

func (n *pushGossiper) GossipSignedOrders(orders []*hu.SignedOrder) error {
	select {
	case n.ordersToGossipChan <- orders:
	case <-n.shutdownChan:
	}
	return nil
}

func (n *pushGossiper) awaitSignedOrderGossip() {
	n.shutdownWg.Add(1)
	go executeFuncAndRecoverPanic(func() {
		var (
			gossipTicker = time.NewTicker(ordersGossipInterval)
		)
		defer func() {
			gossipTicker.Stop()
			n.shutdownWg.Done()
		}()

		for {
			select {
			case <-gossipTicker.C:
				if attempted, err := n.gossipSignedOrders(); err != nil {
					log.Warn(
						"failed to send signed orders",
						"len(orders)", attempted,
						"err", err,
					)
				}
			case orders := <-n.ordersToGossipChan:
				for _, order := range orders {
					n.ordersToGossip[order.OrderHash] = order
				}
				if attempted, err := n.gossipSignedOrders(); err != nil {
					log.Warn(
						"failed to send signed orders",
						"len(orders)", attempted,
						"err", err,
					)
				}
			case <-n.shutdownChan:
				return
			}
		}
	}, "panic in awaitSignedOrderGossip", orderbook.AwaitSignedOrdersGossipPanicsCounter)

}

func (n *pushGossiper) gossipSignedOrders() (int, error) {
	if (time.Since(n.lastOrdersGossiped) < minGossipOrdersBatchInterval) || len(n.ordersToGossip) == 0 {
		return 0, nil
	}
	n.lastOrdersGossiped = time.Now()
	now := time.Now().Unix()
	selectedOrders := []*hu.SignedOrder{}
	for orderHash, order := range n.ordersToGossip {
		if len(selectedOrders) >= maxSignedOrdersGossipBatchSize {
			break
		}
		if order.ExpireAt.Int64() < now {
			n.stats.IncSignedOrdersGossipOrderExpired()
			log.Warn("signed order expired before gossip", "order", order, "now", now)
			delete(n.ordersToGossip, orderHash)
			continue
		}
		selectedOrders = append(selectedOrders, order)
		delete(n.ordersToGossip, orderHash)
	}

	if len(selectedOrders) == 0 {
		return 0, nil
	}

	return len(selectedOrders), n.sendSignedOrders(selectedOrders)
}

func (n *pushGossiper) sendSignedOrders(orders []*hu.SignedOrder) error {
	if len(orders) == 0 {
		return nil
	}

	ordersBytes, err := rlp.EncodeToBytes(orders)
	if err != nil {
		return err
	}
	msg := message.SignedOrdersGossip{
		Orders: ordersBytes,
	}
	msgBytes, err := message.BuildGossipMessage(n.codec, msg)
	if err != nil {
		return err
	}

	log.Trace(
		"gossiping signed orders",
		"len(orders)", len(orders),
		"size(orders)", len(msg.Orders),
	)
	n.stats.IncSignedOrdersGossipSent(int64(len(orders)))
	n.stats.IncSignedOrdersGossipBatchSent()
	return n.client.Gossip(msgBytes)
}

//   #### HANDLER ####

func (h *GossipHandler) HandleSignedOrders(nodeID ids.NodeID, msg message.SignedOrdersGossip) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Trace(
		"AppGossip called with SignedOrdersGossip",
		"peerID", nodeID,
		"bytes(orders)", len(msg.Orders),
	)

	if len(msg.Orders) == 0 {
		log.Warn(
			"AppGossip received empty SignedOrdersGossip Message",
			"peerID", nodeID,
		)
		return nil
	}

	orders := make([]*hu.SignedOrder, 0)
	if err := rlp.DecodeBytes(msg.Orders, &orders); err != nil {
		log.Trace(
			"AppGossip provided invalid orders",
			"peerID", nodeID,
			"err", err,
		)
		return nil
	}

	h.stats.IncSignedOrdersGossipReceived(int64(len(orders)))
	h.stats.IncSignedOrdersGossipBatchReceived()

	tradingAPI := h.vm.limitOrderProcesser.GetTradingAPI()
	if hu.ChainId == 0 { // set once, will need to restart node if we change
		tradingAPI.SetChainIdAndVerifyingSignedOrdersContract()
	}

	// re-gossip orders, but not when we already knew the orders
	ordersToGossip := make([]*hu.SignedOrder, 0)
	for _, order := range orders {
		needToGossip := false
		err := tradingAPI.PlaceOrder(order)
		if err == nil {
			h.stats.IncSignedOrdersGossipReceivedNew()
			needToGossip = true
		} else if err == hu.ErrOrderAlreadyExists {
			h.stats.IncSignedOrdersGossipReceivedKnown()
		} else {
			h.stats.IncSignedOrdersGossipReceiveError()
			log.Error("failed to place order", "err", err)
		}

		if needToGossip {
			ordersToGossip = append(ordersToGossip, order)
		}
	}

	if len(ordersToGossip) > 0 {
		h.vm.gossiper.GossipSignedOrders(ordersToGossip)
	}

	return nil
}
