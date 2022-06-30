package model

type Session string

const (
	Premarket  Session = "premarket"
	Market     Session = "market"
	Postmarket Session = "postmarket"
	All        Session = "all"
)

func (s Session) IsValid() bool {
	return s == Premarket || s == Market || s == Postmarket || s == All
}