package tinkoff

import (
	"fmt"

	sdk "github.com/tinkoff/invest-api-go-sdk"
)

type Streaming struct {
	stream sdk.MarketDataStreamService_MarketDataStreamClient
	// streaming limiter
}

func (sc *Streaming) SubscribeCandles(figiList []string, interval sdk.SubscriptionInterval) error {
	return sc.subscribeCandles(figiList, interval, sdk.SubscriptionAction_SUBSCRIPTION_ACTION_SUBSCRIBE)
}

func (sc *Streaming) UnsubscribeCandles(figiList []string, interval sdk.SubscriptionInterval) error {
	return sc.subscribeCandles(figiList, interval, sdk.SubscriptionAction_SUBSCRIPTION_ACTION_UNSUBSCRIBE)
}

func (sc *Streaming) subscribeCandles(figiList []string, interval sdk.SubscriptionInterval, subscriptionAction sdk.SubscriptionAction) error {
	instruments := make([]*sdk.CandleInstrument, 0, len(figiList))
	for _, figi := range figiList {
		instruments = append(instruments, &sdk.CandleInstrument{
			Figi:     figi,
			Interval: interval,
		})
	}

	return sc.stream.Send(&sdk.MarketDataRequest{
		Payload: &sdk.MarketDataRequest_SubscribeCandlesRequest{
			SubscribeCandlesRequest: &sdk.SubscribeCandlesRequest{
				SubscriptionAction: subscriptionAction,
				Instruments:        instruments,
			},
		},
	})
}

func (sc *Streaming) RecvCandle() (*sdk.Candle, error) {
	resp, err := sc.stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("error while streaming recv, err: %v", err)
	}

	return resp.GetCandle(), nil
}
