package inmemoryDB

import (
	"context"
	"database/sql"
	"errors"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/app/shortenerBL"

	"github.com/google/uuid"

	"sync"
)

var _ shortenerBL.ShortenerStore = &shortnerMapDB{}

type shortnerMapDB struct {
	sync.Mutex
	sht map[uuid.UUID]shortenerBL.Shortener
}

func NewShortenerMapDB() *shortnerMapDB {
	return &shortnerMapDB{
		sht: make(map[uuid.UUID]shortenerBL.Shortener),
	}
}

func (sdb *shortnerMapDB) CreateShort(ctx context.Context, shortner shortenerBL.Shortener) (*shortenerBL.Shortener, error) {
	sdb.Lock()
	defer sdb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	sdb.sht[shortner.ID] = shortner
	return &shortner, nil
}

func (sdb *shortnerMapDB) GetShort(ctx context.Context, uid uuid.UUID) (*shortenerBL.Shortener, error) {
	sdb.Lock()
	defer sdb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	sht, ok := sdb.sht[uid]
	if ok {
		return &sht, nil
	}
	return nil, sql.ErrNoRows
}

func (sdb *shortnerMapDB) DeleteShort(ctx context.Context, uid uuid.UUID) error {
	sdb.Lock()
	defer sdb.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, ok := sdb.sht[uid]; !ok {
		return errors.New("в БД нет такой позиции")
	}

	delete(sdb.sht, uid)
	return nil
}

func (sdb *shortnerMapDB) SearchShort(ctx context.Context, shortLink string) (*shortenerBL.Shortener, error) {
	sdb.Lock()
	defer sdb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	for _, sht := range sdb.sht {
		if sht.ShortLink == shortLink {
			return &sht, nil
		}
	}
	return nil, errors.New("в БД нет данной записи")
}
