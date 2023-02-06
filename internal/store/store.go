package store

import (
	"algotrade_service/internal/model/data"
	"context"
)

type Store interface {
	AppendBars(ctx context.Context, figi, ticker string, timeFrame data.TimeFrame, bars data.Bars)
	GetLastBarsByFigi(ctx context.Context, figi string, timeFrame data.TimeFrame, session data.Session, depth int) (data.Bars, error)
	GetLastBarsByTicker(ctx context.Context, ticker string, timeFrame data.TimeFrame, session data.Session, depth int) (data.Bars, error)

	AppendEvent(ctx context.Context, name string, event data.Event)
	DeleteEvent(ctx context.Context, name string, event data.Event)
	GetEvents(ctx context.Context) ([]data.Event, error)
	
	AppendInfo(ctx context.Context, figi string, info data.Info)
	DeleteInfo(ctx context.Context, figi string)
	GetInfo(ctx context.Context, figi string) (data.Info, error)
}
