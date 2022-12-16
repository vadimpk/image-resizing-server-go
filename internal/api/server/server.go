package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

//
//func NewServer(cfg *config.Config, handler http.Handler) *Server {
//	return &Server{
//		server: &http.Server{
//			Addr:           ":" + cfg.Server.Port,
//			Handler:        handler,
//			ReadTimeout:    cfg.Server.ReadTimeout,
//			WriteTimeout:   cfg.Server.WriteTimeout,
//			MaxHeaderBytes: cfg.Server.MaxHeaderMegabytes << 20,
//		},
//	}
//}

func NewServer(handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:           ":8080",
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 10 << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
