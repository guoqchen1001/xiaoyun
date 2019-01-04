package mysql_test

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/mysql"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetGoodsImage(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

	now := time.Now()
	columns := []string{
		"goods_id", "id", "size", "create_at", "create_by",
	}
	rows := sqlmock.NewRows(columns).
		AddRow(1, "asdfghhj", 102400, now, "admin")
	goodsID := 1

	mockSession.mock.ExpectQuery("SELECT (.+) FROM goods_images WHERE").WithArgs(goodsID).WillReturnRows(rows)

	service := mockSession.mysqlSession.GoodsImageService()
	result, err := service.GetGoodsImage(goodsID)

	if err != nil {
		t.Error(err)
		return
	}

	model := mysql.GoodsImageModel{
		GoodsID:  1,
		ID:       "asdfghhj",
		Size:     102400,
		CreateAt: now,
		CreateBy: "admin",
	}

	image := mysql.ToRootGoodsImage([]mysql.GoodsImageModel{model}, goodsID)

	if !reflect.DeepEqual(result.Image, image.Image) {
		t.Errorf("返回值不符预期， 预期：%+v, 实际%+v", image.Image, result.Image)
		return
	}

}

func TestGoodsImage_NotFound(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

	service := mockSession.mysqlSession.GoodsImageService()

	columns := []string{
		"goods_id", "id", "size", "create_at", "create_by",
	}
	rows := sqlmock.NewRows(columns)

	goodsID := 1
	mockSession.mock.ExpectQuery("SELECT (.+) FROM goods_images WHERE").WithArgs(goodsID).WillReturnRows(rows)

	_, err = service.GetGoodsImage(goodsID)

	if root.ErrorCode(err) != root.ENOFOUND {
		t.Errorf("错误返回值不符预期，预期[%s]，实际[%s]", root.ENOFOUND, err.Error())
		return
	}

}

func TestGoodsImage_QueryError(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

	service := mockSession.mysqlSession.GoodsImageService()

	err = errors.New("selecterr")
	goodsID := 1
	mockSession.mock.ExpectQuery("SELECT (.+) FROM `goods_images` WHERE").WithArgs(goodsID).WillReturnError(err)

	_, err = service.GetGoodsImage(goodsID)
	if root.ErrorCode(err) != "db_query_error" {
		t.Errorf("错误返回值不符合预期，预期[%s]，实际[%s]", "db_prepare_error", root.ErrorCode(err))
	}

}

func TestCreateGoodsImage(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

	service := mockSession.mysqlSession.GoodsImageService()
	result := sqlmock.NewResult(0, 1)

	mock := mockSession.mock

	mock.ExpectBegin()
	stmt := mock.ExpectPrepare("INSERT INTO goods_images")
	stmt.ExpectExec().WillReturnResult(result)
	stmt.ExpectExec().WillReturnResult(result)

	mock.ExpectCommit()

	goodsImage := getGoodsImage(2)

	err = service.CreateGoodsImage(goodsImage)

	if err != nil {
		t.Error(err)
	}

}

func TestCreateGoodsImage_ExecError(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

	service := mockSession.mysqlSession.GoodsImageService()

	err = errors.New("execerror")
	result := sqlmock.NewResult(0, 1)

	mockSession.mock.ExpectBegin()
	stmt := mockSession.mock.ExpectPrepare("INSERT INTO goods_images")
	stmt.ExpectExec().WillReturnResult(result)
	stmt.ExpectExec().WillReturnError(err)
	mockSession.mock.ExpectRollback()

	goodsImage := getGoodsImage(3)

	err = service.CreateGoodsImage(goodsImage)
	if root.ErrorCode(err) != "db_exec_error" {
		t.Error(err)
	}

}

func TestCreateGoodsImage_BeginError(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

	service := mockSession.mysqlSession.GoodsImageService()

	err = errors.New("beginerror")

	mockSession.mock.ExpectBegin().WillReturnError(err)

	goodsImage := getGoodsImage(3)

	err = service.CreateGoodsImage(goodsImage)
	if root.ErrorCode(err) != "db_begin_error" {
		t.Error(err)
	}

}

// getGoodsImage 获取包含指定数量图片的商品图片
func getGoodsImage(num int) *root.GoodsImage {

	var goodsImage root.GoodsImage

	goodsID := rand.Int31n(100)
	goodsImage.GoodsID = int(goodsID)

	for i := 0; i < num; i++ {
		goodsImage.Image = append(
			goodsImage.Image, root.Image{
				ID:       fmt.Sprintf("%d", rand.Int31n(1000)),
				CreateAt: time.Now(),
				Size:     rand.Int63n(10000),
			},
		)
	}

	return &goodsImage

}

func TestCreateGoodsImage_AuthError(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

	var authenticator Authenticator

	authenticator.authFn = func(token string) (*root.User, error) {
		return nil, errors.New("auth_error")
	}

	mockSession.mysqlSession.SetAuthenticator(&authenticator)

	goodsImage := getGoodsImage(2)
	service := mockSession.mysqlSession.GoodsImageService()

	err = service.CreateGoodsImage(goodsImage)
	if root.ErrorCode(err) != "auth_error" {
		t.Errorf("错误值不符合预期，预期[%s]，实际[%s]", "auth_error", root.ErrorCode(err))
	}

}

func TestCreateGoodsImage_PrepareError(t *testing.T) {

	mockSession, err := NewMockSession()
	if err != nil {
		t.Error(err)
	}

	defer mockSession.mockDB.Close()
	defer mockSession.sqlxDB.Close()

	service := mockSession.mysqlSession.GoodsImageService()

	err = errors.New("prepare_error")

	mockSession.mock.ExpectBegin()
	mockSession.mock.ExpectPrepare("INSERT INTO goods_images").WillReturnError(err)

	goodsImage := getGoodsImage(3)

	err = service.CreateGoodsImage(goodsImage)
	if root.ErrorCode(err) != "db_prepare_error" {
		t.Error(err)
	}

}
