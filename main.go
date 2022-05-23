package main

import (
	"log"
	"os"

	"algotrade_service/cmd"
	config "algotrade_service/configs"
	"algotrade_service/controller"
	"algotrade_service/model"
	"algotrade_service/view"
)

const dataFilePath = "./dataFile.json"

const (
	tokenEnv = "TINKOFF_TOKEN"
)

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

	eventController, err := model.NewEventController()
	if err != nil {
		log.Printf("cannot create event controller, err: %v", err)
		os.Exit(1)
	}

	controller, err := controller.NewController(os.Getenv(tokenEnv), cfg, eventController)
	if err != nil {
		log.Printf("cannot cannot create controller, err: %v", err)
		os.Exit(1)
	}

	err = controller.Start()
	if err != nil {
		log.Printf("cannot cannot start controllr, err: %v", err)
		os.Exit(1)
	}

	view := view.NewView(eventController)
	err = view.Start()
	if err != nil {
		log.Printf("cannot cannot start view, err: %v", err)
		os.Exit(1)
	}
	log.Println("service started")

	// provider, err := tinkoff.NewGRPCWrap(
	// 	args.Token,
	// 	cfg.ProviderTinkoff.RateLimitPerSecond,
	// 	cfg.ProviderTinkoff.DaysWithEmptyHistory,
	// 	cfg.ProviderTinkoff.TimeOutRequest,
	// )
	// if err != nil {
	// 	log.Printf("cannot create new provider, err: %v", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(provider)

	// ctx := context.TODO()
	// childCtx := metadata.NewOutgoingContext(ctx, provider.MD)
	// stream, err := provider.Stream.MarketDataStream(childCtx)
	// if err != nil {
	// 	log.Printf("cannot create new stream, err: %v", err)
	// 	os.Exit(1)
	// }

	// err = stream.Send(&sdk.MarketDataRequest{
	// 	Payload: &sdk.MarketDataRequest_SubscribeCandlesRequest{
	// 		SubscribeCandlesRequest: &sdk.SubscribeCandlesRequest{
	// 			SubscriptionAction: sdk.SubscriptionAction_SUBSCRIPTION_ACTION_SUBSCRIBE,
	// 			Instruments: []*sdk.CandleInstrument{
	// 				&sdk.CandleInstrument{
	// 					Figi:     "BBG000B9XRY4",
	// 					Interval: sdk.SubscriptionInterval_SUBSCRIPTION_INTERVAL_ONE_MINUTE,
	// 				},
	// 			},
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	log.Printf("error while streaming send, err: %v", err)
	// 	os.Exit(1)
	// }
	// for {
	// 	resp, err := stream.Recv()
	// 	if err != nil {
	// 		log.Printf("error while streaming recv, err: %v", err)
	// 		os.Exit(1)
	// 	}
	// 	candle := resp.GetCandle()
	// 	log.Println(candle)
	// }
	//sdk.MarketDataRequest_SubscribeInfoRequest)
	// shares, err := provider.GetSharesBase(ctx)
	// if err != nil {
	// 	log.Printf("cannot get shares, err: %v", err)
	// 	os.Exit(1)
	// }

	// to := time.Now()
	// controller := model.Controller{SharesByTicker: make(map[string]*model.Share)}
	// for _, share := range shares {
	// 	if share.Exchange != "SPB_MORNING" {
	// 		continue
	// 	}
	// 	history := make(map[sdk.CandleInterval][]*sdk.HistoricCandle)
	// 	for _, tf := range tinkoff.AvailableTF() {
	// 		candles, err := provider.GetCandles(
	// 			ctx,
	// 			share.Figi,
	// 			tf,
	// 			to,
	// 			cfg.ProviderTinkoff.HistoryDepth,
	// 			cfg.ProviderTinkoff.DayOffset.GetDayOffset(tf),
	// 		)
	// 		if err != nil {
	// 			log.Printf("cannot get candles by ticker: %s, time frame: %s, error: %v", share.Ticker, sdk.CandleInterval_name[int32(tf)], err)
	// 			continue
	// 		}
	// 		history[tf] = candles
	// 	}
	// 	_share := model.Share{
	// 		Info:    *share,
	// 		History: history,
	// 	}
	// 	controller.SharesByTicker[share.Ticker] = &_share
	// 	controller.SaveToFile(dataFilePath)
	// 	log.Printf("history for ticker: %s, was written in file: %s", share.Ticker, dataFilePath)
	// }
	log.Println("service finish")
}
