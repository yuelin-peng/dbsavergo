package client

import (
	"context"
	"fmt"

	"github.com/yuelin-peng/dbsavergo/domain/saver"
	"github.com/yuelin-peng/dbsavergo/domain/saver/do"
	"github.com/yuelin-peng/dbsavergo/infrastructure/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBSaverDirectRPC struct {
	s saver.DBSaverDomain
}

type DBSaverDirectRPCConfig struct {
	WriteURL string
}

func NewDBSaverDirectRPCWithConfig(ctx context.Context, config DBSaverDirectRPCConfig) (*DBSaverDirectRPC, error) {
	db, err := gorm.Open(mysql.Open(config.WriteURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("create db connenction failed, url=%v", config.WriteURL)
	}

	d, err := dao.NewOrderModel(ctx, db)
	if err != nil {
		return nil, err
	}
	s, err := saver.NewSaver(ctx, d)
	if err != nil {
		return nil, err
	}
	return NewDBSaverDirectRPC(s)
}

func NewDBSaverDirectRPC(s saver.DBSaverDomain) (*DBSaverDirectRPC, error) {
	if s == nil {
		return nil, fmt.Errorf("saver domain can't be nil")
	}
	return &DBSaverDirectRPC{
		s: s,
	}, nil
}

func (r *DBSaverDirectRPC) SaveOrder(ctx context.Context, order *Order) error {
	if order == nil {
		return fmt.Errorf("invalid pararm")
	}
	if len(order.OrderNO) == 0 {
		return fmt.Errorf("order number can't be empty")
	}
	if order.Version < 0 {
		return fmt.Errorf("version can't be less than 0, orderNO=%v, Version=%v",
			order.OrderNO, order.Version)
	}
	if order.ModifyTime <= 0 {
		return fmt.Errorf("version can't be less than 0, orderNO=%v, modifyTime=%v",
			order.OrderNO, order.ModifyTime)
	}
	return r.s.Save(ctx, &do.Order{
		OrderNO:    order.OrderNO,
		Version:    order.Version,
		ModifyTime: order.ModifyTime,
	})
}

func (r *DBSaverDirectRPC) QueryOrder(ctx context.Context, orderNO string) (*Order, error) {
	if len(orderNO) == 0 {
		return nil, fmt.Errorf("order number can't be empty")
	}
	o, err := r.s.Query(ctx, orderNO)
	if err != nil {
		return nil, err
	}
	if o == nil {
		return nil, nil
	}
	return &Order{
		OrderNO:    o.OrderNO,
		Version:    o.Version,
		ModifyTime: o.ModifyTime,
	}, nil
}
