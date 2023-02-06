package indicator

import (
	"algotrade_service/internal/model/data"
)

type Indicator interface {
	Calc(bars data.Bars) (float64, error)
	GetHistoryParams() (data.TimeFrame, data.TypeOfPrice, data.Depth)
}

func NewIndicator(rawIndicator string) (Indicator, error) {
	// switch typeOfIndicator {
	// case "rsi":
	// 	return NewRSI(args)
	// case "change":
	// 	return NewChange(args)
	// case "change_percent":
	// 	return NewChangePercent(args)
	// default:
	// 	return nil, fmt.Errorf("unsupported indicator: %s", typeOfIndicator)
	// }

	return nil, nil
}
