package http

import (
	"net/http"
	"net/url"
	root "xiaoyun/pkg"

	"github.com/julienschmidt/httprouter"
)

// GoodsHandler 商品处理器
type GoodsHandler struct {
	*httprouter.Router

	GoodsService root.GoodsService
	log          *root.Log
}

// NewGoodsHandler 创建新的商品处理器
func NewGoodsHandler() *GoodsHandler {
	h := &GoodsHandler{
		Router: httprouter.New(),
	}

	h.GET("/api/goods/:no", h.HandleGetGoods)
	return h
}

// HandleGetGoods 获取商品处理器
func (h *GoodsHandler) HandleGetGoods(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	no := p.ByName("no")

	g, err := h.GoodsService.Goods(no)

	if err != nil {
		if root.ErrorCode(err) == root.ENOFOUND {
			NotFound(w)
		} else {
			Error(w, err, http.StatusInternalServerError, h.log)
		}

	} else if g == nil {
		NotFound(w)
	} else {
		encodeJSON(w, g, h.log)
	}

}

// GoodsService 商品服务的http实现，相当于利用http获取商品数据
type GoodsService struct {
	URL *url.URL
}

// Goods 返回引用
func (s *GoodsService) Goods() (*http.Request, error) {

	s.URL = &url.URL{}
	u := s.URL
	u.Path = "/api/goods/"

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
