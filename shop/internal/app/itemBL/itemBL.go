package itemBL

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type ItemBL struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type ItemFilter struct {
	PriceLeft  *int64
	PriceRight *int64
	Limit      int
	Offset     int
}

type ItemStores interface {
	CreateItem(ctx context.Context, item ItemBL) (*uuid.UUID, error)
	GetItem(ctx context.Context, itemID uuid.UUID) (*ItemBL, error)
	//UpdateItem(ctx context.Context, item ItemBL) (*ItemBL, error)
	//DeleteItem(ctx context.Context, itemID uuid.UUID) error

	ListItems(ctx context.Context, filter ItemFilter) (chan ItemBL, error)
	SearchItems(ctx context.Context, s string) (chan ItemBL, error)
}

type ItemStore struct {
	store ItemStores
}

func NewItemStore(stores ItemStores) *ItemStore {
	return &ItemStore{
		store: stores,
	}
}

func (is *ItemStore) CreateItem(ctx context.Context, item ItemBL) (*ItemBL, error) {
	item.ID = uuid.New()
	timeNow := time.Now().UTC()
	item.CreatedAt = timeNow
	item.UpdatedAt = timeNow

	id, err := is.store.CreateItem(ctx, item)

	if err != nil {
		return nil, fmt.Errorf("create item error: %w", err)
	}

	item.ID = *id

	return &item, nil
}

func (is *ItemStore) GetItem(ctx context.Context, ID uuid.UUID) (*ItemBL, error) {
	item, err := is.store.GetItem(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get item error: %w", err)
	}

	return item, nil
}

func (is *ItemStore) UpdateItem(ctx context.Context, item ItemBL) (*ItemBL, error) {

	updItem, err := is.store.GetItem(ctx, item.ID)
	if err != nil {
		return nil, fmt.Errorf("get item error: %w", err)
	}

	if item.Name != "" {
		updItem.Name = item.Name
	}

	if item.Price != 0 {
		updItem.Price = item.Price
	}

	updItem.UpdatedAt = time.Now().UTC()

	_, err = is.store.CreateItem(ctx, *updItem)
	if err != nil {
		return nil, fmt.Errorf("update item error: %w", err)
	}

	return updItem, nil
}

func (is *ItemStore) DeleteItem(ctx context.Context, ID uuid.UUID) (*ItemBL, error) {

	delItem, err := is.store.GetItem(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get item error: %w", err)
	}

	delItem.DeletedAt = time.Now().UTC()

	_, err = is.store.CreateItem(ctx, *delItem)
	if err != nil {
		return nil, fmt.Errorf("update item error: %w", err)
	}

	return delItem, nil
}

func (is *ItemStore) ListItems(ctx context.Context, filter ItemFilter) (chan ItemBL, error) {
	chin, err := is.store.ListItems(ctx, filter)
	if err != nil {
		return nil, err
	}

	chout := make(chan ItemBL, 100)
	go func() {
		defer close(chout)
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-chin:
				if !ok {
					return
				}
				//здесь может располагается код бизнес-логики
				//item.Permissions = 0755
				//if item.DeletedAt > item.CreatedAt {
				//	continue
				//}
				d := time.Since(item.DeletedAt)
				log.Printf("Deleted_AT = %v %T", d, d)
				chout <- item
			}
		}
	}()
	return chout, nil
}

func (is *ItemStore) SearchItems(ctx context.Context, s string) (chan ItemBL, error) {
	chin, err := is.store.SearchItems(ctx, s)
	if err != nil {
		return nil, err
	}

	chout := make(chan ItemBL, 100)
	go func() {
		defer close(chout)
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-chin:
				if !ok {
					return
				}
				//здесь может располагается код бизнес-логики
				//item.Permissions = 0755
				chout <- item
			}
		}
	}()

	return chout, nil
}
