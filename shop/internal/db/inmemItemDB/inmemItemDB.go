package inmemItemDB

import (
	"context"
	"shop/internal/app/itemBL"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var _ itemBL.ItemStores = &inmemItemDB{}

type inmemItemDB struct {
	mu    sync.Mutex
	items map[uuid.UUID]itemBL.ItemBL
}

func NewinmemoryDB() *inmemItemDB {
	return &inmemItemDB{
		items: make(map[uuid.UUID]itemBL.ItemBL),
	}
}

func (i *inmemItemDB) CreateItem(ctx context.Context, item itemBL.ItemBL) (*uuid.UUID, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	i.items[item.ID] = item

	return &item.ID, nil
}

func (i *inmemItemDB) GetItem(ctx context.Context, ID uuid.UUID) (*itemBL.ItemBL, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	item, ok := i.items[ID]
	if !ok {
		return nil, itemBL.ErrNoRows
	}

	return &item, nil
}

//func (i *inmemItemDB) UpdateItem(ctx context.Context, item itemBL.ItemBL) (*itemBL.ItemBL, error) {
//	i.mu.Lock()
//	defer i.mu.Unlock()
//
//	select {
//	case <-ctx.Done():
//		return nil, ctx.Err()
//	default:
//	}
//
//	_, ok := i.items[item.ID]
//	if !ok {
//		return nil, itemBL.ErrNoRows
//	}
//
//	i.items[item.ID] = item
//
//	return &item, nil
//}

func (i *inmemItemDB) DeleteItem(ctx context.Context, ID uuid.UUID) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, ok := i.items[ID]
	if !ok {
		return itemBL.ErrNoRows
	}

	delete(i.items, ID)
	return nil
}

func (i *inmemItemDB) SearchItems(ctx context.Context, s string) (chan itemBL.ItemBL, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	chout := make(chan itemBL.ItemBL, 100)

	go func() {
		defer close(chout)
		i.mu.Lock()
		defer i.mu.Unlock()
		for _, item := range i.items {
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

func (i *inmemItemDB) ListItems(ctx context.Context, filter itemBL.ItemFilter) (chan itemBL.ItemBL, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	chout := make(chan itemBL.ItemBL, 100)

	go func() {
		defer close(chout)
		i.mu.Lock()
		defer i.mu.Unlock()
		for _, item := range i.items {
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

	//var res []*itemBL.ItemBL
	//
	//log.Println("filter = ", filter)
	//
	//for _, item := range i.items {
	//	res = append(res, &itemBL.ItemBL{
	//		ID:        item.ID,
	//		Name:      item.Name,
	//		Price:     item.Price,
	//		CreatedAt: item.CreatedAt,
	//		UpdatedAt: item.UpdatedAt,
	//	})
	//}
	//return res, nil

	//var res []*ItemDB
	//
	//itemSlice := make([]*ItemDB, 0, len(i.items))
	//for _, item := range i.items {
	//	itemSlice = append(itemSlice, item)
	//}
	//sort.Slice(itemSlice, func(i, j int) bool {
	//	return itemSlice[i].Price < itemSlice[j].Price
	//})
	//
	//for _, item := range itemSlice {
	//	if filter.PriceLeft == nil && filter.PriceRight == nil {
	//		res = itemSlice
	//		break
	//		//res = append(res, &itemBL.ItemBL{
	//		//	ID:        item.ID,
	//		//	Name:      item.Name,
	//		//	Price:     item.Price,
	//		//	CreatedAt: item.CreatedAt,
	//		//	UpdatedAt: item.UpdatedAt,
	//		//})
	//	}
	//	if filter.PriceLeft != nil && filter.PriceRight == nil && item.Price >= *filter.PriceLeft ||
	//		filter.PriceLeft == nil && filter.PriceRight != nil && item.Price <= *filter.PriceRight {
	//		res = append(res, item)
	//		//res = append(res, &itemBL.ItemBL{
	//		//	ID:        item.ID,
	//		//	Name:      item.Name,
	//		//	Price:     item.Price,
	//		//	CreatedAt: item.CreatedAt,
	//		//	UpdatedAt: item.UpdatedAt,
	//		//})
	//	}
	//	if filter.PriceLeft != nil && filter.PriceRight != nil &&
	//		item.Price >= *filter.PriceLeft && item.Price <= *filter.PriceRight {
	//		res = append(res, item)
	//		//res = append(res, &itemBL.ItemBL{
	//		//	ID:        item.ID,
	//		//	Name:      item.Name,
	//		//	Price:     item.Price,
	//		//	CreatedAt: item.CreatedAt,
	//		//	UpdatedAt: item.UpdatedAt,
	//		//})
	//	}
	//}
	//
	//resFiltered := make([]*itemBL.ItemBL, 0, len(res))
	//for idx, item := range res {
	//	if len(resFiltered) == filter.Limit {
	//		break
	//	}
	//	if idx < filter.Offset {
	//		continue
	//	}
	//	resFiltered = append(resFiltered, &itemBL.ItemBL{
	//		ID:        item.ID,
	//		Name:      item.Name,
	//		Price:     item.Price,
	//		CreatedAt: item.CreatedAt,
	//		UpdatedAt: item.UpdatedAt,
	//	})
	//	//resFiltered = append(resFiltered, item)
	//}
	//
	//return resFiltered, nil
	////return res[filter.Offset : filter.Offset+filter.Limit - 1], nil
}
