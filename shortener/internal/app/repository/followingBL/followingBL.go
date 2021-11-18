package followingBL

import (
	"DZ_Backend_dev_Go_level_2/shortener/internal/app/repository/shortenerBL"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Following struct {
	ID           uuid.UUID `json:"id"`
	ShortenerID  uuid.UUID `json:"shortener_id"`
	StatLink     string    `json:"stat_link"`
	IPaddress    string    `json:"ip_address"`
	Count        int       `json:"count"`
	FollowLinkAt time.Time `json:"follow_link_at"`
}

type FollowingStore interface {
	CreateFollow(ctx context.Context, following Following) (*Following, error)
	ReadFollow(ctx context.Context, uid uuid.UUID) (*Following, error)
	UpdateFollow(ctx context.Context, following Following) (*Following, error)
	SearchElement(ctx context.Context, statisticLink string, ipAddress string) (*Following, error)
	GetFollowingList(ctx context.Context, statisticLink string) ([]Following, error)
}

type FollowingBL struct {
	followingStore FollowingStore
}

func NewFollowingBL(followingStr FollowingStore) *FollowingBL {
	return &FollowingBL{
		followingStore: followingStr,
	}
}

func (fwlBL *FollowingBL) CreateFollowing(ctx context.Context, short *shortenerBL.Shortener) (*Following, error) {
	following := Following{
		ID:           uuid.New(),
		ShortenerID:  short.ID,
		StatLink:     short.StatLink,
		FollowLinkAt: time.Now(),
	}

	newFollowing, err := fwlBL.followingStore.CreateFollow(ctx, following)
	if err != nil {
		return nil, fmt.Errorf("create short-link error: %w", err)
	}

	return newFollowing, nil
}

func (fwBL *FollowingBL) Read(ctx context.Context, id uuid.UUID) (*Following, error) {
	readFollowing, err := fwBL.followingStore.ReadFollow(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("read following error: %w", err)
	}
	return readFollowing, nil
}

func (fwlBL *FollowingBL) Update(ctx context.Context, following Following) (*Following, error) {
	_, err := fwlBL.followingStore.ReadFollow(ctx, following.ID)
	if err != nil {
		return nil, fmt.Errorf("search following error: %w", err)
	}

	updateFollowing, err := fwlBL.followingStore.UpdateFollow(ctx, following)
	if err != nil {
		return nil, fmt.Errorf("error update following: %w", err)
	}

	return updateFollowing, nil
}

func (fwl *FollowingBL) SearchFollowing(ctx context.Context, StatLink string, ipAddress string) (*Following, error) {
	following, err := fwl.followingStore.SearchElement(ctx, StatLink, ipAddress)
	if err != nil {
		return nil, err
	}

	return following, nil
}

func (fwl *FollowingBL) GetFollowingList(ctx context.Context, statisticLink string) ([]Following, error) {

	sliceIn, err := fwl.followingStore.GetFollowingList(ctx, statisticLink)
	if err != nil {
		return nil, err
	}

	return sliceIn, nil
}
