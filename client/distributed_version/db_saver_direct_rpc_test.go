package client_test

import (
	"context"
	"fmt"
	reflect "reflect"
	"time"

	"github.com/agiledragon/gomonkey"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	client "github.com/yuelin-peng/dbsavergo/client/distributed_version"
	"github.com/yuelin-peng/dbsavergo/domain/saver"
	"github.com/yuelin-peng/dbsavergo/domain/saver/do"
	"gorm.io/gorm"
)

var _ = Describe("DbSaverDirectRpc-NewDBSaverDirectRPCWithConfig", func() {
	var (
		ctx = context.Background()
		t   = GinkgoT()
	)
	Describe("非法url", func() {
		r, err := client.NewDBSaverDirectRPCWithConfig(ctx, client.DBSaverDirectRPCConfig{
			WriteURL: "invalid url",
		})
		It("空｜有错误", func() {
			assert.NotNil(t, err)
			assert.Nil(t, r)
		})
	})
	Describe("正确url", func() {
		patch := gomonkey.ApplyFunc(gorm.Open, func(_ gorm.Dialector, _ ...gorm.Option) (db *gorm.DB, err error) {
			return &gorm.DB{}, nil
		})
		defer patch.Reset()
		r, err := client.NewDBSaverDirectRPCWithConfig(ctx, client.DBSaverDirectRPCConfig{
			WriteURL: "root:123456@tcp(127.0.0.1:3306)/db_saver?charset=utf8mb4&parseTime=True&loc=Local",
		})
		It("非空｜无错误", func() {
			assert.Nil(t, err)
			assert.NotNil(t, r)
		})
	})
})

var _ = Describe("DbSaverDirectRpc-SaveOrder", func() {
	var (
		ctx   = context.Background()
		t     = GinkgoT()
		order = &client.Order{
			OrderNO:    "abc",
			Version:    1,
			ModifyTime: time.Now().Unix(),
		}
	)
	Describe("入参异常", func() {
		defer GinkgoRecover()
		mockCtrl := gomock.NewController(t)
		s := saver.NewMockDBSaverDomain(mockCtrl)

		Context("入参为nil", func() {
			r, err := client.NewDBSaverDirectRPC(s)
			assert.Nil(t, err)
			assert.NotNil(t, r)
			err = r.SaveOrder(ctx, nil)
			It("失败", func() {
				assert.NotNil(t, err)
			})

		})
		Context("入参参数错误", func() {
			tmp := *order
			tmp.Version = -1
			r, err := client.NewDBSaverDirectRPC(s)
			assert.Nil(t, err)
			assert.NotNil(t, r)
			err = r.SaveOrder(ctx, &tmp)
			It("失败", func() {
				assert.NotNil(t, err)
			})
		})
	})
	Describe("入参正常", func() {
		defer GinkgoRecover()
		mockCtrl := gomock.NewController(t)
		s := saver.NewMockDBSaverDomain(mockCtrl)

		Context("领域正常", func() {
			s.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			r, err := client.NewDBSaverDirectRPC(s)
			assert.Nil(t, err)
			assert.NotNil(t, r)
			err = r.SaveOrder(ctx, order)
			It("无错误", func() {
				assert.Nil(t, err)
			})

		})
		Context("领域执行失败", func() {
			s.EXPECT().Save(gomock.Any(), gomock.Any()).Return(fmt.Errorf("something wrong"))
			r, err := client.NewDBSaverDirectRPC(s)
			assert.Nil(t, err)
			assert.NotNil(t, r)
			err = r.SaveOrder(ctx, order)
			It("失败", func() {
				assert.NotNil(t, err)
			})
		})
	})
})

var _ = Describe("DbSaverDirectRpc-QueryOrder", func() {
	var (
		ctx     = context.Background()
		t       = GinkgoT()
		orderNO = "abc"
		order   = &client.Order{
			OrderNO:    "abc",
			Version:    1,
			ModifyTime: time.Now().Unix(),
		}
	)
	Describe("入参异常", func() {
		defer GinkgoRecover()
		mockCtrl := gomock.NewController(t)
		s := saver.NewMockDBSaverDomain(mockCtrl)

		Context("入参为空", func() {
			r, err := client.NewDBSaverDirectRPC(s)
			assert.Nil(t, err)
			assert.NotNil(t, r)
			n, err := r.QueryOrder(ctx, "")
			It("订单为空|失败", func() {
				assert.NotNil(t, err)
				assert.Nil(t, n)
			})

		})
	})
	Describe("入参正常", func() {
		defer GinkgoRecover()
		mockCtrl := gomock.NewController(t)
		s := saver.NewMockDBSaverDomain(mockCtrl)

		Context("领域执行失败", func() {
			s.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("something wrong"))
			r, err := client.NewDBSaverDirectRPC(s)
			assert.Nil(t, err)
			assert.NotNil(t, r)
			n, err := r.QueryOrder(ctx, orderNO)
			It("订单为空|失败", func() {
				assert.NotNil(t, err)
				assert.Nil(t, n)
			})
		})
		Context("领域执行成功，且有订单", func() {
			s.EXPECT().Query(gomock.Any(), orderNO).Return(&do.Order{
				OrderNO:    order.OrderNO,
				Version:    order.Version,
				ModifyTime: order.ModifyTime,
			}, nil)
			r, err := client.NewDBSaverDirectRPC(s)
			assert.Nil(t, err)
			assert.NotNil(t, r)
			n, err := r.QueryOrder(ctx, orderNO)
			It("订单不为空，且与领域一致|无错误", func() {
				assert.Nil(t, err)
				assert.NotNil(t, n)
				assert.Equal(t, true, reflect.DeepEqual(order, n))
			})
		})
		Context("领域执行成功，且无订单", func() {
			s.EXPECT().Query(gomock.Any(), orderNO).Return(nil, nil)
			r, err := client.NewDBSaverDirectRPC(s)
			assert.Nil(t, err)
			assert.NotNil(t, r)
			n, err := r.QueryOrder(ctx, orderNO)
			It("订单为空|无错误", func() {
				assert.Nil(t, err)
				assert.Nil(t, n)
			})
		})
	})
})
