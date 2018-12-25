package util

import "os"

// FileExist 判断文件是否存在，文件读取失败否会error
func FileExist(path string) (bool, error) {

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
