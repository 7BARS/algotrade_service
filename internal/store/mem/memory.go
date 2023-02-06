package mem

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"algotrade_service/internal/model"
)

type instrument struct {
	info    model.Info
	history map[model.TimeFrame]model.Bars
}

type MemoryStore struct {
	rwmux               *sync.RWMutex
	instrumentsByTicker map[string]instrument
	instrumentsByFIGI   map[string]instrument
	events              map[string]model.Event
}

func NewMemoryStore() (*MemoryStore, error) {
	return &MemoryStore{
		rwmux:               &sync.RWMutex{},
		instrumentsByTicker: make(map[string]instrument),
		instrumentsByFIGI:   make(map[string]instrument),
		events:              make(map[string]model.Event),
	}, nil
}

func (ms *MemoryStore) AppendBars(ctx context.Context, figi string, ticker string, timeFrame model.TimeFrame, bars model.Bars) {
	ms.rwmux.Lock()
	defer ms.rwmux.Unlock()

	_instrument, ok := ms.instrumentsByFIGI[figi]
	if !ok {
		_instrument = instrument{
			history: make(map[model.TimeFrame]model.Bars),
		}
		_instrument.history[timeFrame] = bars
		ms.instrumentsByFIGI[figi] = _instrument
	}
	ms.instrumentsByFIGI[figi] = _instrument
	ms.instrumentsByTicker[ticker] = _instrument
}

func (ms *MemoryStore) GetLastBarsByFigi(ctx context.Context, figi string, timeFrame model.TimeFrame, session model.Session, depth int) (model.Bars, error) {
	ms.rwmux.Lock()
	defer ms.rwmux.RUnlock()

	_instrument, ok := ms.instrumentsByFIGI[figi]
	if !ok {
		return nil, fmt.Errorf("instrument %s not found", figi)
	}

	return getBarsByTimeFrameFromInstrument(_instrument, timeFrame)
}

func (ms *MemoryStore) GetLastBarsByTicker(ctx context.Context, ticker string, timeFrame model.TimeFrame, session model.Session, depth int) (model.Bars, error) {
	ms.rwmux.RLock()
	defer ms.rwmux.RUnlock()

	_instrument, ok := ms.instrumentsByTicker[ticker]
	if !ok {
		return nil, fmt.Errorf("instrument %s not found", ticker)
	}

	return getBarsByTimeFrameFromInstrument(_instrument, timeFrame)
}

func (ms *MemoryStore) AppendEvent(ctx context.Context, name string, event model.Event) {
	ms.rwmux.Lock()
	defer ms.rwmux.Unlock()

	ms.events[name] = event
}

func (ms *MemoryStore) DeleteEvent(ctx context.Context, name string, event model.Event) {
	ms.rwmux.Lock()
	defer ms.rwmux.Unlock()

	delete(ms.events, name)
}

func (ms *MemoryStore) GetEvents(ctx context.Context) ([]model.Event, error) {
	ms.rwmux.RLock()
	defer ms.rwmux.RUnlock()

	events := make([]model.Event, 0, len(ms.events))
	for _, event := range ms.events {
		events = append(events, event)
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp < events[j].Timestamp
	})

	return events, nil
}

func (ms *MemoryStore) AppendInfo(ctx context.Context, figi string, info model.Info) {
	ms.rwmux.Lock()
	defer ms.rwmux.Unlock()

	_instrument, ok := ms.instrumentsByFIGI[figi]
	if !ok {
		_instrument = instrument{}
	}
	_instrument.info = info
	ms.instrumentsByFIGI[figi] = _instrument
}

func (ms *MemoryStore) DeleteInfo(ctx context.Context, figi string) {
	ms.rwmux.Lock()
	defer ms.rwmux.Unlock()

	ms.instrumentsByFIGI[figi] = instrument{}
}

func (ms *MemoryStore) GetInfo(ctx context.Context, figi string) (model.Info, error) {
	ms.rwmux.RLock()
	defer ms.rwmux.RUnlock()

	_instrument, ok := ms.instrumentsByFIGI[figi]
	if !ok {
		return model.Info{}, fmt.Errorf("instrument %s not found", figi)
	}

	return _instrument.info, nil
}

func getBarsByTimeFrameFromInstrument(instrument instrument, timeFrame model.TimeFrame) (model.Bars, error) {
	bars, ok := instrument.history[timeFrame]
	if !ok {
		return nil, fmt.Errorf("bars for instrument %s and time frame %d not found", instrument.info.Figi, timeFrame)
	}

	return bars, nil
}
