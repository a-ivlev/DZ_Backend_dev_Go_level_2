package shortenerBL

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Shortener struct {
	ID         uuid.UUID
	ShortLink  string
	FullLink   string
	StatLink   string
	TotalCount int
	CreatedAt  time.Time
}

type ShortenerStore interface {
	CreateShort(ctx context.Context, short Shortener) (*Shortener, error)
	UpdateShort(ctx context.Context, short Shortener) (*Shortener, error)
	SearchShortLink(ctx context.Context, shortLink string) (*Shortener, error)
	SearchStatLink(ctx context.Context, statisticLink string) (*Shortener, error)
}

type ShortenerBL struct {
	shortenerStore ShortenerStore
}

func NewShotenerBL(shortenerStr ShortenerStore) *ShortenerBL {
	return &ShortenerBL{
		shortenerStore: shortenerStr,
	}
}

func GenerateShortLink(ctx context.Context, id uuid.UUID) string {
	link := strings.Split((id).String(), "-")
	shortLink := fmt.Sprintf("%s%s", strings.ToUpper(link[2]), strings.ToUpper(link[3]))

	return shortLink
}

func GenerateStatLink(ctx context.Context, id uuid.UUID) string {
	link := strings.Split((id).String(), "-")
	str := strings.Join(link, "")

	return str
}

func (sh *ShortenerBL) CreateShort(ctx context.Context, shortener Shortener) (*Shortener, error) {
	shortener.ID = uuid.New()

	var err error
	shortener.ShortLink = GenerateShortLink(ctx, shortener.ID)
	shortener.StatLink = GenerateStatLink(ctx, shortener.ID)

	shortener.CreatedAt = time.Now()

	newShortener, err := sh.shortenerStore.CreateShort(ctx, shortener)
	if err != nil {
		return nil, fmt.Errorf("create short-link error: %w", err)
	}

	return newShortener, nil
}

func (sh *ShortenerBL) ReadShort(ctx context.Context, statisticLink string) (*Shortener, error) {
	readFollowing, err := sh.shortenerStore.SearchStatLink(ctx, statisticLink)
	if err != nil {
		return nil, fmt.Errorf("read shortener error: %w", err)
	}

	return readFollowing, nil
}

func (sh *ShortenerBL) Update(ctx context.Context, shortener Shortener) (*Shortener, error) {
	_, err := sh.shortenerStore.SearchStatLink(ctx, shortener.StatLink)
	if err != nil {
		return nil, fmt.Errorf("search user error: %w", err)
	}
	updateHortener, err := sh.shortenerStore.UpdateShort(ctx, shortener)
	if err != nil {
		return nil, fmt.Errorf("create short-link error: %w", err)
	}
	return updateHortener, nil
}

func (sh *ShortenerBL) GetFullLink(ctx context.Context, shortener Shortener) (*Shortener, error) {
	short, err := sh.shortenerStore.SearchShortLink(ctx, shortener.ShortLink)
	if err != nil {
		return nil, fmt.Errorf("search short-link error: %w", err)
	}

	return short, nil
}
