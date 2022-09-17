package main_test

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	. "github.com/yuelin-peng/dbsavergo/application/consumer"
	"github.com/yuelin-peng/dbsavergo/domain/saver"
)

var _ = Describe("ConsumeBackupOrder-TableConsumer", func() {
	var (
		ctx   = context.Background()
		t     = GinkgoT()
		order = map[string]interface{}{
			"order_no":    "abc",
			"modify_time": time.Now().Unix(),
			"version":     1,
		}
		config = TableConfig{
			TableName:       "t_test",
			OrderNOField:    "order_no",
			ModifyTimeField: "modify_time",
			VersionField:    "version",
		}
	)
	Describe("订单正常", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		s := saver.NewMockDBSaver(mockCtrl)
		c, err := NewConsumer(ctx, s, config)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		Context("领域正常", func() {
			s.EXPECT().Eliminate(gomock.Any(), gomock.Any()).Return(nil)
			err := c.Consume(ctx, order)
			It("无错误", func() {
				assert.Nil(t, err)
			})
		})
		Context("领域异常", func() {
			s.EXPECT().Eliminate(gomock.Any(), gomock.Any()).Return(fmt.Errorf("something wrong"))
			err := c.Consume(ctx, order)
			It("有错误", func() {
				assert.NotNil(t, err)
			})
		})
	})
	Describe("订单异常", func() {
		defer GinkgoRecover()

		Context("订单为空", func() {
			mockCtrl := gomock.NewController(t)
			s := saver.NewMockDBSaver(mockCtrl)
			c, err := NewConsumer(ctx, s, config)
			assert.Nil(t, err)
			assert.NotNil(t, c)
			err = c.Consume(ctx, nil)
			It("有错误", func() {
				assert.NotNil(t, err)
			})

		})
		Context("订单号为空", func() {
			mockCtrl := gomock.NewController(t)
			s := saver.NewMockDBSaver(mockCtrl)
			c, err := NewConsumer(ctx, s, config)
			assert.Nil(t, err)
			assert.NotNil(t, c)
			tmp := order
			delete(tmp, "order_no")
			err = c.Consume(ctx, tmp)
			It("有错误", func() {
				assert.NotNil(t, err)
			})
		})
	})
})
