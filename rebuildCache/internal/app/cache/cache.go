package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheStores interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) ([]byte, error)
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
}

type CacheStore struct {
	cache CacheStores
}

func NewCacheStore(cs CacheStores) *CacheStore {
	return &CacheStore{
		cache: cs,
	}
}

func (cs *CacheStore) getCurrentTags(tags []string) (map[string]int, error) {
	currTags, err := cs.cache.MGet(context.Background(), tags...)
	if err != nil {
		return nil, fmt.Errorf("MGET redis for tags %v: %w", tags, err)
	}
	resultTags := make(map[string]int, len(tags))
	now := int(time.Now().Unix())
	for i, tagKey := range tags {
		tagItem := currTags[i]
		// TODO func getCurrentTags
		log.Println("func getCurrentTags tagItem = ", tagItem, " i = ", i)
		if tagItem == nil {
			err := cs.cache.Set(context.Background(), tagKey, now)
			if err != nil {
				return nil, fmt.Errorf("set to redis key-value: %v-%v", tagKey, now)
			}
			resultTags[tagKey] = now
			continue
		}
		data, ok := tagItem.(string)
		if !ok {
			log.Printf("current tags assertion err for %v with type %T", tagItem,
				tagItem)
			continue
		}
		number, err := strconv.Atoi(data)
		if err != nil {
			return nil, err
		}
		resultTags[tagKey] = number
	}
	return resultTags, nil
}

func (cs *CacheStore) rebuild(mkey string, in interface{}, rebuildCb RebuildFunc) error {
	result, tags, err := rebuildCb()
	if err != nil {
		return fmt.Errorf("rebuild cb: %w", err)
	}
	if reflect.TypeOf(result) != reflect.TypeOf(in) {
		return fmt.Errorf("data type mismatch, expected %s, got %s",
			reflect.TypeOf(in), reflect.TypeOf(result))
	}
	currTags, err := cs.getCurrentTags(tags)
	// TODO func rebuild currTags, err := cs.getCurrentTags(tags)
	log.Println("func rebuild currTags, err := cs.getCurrentTags(tags) currTags = ", currTags)
	if err != nil {
		return fmt.Errorf("get current item tags: %w", err)
	}
	cacheData := CacheItemStore{
		Data: result,
		Tags: currTags,
	}
	rawData, err := json.Marshal(cacheData)
	if err != nil {
		return fmt.Errorf("marshal cache item store: %w", err)
	}
	err = cs.cache.Set(context.Background(), mkey, rawData)
	if err != nil {
		return fmt.Errorf("set raw data: %w", err)
	}
	inVal := reflect.ValueOf(in)
	resultVal := reflect.ValueOf(result)
	rv := reflect.Indirect(inVal)
	rvpresult := reflect.Indirect(resultVal)
	rv.Set(rvpresult)
	return nil
}

func (cs *CacheStore) validateTags(itemTags map[string]int) (bool, error) {
	tags := make([]string, 0, len(itemTags))
	for tagKey := range itemTags {
		tags = append(tags, tagKey)
	}
	curr, err := cs.cache.MGet(context.Background(), tags...)
	if err != nil {
		return false, fmt.Errorf("MGET redis for tags %v: %w", tags, err)
	}
	currentTagsMap := make(map[string]int, len(curr))
	for i, tagItem := range curr {
		data, ok := tagItem.(string)
		if !ok {
			log.Printf("validate tags: type assertion err for value %v with type	%T", tagItem, tagItem)
			continue
		}
		number, err := strconv.Atoi(data)
		if err != nil {
			return false, err
		}
		currentTagsMap[tags[i]] = number
	}
	return reflect.DeepEqual(itemTags, currentTagsMap), nil
}

type CacheItem struct {
	Data json.RawMessage
	Tags map[string]int
}
type CacheItemStore struct {
	Data interface{}
	Tags map[string]int
}

type RebuildFunc func() (interface{}, []string, error)

func (cs *CacheStore) GetCache(mkey string, in interface{}, rebuildCb RebuildFunc) (err error) {
	inKind := reflect.ValueOf(in).Kind()
	if inKind != reflect.Ptr {
		return fmt.Errorf("expect pointer, got %s", inKind)
	}

	itemRaw, err := cs.cache.Get(context.Background(), mkey)
	if err == redis.Nil {
		fmt.Println("record not found in cache")
		return cs.rebuild(mkey, in, rebuildCb)
	} else if err != nil {
		return fmt.Errorf("redis: get info for key %v: %w", mkey, err)
	}

	item := new(CacheItem)
	err = json.Unmarshal(itemRaw, item)
	if err != nil {
		return fmt.Errorf("unmarshal to cache item: %w", err)
	}

	tagsValid, err := cs.validateTags(item.Tags)
	if err != nil {
		return fmt.Errorf("validate item tags error %w", err)
	}

	if tagsValid {
		return json.Unmarshal(item.Data, &in)
	}

	return cs.rebuild(mkey, in, rebuildCb)
}
