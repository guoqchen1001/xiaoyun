package mssql_test

import (
	"errors"
	"reflect"
	"testing"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/mssql"

	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGoods(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	rows := sqlmock.NewRows([]string{
		"fitem_id", "fitem_no", "fitem_name", "fitem_subno", "fitem_size", "funit_no"}).
		AddRow(2, "02", "测试2", "0002", "", "瓶")

	no := "02"

	mockSession.mock.ExpectBegin()
	mockSession.mock.ExpectQuery("^select (.+) from t_bi_master where").WithArgs(no).WillReturnRows(rows)
	mockSession.mock.ExpectRollback()

	service := mockSession.mssqlSession.GoodsService()

	result, err := service.Goods(no)
	if err != nil {
		t.Error(err)
		return
	}

	goods := &root.Goods{
		ID:      2,
		Name:    "测试2",
		No:      "02",
		Barcode: "0002",
		Size:    "",
		Unit:    "瓶",
	}

	if !reflect.DeepEqual(result, goods) {
		t.Errorf("商品查询结果不正确, 期待%v，实际%v ", goods, result)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

}

func TestGoods_NotFound(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	rows := sqlmock.NewRows([]string{
		"fitem_id", "fitem_no", "fitem_name", "fitem_subno", "fitem_size", "funit_no"})

	no := "03"

	mockSession.mock.ExpectBegin()
	mockSession.mock.ExpectQuery("^select (.+) from t_bi_master where").WithArgs(no).WillReturnRows(rows)
	mockSession.mock.ExpectRollback()

	service := mockSession.mssqlSession.GoodsService()

	_, err = service.Goods(no)

	if root.ErrorCode(err) != root.ENOFOUND {
		t.Errorf("错误返回值不正确，期待%s, 实际%s", root.ENOFOUND, root.ErrorCode(err))
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

}

func TestGoods_QueryError(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	no := "01"
	err = errors.New("sqlerrtest")

	mockSession.mock.ExpectBegin()
	mockSession.mock.ExpectQuery("^select (.+) from t_bi_master where").WithArgs(no).WillReturnError(err)
	mockSession.mock.ExpectRollback()

	service := mockSession.mssqlSession.GoodsService()
	_, err = service.Goods(no)

	if root.ErrorCode(err) != root.EDBQUERYERROR {
		t.Errorf(`错误类型不正确，期待包含[%s]，实际[%s]`, root.EDBQUERYERROR, root.ErrorCode(err))
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

}

func TestGoods_TransError(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	err = errors.New("beginerr")
	mockSession.mock.ExpectBegin().WillReturnError(err)

	service := mockSession.mssqlSession.GoodsService()
	_, err = service.Goods("")

	if root.ErrorCode(err) != root.EDBBEGINERROR {
		t.Errorf("错误返回值不正确，期待%s, 实际%s", root.EDBBEGINERROR, root.ErrorCode(err))
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

}

func TestGoods_NilSession(t *testing.T) {

	no := "01"
	service := mssql.GoodsService{}

	_, err := service.Goods(no)
	if root.ErrorCode(err) != root.ESERVICEWITHNILSESSION {
		t.Errorf("错误返回值不正确，期待%s, 实际%s", root.ESERVICEWITHNILSESSION, root.ErrorCode(err))
	}

}

func TestGoods_NilDB(t *testing.T) {

	var db *sqlx.DB
	session := mssql.NewSession(db)

	service := session.GoodsService()

	_, err := service.Goods("01")

	if root.ErrorCode(err) != root.ESERVICEWITHNILDB {
		t.Errorf("错误返回值不正确，期待%s, 实际%s", root.ESERVICEWITHNILDB, root.ErrorCode(err))
	}

}
