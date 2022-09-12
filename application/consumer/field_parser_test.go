package main_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	. "github.com/yuelin-peng/dbsavergo/application/consumer"
)

var _ = Describe("ConsumeBackupOrder-ConsumeBackupOrder", func() {
	var (
		t           = GinkgoT()
		fieldConfig = "order_no"
	)
	Describe("简单字段", func() {
		defer GinkgoRecover()

		Context("有该字段值", func() {
			p, err := NewFieldParser(fieldConfig)
			assert.Nil(t, err)
			assert.NotNil(t, p)
			v, err := p.GetFieldValue(map[string]interface{}{
				"order_no": "abc",
			})
			It("无错误", func() {
				assert.Nil(t, err)
				assert.Equal(t, "abc", v)
			})
		})
		Context("没有该字段值", func() {
			p, err := NewFieldParser(fieldConfig)
			assert.Nil(t, err)
			assert.NotNil(t, p)
			v, err := p.GetFieldValue(map[string]interface{}{
				"not_exist": "abc",
			})
			It("有错误|nil", func() {
				assert.NotNil(t, err)
				assert.Nil(t, v)
			})
		})
	})
	Describe("unix函数", func() {
		defer GinkgoRecover()

		unixFieldConfig := "modify_time.Unix"
		Context("有该字段值", func() {
			p, err := NewFieldParser(unixFieldConfig)
			assert.Nil(t, err)
			assert.NotNil(t, p)
			v, err := p.GetFieldValue(map[string]interface{}{
				"modify_time": "2022-08-09 11:20:13",
			})
			It("无错误", func() {
				assert.Nil(t, err)
				assert.Equal(t, int64(1660015213), v)
			})
		})
		Context("没有该字段值", func() {
			p, err := NewFieldParser(unixFieldConfig)
			assert.Nil(t, err)
			assert.NotNil(t, p)
			v, err := p.GetFieldValue(map[string]interface{}{
				"not_exist": "2022-08-09 11:20:13",
			})
			It("有错误", func() {
				assert.NotNil(t, err)
				assert.Nil(t, v)
			})
		})
	})
})
