package client

import (
	"context"
	"fmt"
	"time"
)

type client struct {
	r             DBSaverServer
	dbMonitor     DBMonitor
	validDuration int64 // in seconds
}

func New(r DBSaverServer) (*client, error) {
	return &client{
		r:             r,
		validDuration: 24 * 60 * 60,
	}, nil
}

func NewDBSaver(r DBSaverServer, dbMonitor DBMonitor) (*client, error) {
	return &client{
		r:             r,
		dbMonitor:     dbMonitor,
		validDuration: 24 * 60 * 60,
	}, nil
}

func (c *client) SaveOrder(ctx context.Context, order *Order) error {
	if err := c.checkOrder(order); err != nil {
		return InvalidParam
	}
	if order.ModifyTime < time.Now().Unix()-c.validDuration {
		return DiscardError
	}
	err := c.r.SaveOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) checkOrder(order *Order) error {
	if order == nil {
		return fmt.Errorf("order can't be nil")
	}
	if len(order.OrderNO) == 0 {
		return fmt.Errorf("order number can't be empty")
	}
	if order.Version < 0 {
		return fmt.Errorf("order version can't be less than 0")
	}
	if order.ModifyTime <= 0 {
		return fmt.Errorf("order modify time can't be less than 0")
	}

	return nil
}

func (c *client) IsNormalWithCreateTime(ctx context.Context, orderNO string, createTime int64) (bool, error) {
	abnormalTimeList, err := c.dbMonitor.GetDBAbnormalDuration(ctx)
	if err == nil { // 如果db监控器正常，且订单发生在异常之后，直接返回正常
		if !c.orderCreatedBeforeDBAbnormal(abnormalTimeList, createTime) {
			return true, nil
		}
	}
	order, err := c.r.QueryOrder(ctx, orderNO)
	if err != nil {
		return true, nil
	}
	if order == nil {
		return true, nil
	}
	return false, nil
}

func (c *client) orderCreatedBeforeDBAbnormal(abnormalTimeList []TimeInterval, createTime int64) bool {
	for _, interval := range abnormalTimeList {
		if createTime < interval.EndTime {
			return true
		}
	}
	return false
}
