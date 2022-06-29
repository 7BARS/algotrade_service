package tinkoff

import (
	"algotrade_service/internal/model"

	sdk "github.com/tinkoff/invest-api-go-sdk"
)

func AvailableIntervals() []sdk.CandleInterval {
	return []sdk.CandleInterval{
		sdk.CandleInterval_CANDLE_INTERVAL_1_MIN,
		sdk.CandleInterval_CANDLE_INTERVAL_5_MIN,
		sdk.CandleInterval_CANDLE_INTERVAL_15_MIN,
		sdk.CandleInterval_CANDLE_INTERVAL_HOUR,
		sdk.CandleInterval_CANDLE_INTERVAL_DAY,
	}
}

func CandleToHistoricCandle(candle *sdk.Candle) sdk.HistoricCandle {
	return sdk.HistoricCandle{
		Open:   candle.GetOpen(),
		High:   candle.GetHigh(),
		Low:    candle.GetLow(),
		Close:  candle.GetClose(),
		Volume: candle.GetVolume(),
		Time:   candle.GetTime(),
	}
}

func TimeFrameToInt(interval sdk.CandleInterval) int {
	switch interval {
	case sdk.CandleInterval_CANDLE_INTERVAL_1_MIN:
		return 1
	case sdk.CandleInterval_CANDLE_INTERVAL_5_MIN:
		return 5
	case sdk.CandleInterval_CANDLE_INTERVAL_15_MIN:
		return 15
	case sdk.CandleInterval_CANDLE_INTERVAL_HOUR:
		return 1
	case sdk.CandleInterval_CANDLE_INTERVAL_DAY:
		return 1
	default:
		return 0
	}
}

const div = 1e+09

func NewCandle(candles []*sdk.Candle) *sdk.Candle {
	candle := &sdk.Candle{
		Figi:     candles[0].GetFigi(),
		Interval: candles[0].GetInterval(),
		Open:     candles[0].GetOpen(),
		Close:    candles[len(candles)-1].Close,
	}
	volume := 0
	for _, c := range candles {
		volume += int(c.GetVolume())
		if c.Low.GetUnits() <= candle.Low.GetUnits() {
			if c.Low.GetNano() <= candle.Low.GetNano() {
				candle.Low = c.Low
			}
		}
		if c.High.GetUnits() >= candle.High.GetUnits() {
			if c.High.GetNano() >= candle.High.GetNano() {
				candle.High = c.High
			}
		}
	}
	candle.Volume = int64(volume)

	return candle
}

func CandleToBar(candle sdk.Candle) model.Bar {
	bar := model.Bar{}

	price := float64(candle.GetOpen().GetNano()) / div
	price += float64(candle.GetOpen().GetUnits())
	bar.Open = price

	price = float64(candle.GetClose().GetNano()) / div
	price += float64(candle.GetClose().GetUnits())
	bar.Close = price

	price = float64(candle.GetHigh().GetNano()) / div
	price += float64(candle.GetHigh().GetUnits())
	bar.High = price

	price = float64(candle.GetLow().GetNano()) / div
	price += float64(candle.GetLow().GetUnits())
	bar.Low = price

	price = float64(candle.GetOpen().GetNano()) / div
	price += float64(candle.GetOpen().GetUnits())
	bar.Open = price

	bar.Volume = candle.GetVolume()

	bar.Timestamp = candle.GetTime().AsTime().Unix()

	return bar
}
