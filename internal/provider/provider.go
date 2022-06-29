package provider

import (
	"algotrade_service/internal/model"
	"context"
	"time"

	sdk "github.com/tinkoff/invest-api-go-sdk"
)

type Unary interface {
	GetCandles(ctx context.Context, figi string, interval sdk.CandleInterval, to time.Time, historyDepth int) ([]*sdk.HistoricCandle, error)
}

type Streaming interface {
	SubscribeCandles(figiList []string, interval sdk.SubscriptionInterval) error
	UnsubscribeCandles(figiList []string, interval sdk.SubscriptionInterval) error
}

type Provider interface {
	AvailableIntervals() []model.TimeFrame
	GetCandles(ctx context.Context, figi string, tf model.TimeFrame, to time.Time, historyDepth int) (model.Bars, error)
}
