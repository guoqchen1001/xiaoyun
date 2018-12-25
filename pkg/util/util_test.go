package util_test

import (
	"io/ioutil"
	"os"
	"xiaoyun/pkg/util"

	"testing"
)

func TestFileExists(t *testing.T) {

	// 创建临时文件
	file, err := ioutil.TempFile("", "test")
	fileName := file.Name()

	if err != nil {
		t.Error(err)
	}

	// 检测文件存在
	if exits, err := util.FileExist(fileName); !exits || (err != nil) {
		t.Error("文件存在检测错误")
		return
	}

	// 关闭和删除闻临时文件
	file.Close()

	err = os.Remove(fileName)
	if err != nil {
		t.Error(err)
		return
	}

	// 检测文件不存在
	if exits, err := util.FileExist(fileName); (err != nil) || exits {
		t.Error("文件不存在检测错误")
		return
	}

}
