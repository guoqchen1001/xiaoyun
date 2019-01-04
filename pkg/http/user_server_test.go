package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/http"
	"xiaoyun/pkg/mock"

	"github.com/julienschmidt/httprouter"
)

func TestUser_LoginOK(t *testing.T) {

	mockService := mock.UserService{}

	mockService.LoginFn = func(root.Credentials) (string, error) {
		return "mock_login", nil
	}

	handler := getMockUserHandler(&mockService)

	credentials := root.Credentials{
		UserName: "0122",
		PassWord: "0122",
	}

	var reqByte []byte

	reqByte, err := json.Marshal(credentials)
	if err != nil {
		t.Error(err)
	}

	reqBody := bytes.NewBuffer(reqByte)

	req := httptest.NewRequest("POST", "/api/user/login", reqBody)
	rr := httptest.NewRecorder()
	params := httprouter.Params{}

	handler.UserHandler.HandleLogin(rr, req, params)

	if rr.Code != 200 {
		t.Errorf("http返回值不符合预期，预期[%d],实际[%d]", 200, rr.Code)
	}

}

func TestUser_LoginBadRequest(t *testing.T) {

	handler := getMockUserHandler(&mock.UserService{})

	reqMap := map[string]string{
		"user_name": "0122",
	}

	reqBytes, err := json.Marshal(reqMap)
	if err != nil {
		t.Error(err)
	}

	reqBody := bytes.NewBuffer(reqBytes)

	req := httptest.NewRequest("POST", "/api/user/login", reqBody)
	rr := httptest.NewRecorder()
	params := httprouter.Params{}

	handler.UserHandler.HandleLogin(rr, req, params)

	if rr.Code != 400 {
		t.Errorf("http状态码不符合预期，预期[%d], 实际[%d]", 400, rr.Code)
	}

}

func TestUser_LoginDecodeError(t *testing.T) {

	handler := getMockUserHandler(&mock.UserService{})

	var reqBytes []byte

	reqBody := bytes.NewBuffer(reqBytes)

	req := httptest.NewRequest("POST", "/api/user/login", reqBody)
	rr := httptest.NewRecorder()
	params := httprouter.Params{}

	handler.UserHandler.HandleLogin(rr, req, params)

	if rr.Code != 400 {
		t.Errorf("http状态码不符合预期，预期[%d], 实际[%d]", 400, rr.Code)
	}

}

func TestUser_LoginNilService(t *testing.T) {

	services := http.Services{}
	handler := http.NewHandler(root.NewLogStdOut())
	handler.Init(services)

	credentials := root.Credentials{
		UserName: "0122",
		PassWord: "0122",
	}

	var reqByte []byte

	reqByte, err := json.Marshal(credentials)
	if err != nil {
		t.Error(err)
	}

	reqBody := bytes.NewBuffer(reqByte)

	req := httptest.NewRequest("POST", "/api/user/login", reqBody)
	rr := httptest.NewRecorder()
	params := httprouter.Params{}

	handler.UserHandler.HandleLogin(rr, req, params)

	if rr.Code != 500 {
		t.Errorf("http状态码不符合预期，预期[%d], 实际[%d]", 500, rr.Code)
	}

}

func TestUser_LoginUserOrPwdInvalid(t *testing.T) {

	mockService := mock.UserService{}
	mockService.LoginFn = func(root.Credentials) (string, error) {
		return "", errors.New("mock_login_eror")
	}

	handler := getMockUserHandler(&mockService)
	credentials := root.Credentials{
		UserName: "0122",
		PassWord: "0122",
	}

	var reqByte []byte

	reqByte, err := json.Marshal(credentials)
	if err != nil {
		t.Error(err)
	}

	reqBody := bytes.NewBuffer(reqByte)

	req := httptest.NewRequest("POST", "/api/user/login", reqBody)
	rr := httptest.NewRecorder()
	params := httprouter.Params{}

	handler.UserHandler.HandleLogin(rr, req, params)

	if rr.Code != 403 {
		t.Errorf("http状态码不符合预期，预期[%d], 实际[%d]", 403, rr.Code)
	}

}

func TestUser_CreateOK(t *testing.T) {

	user := root.User{
		No:       "0001",
		Name:     "测试",
		Password: "0122",
	}

	service := mock.UserService{}
	service.CreateUserFn = func(user *root.User) error {
		return nil
	}

	var reqByte []byte
	reqByte, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}

	reqBody := bytes.NewBuffer(reqByte)

	req := httptest.NewRequest("POST", "/api/user", reqBody)
	rr := httptest.NewRecorder()

	handler := getMockUserHandler(&service)

	handler.UserHandler.HandleCreateUser(rr, req, nil)

	if rr.Code != 200 {
		t.Errorf("http状态码不符合预期，预期[%d], 实际[%d]", 200, rr.Code)
	}

}

func TestUser_CreateDecodeError(t *testing.T) {

	reqBody := bytes.NewBufferString("test")
	req := httptest.NewRequest("POST", "/api/user", reqBody)
	rr := httptest.NewRecorder()

	handler := getMockUserHandler(&mock.UserService{})

	handler.UserHandler.HandleCreateUser(rr, req, nil)

	if rr.Code != 400 {
		t.Errorf("http状态码不符合预期，预期[%d]，实际[%d]", 400, rr.Code)
	}

}

func TestUser_CreateError(t *testing.T) {

	service := mock.UserService{}
	service.CreateUserFn = func(user *root.User) error {
		return errors.New("test_user_create_error")
	}

	handler := getMockUserHandler(&service)

	reqByte, _ := json.Marshal(root.User{})
	reqBody := bytes.NewBuffer(reqByte)
	req := httptest.NewRequest("POST", "/api/user", reqBody)
	rr := httptest.NewRecorder()

	handler.UserHandler.HandleCreateUser(rr, req, nil)

	if rr.Code != 500 {
		t.Errorf("http状态码不符合预期，预期[%d]，实际[%d]", 500, rr.Code)
	}

}

func getMockUserHandler(service *mock.UserService) *http.Handler {

	handler := http.NewHandler(root.NewLogStdOut())

	services := http.Services{
		UserService: service,
	}

	handler.Init(services)
	return handler

}
