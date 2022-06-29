package store

import (
	"algotrade_service/internal/model"
	"context"
)

type Store interface {
	AppendBars(ctx context.Context, figi, ticker string, timeFrame model.TimeFrame, bars model.Bars)
	GetLastBarsByFigi(ctx context.Context, figi string, timeFrame model.TimeFrame, session model.Session, depth int) (model.Bars, error)
	GetLastBarsByTicker(ctx context.Context, ticker string, timeFrame model.TimeFrame, session model.Session, depth int) (model.Bars, error)

	AppendEvent(ctx context.Context, name string, event model.Event)
	DeleteEvent(ctx context.Context, name string, event model.Event)
	GetEvents(ctx context.Context) ([]model.Event, error)

	AppendInfo(ctx context.Context, figi string, info model.Info)
	DeleteInfo(ctx context.Context, figi string)
	GetInfo(ctx context.Context, figi string) (model.Info, error)
}
