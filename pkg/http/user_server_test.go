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

	log := root.NewLogStdOut()

	handler := http.NewHandler(log)
	handler.UserHandler.UserService = &mockService

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

	log := root.NewLogStdOut()

	handler := http.NewHandler(log)

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

	log := root.NewLogStdOut()

	handler := http.NewHandler(log)

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

	log := root.NewLogStdOut()

	handler := http.NewHandler(log)

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

	log := root.NewLogStdOut()

	handler := http.NewHandler(log)
	handler.UserHandler.UserService = &mockService

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

	handler := http.NewUserHandler()
	handler.UserService = &service

	handler.HandleCreateUser(rr, req, nil)

	if rr.Code != 200 {
		t.Errorf("http状态码不符合预期，预期[%d], 实际[%d]", 200, rr.Code)
	}

}
