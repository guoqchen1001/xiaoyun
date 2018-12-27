package mysql

import (
	root "xiaoyun/pkg"
)

// GoodsImageService 商品图片服务
type GoodsImageService struct {
	session *Session
}

// GetGoodsImage 获取商品图片，实现root.GoodsImageService接口
func (s *GoodsImageService) GetGoodsImage(id int) (*root.GoodsImage, error) {

	var customErr root.Error
	customErr.Op = "mssql.GoodsImageService.GetGoodsImage"

	var models []GoodsImageModel
	err := s.session.db.Select(&models, `SELECT 
											goods_id, 
											id, size, 
											create_at,
											create_by 
										FROM goods_image 
										WHERE goods_id = ?`, id)

	if err != nil {
		customErr.Code = root.EDBQUERYERROR
		customErr.Err = err
		return nil, &customErr
	}

	if len(models) == 0 {
		customErr.Code = root.ENOFOUND
		return nil, &customErr
	}

	goodsImage := ToRootGoodsImage(models, id)
	return &goodsImage, nil

}

// CreateGoodsImage 创建商品图片
func (s *GoodsImageService) CreateGoodsImage(goodsImage *root.GoodsImage) error {

	const op = "mysql.GoodsImageService.CreateGoodsImage"
	customErr := root.Error{Op: op}

	token := goodsImage.Token
	user, err := s.session.authenticator.Authenticate(token)
	if err != nil {
		customErr.Err = err
		customErr.Code = root.EAUTHERROR
		return &customErr
	}

	goodsID := goodsImage.GoodsID
	models := []GoodsImageModel{}

	for _, image := range goodsImage.Image {
		model := GoodsImageModel{
			GoodsID:  goodsID,
			ID:       image.ID,
			Size:     image.Size,
			CreateAt: s.session.now,
			CreateBy: user.No,
		}
		models = append(models, model)
	}

	tx, err := s.session.db.Beginx()
	if err != nil {
		customErr.Code = root.EDBBEGINERROR
		return &customErr
	}

	stmt, err := tx.PrepareNamed(`INSERT INTO goods_image(goods_id, id, size, create_at, create_by))
							  VALUES(:goods_id, :id, :size, :create_at, :create_by )`)
	if err != nil {
		customErr.Code = root.EDBPREPAREERROR
		customErr.Err = err
		return &customErr
	}

	defer stmt.Close()

	for _, m := range models {
		_, err := stmt.Exec(m)
		if err != nil {
			tx.Rollback()
			customErr.Code = root.EDBEXECERROR
			customErr.Err = err
			return &customErr
		}
	}

	tx.Commit()
	return nil

}
