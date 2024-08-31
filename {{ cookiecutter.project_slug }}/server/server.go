package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/config"
	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/httplogger"
	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/logutils"
)

type Server struct {
	cfg *config.Config

	failure chan error

	logger *zap.Logger
	server *http.Server
}

func New(cfg *config.Config) (*Server, error) {
	s := &Server{
		cfg:     cfg,
		failure: make(chan error, 1),
		logger:  zap.L(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleHealthcheck)
	mux.Handle("/metrics", promhttp.Handler())
	handler := httplogger.Middleware(s.logger, mux)

	s.server = &http.Server{
		Addr:              cfg.Server.ListenAddress,
		ErrorLog:          logutils.NewHttpServerErrorLogger(s.logger),
		Handler:           handler,
		MaxHeaderBytes:    1024,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	return s, nil
}

func (s *Server) Run() error {
	l := s.logger
	ctx := logutils.ContextWithLogger(context.Background(), l)

	go func() { // run the server
		l.Info("{{ cookiecutter.project_name }} server is going up...",
			zap.String("server_listen_address", s.cfg.Server.ListenAddress),
		)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.failure <- err
		}
		l.Info("{{ cookiecutter.project_name }} server is down")
	}()

	errs := []error{}
	{ // wait until termination or internal failure
		terminator := make(chan os.Signal, 1)
		signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)

		select {
		case stop := <-terminator:
			l.Info("Stop signal received; shutting down...",
				zap.String("signal", stop.String()),
			)
		case err := <-s.failure:
			l.Error("Internal failure; shutting down...",
				zap.Error(err),
			)
			errs = append(errs, err)
		exhaustErrors:
			for { // exhaust the errors
				select {
				case err := <-s.failure:
					l.Error("Extra internal failure",
						zap.Error(err),
					)
					errs = append(errs, err)
				default:
					break exhaustErrors
				}
			}
		}
	}

	{ // stop the server
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			l.Error("{{ cookiecutter.project_name }} server shutdown failed",
				zap.Error(err),
			)
		}
	}

	switch len(errs) {
	default:
		return errors.Join(errs...)
	case 1:
		return errs[0]
	case 0:
		return nil
	}
}
