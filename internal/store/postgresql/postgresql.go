package postgresql

import (
	"algotrade_service/internal/model"
	"context"
	"database/sql"
)

type Postgresql struct {
	db *sql.DB
}

func NewPostgresql(db *sql.DB) *Postgresql {
	return &Postgresql{db: db}
}

func (p *Postgresql) AppendBars(ctx context.Context, figi, ticker string, timeFrame model.TimeFrame, bars model.Bars) {
	panic("not implemented") // TODO: Implement
}


func (p *Postgresql) GetLastBarsByFigi(ctx context.Context, figi string, timeFrame model.TimeFrame, session model.Session, depth int) (model.Bars, error) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgresql) GetLastBarsByTicker(ctx context.Context, ticker string, timeFrame model.TimeFrame, session model.Session, depth int) (model.Bars, error) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgresql) AppendEvent(ctx context.Context, name string, event model.Event) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgresql) DeleteEvent(ctx context.Context, name string, event model.Event) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgresql) GetEvents(ctx context.Context) ([]model.Event, error) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgresql) AppendInfo(ctx context.Context, figi string, info model.Info) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgresql) DeleteInfo(ctx context.Context, figi string) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgresql) GetInfo(ctx context.Context, figi string) (model.Info, error) {
	panic("not implemented") // TODO: Implement
}
