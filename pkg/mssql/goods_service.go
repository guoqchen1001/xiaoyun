package mssql

import (
	"database/sql"
	root "xiaoyun/pkg"
)

// GoodsService mssql中商品服务对象
type GoodsService struct {
	session *Session
}

// Goods 通过ID查找商品信息，实现GoodS接口
func (c *GoodsService) Goods(no string) (*root.Goods, error) {

	var customErr root.Error
	customErr.Op = "mssql.GoodsService.Goods"

	if c.session == nil {
		customErr.Code = "service_with_nil_session"
		return nil, &customErr
	}

	if c.session.db == nil {
		customErr.Code = "service_with_nil_db"
		return nil, &customErr
	}

	tx, err := c.session.db.Beginx()

	if err != nil {
		customErr.Code = "db_begin_error"
		customErr.Err = err
		return nil, &customErr
	}

	defer tx.Rollback()

	gm := &goodsNodel{}
	err = tx.Get(gm, "select fitem_id, fitem_no ,fitem_subno, fitem_name,funit_no, fitem_size from t_bi_master where fitem_no = ?", no)
	if err != nil {

		if err == sql.ErrNoRows {
			customErr.Code = root.ENOFOUND
			return nil, &customErr
		}

		customErr.Code = "db_query_error"
		customErr.Err = err
		return nil, &customErr
	}

	goods := gm.toGoods()

	return goods, nil
}
