package http_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/http"
	"xiaoyun/pkg/mock"

	"github.com/julienschmidt/httprouter"
)

// TestGetGoods_goods 测试获取商品
func TestGetGoods_goods(t *testing.T) {

	mockService := mock.GoodsService{}

	goods := root.Goods{
		No:   "01",
		Name: "测试",
	}
	mockService.GoodsFn = func(no string) (*root.Goods, error) {
		return &goods, nil
	}

	log := root.NewLogStdOut()
	handler := http.NewHandler(log)
	handler.GoodsHandler.GoodsService = &mockService

	req := httptest.NewRequest("GET", "/api/goods/", nil)
	rr := httptest.NewRecorder()
	params := httprouter.Params{}
	handler.GoodsHandler.HandleGetGoods(rr, req, params)

	respGoods := root.Goods{}
	err := json.NewDecoder(rr.Body).Decode(&respGoods)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(goods, respGoods) {
		t.Errorf("输入输出值不相等，输入%+v, 输出%+v", goods, respGoods)
	}

}

func TestGetGoods_NotFound(t *testing.T) {

	mockService := testGetGoods_getMockService()

	// 商品处理器
	log := root.NewLogStdOut()
	handler := http.NewHandler(log)

	handler.GoodsHandler.GoodsService = mockService
	// 注入商品服务为测试数据

	// 创建response
	rr := httptest.NewRecorder()
	// 创建客户端
	req := httptest.NewRequest("GET", "/api/goods/", nil)
	params := httprouter.Params{
		httprouter.Param{
			Key:   "no",
			Value: "xxx",
		},
	}

	// 测试获取数据
	handler.GoodsHandler.HandleGetGoods(rr, req, params)

	if status := rr.Code; status != 404 {
		t.Errorf("状态不正确, %d", status)
	}

}

// TestGetGoods_InternalError 测试内部错误
func TestGetGoods_InternalError(t *testing.T) {

	mockService := testGetGoods_getMockService()

	log := root.NewLogFileOut("app.log")
	handler := http.NewHandler(log)
	handler.GoodsHandler.GoodsService = mockService

	rr := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/goods/", nil)
	parms := httprouter.Params{
		httprouter.Param{
			Key:   "no",
			Value: "yyy",
		},
	}

	handler.GoodsHandler.HandleGetGoods(rr, req, parms)

	if rr.Code != 500 {
		t.Errorf("错误类型不正确： %d", rr.Code)
	}

}

func testGetGoods_getMockService() *mock.GoodsService {
	// 测试数据
	mockService := mock.GoodsService{}
	mockService.GoodsFn = func(no string) (*root.Goods, error) {
		if no == "xxx" {
			return nil, nil
		} else if no == "yyy" {
			return nil, errors.New("测试内部错误")
		} else {
			return &root.Goods{
				No:   "123",
				Name: "测试",
			}, nil
		}
	}

	return &mockService
}
