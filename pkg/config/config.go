package config

import (
	"encoding/json"
	"io"
	"os"
	root "xiaoyun/pkg"
)

// Config 系统服务对象
type Config struct {
	config *root.Config
}

// NewConfig 创建的新的配置对象
func NewConfig(path string) *Config {
	c := Config{}
	c.config = &root.Config{}
	c.config.Path = &root.PathConfig{
		ConfigPath: path,
	}
	return &c
}

// GetConfig 获取系统设置，实现root.Configer接口
func (c *Config) GetConfig() (*root.Config, error) {
	// 配置文件和日志文件
	err := c.getConfFromFile(c.config.Path.ConfigPath)
	if err != nil {
		return c.config, err
	}
	return c.config, nil
}

// getConfFromFile 从文件中获取配置信息
func (c *Config) getConfFromFile(path string) error {

	file, err := getContentFromFile(path)
	defer file.Close()

	if err != nil {
		return err
	}

	err = c.parseFromJSONFile(file)
	if err != nil {
		return err
	}

	return nil
}

// hetContentFromFile 从配置文件中读取内容，需要读取后关闭句柄
func getContentFromFile(path string) (*os.File, error) {

	const op = "config.getContentFromFile"
	var customError root.Error

	file, err := os.Open(path)
	if err != nil {
		customError = root.Error{
			Op:   op,
			Err:  err,
			Code: "config_not_found",
		}
		return file, &customError
	}
	return file, nil
}

// ParseFromJSONFile 从json文件句柄中解析配置信息
func (c *Config) parseFromJSONFile(file io.Reader) error {

	const op = "config.parseFromJSONFile"

	err := json.NewDecoder(file).Decode(&c.config)
	if err != nil {
		customError := root.Error{
			Op:   op,
			Err:  err,
			Code: "config_not_invalid",
		}
		return &customError
	}

	return nil
}
