package saver

import (
	"context"
	"fmt"

	"github.com/yuelin-peng/dbsavergo/domain/infrastructure/dao"
	"github.com/yuelin-peng/dbsavergo/domain/saver/do"
)

type Saver struct {
	orderDAO dao.OrderDAO
}

var (
	InvalidParam    = fmt.Errorf("invalid param")
	StorageAbnormal = fmt.Errorf("storage abnormal")
)

func NewSaver(ctx context.Context, dao dao.OrderDAO) (*Saver, error) {
	if dao == nil {
		return nil, InvalidParam
	}
	return &Saver{
		orderDAO: dao,
	}, nil
}

func (s *Saver) Query(ctx context.Context, orderNO string) (*do.Order, error) {
	if len(orderNO) == 0 {
		return nil, InvalidParam
	}
	return s.orderDAO.QueryByOrderNO(ctx, orderNO)
}

func (s *Saver) Save(ctx context.Context, order *do.Order) error {
	if err := s.checkOrder(order); err != nil {
		return InvalidParam
	}
	affectRows, err := s.orderDAO.SetNX(ctx, order)
	if err == dao.DuplicateEntry {
		var oldOrder *do.Order
		oldOrder, err = s.orderDAO.QueryByOrderNO(ctx, order.OrderNO)
		if err != nil {
			return StorageAbnormal
		}
		if oldOrder == nil {
			return StorageAbnormal
		}
		if oldOrder.IsNewerTo(order) {
			return nil
		}
		affectRows, err = s.orderDAO.SetWithCas(ctx, order, oldOrder)
	}
	if err != nil {
		return StorageAbnormal
	} else if affectRows != 1 {
		return StorageAbnormal
	}

	return nil
}

func (s *Saver) Eliminate(ctx context.Context, order *do.Order) error {
	return nil
}

func (s *Saver) checkOrder(order *do.Order) error {
	if order == nil {
		return fmt.Errorf("[checkOrder] order can't be nil")
	}
	if len(order.OrderNO) == 0 {
		return fmt.Errorf("[checkOrder] orderNO can't be empty")
	}
	if order.ModifyTime <= 0 {
		return fmt.Errorf("[checkOrder] modify time can't be less than zero, ModifyTime=%v", order.ModifyTime)
	}
	if order.Version < 0 {
		return fmt.Errorf("[checkOrder] version can't be less than zero, Version=%v", order.Version)
	}

	return nil
}
