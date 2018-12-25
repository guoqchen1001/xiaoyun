package http

import (
	"net"
	"net/http"
	root "xiaoyun/pkg"
)

// Server http服务
type Server struct {
	ln net.Listener

	Handler *Handler

	configer root.Configer
}

// NewServer 创建web服务器
func NewServer(config root.Configer, handler *Handler) *Server {
	return &Server{
		configer: config,
		Handler:  handler,
	}
}

// Open 打开web服务
func (s *Server) Open() error {

	const op = "http.Server.Open"
	var customError root.Error
	customError.Op = op

	config, err := s.configer.GetConfig()
	if err != nil {
		customError.Err = err
		return &customError
	}

	if config.HTTP == nil {
		customError.Code = root.ECONFIGHTTPNOTFOUND
		return &customError
	}

	ln, err := net.Listen("tcp", config.HTTP.Host)

	if err != nil {
		customError := root.Error{
			Op:  op,
			Err: err,
		}
		return &customError
	}

	s.ln = ln

	http.Serve(s.ln, s.Handler)

	return nil
}

// Close 关闭web服务
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}

// Port 返回服务器端口，仅在服务器端口打开时有效
func (s *Server) Port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}
