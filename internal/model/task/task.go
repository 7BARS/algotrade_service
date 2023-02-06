package task

import "algotrade_service/internal/model/data"

type Task struct {
	name       string
	expression MathExpression
}

type MathExpression struct {
	// indicators []
}

type HistoryParams struct {
	depth       int
	typeOfPrice data.TypeOfPrice
	interval    data.TimeFrame
}

// func (me *MathExpression) GetHistoryParams() HistoryParams {
// 	// return 
// }

type Response struct {
	ticker []data.Ticker
}
