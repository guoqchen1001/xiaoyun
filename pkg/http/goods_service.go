package http

import (
	"net/http"
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
func NewGoodsHandler(service root.GoodsService, log *root.Log) *GoodsHandler {
	h := &GoodsHandler{
		Router:       httprouter.New(),
		GoodsService: service,
		log:          log,
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
