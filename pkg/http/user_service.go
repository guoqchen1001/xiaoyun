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
	log         *root.Log
	UserService root.UserService
}

// NewUserHandler 返回用户处理器
func NewUserHandler(service root.UserService, log *root.Log) *UserHandler {

	h := UserHandler{
		Router:      httprouter.New(),
		UserService: service,
		log:         log,
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

// HandleCreateUser 创建用户
func (h UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	user, err := decodeUser(r)
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.log)
		return
	}

	err = h.UserService.CreateUser(&user)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.log)
		return
	}

	encodeJSON(w, user, h.log)

}

func decodeCredentials(r *http.Request) (root.Credentials, error) {

	var credentials root.Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)

	return credentials, err

}

func decodeUser(r *http.Request) (root.User, error) {

	var user root.User

	err := json.NewDecoder(r.Body).Decode(&user)

	return user, err

}

//loginResponse 登录返回
type loginResponse struct {
	Token string `json:"token"`
}
