package followingBL

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type Following struct {
	ID              uuid.UUID
	ShortenerID     uuid.UUID
	IPaddress       string
	Count           int
	FollowingLinkAt time.Time
}

type FollowingStore interface {
	CreateFollowing(ctx context.Context, following Following) (*Following, error)
	ReadFollowing(ctx context.Context, following Following) (*Following, error)
	UpdateFollowing(ctx context.Context, followingList Following) (*Following, error)
	DeleteFollowing(ctx context.Context, uid uuid.UUID) (*Following, error)
	GetFollowingList(ctx context.Context, followingList Following) (chan Following, error)
}

type FollowingBL struct {
	followingStore FollowingStore
}

func NewFollowingBL(followingStr FollowingStore) *FollowingBL {
	return &FollowingBL{
		followingStore: followingStr,
	}
}

func (fwlBL *FollowingBL) CreateFollowing(ctx context.Context, ShortenerID uuid.UUID) (*Following, error) {
	following := Following{
		ID:          uuid.New(),
		ShortenerID: ShortenerID,
	}

	newFollowing, err := fwlBL.followingStore.CreateFollowing(ctx, following)
	if err != nil {
		return nil, fmt.Errorf("create short-link error: %w", err)
	}

	return newFollowing, nil
}

func (fwBL *FollowingBL) ReadFollowing(ctx context.Context, following Following) (*Following, error) {
	readFollowing, err := fwBL.followingStore.ReadFollowing(ctx, following)
	if err != nil {
		return nil, fmt.Errorf("read following error: %w", err)
	}
	return readFollowing, nil
}

func (fwlBL *FollowingBL) UpdateFollowing(ctx context.Context, following Following) (*Following, error) {
	//IPaddress := ctx.Value("IP_address").(string)
	readFolowing, err := fwlBL.followingStore.ReadFollowing(ctx, following)
	if err != nil {
		return nil, err
	}
	//TODO
	fmt.Println("followingBL -> UpdateFollowingList -> ReadFollowing", readFolowing)

	if readFolowing == nil {
		following = Following{
			ID:              uuid.New(),
			ShortenerID:     following.ShortenerID,
			Count:           1,
			IPaddress:       following.IPaddress,
			FollowingLinkAt: time.Now(),
		}

		var createFollowing *Following
		createFollowing, err = fwlBL.followingStore.CreateFollowing(ctx, following)
		if err != nil {
			log.Println("error new create following")
		}
		return createFollowing, nil
	}

	following = Following{
		Count:           following.Count + 1,
		FollowingLinkAt: time.Now(),
	}

	updateFollowing, err := fwlBL.followingStore.UpdateFollowing(ctx, following)
	if err != nil {
		log.Println("error update following")
	}

	return updateFollowing, nil
}

func (fwl *FollowingBL) GetFollowingList(ctx context.Context, following Following) (chan Following, error) {
	// FIXME: здесь нужно использвоать паттерн Unit of Work
	// бизнес-транзакция
	chin, err := fwl.followingStore.GetFollowingList(ctx, following)
	if err != nil {
		return nil, err
	}
	chout := make(chan Following, 100)
	go func() {
		defer close(chout)
		for {
			select {
			case <-ctx.Done():
				return
			case followingList, ok := <-chin:
				if !ok {
					return
				}
				chout <- followingList
			}
		}
	}()
	return chout, nil
}

//func (fwl *FollowingListBL) Create(ctx context.Context, followingList FollowingList) (*FollowingList, error) {
//	followingList.ID = uuid.New()
//	id, err := fwl.followingListStore.CreateFollowingList(ctx, followingList)
//	if err != nil {
//		return nil, fmt.Errorf("create user error: %w", err)
//	}
//	return &u, nil
//}
//

//
//func (us *Users) Delete(ctx context.Context, uid uuid.UUID) (*User, error) {
//	u, err := us.ustore.Read(ctx, uid)
//	if err != nil {
//		return nil, fmt.Errorf("search user error: %w", err)
//	}
//	return u, us.ustore.Delete(ctx, uid)
//}
