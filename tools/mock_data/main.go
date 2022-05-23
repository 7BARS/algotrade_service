package main

import (
	"context"
	"log"
	"os"
	"time"

	"algotrade_service/cmd"
	config "algotrade_service/configs"
	"algotrade_service/internal/provider/tinkoff"
	"algotrade_service/model"

	sdk "github.com/tinkoff/invest-api-go-sdk"
)

const dataFilePath = "./dataFile.json"

func main() {
	log.Println("service start")
	args, err := cmd.ParseArgs()
	if err != nil {
		log.Printf("cannot parse cmd, err: %v", err)
		os.Exit(1)
	}

	cfg, err := config.GetConfig(args.PathConfig)
	if err != nil {
		log.Printf("cannot get config, err: %v", err)
		os.Exit(1)
	}

	provider, err := tinkoff.NewGRPCWrap(
		args.Token,
		cfg.ProviderTinkoff.RateLimitPerSecond,
		cfg.ProviderTinkoff.DaysWithEmptyHistory,
		cfg.ProviderTinkoff.TimeOutRequest,
	)
	if err != nil {
		log.Printf("cannot create new provider, err: %v", err)
		os.Exit(1)
	}

	ctx := context.TODO()
	shares, err := provider.GetSharesBase(ctx)
	if err != nil {
		log.Printf("cannot get shares, err: %v", err)
		os.Exit(1)
	}

	to := time.Now()
	controller := model.Controller{SharesByTicker: make(map[string]*model.Share)}
	for _, share := range shares {
		if share.Exchange != "SPB_MORNING" {
			continue
		}
		history := make(map[sdk.CandleInterval][]*sdk.HistoricCandle)
		for _, tf := range tinkoff.AvailableIntervals() {
			candles, err := provider.GetCandles(
				ctx,
				share.Figi,
				tf,
				to,
				cfg.ProviderTinkoff.HistoryDepth,
				cfg.ProviderTinkoff.DayOffset.GetDayOffset(tf),
			)
			if err != nil {
				log.Printf("cannot get candles by ticker: %s, time frame: %s, error: %v", share.Ticker, sdk.CandleInterval_name[int32(tf)], err)
				continue
			}
			history[tf] = candles
		}
		_share := model.Share{
			Info:    *share,
			History: history,
		}
		controller.SharesByTicker[share.Ticker] = &_share
		controller.SaveToFile(dataFilePath)
		log.Printf("history for ticker: %s, was written in file: %s", share.Ticker, dataFilePath)
	}
	log.Println("service finish")
}
