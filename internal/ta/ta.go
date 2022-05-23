package ta

import (
	"fmt"
	"strconv"

	talib "github.com/markcheno/go-talib"
	sdk "github.com/tinkoff/invest-api-go-sdk"
)

type priceType string

const (
	open  priceType = "open"
	close priceType = "close"
	high  priceType = "high"
	low   priceType = "low"
	empty priceType = ""

	div = 1e+09
)

type Indicator interface {
	Calc(candles map[sdk.CandleInterval][]*sdk.HistoricCandle) (float64, error)
}

func NewIndicator(typeOfIndicator string, args []string) (Indicator, error) {
	switch typeOfIndicator {
	case "rsi":
		return NewRSI(args)
	case "change":
		return NewChange(args)
	case "change_percent":
		return NewChangePercent(args)
	default:
		return nil, fmt.Errorf("unsupported indicator: %s", typeOfIndicator)
	}
}

type RSI struct {
	pt       priceType
	depth    int
	interval sdk.CandleInterval
}

func NewRSI(args []string) (*RSI, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("unexpected count of arguments: %v", args)
	}

	priceType, err := parseTypeOfPrice(args[1])
	if err != nil {
		return nil, err
	}

	depth, err := parseDepth(args[0])
	if err != nil {
		return nil, err
	}

	interval, err := parseInteval(args[2])
	if err != nil {
		return nil, err
	}

	return &RSI{
		pt:       priceType,
		depth:    depth,
		interval: interval,
	}, nil
}

func (r *RSI) Calc(candles map[sdk.CandleInterval][]*sdk.HistoricCandle) (float64, error) {
	prices := pricesFromCandles(candles[r.interval], r.pt, r.depth+1)
	if len(prices) < r.depth {
		return 0, fmt.Errorf("not enough candles: %d, but need: %d", len(prices), r.depth)
	}
	rsi := talib.Rsi(prices, r.depth)
	return rsi[len(rsi)-1], nil
}

type Change struct {
	pt       priceType
	interval sdk.CandleInterval
}

func NewChange(args []string) (*Change, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("unexpected count of arguments: %v", args)
	}

	priceType, err := parseTypeOfPrice(args[1])
	if err != nil {
		return nil, err
	}

	interval, err := parseInteval(args[2])
	if err != nil {
		return nil, err
	}

	return &Change{
		pt:       priceType,
		interval: interval,
	}, nil
}

func (c *Change) Calc(candles map[sdk.CandleInterval][]*sdk.HistoricCandle) (float64, error) {
	cur, prev := getCurAndPrev(candles[c.interval], c.pt, 2)
	return change(cur, prev), nil
}

type ChangePercent struct {
	pt       priceType
	interval sdk.CandleInterval
}

func NewChangePercent(args []string) (*ChangePercent, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("unexpected count of arguments: %v", args)
	}

	priceType, err := parseTypeOfPrice(args[0])
	if err != nil {
		return nil, err
	}

	interval, err := parseInteval(args[1])
	if err != nil {
		return nil, err
	}

	return &ChangePercent{
		pt:       priceType,
		interval: interval,
	}, nil
}

func (cp *ChangePercent) Calc(candles map[sdk.CandleInterval][]*sdk.HistoricCandle) (float64, error) {
	cur, prev := getCurAndPrev(candles[cp.interval], cp.pt, 2)
	return change(cur, prev) / cur * 100, nil
}

func change(cur, prev float64) float64 {
	return cur - prev
}

func getCurAndPrev(candles []*sdk.HistoricCandle, priceType priceType, depth int) (cur float64, prev float64) {
	prices := pricesFromCandles(candles, priceType, depth)
	return prices[0], prices[1]
}

func pricesFromCandles(candles []*sdk.HistoricCandle, priceType priceType, depth int) []float64 {
	prices := make([]float64, 0, depth)
	switch priceType {
	case open:
		for i := len(candles) - 1; i >= len(candles)-depth; i-- {
			p := float64(candles[i].GetOpen().GetNano()) / div
			p += float64(candles[i].GetOpen().GetUnits())
			prices = append(prices, p)
		}
	case close:
		for i := len(candles) - 1; i >= len(candles)-depth; i-- {
			p := float64(candles[i].GetClose().GetNano()) / div
			p += float64(candles[i].GetClose().GetUnits())
			prices = append(prices, p)
		}
	case high:
		for i := len(candles) - 1; i >= len(candles)-depth; i-- {
			p := float64(candles[i].GetHigh().GetNano()) / div
			p += float64(candles[i].GetHigh().GetUnits())
			prices = append(prices, p)
		}
	case low:
		for i := len(candles) - 1; i >= len(candles)-depth; i-- {
			p := float64(candles[i].GetLow().GetNano()) / div
			p += float64(candles[i].GetLow().GetUnits())
			prices = append(prices, p)
		}
	}

	return prices
}

func parseTypeOfPrice(raw string) (priceType, error) {
	switch raw {
	case string(open):
		return open, nil
	case string(close):
		return close, nil
	case string(high):
		return high, nil
	case string(low):
		return low, nil
	default:
		return empty, fmt.Errorf("unexpected first argument: %s", raw)
	}
}

func parseDepth(raw string) (int, error) {
	depth, err := strconv.Atoi(raw)
	if err != err {
		return 0, fmt.Errorf("cannot parse depth: %s, error: %v", raw, err)
	}

	return depth, nil
}

func parseInteval(raw string) (sdk.CandleInterval, error) {
	switch raw {
	case "1":
		return sdk.CandleInterval_CANDLE_INTERVAL_1_MIN, nil
	case "5":
		return sdk.CandleInterval_CANDLE_INTERVAL_5_MIN, nil
	case "15":
		return sdk.CandleInterval_CANDLE_INTERVAL_15_MIN, nil
	case "60":
		return sdk.CandleInterval_CANDLE_INTERVAL_HOUR, nil
	case "1440":
		return sdk.CandleInterval_CANDLE_INTERVAL_DAY, nil
	default:
		return sdk.CandleInterval_CANDLE_INTERVAL_UNSPECIFIED, fmt.Errorf("unsupported interval: %s", raw)
	}
}
