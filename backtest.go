package main

import (
	oandaapi "oanda-api/api"
	"testing"
)

type IStrategy interface {
	Process(candle *oandaapi.Candlestick) *Action
}

type Action struct {
	Type  string
	Units int64
	Price float64
}

type OrderBook struct {
}

type PositionBook struct {
}

type Results struct {
	Actions [
}

func BackTest(
	strategy IStrategy,
	orderBook *OrderBook,
	positionBook *PositionBook,
	candles []*oandaapi.Candlestick,
) *Results {
	actions := []*Action{}

	for _, candle := range candles {
		action := strategy.Process(candle)
		actions = append(actions, action)

		orderBook.execute(action)
	}

	return &Results{}
}

type BasicStrategy struct {
}

func (strat *BasicStrategy) Process(candle *oandaapi.Candlestick) *Action {
	return nil
}

func TestBackTest(t *testing.T) {
	candles := []*oandaapi.Candlestick{}
	orderBook := OrderBook{}
	positionBook := PositionBook{}

	strat := BasicStrategy{}
	actions := BackTest(&strat, &orderBook, &positionBook, candles)
}
