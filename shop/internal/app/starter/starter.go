package starter

import (
	"context"
	"shop/internal/app/itemBL"
	"sync"
)

type App struct {
	*itemBL.ItemStore
}

func NewApp(itemDB itemBL.ItemStores) *App {
	return &App{
		itemBL.NewItemStore(itemDB),
	}
}

type APIserver interface {
	Start()
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, srv APIserver) {
	defer wg.Done()
	srv.Start()
	<-ctx.Done()
	srv.Stop()
}
