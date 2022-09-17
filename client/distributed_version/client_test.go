package client_test

import (
	"context"
	"fmt"
	"time"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	client "github.com/yuelin-peng/dbsavergo/client/distributed_version"
)

var _ = Describe("Client-Save", func() {
	var (
		t     = GinkgoT()
		ctx   = context.Background()
		order = &client.Order{
			OrderNO:    "abc",
			Version:    1,
			ModifyTime: time.Now().Unix(),
		}
	)
	Describe("happy case:一切正常", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		r := client.NewMockDBSaverServer(mockCtrl)
		r.EXPECT().SaveOrder(gomock.Any(), order).Return(nil)
		c, err := client.New(r)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		err = c.SaveOrder(ctx, order)
		It("无错误", func() {
			assert.Nil(t, err)
		})
	})
	Describe("happy case:rpc保存失败", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		r := client.NewMockDBSaverServer(mockCtrl)
		r.EXPECT().SaveOrder(gomock.Any(), order).Return(fmt.Errorf("save failed"))
		c, err := client.New(r)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		err = c.SaveOrder(ctx, order)
		It("保存失败", func() {
			assert.NotNil(t, err)
		})
	})
	Describe("happy case:过期", func() {
		defer GinkgoRecover()

		tmp := *order
		tmp.ModifyTime = 10
		mockCtrl := gomock.NewController(t)
		r := client.NewMockDBSaverServer(mockCtrl)
		c, err := client.New(r)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		err = c.SaveOrder(ctx, &tmp)
		It("丢弃", func() {
			assert.Equal(t, client.DiscardError, err)
		})
	})
	Describe("happy case:参数错误", func() {
		defer GinkgoRecover()

		tmp := *order
		tmp.ModifyTime = 0
		mockCtrl := gomock.NewController(t)
		r := client.NewMockDBSaverServer(mockCtrl)
		c, err := client.New(r)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		err = c.SaveOrder(ctx, &tmp)
		It("参数错误", func() {
			assert.Equal(t, client.InvalidParam, err)
		})
	})
})

var _ = Describe("Client-IsNormal", func() {
	var (
		t          = GinkgoT()
		ctx        = context.Background()
		orderNO    = "abc"
		createTime = int64(1662134400) // 2022-09-03 00:00:00
		order      = &client.Order{
			OrderNO:    "abc",
			Version:    1,
			ModifyTime: time.Now().Unix(),
		}
	)
	Describe("happy case:正常情况", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		r := client.NewMockDBSaverServer(mockCtrl)
		dbMonitor := client.NewMockDBMonitor(mockCtrl)
		dbMonitor.EXPECT().GetDBAbnormalDuration(gomock.Any()).Return(nil, nil)
		c, err := client.NewDBSaver(r, dbMonitor)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		normal, err := c.IsNormalWithCreateTime(ctx, orderNO, createTime)
		It("正常|无错误", func() {
			assert.Nil(t, err)
			assert.Equal(t, true, normal)
		})
	})

	Describe("happy case:订单在db异常之前创建,且未同步", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		r := client.NewMockDBSaverServer(mockCtrl)
		r.EXPECT().QueryOrder(gomock.Any(), gomock.Any()).Return(order, nil)
		dbMonitor := client.NewMockDBMonitor(mockCtrl)
		dbMonitor.EXPECT().GetDBAbnormalDuration(gomock.Any()).Return([]client.TimeInterval{
			client.TimeInterval{StartTime: createTime, EndTime: time.Now().Unix()},
		}, nil)
		c, err := client.NewDBSaver(r, dbMonitor)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		normal, err := c.IsNormalWithCreateTime(ctx, orderNO, createTime)
		It("订单异常｜无错误", func() {
			assert.Nil(t, err)
			assert.Equal(t, false, normal)
		})
	})

	Describe("happy case:订单在db异常之前创建,但已同步", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		r := client.NewMockDBSaverServer(mockCtrl)
		r.EXPECT().QueryOrder(gomock.Any(), gomock.Any()).Return(nil, nil)
		dbMonitor := client.NewMockDBMonitor(mockCtrl)
		dbMonitor.EXPECT().GetDBAbnormalDuration(gomock.Any()).Return([]client.TimeInterval{
			client.TimeInterval{StartTime: createTime, EndTime: time.Now().Unix()},
		}, nil)
		c, err := client.NewDBSaver(r, dbMonitor)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		normal, err := c.IsNormalWithCreateTime(ctx, orderNO, createTime)
		It("订单正常｜无错误", func() {
			assert.Nil(t, err)
			assert.Equal(t, true, normal)
		})
	})

	Describe("happy case:订单在db异常之后创建,但未同步", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		r := client.NewMockDBSaverServer(mockCtrl)
		r.EXPECT().QueryOrder(gomock.Any(), gomock.Any()).Return(order, nil)
		dbMonitor := client.NewMockDBMonitor(mockCtrl)
		dbMonitor.EXPECT().GetDBAbnormalDuration(gomock.Any()).Return([]client.TimeInterval{
			client.TimeInterval{StartTime: createTime - 1, EndTime: createTime},
		}, nil)
		c, err := client.NewDBSaver(r, dbMonitor)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		normal, err := c.IsNormalWithCreateTime(ctx, orderNO, createTime)
		It("订单正常｜无错误", func() {
			assert.Nil(t, err)
			assert.Equal(t, true, normal)
		})
	})

	Describe("happy case:订单在db异常之后创建,但已同步", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		r := client.NewMockDBSaverServer(mockCtrl)
		r.EXPECT().QueryOrder(gomock.Any(), gomock.Any()).Return(nil, nil)
		dbMonitor := client.NewMockDBMonitor(mockCtrl)
		dbMonitor.EXPECT().GetDBAbnormalDuration(gomock.Any()).Return([]client.TimeInterval{
			client.TimeInterval{StartTime: createTime - 1, EndTime: createTime},
		}, nil)
		c, err := client.NewDBSaver(r, dbMonitor)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		normal, err := c.IsNormalWithCreateTime(ctx, orderNO, createTime)
		It("订单正常｜无错误", func() {
			assert.Nil(t, err)
			assert.Equal(t, true, normal)
		})
	})

	Describe("happy case:DB监控异常", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		dbMonitor := client.NewMockDBMonitor(mockCtrl)
		dbMonitor.EXPECT().GetDBAbnormalDuration(gomock.Any()).Return(nil, fmt.Errorf("monitor abnormal"))
		Context("订单未同步", func() {
			r := client.NewMockDBSaverServer(mockCtrl)
			r.EXPECT().QueryOrder(gomock.Any(), gomock.Any()).Return(order, nil)
			c, err := client.NewDBSaver(r, dbMonitor)
			assert.Nil(t, err)
			assert.NotNil(t, c)
			normal, err := c.IsNormalWithCreateTime(ctx, orderNO, createTime)
			It("订单异常｜无错误", func() {
				assert.Nil(t, err)
				assert.Equal(t, false, normal)
			})
		})
		dbMonitor.EXPECT().GetDBAbnormalDuration(gomock.Any()).Return(nil, fmt.Errorf("monitor abnormal"))
		Context("订单已同步", func() {
			r := client.NewMockDBSaverServer(mockCtrl)
			r.EXPECT().QueryOrder(gomock.Any(), gomock.Any()).Return(nil, nil)
			c, err := client.NewDBSaver(r, dbMonitor)
			assert.Nil(t, err)
			assert.NotNil(t, c)
			normal, err := c.IsNormalWithCreateTime(ctx, orderNO, createTime)
			It("订单正常｜无错误", func() {
				assert.Nil(t, err)
				assert.Equal(t, true, normal)
			})
		})
	})

	Describe("happy case:DB监控正常,但订单创建在异常之前", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		dbMonitor := client.NewMockDBMonitor(mockCtrl)
		dbMonitor.EXPECT().GetDBAbnormalDuration(gomock.Any()).Return([]client.TimeInterval{
			client.TimeInterval{StartTime: createTime - 1, EndTime: createTime + 1},
		}, nil)
		Context("db saver server异常", func() {
			r := client.NewMockDBSaverServer(mockCtrl)
			r.EXPECT().QueryOrder(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("invalid"))
			c, err := client.NewDBSaver(r, dbMonitor)
			assert.Nil(t, err)
			assert.NotNil(t, c)
			normal, err := c.IsNormalWithCreateTime(ctx, orderNO, createTime)
			It("订单正常｜无错误", func() {
				assert.Nil(t, err)
				assert.Equal(t, true, normal)
			})
		})
	})
})
