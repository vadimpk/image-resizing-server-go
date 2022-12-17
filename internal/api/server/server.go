package server

import (
	"context"
	"log"
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

func (s *Server) Run() {
	go func() {
		log.Println("Starting http server...")
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("couldn't start the server: [%s]\n", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := s.server.Shutdown(ctx); err != nil {
			log.Printf("couldn't shutdown http server: [%s]\n", err)
		}
	}()
	return done
}
