package dao

import (
	"context"
	"fmt"
	"time"

	I "github.com/yuelin-peng/dbsavergo/domain/infrastructure/dao"
	"github.com/yuelin-peng/dbsavergo/domain/saver/do"
	"gorm.io/gorm"
)

type OrderModel struct {
	db *gorm.DB
}

type TOrder struct {
	ID         int64  `gorm:"column:id;primary_key"`
	OrderNO    string `gorm:"column:order_no"`
	Version    int64  `gorm:"column:version"`
	CreateTime int64  `gorm:"column:create_time"`
	ModifyTime int64  `gorm:"column:modify_time"`
}

func NewOrderModel(ctx context.Context, db *gorm.DB) (*OrderModel, error) {
	if db == nil {
		return nil, InvalidConn
	}

	return &OrderModel{
		db: db,
	}, nil
}

func (m *OrderModel) TableName() string {
	return "t_order"
}

func (m *OrderModel) QueryByOrderNO(ctx context.Context, orderNO string) (*do.Order, error) {
	if len(orderNO) == 0 {
		return nil, InvalidParam
	}
	var order do.Order
	db := m.db.Table(m.TableName()).Where("order_no = ?", orderNO).First(&order)
	if db.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if db.Error != nil {
		return nil, db.Error
	}
	return &order, nil
}

func (m *OrderModel) SetNX(ctx context.Context, o *do.Order) (int, error) {
	if err := m.checkOrderForSetNX(o); err != nil {
		return 0, InvalidParam
	}
	order := m.toOrderModel(o)
	db := m.db.Table(m.TableName()).Create(order)
	if isDuplicateEntryError(db.Error) {
		return 0, I.DuplicateEntry
	}
	return int(db.RowsAffected), db.Error
}

func (m *OrderModel) SetWithCas(ctx context.Context, n *do.Order, o *do.Order) (int, error) {
	if o == nil {
		return 0, CasNotBeNil
	}
	if err := m.checkOrderForSetNX(o); err != nil {
		return 0, InvalidParam
	}
	if err := m.checkOrderForSetNX(n); err != nil {
		return 0, InvalidParam
	}
	order := m.toOrderModel(n)
	db := m.db.Table(m.TableName()).Where(map[string]interface{}{
		"order_no": o.OrderNO,
		"version":  o.Version,
	}).Updates(order)
	return int(db.RowsAffected), db.Error
}

func (m *OrderModel) toOrderModel(o *do.Order) *TOrder {
	return &TOrder{
		OrderNO:    o.OrderNO,
		Version:    o.Version,
		CreateTime: time.Now().Unix(),
		ModifyTime: time.Now().Unix(),
	}
}
func (m *OrderModel) checkOrderForSetNX(o *do.Order) error {
	if o == nil {
		return fmt.Errorf("order can't be nil")
	}
	if len(o.OrderNO) == 0 {
		return fmt.Errorf("order no can't be empty")
	}
	if o.Version < 0 {
		return fmt.Errorf("version can't be less than 0, OrderNO=%v, Version=%v",
			o.OrderNO, o.Version)
	}
	return nil
}
