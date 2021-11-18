package inmemoryDB

import (
	"DZ_Backend_dev_Go_level_2/shortener/internal/app/repository/followingBL"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"strings"
	"sync"
)

var ErrorNotFound = errors.New("error not found")

var _ followingBL.FollowingStore = &followingMapDB{}

type followingMapDB struct {
	sync.Mutex
	followingDB map[uuid.UUID]followingBL.Following
}

func NewFollowingMapDB() *followingMapDB {
	return &followingMapDB{
		followingDB: make(map[uuid.UUID]followingBL.Following),
	}
}

func (fldb *followingMapDB) CreateFollow(ctx context.Context, following followingBL.Following) (*followingBL.Following, error) {
	fldb.Lock()
	defer fldb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	fldb.followingDB[following.ID] = following
	return &following, nil
}

func (fldb *followingMapDB) ReadFollow(ctx context.Context, uid uuid.UUID) (*followingBL.Following, error) {
	fldb.Lock()
	defer fldb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	u, ok := fldb.followingDB[uid]
	if ok {
		return &u, nil
	}
	return nil, sql.ErrNoRows
}

func (fldb *followingMapDB) UpdateFollow(ctx context.Context, following followingBL.Following) (*followingBL.Following, error) {
	fldb.Lock()
	defer fldb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if _, ok := fldb.followingDB[following.ID]; !ok {
		return nil, ErrorNotFound
	}

	fldb.followingDB[following.ID] = following

	return &following, nil
}

//func (fldb *followingMapDB) DeleteFollow(ctx context.Context, uid uuid.UUID) (*followingBL.Following, error) {
//	fldb.Lock()
//	defer fldb.Unlock()
//
//	select {
//	case <-ctx.Done():
//		return nil, ctx.Err()
//	default:
//	}
//
//	if _, ok := fldb.followingDB[uid]; !ok {
//		return nil, errors.New("в БД нет такой позиции")
//	}
//
//	deleteFollowingList := fldb.followingDB[uid]
//	delete(fldb.followingDB, uid)
//	return &deleteFollowingList, nil
//}


func (fldb *followingMapDB) SearchElement(ctx context.Context, statisticLink string, ipAddress string) (*followingBL.Following, error) {
	fldb.Lock()
	defer fldb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	for _, elem := range fldb.followingDB {
		//if elem.StatisticLink == statisticLink && elem.IPaddress == ipAddress {
		if strings.Contains(elem.StatLink, statisticLink) && elem.IPaddress == ipAddress {
			return &elem, nil
		}
	}

	return nil, ErrorNotFound
}

func (fldb *followingMapDB) GetFollowingList(ctx context.Context, statisticLink string) ([]followingBL.Following, error) {
	sliceOut := make([]followingBL.Following, 0, 100)

		fldb.Lock()
		defer fldb.Unlock()
		for _, followingItem := range fldb.followingDB {
			// if followingItem.StatisticLink == statisticLink {
			if strings.Contains(followingItem.StatLink, statisticLink) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				//case <-time.After(2 * time.Second):
				//	return
				//case chout <- followingItem:
				default:
				}
				sliceOut = append(sliceOut, followingItem)
			}
		}

	return sliceOut, nil
}

//func (fldb *followingMapDB) ReadFollowing(ctx context.Context, following followingBL.Following) (*followingBL.Following, error) {
//	fldb.Lock()
//	defer fldb.Unlock()
//
//	select {
//	case <-ctx.Done():
//		return nil, ctx.Err()
//	default:
//	}
//
//	for _, readFollowing := range fldb.followingDB {
//		if following.ShortenerID == readFollowing.ShortenerID && following.IPaddress == readFollowing.IPaddress {
//			return &readFollowing, nil
//		}
//	}
//
//	return nil, sql.ErrNoRows
//}