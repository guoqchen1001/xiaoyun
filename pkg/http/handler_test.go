package http_test

import (
	"testing"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/http"
)

func TestHandler_new(t *testing.T) {
	log := root.NewLogStdOut()
	_ = http.NewHandler(log)

}
