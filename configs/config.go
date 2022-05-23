package config

import (
	"time"

	sdk "github.com/tinkoff/invest-api-go-sdk"
)

type Config struct {
	ProviderTinkoff ProviderTinkoff `yaml:"provider_tinkoff"`
}

type ProviderTinkoff struct {
	TimeOutRequest       time.Duration `yaml:"timeout_request"`
	RateLimitPerSecond   int           `yaml:"rate_limit_per_second"`
	DaysWithEmptyHistory int           `yaml:"days_with_empty_history"`
	DayOffset            DayOffset     `yaml:"day_offset"`
	HistoryDepth         int           `yaml:"history_depth"`
}

type DayOffset struct {
	DaysForHistory1    int `yaml:"days_for_history_1"`
	DaysForHistory5    int `yaml:"days_for_history_5"`
	DaysForHistory15   int `yaml:"days_for_history_15"`
	DaysForHistory60   int `yaml:"days_for_history_60"`
	DaysForHistory1440 int `yaml:"days_for_history_1440"`
}

func (do *DayOffset) GetDayOffset(interval sdk.CandleInterval) int {
	switch interval {
	case sdk.CandleInterval_CANDLE_INTERVAL_1_MIN:
		return do.DaysForHistory1
	case sdk.CandleInterval_CANDLE_INTERVAL_5_MIN:
		return do.DaysForHistory5
	case sdk.CandleInterval_CANDLE_INTERVAL_15_MIN:
		return do.DaysForHistory15
	case sdk.CandleInterval_CANDLE_INTERVAL_HOUR:
		return do.DaysForHistory60
	case sdk.CandleInterval_CANDLE_INTERVAL_DAY:
		return do.DaysForHistory1440
	default:
		return 0
	}
}
