package root

// Crypto 加密接口
type Crypto interface {
	Salt(s string) (string, error)               // 加密
	Compare(hash string, s string) (bool, error) // 比较
}
