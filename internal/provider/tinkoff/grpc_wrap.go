package tinkoff

import (
	"context"
	"crypto/tls"
	"math/rand"
	"time"

	sdk "github.com/tinkoff/invest-api-go-sdk"
	"go.uber.org/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	address = "invest-public-api.tinkoff.ru:443"
)

type GRPCWrap struct {
	Stream      sdk.MarketDataStreamServiceClient
	instruments sdk.InstrumentsServiceClient
	marketData  sdk.MarketDataServiceClient
	operations  sdk.OperationsServiceClient
	users       sdk.UsersServiceClient

	conn *grpc.ClientConn
	MD   metadata.MD

	limiter            ratelimit.Limiter
	magicNumberDaysAgo int
	coolDownTime       time.Duration
}

func NewGRPCWrap(token string, limitPerSecond, magicNumberDaysAgo int, coolDownTime time.Duration) (*GRPCWrap, error) {
	rand.Seed(time.Now().UnixNano())

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{"Authorization": "Bearer " + token})
	return &GRPCWrap{
		Stream:      sdk.NewMarketDataStreamServiceClient(conn),
		instruments: sdk.NewInstrumentsServiceClient(conn),
		marketData:  sdk.NewMarketDataServiceClient(conn),
		operations:  sdk.NewOperationsServiceClient(conn),
		users:       sdk.NewUsersServiceClient(conn),

		conn: conn,
		MD:   md,

		limiter:            ratelimit.New(limitPerSecond),
		magicNumberDaysAgo: magicNumberDaysAgo,
		coolDownTime:       coolDownTime,
	}, nil
}

func (gw *GRPCWrap) NewStreaming(ctx context.Context) (*Streaming, error) {
	childCtx := metadata.NewOutgoingContext(ctx, gw.MD)
	stream, err := gw.Stream.MarketDataStream(childCtx)
	if err != nil {
		return nil, err
	}

	return &Streaming{
		stream: stream,
	}, nil
}
