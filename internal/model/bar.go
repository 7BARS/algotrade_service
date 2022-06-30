package model

type Bar struct {
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Volume    int64
	Timestamp int64
}

func (b Bar) IsValid() bool {
	return b.Open != 0 && b.Close != 0 && b.High != 0 && b.Low != 0 && b.Volume != 0 && b.Timestamp != 0
}

type Bars []Bar

func (b Bars) IsValid() bool {
	for _, bar := range b {
		if !bar.IsValid() {
			return false
		}
	}
	return true
}

type BarMsg struct {
	Bar    Bar
	Ticker string
	Figi   string
}
