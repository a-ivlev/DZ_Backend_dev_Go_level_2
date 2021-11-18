package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/api/router"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/api/server"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/app/shortenerBL"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/app/starter"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/db/inmemoryDB"
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
