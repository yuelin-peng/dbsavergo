package client

import "context"

type Order struct {
	OrderNO    string
	Version    int64
	ModifyTime int64
}

//go:generate mockgen -destination=mock_db_saver_server.go  -source=rpc.go -package=client DBSaverServer

type DBSaverServer interface {
	SaveOrder(ctx context.Context, order *Order) error
	QueryOrder(ctx context.Context, orderNO string) (*Order, error)
}
