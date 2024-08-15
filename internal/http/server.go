package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}

	Address string
}

func NewServer(ctx context.Context, address string) *Server {
	return &Server{
		ctx:         ctx,
		Address:     address,
		idleConnsCh: make(chan struct{}),
	}
}

func (s *Server) Run() error {
	handlers := http.NewServeMux()
	handlers.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Alive"))
	})

	srv := &http.Server{
		Addr:         s.Address,
		Handler:      handlers,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	go s.ListenCtxForGT(srv)

	log.Println("server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP got err shutdown %v", err)
		return
	}

	log.Println("HTTP processed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}
