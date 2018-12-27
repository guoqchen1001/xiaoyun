package http

import (
	"encoding/json"
	"net/http"
	root "xiaoyun/pkg"

	"github.com/julienschmidt/httprouter"
)

// UserHandler 获取token
type UserHandler struct {
	*httprouter.Router
	log           *root.Log
	UserService   root.UserService
	authenticator root.Authenticator
}

// NewUserHandler 返回用户处理器
func NewUserHandler() *UserHandler {

	h := UserHandler{
		Router: httprouter.New(),
	}

	h.POST("/api/user/login", h.HandleLogin)
	return &h
}

// HandleLogin 用户登录
func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var customErr root.Error
	customErr.Op = "http.UserHandler.handleLogin"

	credentials, err := decodeCredentials(r)
	if err != nil {
		customErr.Err = err
		customErr.Code = "json_decode_error"
		Error(w, &customErr, http.StatusBadRequest, h.log)
		return
	}

	if credentials.PassWord == "" || credentials.UserName == "" {
		customErr.Err = err
		Error(w, &customErr, http.StatusBadRequest, h.log)
		return
	}

	if h.UserService == nil {
		customErr.Code = root.ENILSERVICE
		Error(w, &customErr, http.StatusInternalServerError, h.log)
		return
	}

	cookieStr, err := h.UserService.Login(credentials)
	if err != nil {
		customErr.Err = err
		customErr.Message = "用户名或密码错误"
		Error(w, &customErr, http.StatusForbidden, h.log)
		return
	}

	response := loginResponse{
		Token: cookieStr,
	}

	httpCookie := http.Cookie{
		Value: cookieStr,
	}

	JSONWithCookie(w, response, httpCookie, h.log)

}

func decodeCredentials(r *http.Request) (root.Credentials, error) {

	var credentials root.Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		return credentials, err
	}

	return credentials, nil

}

//loginResponse 登录返回
type loginResponse struct {
	Token string `json:"token"`
}