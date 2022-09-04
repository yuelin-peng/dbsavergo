package dao_test

import (
	"context"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	I "github.com/yuelin-peng/dbsavergo/domain/infrastructure/dao"
	"github.com/yuelin-peng/dbsavergo/domain/saver/do"
	"github.com/yuelin-peng/dbsavergo/infrastructure/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ = Describe("OrderModel-SetNX", func() {
	var (
		order = &do.Order{
			OrderNO:    "abc",
			Version:    1,
			ModifyTime: time.Now().Unix(),
		}
		ctx = context.Background()
		t   = GinkgoT()
	)

	Describe("[SetNX]happy case:写入数据正常", func() {
		defer GinkgoRecover()

		gormDB, mock := createMockGormDB(t)
		d := createMockOrderModel(t, gormDB)
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `" + d.TableName() + "`.+").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit() // 有顺序要求
		affectedRows, err := d.SetNX(ctx, order)
		It("影响行数=1｜无错误", func() {
			assert.Nil(t, err)
			assert.Equal(t, 1, affectedRows)
		})
	})
	Describe("[SetNX]sad case:数据已存在", func() {
		defer GinkgoRecover()

		gormDB, mock := createMockGormDB(t)
		d := createMockOrderModel(t, gormDB)
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `" + d.TableName() + "`.+").WillReturnError(I.DuplicateEntry)
		mock.ExpectRollback() // 有顺序要求
		affectedRows, err := d.SetNX(ctx, order)
		It("影响行数=0｜重复", func() {
			assert.Equal(t, I.DuplicateEntry, err)
			assert.Equal(t, 0, affectedRows)
		})
	})
	Describe("[SetNX]sad case:参数错误", func() {
		defer GinkgoRecover()

		tmpOrder := *order
		tmpOrder.OrderNO = ""
		gormDB, _ := createMockGormDB(t)
		d := createMockOrderModel(t, gormDB)
		affectedRows, err := d.SetNX(ctx, &tmpOrder)
		It("影响行数=0｜参数错误", func() {
			assert.Equal(t, dao.InvalidParam, err)
			assert.Equal(t, 0, affectedRows)
		})
	})
	Describe("[SetNX]sad case:其他异常", func() {
		defer GinkgoRecover()

		gormDB, mock := createMockGormDB(t)
		d := createMockOrderModel(t, gormDB)
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `" + d.TableName() + "`.+").WillReturnError(dao.DBAbnormal)
		mock.ExpectRollback() // 有顺序要求
		affectedRows, err := d.SetNX(ctx, order)
		It("影响行数=0｜其他异常", func() {
			assert.NotNil(t, err)
			assert.Equal(t, 0, affectedRows)
		})
	})
})

var _ = Describe("OrderModel-SetWithCas", func() {
	var (
		order = &do.Order{
			OrderNO:    "abc",
			Version:    1,
			ModifyTime: time.Now().Unix(),
		}
		ctx = context.Background()
		t   = GinkgoT()
	)

	Describe("[SetWithCas]happy case:更新数据正常", func() {
		defer GinkgoRecover()

		oldOrder := *order
		oldOrder.Version = 1
		gormDB, mock := createMockGormDB(t)
		d := createMockOrderModel(t, gormDB)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `" + d.TableName() + "`.+").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit() // 有顺序要求
		affectedRows, err := d.SetWithCas(ctx, order, &oldOrder)
		It("影响行数=1｜无错误", func() {
			assert.Nil(t, err)
			assert.Equal(t, 1, affectedRows)
		})
	})
	Describe("[SetWithCas]sad case:数据不存在", func() {
		defer GinkgoRecover()

		oldOrder := *order
		oldOrder.Version = 1
		gormDB, mock := createMockGormDB(t)
		d := createMockOrderModel(t, gormDB)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `" + d.TableName() + "`.+").WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectCommit() // 有顺序要求
		affectedRows, err := d.SetWithCas(ctx, order, &oldOrder)
		It("影响行数=0｜无错误", func() {
			assert.Nil(t, err)
			assert.Equal(t, 0, affectedRows)
		})
	})
	Describe("[SetWithCas]sad case:cas未空", func() {
		defer GinkgoRecover()

		gormDB, _ := createMockGormDB(t)
		d := createMockOrderModel(t, gormDB)
		affectedRows, err := d.SetWithCas(ctx, order, nil)
		It("影响行数=0｜cas不能空", func() {
			assert.Equal(t, dao.CasNotBeNil, err)
			assert.Equal(t, 0, affectedRows)
		})
	})
	Describe("[SetWithCas]sad case:参数错误", func() {
		defer GinkgoRecover()

		oldOrder := *order
		oldOrder.Version = -1
		gormDB, _ := createMockGormDB(t)
		d := createMockOrderModel(t, gormDB)
		affectedRows, err := d.SetWithCas(ctx, order, &oldOrder)
		It("影响行数=0｜参数错误", func() {
			assert.Equal(t, dao.InvalidParam, err)
			assert.Equal(t, 0, affectedRows)
		})
	})
	Describe("[SetWithCas]sad case:其他异常", func() {
		defer GinkgoRecover()

		oldOrder := *order
		oldOrder.Version = 1
		gormDB, mock := createMockGormDB(t)
		d := createMockOrderModel(t, gormDB)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `" + d.TableName() + "`.+").WillReturnError(dao.DBAbnormal)
		mock.ExpectRollback() // 有顺序要求
		affectedRows, err := d.SetWithCas(ctx, order, &oldOrder)
		It("影响行数=0｜其他异常", func() {
			assert.Equal(t, dao.DBAbnormal, err)
			assert.Equal(t, 0, affectedRows)
		})
	})
})

func createMockGormDB(t ginkgo.GinkgoTInterface) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	assert.NotNil(t, mock)
	assert.NotNil(t, db)
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
		AddRow("5.7.6"))
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, gormDB)

	return gormDB, mock
}

func createMockOrderModel(t ginkgo.GinkgoTInterface, db *gorm.DB) *dao.OrderModel {
	d, err := dao.NewOrderModel(context.Background(), db)
	assert.Nil(t, err)
	assert.NotNil(t, d)
	return d
}
