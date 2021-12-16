package msghandler

import (
	"context"
	"log"
	"sync"
)

type Request struct {
	ID      string
	Message string
	Headers map[string]string
}

func Handler(ctx context.Context, ch chan Request, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case r := <-ch:
			log.Printf("id=%s body=%s metadata=%+v",
				r.ID, r.Message, r.Headers)
		}
	}
}
