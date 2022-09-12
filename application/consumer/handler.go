package main

import (
	"context"
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/imdario/mergo"
	"github.com/siddontang/go/log"
	"github.com/yuelin-peng/dbsavergo/domain/saver"
)

type BinLogHandler struct {
	canal.DummyEventHandler
	c   *TableConsumer
	ctx context.Context
}

func NewBinLogHandlerWithConfig(ctx context.Context, dbSaverConfig DBConfig, tableConfig TableConfig) (*BinLogHandler, error) {
	c, err := NewConsumerWithConfig(ctx, dbSaverConfig, tableConfig)
	if err != nil {
		return nil, err
	}
	return &BinLogHandler{
		c:   c,
		ctx: ctx,
	}, nil
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
		c:   c,
		ctx: ctx,
	}, nil
}

func (h *BinLogHandler) OnRow(e *canal.RowsEvent) error {
	log.Infof("get row action=%v, rows=%v, table=%v", e.Action, e.Rows, e.Table)

	value, err := h.rowsToMap(e)
	if err != nil {
		return err
	}
	return h.c.Consume(h.ctx, value)
}

func (h *BinLogHandler) rowsToMap(e *canal.RowsEvent) (map[string]interface{}, error) {
	if e == nil {
		return nil, fmt.Errorf("bin log rows event can't be nil")
	}
	if len(e.Rows) == 0 {
		return nil, fmt.Errorf("bin log rows can't be empty, event=%v", *e)
	}
	switch e.Action {
	case canal.InsertAction, canal.DeleteAction:
		if len(e.Rows) != 1 {
			return nil, fmt.Errorf("bin log rows number for insert is not 1, event=%v", *e)
		}
		return h.rowToMap(e.Table.Columns, e.Rows[0])
	case canal.UpdateAction:
		if len(e.Rows) != 2 {
			return nil, fmt.Errorf("bin log rows number for update is not 2, event=%v", *e)
		}
		oldRow, err := h.rowToMap(e.Table.Columns, e.Rows[0])
		if err != nil {
			return nil, err
		}
		newRow, err := h.rowToMap(e.Table.Columns, e.Rows[1])
		if err != nil {
			return nil, err
		}
		if err = mergo.Merge(&newRow, oldRow); err != nil {
			return nil, err
		}
		return newRow, nil
	default:
		return nil, fmt.Errorf("invalid action=%v", e.Action)
	}
}

func (h *BinLogHandler) rowToMap(columns []schema.TableColumn, row []interface{}) (map[string]interface{}, error) {
	if len(columns) < len(row) {
		return nil, fmt.Errorf("rows event's scheme is not match rows, columns=%v, rows=%v",
			columns, row)
	}
	if len(row) == 0 {
		return nil, fmt.Errorf("row can't be empty, columns=%v, rows=%v",
			columns, row)
	}
	result := make(map[string]interface{}, len(row))
	for i, _ := range row {
		result[columns[i].Name] = row[i]
	}
	return result, nil
}

func (h *BinLogHandler) String() string {
	return "BinLogHandler"
}
