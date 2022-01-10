package handlers

import (
	"context"
	"errors"
	"fmt"
	"shop/internal/app/itemBL"
	"shop/internal/app/starter"
	"time"

	"github.com/google/uuid"
)

type Handlers struct {
	app *starter.App
}

func NewHandlers(app *starter.App) *Handlers {
	h := &Handlers{
		app: app,
	}

	return h
}

type ItemHendler struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DelatedAT time.Time `json:"delated_at"`
}

func (h *Handlers) CreateItemHandler(ctx context.Context, item ItemHendler) (ItemHendler, error) {

	itemBL := itemBL.ItemBL{
		Name:  item.Name,
		Price: item.Price,
	}

	newItem, err := h.app.CreateItem(ctx, itemBL)
	if err != nil {
		return ItemHendler{}, fmt.Errorf("error when creating: %w", err)
	}

	return ItemHendler{
		ID:        newItem.ID,
		Name:      newItem.Name,
		Price:     newItem.Price,
		CreatedAt: newItem.CreatedAt,
		UpdatedAt: newItem.UpdatedAt,
	}, nil
}

func (h *Handlers) GetItemHandler(ctx context.Context, uid uuid.UUID) (*ItemHendler, error) {
	if (uid == uuid.UUID{}) {
		return &ItemHendler{}, fmt.Errorf("bad request: uid is empty")
	}

	getItem, err := h.app.GetItem(ctx, uid)
	if err != nil {
		if errors.Is(err, itemBL.ErrNoRows) {
			return &ItemHendler{}, itemBL.ErrItemNotFound
		}
		return &ItemHendler{}, fmt.Errorf("error when reading: %w", err)
	}

	return &ItemHendler{
		ID:        getItem.ID,
		Name:      getItem.Name,
		Price:     getItem.Price,
		CreatedAt: getItem.CreatedAt,
		UpdatedAt: getItem.UpdatedAt,
	}, nil
}

func (h *Handlers) UpdateItemHandler(ctx context.Context, item ItemHendler) (*ItemHendler, error) {

	updItemBL := itemBL.ItemBL{
		ID:    item.ID,
		Name:  item.Name,
		Price: item.Price,
	}

	updItem, err := h.app.UpdateItem(ctx, updItemBL)
	if err != nil {
		if errors.Is(err, itemBL.ErrNoRows) {
			return &ItemHendler{}, itemBL.ErrItemNotFound
		}
		return &ItemHendler{}, fmt.Errorf("error when updating: %w", err)
	}

	return &ItemHendler{
		ID:        updItem.ID,
		Name:      updItem.Name,
		Price:     updItem.Price,
		CreatedAt: updItem.CreatedAt,
		UpdatedAt: updItem.UpdatedAt,
	}, nil
}

func (h *Handlers) DeleteItemHandler(ctx context.Context, uid uuid.UUID) (*ItemHendler, error) {
	if (uid == uuid.UUID{}) {
		return nil, fmt.Errorf("bad request: uid is empty")
	}

	delItem, err := h.app.DeleteItem(ctx, uid)
	if err != nil {
		if errors.Is(err, itemBL.ErrNoRows) {
			return &ItemHendler{}, itemBL.ErrItemNotFound
		}
		return &ItemHendler{}, fmt.Errorf("error when deleting: %w", err)
	}

	return &ItemHendler{
		ID:        delItem.ID,
		Name:      delItem.Name,
		Price:     delItem.Price,
		CreatedAt: delItem.CreatedAt,
		UpdatedAt: delItem.UpdatedAt,
		DelatedAT: delItem.DeletedAt,
	}, nil
}

func (h *Handlers) ListItemHandler(ctx context.Context, filter itemBL.ItemFilter, f func(ItemHendler) error) error {

	ch, err := h.app.ListItems(ctx, filter)
	if err != nil {
		return fmt.Errorf("error when reading: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-ch:
			if !ok {
				return nil
			}
			if err := f(ItemHendler{
				ID:        item.ID,
				Name:      item.Name,
				Price:     item.Price,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}); err != nil {
				return err
			}
		}
	}
}

func (h *Handlers) SearchItemsHandler(ctx context.Context, s string, f func(ItemHendler) error) error {

	ch, err := h.app.SearchItems(ctx, s)
	if err != nil {
		return fmt.Errorf("error when reading: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-ch:
			if !ok {
				return nil
			}
			if err := f(ItemHendler{
				ID:        item.ID,
				Name:      item.Name,
				Price:     item.Price,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}); err != nil {
				return err
			}
		}
	}
}
