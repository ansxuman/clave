package backend

import (
	"context"
	"fmt"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	ctx    context.Context
	cancel cmap.ConcurrentMap[string, func()]
}

func NewApp() *App {
	return &App{
		cancel: cmap.New[func()](),
	}
}

func (a *App) OnStartup(ctx context.Context, options application.ServiceOptions) error {
	a.ctx = ctx
	fmt.Println("ClaveBackendService is starting up")
	return nil
}

func (a *App) OnShutdown() error {
	fmt.Println("ClaveBackendService is shutting down")
	return nil
}
