package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"algotrade_service/internal/ta"
)

type Event struct {
	Name           string  `json:"name"`
	ComparisonSign string  `json:"comparison_sign"`
	Value          float64 `json:"value"`

	Indicator ta.Indicator
}

type EventStore struct {
	EventsAND map[string][]Event `json:"events_and"`
	// EventsOR []Event

}

type EventController struct {
	rg *regexp.Regexp

	rwmux  *sync.RWMutex
	events []EventStore
}

func NewEventController() (*EventController, error) {
	rg, err := regexp.Compile(`(([a-z]*)\.([a-z_]*))\(([^)]+)\)[\s|\S](>|<|==|!=|>=|<=)[\s|\S](\d*)([\s|\S]and|or|)`)
	if err != nil {
		return nil, err
	}
	return &EventController{
		rg:    rg,
		rwmux: &sync.RWMutex{},
	}, nil
}

func (ec *EventController) AddNewEventFromRaw(raw string) error {
	// raw := "any.rsi(14, close, 1440) > 70 and any.rsi(4, close, 5) < 30 or any.rsi(3, close, 1) < 20 and aapl.change_percent(close) > 2"
	strSubmatch := ec.rg.FindAllStringSubmatch(raw, 1000)
	lenSubmatch := len(strSubmatch)
	if lenSubmatch == 0 {
		return fmt.Errorf("cannot parse event: %s", raw)
	}
	event := EventStore{
		EventsAND: make(map[string][]Event),
	}
	for i := 0; i < lenSubmatch; i++ {
		booleanOperator := strSubmatch[i][7]
		booleanOperator = strings.ReplaceAll(booleanOperator, " ", "")
		if i != lenSubmatch-1 && booleanOperator != "and" {
			return fmt.Errorf("unsupported boolean operator")
		}

		typeOfIndicator := strSubmatch[i][3]
		argsRaw := strSubmatch[i][4]
		argsRaw = strings.ReplaceAll(argsRaw, " ", "")
		args := strings.Split(argsRaw, ",")
		indicator, err := ta.NewIndicator(typeOfIndicator, args)
		if err != nil {
			return fmt.Errorf("cannot create new indicator, raw: %s, error: %v", raw, err)
		}

		valueRaw := strSubmatch[i][6]
		value, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			return fmt.Errorf("cannot parse value, raw: %s, error: %v", valueRaw, err)
		}

		ticker := strSubmatch[i][2]
		ticker = strings.ToUpper(ticker)
		comparisonSign := strSubmatch[i][5]
		event.EventsAND[ticker] = append(event.EventsAND[ticker],
			Event{
				Name:           typeOfIndicator + "(" + argsRaw + ")",
				ComparisonSign: comparisonSign,
				Indicator:      indicator,
				Value:          value,
			})
	}
	ec.addNewEvent(event)

	return nil
}

func (ec *EventController) addNewEvent(event EventStore) {
	ec.rwmux.Lock()
	defer ec.rwmux.Unlock()

	ec.events = append(ec.events, event)
}

func (ec *EventController) GetEvents() []EventStore {
	ec.rwmux.RLock()
	defer ec.rwmux.RUnlock()

	return ec.events
}
