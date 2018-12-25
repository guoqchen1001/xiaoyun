package http_test

import (
	"testing"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/http"
)

type Config struct{}

func (c *Config) GetConfig() (*root.Config, error) {
	return &root.Config{
		HTTP: &root.HTTPConfig{
			Host: ":2222",
		},
	}, nil
}
func TestServer(t *testing.T) {
	log := root.NewLogStdOut()
	handler := http.NewHandler(log)
	c := Config{}
	http.NewServer(&c, handler)

}
