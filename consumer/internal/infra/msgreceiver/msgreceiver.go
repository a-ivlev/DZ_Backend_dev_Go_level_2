package msgreceiver

import (
	"context"
	"lesson07/consumer/internal/infra/msghandler"
	"log"
	"os"
	"sync"

	"gocloud.dev/pubsub"
)

func Run(ctx context.Context, chout chan msghandler.Request, wg *sync.WaitGroup) {
	defer wg.Done()

	// pubsub.OpenSubscription creates a *pubsub.Subscription from a URL.
	// This URL will Dial the RabbitMQ server at the URL in the environment
	// variable RABBIT_SERVER_URL and open the queue "myqueue".
	subs, err := pubsub.OpenSubscription(ctx, os.Getenv("Q1"))
	if err != nil {
		log.Fatal(err)
	}
	defer subs.Shutdown(ctx)

	for {
		msg, err := subs.Receive(ctx)
		if err != nil {
			// Errors from Receive indicate that Receive will no longer succeed.
			log.Printf("Receiving message: %v", err)
			break
		}

		chout <- msghandler.Request{
			Message: string(msg.Body),
			Headers: msg.Metadata,
		}

		msg.Ack()
	}
}
