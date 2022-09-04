package saver_test

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"

	"github.com/yuelin-peng/dbsavergo/domain/infrastructure/dao"
	"github.com/yuelin-peng/dbsavergo/domain/saver"
	"github.com/yuelin-peng/dbsavergo/domain/saver/do"
)

var _ = Describe("Save", func() {
	var (
		t          = GinkgoT()
		ctx        = context.Background()
		savedOrder = &do.Order{}
		order      = &do.Order{
			OrderNO:    "abc",
			ModifyTime: time.Now().Unix(),
			Version:    1,
		}
	)

	Describe("happy case:订单首次保存成功", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		dao := dao.NewMockOrderDAO(mockCtrl)
		dao.EXPECT().SetNX(gomock.Any(), gomock.Any()).Return(1, nil).
			Do(func(ctx context.Context, newOrder *do.Order) {
				savedOrder = newOrder
			})
		err := createSaver(t, order, dao).Save(ctx)
		It("无错误｜保存信息与入参一致", func() {
			assert.Nil(t, err)
			assert.Equal(t, true, order.IsEqualForReqInfo(savedOrder))
		})
	})
	Describe("happy case:订单非首次保存成功", func() {
		defer GinkgoRecover()

		oldOrder := order
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		d := dao.NewMockOrderDAO(mockCtrl)
		d.EXPECT().SetNX(gomock.Any(), gomock.Any()).Return(0, dao.DuplicateEntry)
		d.EXPECT().QueryByOrderNO(gomock.Any(), gomock.Any()).Return(oldOrder, nil)
		d.EXPECT().SetWithCas(gomock.Any(), gomock.Any(), oldOrder).Return(1, nil).
			Do(func(ctx context.Context, newOrder *do.Order, oldOrder *do.Order) {
				savedOrder = newOrder
			})
		err := createSaver(t, order, d).Save(ctx)
		It("无错误｜保存信息与入参一致|保存信息版本号不小于原版本号", func() {
			assert.Nil(t, err)
			assert.Equal(t, true, order.IsEqualForReqInfo(savedOrder))
			assert.Equal(t, false, oldOrder.IsNewerTo(savedOrder))
		})
	})
	Describe("happy case:订单版本号低于当前数据", func() {
		defer GinkgoRecover()

		order.Version = 1
		oldOrder := *order
		oldOrder.Version = 2
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		d := dao.NewMockOrderDAO(mockCtrl)
		d.EXPECT().SetNX(gomock.Any(), gomock.Any()).Return(0, dao.DuplicateEntry)
		d.EXPECT().QueryByOrderNO(gomock.Any(), gomock.Any()).Return(&oldOrder, nil)
		err := createSaver(t, order, d).Save(ctx)
		It("无错误｜不保存", func() {
			assert.Nil(t, err)
		})
	})
	Describe("bad case:参数异常", func() {
		defer GinkgoRecover()

		invalidOrder := *order
		invalidOrder.Version = -1
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		d := dao.NewMockOrderDAO(mockCtrl)
		err := createSaver(t, &invalidOrder, d).Save(ctx)
		It("参数异常｜订单未保存", func() {
			assert.Equal(t, err, saver.InvalidParam)
		})
	})
	Describe("bad case:存储异常", func() {
		defer GinkgoRecover()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		d := dao.NewMockOrderDAO(mockCtrl)
		d.EXPECT().SetNX(gomock.Any(), gomock.Any()).Return(0, fmt.Errorf("db failed"))
		err := createSaver(t, order, d).Save(ctx)
		It("存储异常｜订单未保存", func() {
			assert.Equal(t, saver.StorageAbnormal, err)
		})
	})
})

func createSaver(t ginkgo.GinkgoTInterface, order *do.Order, dao dao.OrderDAO) *saver.Saver {
	s, err := saver.NewSaver(context.Background(), order, dao)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	return s
}
