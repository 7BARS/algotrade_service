package data

type TimeFrame int

func (tf TimeFrame) IsValid() bool {
	return tf >= 1 && tf <= 1440
}
