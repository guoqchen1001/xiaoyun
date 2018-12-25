package root

// Goods 商品信息
type Goods struct {
	ID         int
	No         string
	Barcode    string
	Name       string
	Size       string
	Unit       string
	CategoryNo string
	Brand      string
}

// GoodsService 商品服务
type GoodsService interface {
	Goods(no string) (*Goods, error)
}
