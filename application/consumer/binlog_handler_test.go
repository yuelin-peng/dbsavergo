package main_test

import (
	"context"
	"fmt"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	. "github.com/yuelin-peng/dbsavergo/application/consumer"
	"github.com/yuelin-peng/dbsavergo/domain/saver"
)

var _ = Describe("ConsumeBackupOrder-Handler", func() {
	var (
		ctx    = context.Background()
		t      = GinkgoT()
		config = TableConfig{
			TableName:       "t_test",
			OrderNOField:    "order_no",
			ModifyTimeField: "modify_time",
			VersionField:    "version",
		}
		e = &canal.RowsEvent{
			Action: canal.InsertAction,
			Rows: [][]interface{}{
				{1, "abc", time.Now().Unix(), 1},
			},
			Table: &schema.Table{
				Columns: []schema.TableColumn{
					{Name: "id"},
					{Name: "order_no"},
					{Name: "modify_time"},
					{Name: "version"},
				},
			},
		}
	)
	Describe("插入场景", func() {
		defer GinkgoRecover()

		Describe("行数据正常", func() {
			Context("消费者正常", func() {
				mockCtrl := gomock.NewController(t)
				s := saver.NewMockDBSaver(mockCtrl)
				s.EXPECT().Eliminate(gomock.Any(), gomock.Any()).Return(nil)
				h, err := NewBinLogHandler(ctx, s, config)
				assert.Nil(t, err)
				assert.NotNil(t, h)
				err = h.OnRow(e)
				It("无错误", func() {
					assert.Nil(t, err)
				})
			})
			Context("消费者异常", func() {
				mockCtrl := gomock.NewController(t)
				s := saver.NewMockDBSaver(mockCtrl)
				s.EXPECT().Eliminate(gomock.Any(), gomock.Any()).Return(fmt.Errorf("something error"))
				h, err := NewBinLogHandler(ctx, s, config)
				assert.Nil(t, err)
				assert.NotNil(t, h)
				err = h.OnRow(e)
				It("有错误", func() {
					assert.NotNil(t, err)
				})
			})
		})
		Describe("行数据异常", func() {
			mockCtrl := gomock.NewController(t)
			s := saver.NewMockDBSaver(mockCtrl)
			h, err := NewBinLogHandler(ctx, s, config)
			assert.Nil(t, err)
			assert.NotNil(t, h)
			tmp := *e
			tmp.Rows = append(tmp.Rows, []interface{}{1, 2})
			err = h.OnRow(&tmp)
			It("有错误", func() {
				assert.NotNil(t, err)
			})
		})
	})
	Describe("修改场景", func() {
		defer GinkgoRecover()

		updateEvent := &canal.RowsEvent{
			Action: canal.UpdateAction,
			Rows: [][]interface{}{
				{1, "abc", time.Now().Unix(), 1},
				{1, "abc", time.Now().Unix(), 2},
			},
			Table: &schema.Table{
				Columns: []schema.TableColumn{
					{Name: "id"},
					{Name: "order_no"},
					{Name: "modify_time"},
					{Name: "version"},
				},
			},
		}

		Describe("行数据正常", func() {
			Context("消费者正常", func() {
				mockCtrl := gomock.NewController(t)
				s := saver.NewMockDBSaver(mockCtrl)
				s.EXPECT().Eliminate(gomock.Any(), gomock.Any()).Return(nil)
				h, err := NewBinLogHandler(ctx, s, config)
				assert.Nil(t, err)
				assert.NotNil(t, h)
				err = h.OnRow(updateEvent)
				It("无错误", func() {
					assert.Nil(t, err)
				})
			})
			Context("消费者异常", func() {
				mockCtrl := gomock.NewController(t)
				s := saver.NewMockDBSaver(mockCtrl)
				s.EXPECT().Eliminate(gomock.Any(), gomock.Any()).Return(fmt.Errorf("something error"))
				h, err := NewBinLogHandler(ctx, s, config)
				assert.Nil(t, err)
				assert.NotNil(t, h)
				err = h.OnRow(updateEvent)
				It("有错误", func() {
					assert.NotNil(t, err)
				})
			})
		})
		Describe("行数据异常", func() {
			mockCtrl := gomock.NewController(t)
			s := saver.NewMockDBSaver(mockCtrl)
			h, err := NewBinLogHandler(ctx, s, config)
			assert.Nil(t, err)
			assert.NotNil(t, h)
			tmp := *updateEvent
			tmp.Rows = append(tmp.Rows, []interface{}{1, 2})
			err = h.OnRow(&tmp)
			It("有错误", func() {
				assert.NotNil(t, err)
			})
		})
	})
	Describe("删除场景", func() {
		defer GinkgoRecover()

		deleteEvent := &canal.RowsEvent{
			Action: canal.DeleteAction,
			Rows: [][]interface{}{
				{1, "abc", time.Now().Unix(), 1},
			},
			Table: &schema.Table{
				Columns: []schema.TableColumn{
					{Name: "id"},
					{Name: "order_no"},
					{Name: "modify_time"},
					{Name: "version"},
				},
			},
		}

		Describe("行数据正常", func() {
			Context("消费者正常", func() {
				mockCtrl := gomock.NewController(t)
				s := saver.NewMockDBSaver(mockCtrl)
				s.EXPECT().Eliminate(gomock.Any(), gomock.Any()).Return(nil)
				h, err := NewBinLogHandler(ctx, s, config)
				assert.Nil(t, err)
				assert.NotNil(t, h)
				err = h.OnRow(deleteEvent)
				It("无错误", func() {
					assert.Nil(t, err)
				})
			})
			Context("消费者异常", func() {
				mockCtrl := gomock.NewController(t)
				s := saver.NewMockDBSaver(mockCtrl)
				s.EXPECT().Eliminate(gomock.Any(), gomock.Any()).Return(fmt.Errorf("something error"))
				h, err := NewBinLogHandler(ctx, s, config)
				assert.Nil(t, err)
				assert.NotNil(t, h)
				err = h.OnRow(deleteEvent)
				It("有错误", func() {
					assert.NotNil(t, err)
				})
			})
		})
		Describe("行数据异常", func() {
			mockCtrl := gomock.NewController(t)
			s := saver.NewMockDBSaver(mockCtrl)
			h, err := NewBinLogHandler(ctx, s, config)
			assert.Nil(t, err)
			assert.NotNil(t, h)
			tmp := *deleteEvent
			tmp.Rows = append(tmp.Rows, []interface{}{1, 2})
			err = h.OnRow(&tmp)
			It("有错误", func() {
				assert.NotNil(t, err)
			})
		})
	})
	Describe("其他场景", func() {
		defer GinkgoRecover()

		otherEvent := &canal.RowsEvent{
			Action: "other",
			Rows: [][]interface{}{
				{1, "abc", time.Now().Unix(), 1},
			},
			Table: &schema.Table{
				Columns: []schema.TableColumn{
					{Name: "id"},
					{Name: "order_no"},
					{Name: "modify_time"},
					{Name: "version"},
				},
			},
		}

		mockCtrl := gomock.NewController(t)
		s := saver.NewMockDBSaver(mockCtrl)
		h, err := NewBinLogHandler(ctx, s, config)
		assert.Nil(t, err)
		assert.NotNil(t, h)
		err = h.OnRow(otherEvent)
		It("有错误", func() {
			assert.NotNil(t, err)
		})
	})
})
