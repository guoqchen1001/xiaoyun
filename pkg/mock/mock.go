package mock

import (
	root "xiaoyun/pkg"
)

// GoodsService 商品服务接口
type GoodsService struct {
	GoodsFn      func(no string) (*root.Goods, error)
	GoodsInvoked bool
}

// Goods 商品查询的mock实现
func (s *GoodsService) Goods(no string) (*root.Goods, error) {
	s.GoodsInvoked = true
	return s.GoodsFn(no)
}
