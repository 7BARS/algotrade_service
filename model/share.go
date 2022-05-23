package model

import (
	"fmt"
	"sync"
	"time"

	"algotrade_service/internal/provider/tinkoff"

	sdk "github.com/tinkoff/invest-api-go-sdk"
	"go.uber.org/atomic"
)

// Share has a lot of hardcode about interval :(, but it's proof of concept
type Share struct {
	chOnUpdateByTicker chan string

	Info       sdk.Share `json:"info"`
	lastUpdate atomic.Time

	rwmux       *sync.RWMutex
	History     map[sdk.CandleInterval][]*sdk.HistoricCandle `json:"history"`
	TempCandles map[sdk.CandleInterval]*sdk.Candle
}

func NewShare(info sdk.Share, history map[sdk.CandleInterval][]*sdk.HistoricCandle, chOnUpdateByTicker chan string) *Share {
	return &Share{
		chOnUpdateByTicker: chOnUpdateByTicker,
		Info:               info,
		History:            history,
		TempCandles:        make(map[sdk.CandleInterval]*sdk.Candle),
		rwmux:              &sync.RWMutex{},
	}
}

func (s *Share) CloseOpenCandle() {
	curHor, curMin, _ := time.Now().Clock()
	lastHour, lastMin, _ := s.lastUpdate.Load().Clock()
	s.closeOpenCandle(lastMin, curMin, 5, sdk.CandleInterval_CANDLE_INTERVAL_5_MIN)
	s.closeOpenCandle(lastMin, curMin, 15, sdk.CandleInterval_CANDLE_INTERVAL_15_MIN)
	s.closeOpenCandleHour(lastHour, curHor, sdk.CandleInterval_CANDLE_INTERVAL_HOUR)

	s.onChange()
}

func (s *Share) AddMinuteCandle(candle sdk.Candle) {
	t := candle.GetTime().AsTime()
	hour, min, sec := t.Clock()
	if sec != 0 {
		return
	}

	s.updateMinuteCandle(candle, min, 5, sdk.CandleInterval_CANDLE_INTERVAL_5_MIN)
	s.updateMinuteCandle(candle, min, 15, sdk.CandleInterval_CANDLE_INTERVAL_15_MIN)
	s.updateHourCandle(candle, hour, sdk.CandleInterval_CANDLE_INTERVAL_HOUR)
	s.onChange()

	s.lastUpdate.Store(t)
}

func (s *Share) updateMinuteCandle(candle sdk.Candle, min, mod int, interval sdk.CandleInterval) bool {
	if min%mod != 0 {
		s.updateTempCandle(interval, &candle)
		return false
	}

	_candle := s.getTempCandle(interval)
	hc := tinkoff.CandleToHistoricCandle(_candle)
	hc.IsComplete = true
	s.updateHistoricCandle(interval, &hc)
	s.removeTempCandle(interval)

	return true
}

func (s *Share) updateHourCandle(candle sdk.Candle, hour int, interval sdk.CandleInterval) bool {
	lastUpdate := s.lastUpdate.Load()
	if lastUpdate.IsZero() || lastUpdate.Hour() == hour {
		s.updateTempCandle(interval, &candle)
		return false
	}

	_candle := s.getTempCandle(interval)
	hc := tinkoff.CandleToHistoricCandle(_candle)
	hc.IsComplete = true
	s.History[interval] = append(s.History[interval], &hc)
	s.TempCandles[interval] = nil

	return true
}

func (s *Share) onChange() {
	s.chOnUpdateByTicker <- s.Info.Ticker
}

func (s *Share) updateHistoricCandle(interval sdk.CandleInterval, historicCandle *sdk.HistoricCandle) {
	s.rwmux.Lock()
	defer s.rwmux.Unlock()

	if s.History[interval] != nil {
		s.History[interval] = append(s.History[interval][:0], s.History[interval][:1]...)
	}
	s.History[interval] = append(s.History[interval], historicCandle)
}

func (s *Share) updateTempCandle(interval sdk.CandleInterval, candle *sdk.Candle) {
	s.rwmux.Lock()
	defer s.rwmux.Unlock()

	_candle := s.TempCandles[interval]
	if _candle == nil {
		s.TempCandles[interval] = candle
		return
	}
	_candle = tinkoff.NewCandle([]*sdk.Candle{candle, _candle})
	s.TempCandles[interval] = _candle
}

func (s *Share) removeTempCandle(interval sdk.CandleInterval) {
	s.rwmux.Lock()
	defer s.rwmux.Unlock()

	s.TempCandles[interval] = nil
}

func (s *Share) getTempCandle(interval sdk.CandleInterval) *sdk.Candle {
	s.rwmux.Lock()
	defer s.rwmux.Unlock()

	return s.TempCandles[interval]
}

func (s *Share) closeOpenCandle(lastMin, curMin, mod int, interval sdk.CandleInterval) bool {
	offset := lastMin % mod
	roundedMin := lastMin + offset
	if roundedMin >= curMin {
		return false
	}

	candle := s.getTempCandle(interval)
	hc := tinkoff.CandleToHistoricCandle(candle)
	s.updateHistoricCandle(interval, &hc)
	s.removeTempCandle(interval)

	return true
}

func (s *Share) closeOpenCandleHour(lastHour, curHor int, interval sdk.CandleInterval) bool {
	if lastHour != curHor {
		return false
	}

	candle := s.getTempCandle(interval)
	hc := tinkoff.CandleToHistoricCandle(candle)
	s.updateHistoricCandle(interval, &hc)
	s.removeTempCandle(interval)

	return true
}

func (s *Share) CheckEvent(event *Event) (bool, error) {
	s.rwmux.RLock()
	value, err := event.Indicator.Calc(s.History)
	s.rwmux.RUnlock()
	if err != nil {
		return false, fmt.Errorf("cannot calc indicator, error: %v", err)
	}

	switch event.ComparisonSign {
	case ">":
		if value > event.Value {
			return true, nil
		}
		return false, nil
	case "<":
		if value < event.Value {
			return true, nil
		}
		return false, nil
	case ">=":
		if value >= event.Value {
			return true, nil
		}
		return false, nil
	case "<=":
		if value <= event.Value {
			return true, nil
		}
		return false, nil
	case "==":
		if value == event.Value {
			return true, nil
		}
		return false, nil
	case "!=":
		if value != event.Value {
			return true, nil
		}
		return false, nil
	default:
		return false, fmt.Errorf("sign does not determine")
	}
}
