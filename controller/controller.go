package controller

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	config "algotrade_service/configs"
	"algotrade_service/internal/provider/tinkoff"
	"algotrade_service/model"

	sdk "github.com/tinkoff/invest-api-go-sdk"
)

type Controller struct {
	cfg   *config.Config
	token string

	eventController *model.EventController

	chOnUpdateByTicker chan string
	rwmux              *sync.RWMutex
	SharesByTicker     map[string]*model.Share
	SharesByFIGI       map[string]*model.Share

	provider *tinkoff.GRPCWrap
}

func NewController(token string, cfg *config.Config, eventController *model.EventController) (*Controller, error) {
	provider, err := tinkoff.NewGRPCWrap(
		token,
		cfg.ProviderTinkoff.RateLimitPerSecond,
		cfg.ProviderTinkoff.DaysWithEmptyHistory,
		cfg.ProviderTinkoff.TimeOutRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create new provider, err: %v", err)
	}

	return &Controller{
		token:              token,
		cfg:                cfg,
		rwmux:              &sync.RWMutex{},
		chOnUpdateByTicker: make(chan string, 100),
		SharesByTicker:     make(map[string]*model.Share),
		SharesByFIGI:       make(map[string]*model.Share),
		eventController:    eventController,
		provider:           provider,
	}, nil
}

const (
	tokenEnv = "TINKOFF_TOKEN"
)

func (c *Controller) Start() error {
	ctx := context.TODO()
	shares, err := c.provider.GetSharesBase(ctx)
	if err != nil {
		return fmt.Errorf("cannot get shares, err: %v", err)
	}

	streaming, err := c.provider.NewStreaming(ctx)
	if err != nil {
		return fmt.Errorf("cannot get shares, err: %v", err)
	}
	go c.runCandleStreaming(streaming)
	go c.run(ctx, shares, streaming)
	go c.runEventTrigger(streaming)

	return nil
}

func (c *Controller) Stop() {

}

func (c *Controller) run(ctx context.Context, shares []*sdk.Share, streaming *tinkoff.Streaming) error {
	to := time.Now()
	figi := []string{}
	for _, share := range shares {
		figi = append(figi, share.Figi)
	}
	streaming.UnsubscribeCandles(figi, sdk.SubscriptionInterval_SUBSCRIPTION_INTERVAL_ONE_MINUTE)

	for _, share := range shares {
		if share.Exchange != "SPB_MORNING" {
			continue
		}
		history := make(map[sdk.CandleInterval][]*sdk.HistoricCandle)
		for _, interval := range tinkoff.AvailableIntervals() {
			candles, err := c.provider.GetCandles(
				ctx,
				share.Figi,
				interval,
				to,
				c.cfg.ProviderTinkoff.HistoryDepth,
				c.cfg.ProviderTinkoff.DayOffset.GetDayOffset(interval),
			)
			if err != nil {
				log.Printf("cannot get candles by ticker: %s, time frame: %s, error: %v", share.Ticker, sdk.CandleInterval_name[int32(interval)], err)
				continue
			}
			history[interval] = candles
		}
		_share := model.NewShare(*share, history, c.chOnUpdateByTicker)
		c.rwmux.Lock()
		c.SharesByTicker[share.Ticker] = _share
		c.SharesByFIGI[share.Figi] = _share
		c.rwmux.Unlock()

		streaming.SubscribeCandles([]string{share.Figi}, sdk.SubscriptionInterval_SUBSCRIPTION_INTERVAL_ONE_MINUTE)
		log.Printf("ticker: %s, has full history and subscribed to streaming", share.Ticker)
	}

	return nil
}

func (c *Controller) runCandleStreaming(streaming *tinkoff.Streaming) {
	for {
		candle, err := streaming.RecvCandle()
		if err != nil {
			log.Println(err)
		}
		if candle == nil || candle.Figi == "" {
			continue
		}
		c.rwmux.RLock()
		c.SharesByFIGI[candle.Figi].AddMinuteCandle(*candle)
		log.Printf("new update candle for ticker: %s", c.SharesByFIGI[candle.Figi].Info.Ticker)
		c.rwmux.RUnlock()
	}
}

func (c *Controller) runEventTrigger(streaming *tinkoff.Streaming) {
	for {
		select {
		case ticker := <-c.chOnUpdateByTicker:
			c.checkEvent(ticker)
		}
	}
}

const (
	specialTickerANY = "ANY"
)

func (c *Controller) checkEvent(ticker string) {
	c.rwmux.RLock()
	defer c.rwmux.RUnlock()

	// in first line check specific symbols
	// than any
	for _, events := range c.eventController.GetEvents() {
		trigger := true
		for ticker, eventAND := range events.EventsAND {
			if ticker == specialTickerANY {
				continue
			}
			if !c.checkEventByShare(eventAND, ticker) {
				trigger = false
				break
			}
		}
		if !trigger {
			continue
		}
		for _ticker := range c.SharesByTicker {
			trigger := true
			for ticker, eventAND := range events.EventsAND {
				if ticker != specialTickerANY {
					continue
				}
				if !c.checkEventByShare(eventAND, _ticker) {
					trigger = false
					break
				}
			}
			if trigger {
				log.Printf("want to buy ticker: %v", _ticker)
			}
		}

	}
}

func (c *Controller) checkEventByShare(events []model.Event, ticker string) bool {
	share, ok := c.SharesByTicker[ticker]
	if !ok {
		return false
	}
	for _, event := range events {
		ok, err := share.CheckEvent(&event)
		if err != nil {
			log.Printf("cannot check event by ticker: %s, error: %v", ticker, err)
			return false
		}
		if !ok {
			log.Printf("event: %s does not match to ticker: %s", event.Name, ticker)
			return false
		}
	}

	return true
}
