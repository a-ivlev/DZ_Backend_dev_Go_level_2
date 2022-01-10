package routerchi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"shop/internal/api/handlers"
	"shop/internal/app/itemBL"
	"shop/internal/app/starter"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

var _ itemBL.ItemStores = &mokeDB{}

var uid1, _ = uuid.Parse("29ec48fa-31fc-4513-a4e9-95220dcf3696")
var createAT1, _ = time.Parse(time.RFC3339, "2022-01-10 12:30:45Z03:00")
var updateAT1, _ = time.Parse(time.RFC3339, "2022-01-10 12:30:45Z03:00")

var uid2, _ = uuid.Parse("25405eb2-d939-4ead-bc50-3c52f6bf047c")
var createAT2, _ = time.Parse(time.RFC3339, "2022-01-10 12:35:01.000000000 +0000 UTC")
var updateAT2, _ = time.Parse(time.RFC3339, "2022-01-10 12:35:01.000000000 +0000 UTC")

type mokeDB struct {
	db map[uuid.UUID]itemBL.ItemBL
}

func NewMokeDB() *mokeDB {
	return &mokeDB{
		db: map[uuid.UUID]itemBL.ItemBL{
			uid1: {
				ID:        uid1,
				Name:      "test-1",
				Price:     10,
				CreatedAt: createAT1,
				UpdatedAt: updateAT1,
			},
			uid2: {
				ID:        uid2,
				Name:      "trst-2",
				Price:     20,
				CreatedAt: createAT2,
				UpdatedAt: updateAT2,
			},
		},
	}
}

func (m *mokeDB) CreateItem(ctx context.Context, item itemBL.ItemBL) (*uuid.UUID, error) {
	if item.ID == uuid.Nil {
		return &uuid.Nil, errors.New("не задан uuid")
	}
	if item.Name == "" {
		return &uuid.Nil, errors.New("не задано поле name")
	}
	if item.Price == 0 {
		return &uuid.Nil, errors.New("не задано поле Price")
	}

	if item.Name == "test-1" {
		return &uid1, nil
	}

	if item.Name == "test-2" {
		return &uid2, nil
	}

	return &uuid.Nil, errors.New("create error")
}

func (m *mokeDB) GetItem(ctx context.Context, ID uuid.UUID) (*itemBL.ItemBL, error) {

	item, ok := m.db[ID]
	if !ok {
		return nil, itemBL.ErrNoRows
	}

	for itm, _ := range m.db {
		fmt.Println("itm ", itm)
	}

	fmt.Println("CreateAT1 ", createAT1)

	fmt.Println("get item ", item)

	return &item, nil
}

func (m *mokeDB) DeleteItem(ctx context.Context, ID uuid.UUID) error {
	_, ok := m.db[ID]
	if !ok {
		return itemBL.ErrNoRows
	}

	return nil
}

func (m *mokeDB) SearchItems(ctx context.Context, s string) (chan itemBL.ItemBL, error) {

	chout := make(chan itemBL.ItemBL, 100)

	go func() {
		defer close(chout)
		for _, item := range m.db {
			if strings.Contains(item.Name, s) {
				select {
				case <-ctx.Done():
					return
				case <-time.After(2 * time.Second):
					return
				case chout <- item:
				}
			}
		}
	}()

	return chout, nil
}

func (m *mokeDB) ListItems(ctx context.Context, filter itemBL.ItemFilter) (chan itemBL.ItemBL, error) {

	chout := make(chan itemBL.ItemBL, 100)

	go func() {
		defer close(chout)
		for _, item := range m.db {
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
				return
			case chout <- item:
			}
		}
	}()

	return chout, nil
}

func TestRouterChi(t *testing.T) {
	//itemDB := inmemItemDB.NewinmemoryDB()
	itemDB := NewMokeDB()
	app := starter.NewApp(itemDB)
	h := handlers.NewHandlers(app)
	rt := NewRouterChi(h)

	//// Тестирование отдельного хендлера работает не корректно, тест проходит даже если указан не правильный метод.
	//w := &httptest.ResponseRecorder{}
	//r := httptest.NewRequest("PUT", "/item", strings.NewReader(`{"name": "test 1", "price": 10}`))
	//// ОБЯЗАТЕЛЬНО нужно задавать заголовок! Без заголовка не работает.
	//r.Header.Set("Content-Type", "application/json")
	//rt.CreateItem(w, r)
	//
	//expected := http.StatusOK
	//if w.Code != expected {
	//	t.Errorf("status wrong, expected %d got %d", expected, w.Code)
	//}

	// Рабочий вариант с моканием сервера.
	hts := httptest.NewServer(rt)

	cli := hts.Client()
	req, _ := http.NewRequest("PUT", hts.URL+"/item", strings.NewReader(`{"name": "test-1", "price": 10}`))
	// Обязательно нужно задавать заголовки.
	req.Header.Set("Content-Type", "application/json")
	resp, _ := cli.Do(req)

	//b, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(b))

	createItem := &Item{}
	_ = json.NewDecoder(resp.Body).Decode(createItem)
	fmt.Println("Тест /item хендлера на запись нового item:\n", *createItem)

	expected := http.StatusOK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status wrong, expected %d got %d", expected, resp.StatusCode)
	}

	req, _ = http.NewRequest("POST", hts.URL+"/item/"+fmt.Sprint(createItem.ID), strings.NewReader(``))
	// Обязательно нужно задавать заголовки.
	req.Header.Set("Content-Type", "application/json")
	resp, _ = cli.Do(req)

	expected = http.StatusOK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status wrong, expected %d got %d", expected, resp.StatusCode)
	}

	getItem := &Item{}
	_ = json.NewDecoder(resp.Body).Decode(getItem)
	fmt.Println("Тест /item/{id} хендлера на получение item:\n", *getItem)

	req, _ = http.NewRequest("PUT", hts.URL+"/item/"+fmt.Sprint(createItem.ID), strings.NewReader(`{"name": "test-update", "price": 25}`))
	// Обязательно нужно задавать заголовки.
	req.Header.Set("Content-Type", "application/json")
	resp, _ = cli.Do(req)

	expected = http.StatusOK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status wrong, expected %d got %d", expected, resp.StatusCode)
	}

	updItem := &Item{}
	_ = json.NewDecoder(resp.Body).Decode(updItem)
	fmt.Println("Тест /item/{id} хендлера на обновление item:\n", *updItem)

	req, _ = http.NewRequest("DELETE", hts.URL+"/item/"+fmt.Sprint(createItem.ID), strings.NewReader(``))
	// Обязательно нужно задавать заголовки.
	//req.Header.Set("Content-Type", "application/json")
	resp, _ = cli.Do(req)

	expected = http.StatusOK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status wrong, expected %d got %d", expected, resp.StatusCode)
	}

	delItem := &Item{}
	_ = json.NewDecoder(resp.Body).Decode(delItem)
	fmt.Println("Тест /item/{id} хендлера на удаление item:\n", *delItem)

}
