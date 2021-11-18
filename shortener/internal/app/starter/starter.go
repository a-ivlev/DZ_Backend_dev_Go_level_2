package starter

import (
	"DZ_Backend_dev_Go_level_2/shortener/internal/app/redirectBL"
	"context"
	"sync"
)

type App struct {
	redirectBL *redirectBL.Redirect
}

func NewApp(redirect *redirectBL.Redirect) *App {
	app := &App{
		redirectBL: redirect,
	}
	return app
}

type APIServer interface {
	Start(redirectBL *redirectBL.Redirect)
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs APIServer) {
	defer wg.Done()
	hs.Start(a.redirectBL)
	<-ctx.Done()
	hs.Stop()
}
