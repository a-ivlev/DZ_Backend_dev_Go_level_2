package starter

import (
	"context"
	"sync"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/app/shortenerBL"
)

type App struct {
	shortenerBL *shortenerBL.ShortenerBL
}

func NewApp(shortenerStore shortenerBL.ShortenerStore) *App {
	app := &App{
		shortenerBL: shortenerBL.NewShotenerBL(shortenerStore),
	}
	return app
}

type APIServer interface {
	Start(shortenerBL *shortenerBL.ShortenerBL)
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs APIServer) {
	defer wg.Done()
	hs.Start(a.shortenerBL)
	<-ctx.Done()
	hs.Stop()
}
