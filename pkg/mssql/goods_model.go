package mssql

import (
	"xiaoyun/pkg"
)

type goodsNodel struct {
	ID      int    `db:"fitem_id"`
	No      string `db:"fitem_no"`
	Barcode string `db:"fitem_subno"`
	Name    string `db:"fitem_name"`
	Size    string `db:"fitem_size"`
	Unit    string `db:"funit_no"`
}

// toRootGoods 转换为root.goods
func (gm *goodsNodel) toGoods() *root.Goods {
	g := root.Goods{}
	g.ID = gm.ID
	g.No = gm.No
	g.Name = gm.Name
	g.Barcode = gm.Barcode
	g.Size = gm.Size
	g.Unit = gm.Unit

	return &g

}
