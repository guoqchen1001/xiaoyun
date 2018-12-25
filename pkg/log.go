package root

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Log 日志对象
type Log struct {
	Logger *logrus.Logger
}

// NewLogFileOut  创建日志对象
func NewLogFileOut(path string) *Log {

	var l Log

	log := logrus.New()

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Info("打开或创建日志文件失败，将只用标准输入输出流来操作日志")
	} else {
		log.Out = file
	}

	l.Logger = log

	return &l
}

// NewLogStdOut 创建标准输入输出的log对象，用于测试
func NewLogStdOut() *Log {
	l := Log{}
	log := logrus.New()
	l.Logger = log
	return &l
}

// NewLogMultiOut 创建日志输出到文件和标准输出
func NewLogMultiOut(path string) *Log {

	var l Log

	log := logrus.New()

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Info("打开或创建日志文件失败，将只用标准输入输出流来操作日志")
	} else {
		writers := []io.Writer{
			os.Stdout,
			file,
		}
		multiWriter := io.MultiWriter(writers...)

		log.Out = multiWriter
	}

	l.Logger = log

	return &l

}
