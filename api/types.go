package oandaapi

import (
	"math"
	"time"
)

const (
	GRANULARITY_S5  = "S5"
	GRANULARITY_S10 = "S10"
	GRANULARITY_S30 = "S30"
	GRANULARITY_M1  = "M1"
	GRANULARITY_M5  = "M5"
	GRANULARITY_M10 = "M10"
	GRANULARITY_M30 = "M30"
	GRANULARITY_H1  = "H1"
)

type Instrument struct {
	Name                        string                              `json:"name"`
	Type                        string                              `json:"type"`
	DisplayName                 string                              `json:"displayName"`
	PipLocation                 int64                               `json:"pipLocation"`
	DisplayPrecision            int64                               `json:"displayPrecision"`
	TradeUnitsPrecision         int64                               `json:"tradeUnitsPrecision"`
	MinimumTradeSize            string                              `json:"minimumTradeSize"`
	MaximumTrailingStopDistance string                              `json:"maximumTrailingStopDistance"`
	MinimumTrailingStopDistance string                              `json:"minimumTrailingStopDistance"`
	MaximumPositionSize         string                              `json:"maximumPositionSize"`
	MaximumOrderUnits           string                              `json:"maximumOrderUnits"`
	MarginRate                  string                              `json:"marginRate"`
	Commission                  InstrumentCommission                `json:"instrumentCommission"`
	GuaranteedStopLossMode      GuaranteedStopLossModeForInstrument `json:"guaranteedStopLossMode"`
	Financing                   Financing                           `json:"financing"`
}

type InstrumentCommission struct {
	Commission        string `json:"commission"`
	UnitsTraded       string `json:"unitsTraded"`
	MinimumCommission string `json:"minimumCommission"`
}

type GuaranteedStopLossModeForInstrument = string

type Financing struct {
	LongRate            string `json:"longRate"`
	ShortRate           string `json:"shortRate"`
	FinancingDaysOfWeek []*FinancingDayOfWeek
}

type FinancingDayOfWeek struct {
	DayOfWeek   string `json:"dayOfWeek"`
	DaysCharged int64  `json:"daysCharged"`
}

type Candlestick struct {
	Time     time.Time       `json:"time"`
	Bid      CandlestickData `json:"bid"`
	Mid      CandlestickData `json:"mid"`
	Ask      CandlestickData `json:"ask"`
	Volume   int64           `json:"volume"`
	Complete bool            `json:"complete"`
}

type CandlestickData struct {
	O float64 `json:"o,string"`
	H float64 `json:"h,string"`
	L float64 `json:"l,string"`
	C float64 `json:"c,string"`
}

// DetectEngulfing checks if a candlestick is an engulfing pattern
func DetectEngulfing(prev, curr CandlestickData) int {
	// Bullish Engulfing
	if prev.O < prev.O && curr.C > curr.O &&
		curr.O <= prev.C && curr.C >= prev.O {
		return 1
	}

	// Bearish Engulfing
	if prev.C > prev.O && curr.C < curr.O &&
		curr.O >= prev.C && curr.C <= prev.O {
		return -1
	}

	return 0
}

// DetectDoji checks if a candlestick is a Doji pattern
func DetectDoji(candle CandlestickData) bool {
	// A Doji forms when the open and close prices are almost equal
	return math.Abs(candle.C-candle.O) <= (candle.H-candle.L)*0.1
}

// DetectHammer checks if a candlestick is a Hammer pattern
func DetectHammer(candle CandlestickData) bool {
	// A Hammer forms when there's a small body, a long lower shadow, and little or no upper shadow
	bodySize := math.Abs(candle.C - candle.O)
	lowerShadow := candle.O - candle.L
	upperShadow := candle.H - candle.C

	return bodySize <= (candle.H-candle.L)*0.25 && lowerShadow >= bodySize*2 && upperShadow <= bodySize
}

type Event struct {
	Index     int64
	Timestamp time.Time
	Events    []Tag
}

type Tag struct {
	Name  string
	Value string
}

// AnalyzeCandles analyzes a list of candles to detect patterns
func AnalyzeCandles(candles []*Candlestick) []*Event {
	events := []*Event{}

	for i := 1; i < len(candles); i++ {
		event := &Event{
			Timestamp: candles[i].Time,
			Index:     int64(i),
		}

		if v := DetectEngulfing(candles[i-1].Mid, candles[i].Mid); v != 0 {
			patternType := "bearish"
			if v == 1 {
				patternType = "bullish"
			}

			event.Events = append(event.Events, Tag{
				Name:  "engulfing",
				Value: patternType,
			})

			// detected = append(detected, fmt.Sprintf("%s (%s)", "engulfing", patternType))
			// fmt.Printf("Engulfing pattern detected at candle %d - %s\n", i, patternType)
		}
		if DetectDoji(candles[i].Mid) {
			// fmt.Printf("Doji pattern detected at candle %d\n", i)

			// detected = append(detected, "doji")

			event.Events = append(event.Events, Tag{
				Name: "doji",
			})
		}
		if DetectHammer(candles[i].Mid) {
			// fmt.Printf("Hammer pattern detected at candle %d\n", i)
			// detected = append(detected, "hammer")

			event.Events = append(event.Events, Tag{
				Name: "hammer",
			})
		}

		if len(event.Events) > 0 {
			events = append(events, event)
		}

		// if len(detected) > 0 {
		// 	fmt.Printf("(%d) detected patterns %v\n", i, strings.Join(detected, ", "))
		// }
	}

	return events
}
