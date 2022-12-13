package orderbook

import (
	"strings"
	"time"

	"github.com/e-zhydzetski/strips-tt/orderbook/queue"
	"github.com/e-zhydzetski/strips-tt/orderbook/tree"
)

// The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func lowToHighPrice(a Order, b Order) int {
	if a.Price < b.Price {
		return -1
	}
	if a.Price > b.Price {
		return 1
	}
	if a.AcceptTime.Before(b.AcceptTime) {
		return -1
	}
	if a.AcceptTime.After(b.AcceptTime) {
		return 1
	}
	return strings.Compare(a.ID, b.ID)
}

// The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func highToLowPrice(a Order, b Order) int {
	if a.Price > b.Price {
		return -1
	}
	if a.Price < b.Price {
		return 1
	}
	if a.AcceptTime.Before(b.AcceptTime) {
		return -1
	}
	if a.AcceptTime.After(b.AcceptTime) {
		return 1
	}
	return strings.Compare(a.ID, b.ID)
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		limitBids:  tree.New[Order, Value](highToLowPrice),
		limitAsks:  tree.New[Order, Value](lowToHighPrice),
		marketBids: queue.New[Order](),
		marketAsks: queue.New[Order](),
		events:     NewEvents(100),
	}
}

type OrderBook struct {
	limitBids  *tree.Tree[Order, Value] // TODO tree of queues, or skip list
	limitAsks  *tree.Tree[Order, Value] // TODO tree of queues, or skip list
	marketBids *queue.Queue[Order]
	marketAsks *queue.Queue[Order]
	events     *Events
}

//nolint:dupl // TODO refactor
func (o *OrderBook) Ask(id string, value Value, price PriceLimit) {
	now := time.Now()
	o.events.Emit(OrderAccepted{
		ID:         id,
		Type:       OTAsk,
		Value:      value,
		Price:      price,
		AcceptTime: now,
	})
	o.marketBids.Iterate(func(order *Order) queue.IteratorAction {
		if order.Value > value {
			order.Value -= value
			value = 0
			o.events.Emit(OrderChanged{
				ID:    order.ID,
				Value: order.Value,
			})
			return queue.IAStop
		}
		// order.Value <= value
		value -= order.Value
		o.events.Emit(OrderExecuted{
			ID: order.ID,
		})
		return queue.IARemoveAndContinue
	})
	if value == 0 {
		o.events.Emit(OrderExecuted{
			ID: id,
		})
		return
	}
	o.limitBids.Iterate(func(order Order, remainedValue *Value) tree.IteratorAction {
		if !price.IsMarket() {
			if order.Price < price {
				return tree.IAStop
			}
		}

		if *remainedValue > value {
			*remainedValue -= value
			value = 0
			o.events.Emit(OrderChanged{
				ID:    order.ID,
				Value: order.Value,
			})
			return tree.IAStop
		}
		// remainedValue <= value
		value -= *remainedValue
		o.events.Emit(OrderExecuted{
			ID: order.ID,
		})
		return tree.IARemoveAndContinue
	})
	if value == 0 {
		o.events.Emit(OrderExecuted{
			ID: id,
		})
		return
	}
	o.events.Emit(OrderChanged{
		ID:    id,
		Value: value,
	})

	newOrder := Order{
		ID:         id,
		Type:       OTAsk,
		Value:      value,
		Price:      price,
		AcceptTime: now,
	}
	if price.IsMarket() {
		o.marketAsks.Add(newOrder)
	} else {
		o.limitAsks.Set(newOrder, value)
	}
	// o.events.PrintAll()
}

//nolint:dupl // TODO refactor
func (o *OrderBook) Bid(id string, value Value, price PriceLimit) {
	now := time.Now()
	o.events.Emit(OrderAccepted{
		ID:         id,
		Type:       OTBid,
		Value:      value,
		Price:      price,
		AcceptTime: now,
	})
	o.marketAsks.Iterate(func(order *Order) queue.IteratorAction {
		if order.Value > value {
			order.Value -= value
			value = 0
			o.events.Emit(OrderChanged{
				ID:    order.ID,
				Value: order.Value,
			})
			return queue.IAStop
		}
		// order.Value <= value
		value -= order.Value
		o.events.Emit(OrderExecuted{
			ID: order.ID,
		})
		return queue.IARemoveAndContinue
	})
	if value == 0 {
		o.events.Emit(OrderExecuted{
			ID: id,
		})
		return
	}
	o.limitAsks.Iterate(func(order Order, remainedValue *Value) tree.IteratorAction {
		if !price.IsMarket() {
			if order.Price < price {
				return tree.IAStop
			}
		}

		if *remainedValue > value {
			*remainedValue -= value
			value = 0
			o.events.Emit(OrderChanged{
				ID:    order.ID,
				Value: order.Value,
			})
			return tree.IAStop
		}
		// remainedValue <= value
		value -= *remainedValue
		o.events.Emit(OrderExecuted{
			ID: order.ID,
		})
		return tree.IARemoveAndContinue
	})
	if value == 0 {
		o.events.Emit(OrderExecuted{
			ID: id,
		})
		return
	}
	o.events.Emit(OrderChanged{
		ID:    id,
		Value: value,
	})

	newOrder := Order{
		ID:         id,
		Type:       OTBid,
		Value:      value,
		Price:      price,
		AcceptTime: now,
	}
	if price.IsMarket() {
		o.marketBids.Add(newOrder)
	} else {
		o.limitBids.Set(newOrder, value)
	}
	// o.events.PrintAll()
}
