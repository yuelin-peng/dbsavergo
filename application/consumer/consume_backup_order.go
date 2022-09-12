package main

import (
	"context"
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/siddontang/go-log/log"
	"github.com/yuelin-peng/dbsavergo/domain/saver"
	"github.com/yuelin-peng/dbsavergo/domain/saver/do"
	"github.com/yuelin-peng/dbsavergo/infrastructure/dao"
	"github.com/yuelin-peng/dbsavergo/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	return
}

// func ConsumeBinLog(ctx context.Context, tableName string) {
// 	cfg := canal.NewDefaultConfig()
// 	cfg.Addr = fmt.Sprintf("%s:3306", *testHost)
// 	cfg.User = "root"
// 	cfg.HeartbeatPeriod = 200 * time.Millisecond
// 	cfg.ReadTimeout = 300 * time.Millisecond
// 	cfg.Dump.ExecutionPath = "mysqldump"
// 	cfg.Dump.TableDB = "test"
// 	cfg.Dump.Tables = []string{"canal_test"}
// 	c, err := canal.NewCanal(cfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// Register a handler to handle RowsEvent
// 	c.SetEventHandler(&BinLogHandler{})
//
// 	// Start canal
// 	c.Run()
// 	canal.Start()
// }

type BinLogHandler struct {
	canal.DummyEventHandler
	c *TableConsumer
}

func NewBinLogHandler(ctx context.Context, s saver.DBSaver, tableConfig TableConfig) (*BinLogHandler, error) {
	c, err := NewConsumer(ctx, s, tableConfig)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("new consumer failed, consumer is nil")
	}
	return &BinLogHandler{
		c: c,
	}, nil
}

func (h *BinLogHandler) OnRow(e *canal.RowsEvent) error {
	log.Infof("%s %v\n", e.Action, e.Rows)
	return nil
}

func (h *BinLogHandler) String() string {
	return "BinLogHandler"
}

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
	s      saver.DBSaver
	config TableConfig
}

func NewConsumerWithConfig(ctx context.Context, config DBConfig, tableConfig TableConfig) (*TableConsumer, error) {
	db, err := gorm.Open(mysql.Open(config.URL), &gorm.Config{})
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
	return &TableConsumer{
		s:      s,
		config: tableConfig,
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
	order := &do.Order{}
	if orderNO, ok := value[c.config.OrderNOField]; ok {
		order.OrderNO = orderNO.(string)
	} else if !ok {
		return nil, fmt.Errorf("missing order number field=%s", c.config.OrderNOField)
	}
	if modifyTime, ok := value[c.config.ModifyTimeField]; ok {
		order.ModifyTime = util.Interface2Int64(modifyTime)
	} else if !ok {
		return nil, fmt.Errorf("missing modify number field=%s", c.config.ModifyTimeField)
	}
	if version, ok := value[c.config.VersionField]; ok {
		order.Version = util.Interface2Int64(version)
	} else if !ok {
		return nil, fmt.Errorf("missing version field=%s", c.config.VersionField)
	}
	return order, nil
}
