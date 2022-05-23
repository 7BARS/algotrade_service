package tinkoff

import (
	"context"

	sdk "github.com/tinkoff/invest-api-go-sdk"
	"google.golang.org/grpc/metadata"
)

func (gw *GRPCWrap) GetSharesBase(ctx context.Context) ([]*sdk.Share, error) {
	ir := sdk.InstrumentsRequest{
		InstrumentStatus: sdk.InstrumentStatus_INSTRUMENT_STATUS_BASE,
	}
	childCtx := metadata.NewOutgoingContext(ctx, gw.MD)
	gw.limiter.Take()
	r, err := gw.instruments.Shares(childCtx, &ir)
	if err != nil {
		return nil, err
	}

	return r.GetInstruments(), nil
}
