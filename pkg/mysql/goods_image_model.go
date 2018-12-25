package mysql

import (
	"time"
	root "xiaoyun/pkg"
)

// GoodsImageModel 商品图片
type GoodsImageModel struct {
	GoodsID  int       `db:"goods_id"`
	ID       string    `db:"id"`
	Size     int64     `db:"size"`
	CreateAt time.Time `db:"create_at"`
	CreateBy string    `db:"create_by"`
}

// ToRootGoodsImage 转换为GoodsImage
func ToRootGoodsImage(gms []GoodsImageModel, goodsID int) root.GoodsImage {

	var goodsImage root.GoodsImage
	for _, gm := range gms {
		if gm.GoodsID == goodsID {
			image := root.Image{
				ID:   gm.ID,
				Size: gm.Size,
			}
			goodsImage.Image = append(goodsImage.Image, image)
		}
	}

	goodsImage.GoodsID = goodsID
	return goodsImage
}
