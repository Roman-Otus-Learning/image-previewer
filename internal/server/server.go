package server

import (
	"context"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"

	"github.com/Roman-Otus-Learning/image-previewer/internal/app"
)

type Server struct {
	server *http.Server
}

func CreateHTTPServer(addr string, app app.App) *Server {
	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: &Handler{app},
		},
	}
}

func (s *Server) Start(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		log.Info().Msg("starting http server on " + s.server.Addr)

		if err := s.server.ListenAndServe(); err != nil {
			log.Error().Err(err).Send()
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	log.Info().Msg("stopping http server")

	if err := s.server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Send()
	}

	log.Info().Msg("http server stopped")
}
