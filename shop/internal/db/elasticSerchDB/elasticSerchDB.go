package elasticSerchDB

import (
	"bytes"
	"context"
	"encoding/json"
	"shop/internal/app/itemBL"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/google/uuid"
)

var _ itemBL.ItemStores = &elasticDB{}

type elasticDB struct {
	client *elasticsearch.Client
}

func NewElasticDB() *elasticDB {
	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil
	}
	return &elasticDB{
		client: esClient,
	}
}

func (es *elasticDB) CreateItem(ctx context.Context, item itemBL.ItemBL) (*uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	uid := item.ID.String()
	jsonString, _ := json.Marshal(item)

	request := esapi.IndexRequest{
		Index:      "items",
		DocumentID: uid,
		Body:       strings.NewReader(string(jsonString)),
	}
	request.Do(context.Background(), es.client)

	return &item.ID, nil
}

func (es *elasticDB) GetItem(ctx context.Context, ID uuid.UUID) (*itemBL.ItemBL, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	request := esapi.GetRequest{Index: "items", DocumentID: ID.String()}
	response, _ := request.Do(context.Background(), es.client)

	var item = itemBL.ItemBL{}

	var results map[string]interface{}
	json.NewDecoder(response.Body).Decode(&results)

	res := results["_source"].(map[string]interface{})

	item.ID, _ = uuid.Parse(res["id"].(string))
	item.Name = res["name"].(string)
	item.Price = int64(res["price"].(float64))
	item.CreatedAt, _ = time.Parse(time.RFC3339Nano, res["created_at"].(string))
	if res["updated_at"] != nil {
		item.UpdatedAt, _ = time.Parse(time.RFC3339Nano, res["updated_at"].(string))
	}
	if res["deleted_at"] != nil {
		item.DeletedAt, _ = time.Parse(time.RFC3339Nano, res["deleted_at"].(string))
	}

	return &item, nil
}

func (es *elasticDB) DeleteItem(ctx context.Context, ID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	request := esapi.GetRequest{Index: "items", DocumentID: ID.String()}
	response, _ := request.Do(context.Background(), es.client)

	var results map[string]interface{}
	json.NewDecoder(response.Body).Decode(&results)

	res := results["_source"].(map[string]interface{})

	var item = &itemBL.ItemBL{}
	item.ID, _ = uuid.Parse(res["id"].(string))
	item.Name = res["name"].(string)
	item.Price = int64(res["price"].(float64))
	item.CreatedAt, _ = time.Parse(time.RFC3339Nano, res["created_at"].(string))
	if res["updated_at"] != nil {
		item.UpdatedAt, _ = time.Parse(time.RFC3339Nano, res["updated_at"].(string))
	}
	if res["deleted_at"] != nil {
		return itemBL.ErrDeleted
	}

	item.DeletedAt = time.Now()

	jsonString, _ := json.Marshal(item)

	req := esapi.IndexRequest{
		Index:      "items",
		DocumentID: item.ID.String(),
		Body:       strings.NewReader(string(jsonString)),
	}
	req.Do(context.Background(), es.client)

	return nil
}

func (es *elasticDB) ListItems(ctx context.Context, filter itemBL.ItemFilter) (chan itemBL.ItemBL, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var buffer bytes.Buffer
	var item = &itemBL.ItemBL{}
	chout := make(chan itemBL.ItemBL, 100)

	go func() {
		defer close(chout)

		response, _ := es.client.Search(es.client.Search.WithIndex("items"),
			es.client.Search.WithBody(&buffer))

		var result map[string]interface{}
		json.NewDecoder(response.Body).Decode(&result)

		for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
			res := hit.(map[string]interface{})["_source"].(map[string]interface{})

			item.ID, _ = uuid.Parse(res["id"].(string))
			item.Name = res["name"].(string)
			item.Price = int64(res["price"].(float64))
			item.CreatedAt, _ = time.Parse(time.RFC3339Nano, res["created_at"].(string))
			item.UpdatedAt, _ = time.Parse(time.RFC3339Nano, res["updated_at"].(string))
			item.DeletedAt, _ = time.Parse(time.RFC3339Nano, res["deleted_at"].(string))

			select {
			case <-ctx.Done():
				return
			case chout <- *item:
			}
		}
	}()

	return chout, nil
}

func (es *elasticDB) SearchItems(ctx context.Context, value string) (chan itemBL.ItemBL, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var buffer bytes.Buffer
	var item = &itemBL.ItemBL{}

	chout := make(chan itemBL.ItemBL, 100)

	go func() {
		defer close(chout)
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"name": value,
				},
			},
		}
		json.NewEncoder(&buffer).Encode(query)
		response, _ := es.client.Search(es.client.Search.WithIndex("items"),
			es.client.Search.WithBody(&buffer))

		var result map[string]interface{}
		json.NewDecoder(response.Body).Decode(&result)

		for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
			res :=
				hit.(map[string]interface{})["_source"].(map[string]interface{})

			item.ID, _ = uuid.Parse(res["id"].(string))
			item.Name = res["name"].(string)
			item.Price = int64(res["price"].(float64))
			item.CreatedAt, _ = time.Parse(time.RFC3339Nano, res["created_at"].(string))
			if res["updated_at"] != nil {
				item.UpdatedAt, _ = time.Parse(time.RFC3339Nano, res["updated_at"].(string))
			}
			if res["deleted_at"] != nil {
				item.DeletedAt, _ = time.Parse(time.RFC3339Nano, res["deleted_at"].(string))
			}

			chout <- *item
		}
	}()

	return chout, nil
}
