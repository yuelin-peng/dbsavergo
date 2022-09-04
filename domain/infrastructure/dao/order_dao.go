package dao

import (
	"context"
	"fmt"

	do "github.com/yuelin-peng/dbsavergo/domain/saver/do"
)

//go:generate mockgen -destination=mock_order_dao.go -source=order_dao.go -package=dao OrderDAO

var (
	DuplicateEntry = fmt.Errorf("duplicate entry")
)

type OrderDAO interface {
	SetNX(ctx context.Context, newOrder *do.Order) (int, error)
	SetWithCas(ctx context.Context, newOrder *do.Order, oldOrder *do.Order) (int, error)
	QueryByOrderNO(ctx context.Context, orderNO string) (*do.Order, error)
}
