package provider

import (
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
