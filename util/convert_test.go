package util_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/yuelin-peng/dbsavergo/util"
)

var _ = Describe("InterfaceToInt64", func() {
	var (
		t = GinkgoT()
	)
	Describe("int64", func() {
		v := util.InterfaceToInt64(int64(1))
		It("相等", func() {
			assert.Equal(t, int64(1), v)
		})
	})
})
