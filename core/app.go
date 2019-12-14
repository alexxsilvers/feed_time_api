package core

import (
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
	Closer
}

func NewApp() (*App, error) {
	a := &App{
		publicServer: nil,
		Closer:       closer.New(),
	}

	return a, nil
}

func (a *App) Run() error {
	// todo
	println("olololo")

	a.Closer.Wait()

	a.Closer.CloseAll()

	return nil
}
