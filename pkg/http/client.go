package http

// Client 客户端
type Client struct {
	goodsService GoodsService
}

//NewClient 创建信的客户端
func NewClient() *Client {
	c := &Client{}
	return c
}

// GoodsService 返回goodsService
func (c *Client) GoodsService() *GoodsService {
	return &c.goodsService
}
