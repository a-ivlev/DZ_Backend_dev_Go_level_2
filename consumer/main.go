package main

import (
	"context"
	"lesson07/consumer/internal/infra/msghandler"
	"lesson07/consumer/internal/infra/msgreceiver"
	"os"
	"os/signal"
	"sync"

	// "github.com/streadway/amqp"

	_ "gocloud.dev/pubsub/rabbitpubsub"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	ch := make(chan msghandler.Request, 100)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go msgreceiver.Run(ctx, ch, wg)
	go msghandler.Handler(ctx, ch, wg)
	wg.Wait()
	stop()
}
