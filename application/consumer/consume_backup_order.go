package main

import (
	"context"
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/siddontang/go-log/log"
)

var (
	cfg           *canal.Config
	dbSaverConfig DBConfig
	tableConfig   = TableConfig{
		TableName:       "test.canal_test",
		OrderNOField:    "order_no",
		ModifyTimeField: "modify_time",
		VersionField:    "modify_time.Unix",
	}
)

func init() {
	cfg = canal.NewDefaultConfig()
	cfg.Addr = "127.0.0.1:3306"
	cfg.User = "root"
	cfg.Password = "peng823971@"
	cfg.Dump.TableDB = "test"
	cfg.Dump.Tables = []string{"canal_test"}
	cfg.IncludeTableRegex = make([]string, 1)
	cfg.IncludeTableRegex[0] = ".*\\.canal_test"

	dbSaverConfig.URL = "root:peng823971@@tcp(127.0.0.1:3306)/db_saver?charset=utf8mb4&parseTime=True&loc=Local"
}

func main() {
	ctx := context.Background()

	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// create handler
	h, err := NewBinLogHandlerWithConfig(ctx, dbSaverConfig, tableConfig)
	if err != nil {
		log.Fatal(err)
	} else if h == nil {
		log.Fatal(fmt.Errorf("create handler failed without error message"))
	}

	// Register a handler to handle RowsEvent
	c.SetEventHandler(h)

	// Start canal
	fmt.Println(c.Run())
}
