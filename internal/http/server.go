package http

import (
	"context"
	"fmt"
	muxhandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang-skeleton/internal/domain/adaptors"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/pkg/container"
	"net/http"
)

type Server struct {
	conf   *Conf
	server *http.Server
	logger adaptors.Logger
	router *mux.Router
}

func (s *Server) Run() error {
	s.logger.Info(fmt.Sprintf(`server starting on %s`, s.conf.Host))
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Ready() chan bool {
	return nil
}

func (s *Server) Stop() error {
	c, fn := context.WithTimeout(context.Background(), s.conf.Timeouts.ShutdownWait)
	defer fn()
	return s.server.Shutdown(c)
}

func (s *Server) Init(con container.Container) error {
	s.conf = con.GetGlobalConfig(application.ModuleHTTPServer).(*Conf)
	s.logger = con.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(`http-server`))
	s.router = con.Resolve(application.ModuleHTTPRouter).(*Router).router
	s.server = &http.Server{
		Addr: s.conf.Host,
		Handler: muxhandlers.RecoveryHandler(
			muxhandlers.PrintRecoveryStack(true),
			muxhandlers.RecoveryLogger(s.logger),
		)(s.router),
		ReadTimeout:       s.conf.Timeouts.Read,
		ReadHeaderTimeout: s.conf.Timeouts.ReadHeader,
		WriteTimeout:      s.conf.Timeouts.Write,
		IdleTimeout:       s.conf.Timeouts.Idle,
	}

	return nil
}
