package main

import (
	"context"
	"log"

	// "github.com/streadway/amqp"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

func main() {
	ctx := context.Background()

	// pubsub.OpenTopic creates a *pubsub.Topic from a URL.
	// This URL will Dial the RabbitMQ server at the URL in the environment
	// variable RABBIT_SERVER_URL and open the exchange "myexchange".
	topic, err := pubsub.OpenTopic(ctx, "rabbit://ex1")
	if err != nil {
		log.Fatal(err)
	}
	defer topic.Shutdown(ctx)

	err = topic.Send(ctx, &pubsub.Message{
		Body: []byte("Hello, World!"),
		// Metadata is optional and can be nil.
		Metadata: map[string]string{
			// These are examples of metadata.
			// There is nothing special about the key names.
			"language":   "en",
			"importance": "high",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
