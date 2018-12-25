package http

import (
	"encoding/json"
	"net/http"
	"strings"
	root "xiaoyun/pkg"
)

// Handler 处理器
type Handler struct {
	Log               *root.Log
	GoodsHandler      *GoodsHandler
	GoodsImageHandler *GoodsImageHandler
}

// NewHandler 创建新的处理器
func NewHandler(log *root.Log) *Handler {
	h := &Handler{Log: log}
	h.InitHandler()
	return h
}

// InitHandler 初始化处理器函数
func (h *Handler) InitHandler() {

	h.GoodsHandler = NewGoodsHandler()
	h.GoodsHandler.log = h.Log

	h.GoodsImageHandler = NewGoodsImageHandler()
	h.GoodsImageHandler.log = h.Log
}

// ServeHTTP 开启http服务
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/api/goods") {
		h.GoodsHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/image/goods") {
		h.GoodsImageHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// Error api错误处理
func Error(w http.ResponseWriter, err error, code int, log *root.Log) {

	// 如果错误码 not_found则跳转到not_found
	if root.ErrorCode(err) == root.ENOFOUND {
		NotFound(w)
		return
	}

	customErr := root.Error{}
	// 记录错误日志
	log.Logger.Error(err)
	// 隐藏服务器内部错误
	if code == http.StatusInternalServerError {
		customErr.Code = root.EINTERNAL
	} else {
		customErr.Err = err
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: root.ErrorMessage(&customErr)})
}

// errorResponse 通用错误返回
type errorResponse struct {
	Err string `json:"err,omitempty"`
}

// encodeJson json解析，解析错误时返回内部错误
func encodeJSON(w http.ResponseWriter, v interface{}, log *root.Log) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, log)
	}
}

// NotFound 未找到记录处理.
func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{}` + "\n"))
}
