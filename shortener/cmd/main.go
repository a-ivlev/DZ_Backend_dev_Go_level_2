package main

import (
	"context"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/internal/api/router"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/internal/api/server"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/internal/app/shortenerBL"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/internal/app/starter"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/internal/db/inmemoryDB"
	"os"
	"os/signal"
	"sync"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	shortenerdb := inmemoryDB.NewShortenerMapDB()
	shortenerBL := shortenerBL.NewShotenerBL(shortenerdb)
	app := starter.NewApp(shortenerdb)
	router := router.NewRouter(shortenerBL)
	srv := server.NewServer(":8035", router)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go app.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
