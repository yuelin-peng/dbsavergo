package saver

import (
	"context"

	"github.com/yuelin-peng/dbsavergo/domain/saver/do"
)

//go:generate mockgen -destination=mock_interface.go -source=interface.go -package=saver DBSaver

type DBSaver interface {
	Query(ctx context.Context, orderNO string) (*do.Order, error)
	Save(ctx context.Context, order *do.Order) error
	Eliminate(ctx context.Context, order *do.Order) error
}
