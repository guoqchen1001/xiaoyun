package root

import "time"

// Image 图片
type Image struct {
	ID       string
	Size     int64
	CreateAt time.Time
}

//GoodsImage 商品图片
type GoodsImage struct {
	Token   string `json:"-"`
	GoodsID int
	Image   []Image
}

// GoodsImageService 商品图片服务
type GoodsImageService interface {
	GetGoodsImage(gooodsID int) (*GoodsImage, error)
	CreateGoodsImage(*GoodsImage) error
}
