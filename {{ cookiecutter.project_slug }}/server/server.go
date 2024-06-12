package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/config"
	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/httplogger"
	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/logutils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	cfg *config.Config
	log *zap.Logger
}

func New(cfg *config.Config) (*Server, error) {
	l := zap.L()

	return &Server{
		cfg: cfg,
		log: l,
	}, nil
}

func (s *Server) Run() error {
	l := s.log
	ctx := logutils.ContextWithLogger(context.Background(), l)

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleHealthcheck)
	mux.Handle("/metrics", promhttp.Handler())
	handler := httplogger.Middleware(l, mux)

	srv := &http.Server{
		Addr:              s.cfg.Server.ListenAddress,
		ErrorLog:          logutils.NewHttpServerErrorLogger(l),
		Handler:           handler,
		MaxHeaderBytes:    1024,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	failure := make(chan error, 1)

	go func() { // run the server
		l.Info("{{ cookiecutter.project_name }} server is going up...",
			zap.String("server_listen_address", s.cfg.Server.ListenAddress),
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			failure <- err
		}
		l.Info("{{ cookiecutter.project_name }} server is down")
	}()

	{ // wait until termination or internal failure
		terminator := make(chan os.Signal, 1)
		signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)

		select {
		case stop := <-terminator:
			l.Info("Stop signal received; shutting down...",
				zap.String("signal", stop.String()),
			)
		case err := <-failure:
			l.Error("Internal failure; shutting down...",
				zap.Error(err),
			)
		}
	}

	{ // stop the server
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			l.Error("{{ cookiecutter.project_name }} server shutdown failed",
				zap.Error(err),
			)
		}
	}

	return nil
}
