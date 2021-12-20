package main

import (
	"context"
	"fmt"
	"log"
	"rebuildCache/internal/api"
	"rebuildCache/internal/db/redisDB"
	"time"
)

func main() {
	const (
		host = "localhost"
		port = "6379"
		url  = "https://habr.com/ru/rss/hub/go"
		ttl  = 30 * time.Second
	)
	client, err := redisDB.NewRedisClient(host, port, ttl)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	const (
		mkey                = "rebuild_cache_key"
		customTagOne        = "python"
		customTagTwo        = "go"
		customTagHabr       = "Habr"
		customTagGeekBrains = "GeekBrains"
	)
	tags := []string{customTagHabr, customTagGeekBrains}
	/**/
	// comment it if you dont want delete tags before work
	for _, v := range append(tags, mkey) {
		client.Client.Del(context.Background(), v)
	}
	/**/
	rebuild := func() (interface{}, []string, error) {
		posts, err := api.FetchContent(url)
		if err != nil {
			return nil, nil, err
		}
		// for lesson example we use here hardcode tags
		return posts, tags, nil
	}
	fmt.Println("FIRST call")
	posts := api.RSS{}
	err = client.GetCache(mkey, &posts, rebuild)
	log.Printf("FIRST result: posts: %v, error: %v\n\n", len(posts.Items), err)
	fmt.Println("SECOND call")
	posts = api.RSS{}
	err = client.GetCache(mkey, &posts, rebuild)
	log.Printf("SECOND result: posts: %v, error: %v\n\n", len(posts.Items), err)
	fmt.Printf("increment tag: %v\n", customTagOne)
	client.Client.Incr(context.Background(), customTagHabr)

	fmt.Println("THIRD call")
	posts = api.RSS{}
	err = client.GetCache(mkey, &posts, rebuild)
	log.Printf("THIRD result: posts: %v, error: %v\n\n", len(posts.Items), err)
}
