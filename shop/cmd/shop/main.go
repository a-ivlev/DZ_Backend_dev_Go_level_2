package main

import (
	"context"
	"os"
	"os/signal"
	"shop/internal/api/handlers"
	"shop/internal/api/routerchi"
	"shop/internal/api/server"
	"shop/internal/db/elasticSerchDB"
	"sync"

	"shop/internal/app/starter"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	//db := inmemItemDB.NewinmemoryDB()
	db := elasticSerchDB.NewElasticDB()
	a := starter.NewApp(db)

	h := handlers.NewHandlers(a)
	r := routerchi.NewRouterChi(h)
	srv := server.NewServer(":8000", r)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
