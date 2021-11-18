package inmemoryDB

import (
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/app/followingBL"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

var ErrorNotFoundElement = errors.New("в БД нет такой позиции")

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

func (fldb *followingMapDB) CreateFollowing(ctx context.Context, followingList followingBL.Following) (*followingBL.Following, error) {
	fldb.Lock()
	defer fldb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	fldb.followingDB[followingList.ID] = followingList
	return &followingList, nil
}

func (fldb *followingMapDB) ReadFollowing(ctx context.Context, following followingBL.Following) (*followingBL.Following, error) {
	fldb.Lock()
	defer fldb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	for _, readFollowing := range fldb.followingDB {
		if following.ShortenerID == readFollowing.ShortenerID && following.IPaddress == readFollowing.IPaddress {
			return &readFollowing, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (fldb *followingMapDB) UpdateFollowing(ctx context.Context, following followingBL.Following) (*followingBL.Following, error) {
	fldb.Lock()
	defer fldb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if _, ok := fldb.followingDB[following.ID]; !ok {
		return nil, ErrorNotFoundElement
	}

	fldb.followingDB[following.ID] = following

	return &following, nil
}

func (fldb *followingMapDB) GetFollowingList(ctx context.Context, following followingBL.Following) (chan followingBL.Following, error) {
	chout := make(chan followingBL.Following, 100)

	go func() {
		defer close(chout)
		fldb.Lock()
		defer fldb.Unlock()
		for _, followingItem := range fldb.followingDB {
			if followingItem.ShortenerID == following.ShortenerID {
				select {
				case <-ctx.Done():
					return
				case <-time.After(2 * time.Second):
					return
				case chout <- followingItem:
				}
			}
		}
	}()

	return chout, nil

	//fldb.Lock()
	//defer fldb.Unlock()
	//
	//select {
	//case <-ctx.Done():
	//	return nil, ctx.Err()
	//default:
	//}
	//
	//for _, readFollowing := range fldb.followingListDB {
	//	if following.ShortenerID == readFollowing.ShortenerID {
	//		return &readFollowing, nil
	//	}
	//}
	//
	//return nil, sql.ErrNoRows
}

func (fldb *followingMapDB) DeleteFollowing(ctx context.Context, uid uuid.UUID) (*followingBL.Following, error) {
	fldb.Lock()
	defer fldb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if _, ok := fldb.followingDB[uid]; !ok {
		return nil, errors.New("в БД нет такой позиции")
	}

	deleteFollowingList := fldb.followingDB[uid]
	delete(fldb.followingDB, uid)
	return &deleteFollowingList, nil
}

func (fldb *followingMapDB) SearchShort(ctx context.Context, IPaddres string) (*followingBL.Following, error) {
	fldb.Lock()
	defer fldb.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	for _, elem := range fldb.followingDB {
		if elem.IPaddress == IPaddres {
			return &elem, nil
		}
	}
	return nil, errors.New("в БД нет данной записи")
}
