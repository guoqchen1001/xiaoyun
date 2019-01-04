package http

import (
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/util"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/nu7hatch/gouuid"
)

// GoodsImageHandler 商品图片处理器
type GoodsImageHandler struct {
	*httprouter.Router
	GoodsImageService root.GoodsImageService
	log               *root.Log
}

// NewGoodsImageHandler 生成商品图片处理器
func NewGoodsImageHandler(service root.GoodsImageService, log *root.Log) *GoodsImageHandler {
	h := GoodsImageHandler{
		Router:            httprouter.New(),
		GoodsImageService: service,
		log:               log,
	}

	h.GET("/api/image/goods/:id", h.handleGoodsImage)
	h.POST("/api/image/goods/:id", h.handleUploadGoodsImage)
	return &h

}

// handleGoodsImage 获取商品图片
func (h *GoodsImageHandler) handleGoodsImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	const op = "http.GoodsImageHandler.GoodsImage"

	id := p.ByName("id")

	GoodsID, err := strconv.Atoi(id)
	if err != nil {
		customErr := root.Error{
			Message: "商品ID不合法",
			Err:     err,
		}
		Error(w, &customErr, http.StatusNotAcceptable, h.log)
		return
	}

	goodsImage, err := h.GoodsImageService.GetGoodsImage(GoodsID)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.log)
		return
	}

	encodeJSON(w, goodsImage, h.log)

}

// handleUploadGoodsImage 上传商品图片
func (h *GoodsImageHandler) handleUploadGoodsImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// 解析token
	r.ParseForm()
	token := ""
	if len(r.Form["token"]) > 0 {
		token = r.Form["token"][0]
	}

	// 解析商品id
	id := p.ByName("id")
	goodsID, err := strconv.Atoi(id)
	if err != nil {
		customErr := root.Error{
			Message: "商品ID不合法",
			Err:     err,
		}
		Error(w, &customErr, http.StatusNotAcceptable, h.log)
		return
	}

	// 解析静态文件,内存存入最大32兆
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploaffile")
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.log)
		return
	}
	defer file.Close()

	// 文件uuid
	fileUUID, err := getUploadFileUUID()
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.log)
		return
	}

	// 保存图片文件到本地
	err = saveFile(file, fileUUID)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.log)
		return
	}

	// 保存图片信息
	goodsImage := root.GoodsImage{}
	goodsImage.Token = token

	goods := root.Goods{}
	goods.ID = goodsID

	image := root.Image{}
	image.Size = handler.Size
	image.ID = fileUUID

	goodsImage.Image = append(goodsImage.Image, image)

	err = h.GoodsImageService.CreateGoodsImage(&goodsImage)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.log)
	}

	encodeJSON(w, goodsImage, h.log)

}

// getUploadFileUUID 获取上传文件的最终保存名城
func getUploadFileUUID() (string, error) {

	u, err := uuid.NewV4()
	if err != nil {
		return "", nil
	}
	return u.String(), nil
}

// 保存上传文件
func saveFile(file io.Reader, fileName string) error {

	// 保存文件目录
	const dir = "upload"
	exists, err := util.FileExist(dir)
	if err != nil {
		return err
	}

	if !exists {
		err = os.Mkdir(dir, 666)
		if err != nil {
			return err
		}
	}

	// 写入文件
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	saveFileName := path.Join(currentDir, dir, fileName)
	f, err := os.OpenFile(saveFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 666)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}
