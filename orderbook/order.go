package orderbook

import (
	"time"
)

type Value uint64

type PriceLimit uint64

func (p PriceLimit) IsMarket() bool {
	return p == PLMarket
}

const PLMarket = PriceLimit(0)

type OrderType byte

const (
	OTBid OrderType = iota
	OTAsk
)

type Order struct {
	ID         string
	Type       OrderType
	Value      Value
	Price      PriceLimit
	AcceptTime time.Time
}
