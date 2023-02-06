package data

import (
	"algotrade_service/internal/model/data"

)

type History interface {
	Get(/*arguments*/) data.Bars
}
type HistoryChange interface {
	Subscribe() <-chan data.Ticker
}

type Info interface {
	Get(/*arguments*/) data.Info
}

// type Task interface {
// 	Add()
// 	Delete()
// }