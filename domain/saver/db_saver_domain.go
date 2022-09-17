package saver

import (
	"context"

	"github.com/yuelin-peng/dbsavergo/domain/saver/do"
)

//go:generate mockgen -destination=mock_db_saver_domain.go -source=db_saver_domain.go -package=saver DBSaverDomain

type DBSaverDomain interface {
	Save(ctx context.Context, order *do.Order) error
	Eliminate(ctx context.Context, order *do.Order) (*do.Order, error)
	Query(ctx context.Context, orderNO string) (*do.Order, error)
}
