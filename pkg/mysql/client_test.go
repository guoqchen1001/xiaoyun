package mysql_test

import (
	"errors"
	"testing"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/mysql"
)

type Configer struct {
	configFn func() (*root.Config, error)
}

func (c *Configer) GetConfig() (*root.Config, error) {
	return c.configFn()
}

type Client struct {
	Client *mysql.Client
}

func NewClient() *Client {

	var configer Configer

	configer.configFn = func() (*root.Config, error) {
		var config root.Config
		config.Mysql = &root.MysqlConfig{
			Parm: "",
		}
		return &config, nil
	}

	client := mysql.NewClient(&configer)

	return &Client{
		Client: client,
	}

}

func TestClient_Open(t *testing.T) {

	client := NewClient()

	err := client.Client.Open()
	if err != nil {
		t.Error(err)
	}
	defer client.Client.Close()

}

func TestClient_OpenNilMysql(t *testing.T) {

	var configer Configer

	configer.configFn = func() (*root.Config, error) {
		var config root.Config
		return &config, nil
	}

	Client := mysql.NewClient(&configer)

	err := Client.Open()
	if root.ErrorCode(err) != root.ECONFIGMYSQLNOTFOUND {
		t.Error(err)
	}

	defer Client.Close()
}

func TestClient_OpenConfigError(t *testing.T) {

	var configer Configer

	configer.configFn = func() (*root.Config, error) {
		return nil, errors.New("get_config_error")
	}

	Client := mysql.NewClient(&configer)

	err := Client.Open()
	if root.ErrorCode(err) != root.ECONFIGNOTINVALID {
		t.Error(err)
	}

	defer Client.Close()
}

func TestClient_OpenError(t *testing.T) {

	var configer Configer

	configer.configFn = func() (*root.Config, error) {
		var config root.Config
		config.Mysql = &root.MysqlConfig{
			Dialects: "test_not_exists_dialects",
			Parm:     "",
		}
		return &config, nil
	}

	Client := mysql.NewClient(&configer)

	err := Client.Open()

	if root.ErrorCode(err) != root.EDBOPENERROR {
		t.Error(err)
	}

	defer Client.Close()

}

func TestClient_Connect(t *testing.T) {
	client := NewClient()
	client.Client.Connect()

}

func TestClient_Migrate(t *testing.T) {

	client := NewClient()

	log := root.NewLogStdOut()
	err := client.Client.MigrateUp(log)
	if err != nil {
		t.Error(err)
	}

	err = client.Client.Open()
	if err != nil {
		t.Error(err)
	}

	err = client.Client.MigrateUp(log)
	if err != nil {
		t.Error(err)
	}

}
