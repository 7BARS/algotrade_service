package tinkoff

import (
	"context"
	"log"
	"strings"
	"time"

	sdk "github.com/tinkoff/invest-api-go-sdk"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (gw *GRPCWrap) GetCandles(ctx context.Context, figi string, interval sdk.CandleInterval, to time.Time, historyDepth, dayOffset int) ([]*sdk.HistoricCandle, error) {
	var (
		daysWithoutCandles = 0
		from               = to.AddDate(0, 0, -dayOffset)
		candles            = make([]*sdk.HistoricCandle, 0, historyDepth)
	)
	for {
		if len(candles) > historyDepth {
			break
		}

		if daysWithoutCandles == gw.magicNumberDaysAgo {
			break
		}

		cs, err := gw.getCandles(ctx, figi, from, to, interval)
		if err != nil {
			if strings.Contains(err.Error(), "ResourceExhausted") {
				log.Printf("resource exhausted, figi: %s, time frame: %s, error: %v",
					figi, sdk.CandleInterval_name[int32(interval)], err)
				log.Printf("sleep for %v", gw.coolDownTime)
				time.Sleep(gw.coolDownTime)
				continue
			}
			return nil, err
		}

		if len(cs) == 0 {
			daysWithoutCandles++
			continue
		}

		daysWithoutCandles = 0
		candles = append(candles, cs...)
		from = from.AddDate(0, 0, dayOffset)
		to = to.AddDate(0, 0, dayOffset)
	}

	return candles, nil
}

func (gw *GRPCWrap) getCandles(ctx context.Context, figi string, from time.Time, to time.Time, interval sdk.CandleInterval) ([]*sdk.HistoricCandle, error) {
	candleReq := sdk.GetCandlesRequest{
		Figi:     figi,
		From:     timestamppb.New(from),
		To:       timestamppb.New(to),
		Interval: interval,
	}
	childCtx := metadata.NewOutgoingContext(ctx, gw.MD)
	gw.limiter.Take()
	r, err := gw.marketData.GetCandles(childCtx, &candleReq)
	if err != nil {
		return nil, err
	}

	return r.GetCandles(), nil
}
