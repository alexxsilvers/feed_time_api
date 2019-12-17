package core

import (
	"context"
	"net/http"
	"syscall"
	"time"

	"github.com/alexxsilvers/feed_time_api/core/logger"
	"github.com/pkg/errors"

	"github.com/alexxsilvers/feed_time_api/config"
	"github.com/go-chi/chi"

	"github.com/alexxsilvers/feed_time_api/core/closer"
)

type Closer interface {
	Add(f func() error)
	Wait()
	CloseAll()
}

type App struct {
	publicServer chi.Router
	lis          *listeners
	Closer
}

func NewApp() (*App, error) {
	config.ParseConfig()

	a := &App{
		publicServer: nil,
		Closer:       closer.New(syscall.SIGTERM, syscall.SIGINT),
	}

	lis, err := newListeners(config.GetPublicPort())
	if err != nil {
		return nil, errors.Wrap(err, "create listeners")
	}
	a.lis = lis

	return a, nil
}

func (a *App) RegisterPublicServer(r chi.Router) {
	a.publicServer = r
}

func (a *App) Run() error {
	a.runPublicHTTP()

	logger.Infof(context.Background(), "app started, public port: %d", config.GetPublicPort())

	a.Closer.Wait()

	a.Closer.CloseAll()

	return nil
}

func (a *App) runPublicHTTP() {
	publicServer := &http.Server{
		Handler: a.publicServer,
	}

	go func() {
		err := errors.Wrap(publicServer.Serve(a.lis.http), "http.public")
		if err != nil && errors.Cause(err) != http.ErrServerClosed {
			logger.Error(context.Background(), err)
			a.Closer.CloseAll()
		}
	}()

	a.Closer.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GetGracefulTimeoutInSeconds())*time.Second)
		defer cancel()

		logger.Warn(context.Background(), "http.public: waiting stop of traffic")
		time.Sleep(time.Duration(config.GetGracefulDelayInSeconds()) * time.Second)
		logger.Warn(context.Background(), "http.public: shutting down")

		publicServer.SetKeepAlivesEnabled(false)
		err := errors.Wrap(publicServer.Shutdown(ctx), "http.public: error during shutdown")
		if err != nil {
			return err
		}
		logger.Warn(context.Background(), "http.public: gracefully stopped")
		return nil
	})
}
