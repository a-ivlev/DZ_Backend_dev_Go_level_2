package shortenerBL

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/app/followingBL"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/db/inmemoryDB"

	"github.com/google/uuid"
)

type Shortener struct {
	ID            uuid.UUID
	ShortLink     string
	FullLink      string
	StatisticLink string
	Count         int
	CreatedAt     time.Time
}

type ShortenerStore interface {
	CreateShort(ctx context.Context, short Shortener) (*Shortener, error)
	GetShort(ctx context.Context, uid uuid.UUID) (*Shortener, error)
	DeleteShort(ctx context.Context, uid uuid.UUID) error
	SearchShort(ctx context.Context, shortLink string) (*Shortener, error)
}

type ShortenerBL struct {
	shortenerStore ShortenerStore
	followingBL    followingBL.FollowingBL
}

func NewShotenerBL(shortenerStr ShortenerStore) *ShortenerBL {
	followingStr := inmemoryDB.NewFollowingMapDB()
	return &ShortenerBL{
		shortenerStore: shortenerStr,
		followingBL:    *followingBL.NewFollowingBL(followingStr),
	}
}

func generateShortLink(id uuid.UUID) (string, error) {
	srvHost := "http://localhost" // ctx.Value("SRV_HOST")
	link := strings.Split((id).String(), "-")
	shortLink := fmt.Sprintf("%s/%s", srvHost, strings.ToUpper(link[0]))

	return shortLink, nil
}

func generateStatisticLink(id uuid.UUID) (string, error) {
	srvHost := "http://localhost" // ctx.Value("SRV_HOST")

	link := strings.Split((id).String(), "-")
	str := strings.Join(link, "")
	statisticLink := fmt.Sprintf("%s/%s", srvHost, str)

	return statisticLink, nil
}

func (sh *ShortenerBL) CreateShort(ctx context.Context, shortener Shortener) (*Shortener, error) {
	shortener.ID = uuid.New()
	var err error
	shortener.ShortLink, err = generateShortLink(shortener.ID)
	if err != nil {
		log.Printf("error func GenerateShortLink: %v", err)
	}

	shortener.StatisticLink, err = generateStatisticLink(shortener.ID)
	if err != nil {
		log.Printf("error func GenerateShortLink: %v", err)
	}

	createFollowing, err := sh.followingBL.CreateFollowing(ctx, shortener.ID)
	if err != nil {
		return nil, fmt.Errorf("create folloeing error: %w", err)
	}
	//TODO
	fmt.Println("CreateShort -> createFollowing", createFollowing)

	//shortner.ClientID = ctx.Value("ClientID").(uuid.UUID)
	shortener.CreatedAt = time.Now()

	newShortener, err := sh.shortenerStore.CreateShort(ctx, shortener)
	if err != nil {
		return nil, fmt.Errorf("create short-link error: %w", err)
	}

	return newShortener, nil
}

func (sh *ShortenerBL) GetFullLink(ctx context.Context, shortener Shortener) (*Shortener, error) {
	short, err := sh.shortenerStore.SearchShort(ctx, shortener.ShortLink)
	if err != nil {
		return nil, err
	}

	following := &followingBL.Following{
		IPaddress:   ctx.Value("IP_addres").(string),
		ShortenerID: shortener.ID,
	}
	//readFollowing, err := sh.followingStore.ReadFollowing(ctx, *following)
	//if err != nil {
	//	log.Println("error func ReadFollowingList: ", err)
	//}

	_, err = sh.followingBL.UpdateFollowing(ctx, *following)
	if err != nil {
		log.Println("error func UpdateFollowingList: ", err)
	}

	return short, nil
}

//func (sh *ShortherBL) Redirect(ctx context.Context, shortLink string) (*Shortner, error) {
//	short, err := sh.shortner.SearchShort(ctx, shortLink)
//	if err != nil {
//		return nil, err
//	}
//
//	return short, nil
//}
