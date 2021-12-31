package redisDB

import (
	"context"
	"encoding/json"
	"fmt"
	"simple-redis/internal/db/redisDB/client"
	"time"
)

func WithStructWork(client *client.RedisClient) error {
	jsonKey := "json_key"
	// comment it if you want data from previous launch
	/**/
	err := client.Del(context.Background(), []string{jsonKey}...).Err()
	if err != nil {
		return err
	}
	/**/
	type exampleStruct struct {
		FieldOne string `json:"field_one"`
		FieldTwo string `json:"field_two"`
	}
	s := exampleStruct{
		FieldOne: "one",
		FieldTwo: "two",
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	ttl := 5 * time.Second
	err = client.Set(context.Background(), jsonKey, data, ttl).Err()
	if err != nil {
		return err
	}
	item, err := client.GetRecord(jsonKey)
	if err != nil {
		return err
	}
	fmt.Printf("GetRecord for key %q `%s`\n", jsonKey, item)

	e := new(exampleStruct)
	if err := json.Unmarshal(item, e); err != nil {
		return err
	}
	fmt.Printf("example struct: %+v\n", e)
	return nil
}
