package service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/yuelin-peng/dbsavergo/application/service"
	"github.com/yuelin-peng/dbsavergo/kitex_gen/db_saver_service"
)

var _ = Describe("SaveOrderService", func() {
	var (
		t   = GinkgoT()
		ctx = context.Background()
		req = &db_saver_service.SaveOrderRequest{
			OrderNO: "abc",
			Version: 1,
		}
	)
	Describe("happy case:正常情况", func() {
		resp, err := service.SaveOrder(ctx, req)
		It("无错误｜返回成功", func() {
			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, service.SUCCESS.RetCode, resp.RetCode)
			assert.Equal(t, service.SUCCESS.RetMsg, resp.RetMsg)
			assert.Equal(t, req.OrderNO, resp.OrderNO)
			assert.Equal(t, req.Version, resp.Version)
		})
	})
	Describe("sad case:请求参数异常", func() {
		tmp := *req
		tmp.Version = -1
		resp, err := service.SaveOrder(ctx, &tmp)
		It("无错误｜返回失败", func() {
			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, service.InvalidParam.RetCode, resp.RetCode)
			assert.Equal(t, service.InvalidParam.RetMsg, resp.RetMsg)
		})
	})
	Describe("sad case:创建数据操作对象失败", func() {
		It("无错误｜返回失败", func() {

		})
	})
	Describe("sad case:创建领域对象失败", func() {
		It("无错误｜返回失败", func() {

		})
	})
	Describe("sad case:领域执行失败", func() {
		It("无错误｜返回失败", func() {

		})
	})
})
