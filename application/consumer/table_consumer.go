package main

import (
	"context"
	"fmt"

	"github.com/siddontang/go/log"
	"github.com/yuelin-peng/dbsavergo/domain/saver"
	"github.com/yuelin-peng/dbsavergo/domain/saver/do"
	"github.com/yuelin-peng/dbsavergo/infrastructure/dao"
	"github.com/yuelin-peng/dbsavergo/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	URL string
}

type TableConfig struct {
	TableName       string
	OrderNOField    string
	ModifyTimeField string
	VersionField    string
}

type TableConsumer struct {
	s                     saver.DBSaver
	orderNOFieldParser    *FieldParser
	modifyTimeFieldParser *FieldParser
	versionFieldParser    *FieldParser
}

func NewConsumerWithConfig(ctx context.Context, config DBConfig, tableConfig TableConfig) (*TableConsumer, error) {
	db, err := gorm.Open(mysql.Open(config.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("create db connenction failed, url=%v", config.URL)
	}

	d, err := dao.NewOrderModel(ctx, db)
	if err != nil {
		return nil, err
	}
	s, err := saver.NewSaver(ctx, d)
	if err != nil {
		return nil, err
	}
	return NewConsumer(ctx, s, tableConfig)
}

func NewConsumer(ctx context.Context, s saver.DBSaver, tableConfig TableConfig) (*TableConsumer, error) {
	if s == nil {
		return nil, fmt.Errorf("NewConsumer failed, saver can't be nil")
	}
	orderNOFieldParser, err := NewFieldParser(tableConfig.OrderNOField)
	if err != nil {
		return nil, err
	}
	modifyTimeFieldParser, err := NewFieldParser(tableConfig.ModifyTimeField)
	if err != nil {
		return nil, err
	}
	versionFieldParser, err := NewFieldParser(tableConfig.VersionField)
	if err != nil {
		return nil, err
	}
	return &TableConsumer{
		s:                     s,
		orderNOFieldParser:    orderNOFieldParser,
		modifyTimeFieldParser: modifyTimeFieldParser,
		versionFieldParser:    versionFieldParser,
	}, nil
}

func (c *TableConsumer) Consume(ctx context.Context, value map[string]interface{}) error {
	order, err := c.convertToOrder(value)
	if err != nil {
		return err
	}
	return c.s.Eliminate(ctx, order)
}

func (c *TableConsumer) convertToOrder(value map[string]interface{}) (*do.Order, error) {
	if len(value) == 0 {
		return nil, fmt.Errorf("value can't be nil")
	}
	log.Infof("convertToOrder:value=%v", value)

	order := &do.Order{}
	if orderNO, err := c.orderNOFieldParser.GetFieldValue(value); err == nil {
		order.OrderNO = orderNO.(string)
	} else if err != nil {
		return nil, fmt.Errorf("missing order number, err=%v", err)
	}
	if modifyTime, err := c.modifyTimeFieldParser.GetFieldValue(value); err == nil {
		order.ModifyTime = util.Interface2Int64(modifyTime)
	} else if err != nil {
		return nil, fmt.Errorf("missing modify number err=%v", err)
	}
	if version, err := c.versionFieldParser.GetFieldValue(value); err == nil {
		order.Version = util.Interface2Int64(version)
	} else if err != nil {
		return nil, fmt.Errorf("missing version field, err=%s", err)
	}
	return order, nil
}
