package inmemoryDB

import (
	"CourseProjectBackendDevGoLevel-1/shortener/internal/app/repository/shortenerBL"
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var ErrorNotFoun = errors.New("error not found")

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

func (sdb *shortnerMapDB) UpdateShort(ctx context.Context, shortner shortenerBL.Shortener) (*shortenerBL.Shortener, error) {
	sdb.Lock()
	defer sdb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if _, ok := sdb.sht[shortner.ID]; !ok {
		return nil, ErrorNotFoun
	}

	sdb.sht[shortner.ID] = shortner
	return &shortner, nil
}

func (sdb *shortnerMapDB) SearchShortLink(ctx context.Context, shortLink string) (*shortenerBL.Shortener, error) {
	sdb.Lock()
	defer sdb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	for _, elem := range sdb.sht {
		if strings.Contains(elem.ShortLink, shortLink) {
			return &elem, nil
		}
	}

	return nil, ErrorNotFoun
}

func (sdb *shortnerMapDB) SearchStatLink(ctx context.Context, statisticLink string) (*shortenerBL.Shortener, error) {
	sdb.Lock()
	defer sdb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	for _, elem := range sdb.sht {
		// if elem.StatisticLink == string(statisticLink) {
		if strings.Contains(elem.StatLink, statisticLink) {
			return &elem, nil
		}
	}

	return nil, sql.ErrNoRows
}
