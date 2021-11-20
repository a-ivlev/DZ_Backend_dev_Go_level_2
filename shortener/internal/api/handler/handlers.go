package handler

import (
	"DZ_Backend_dev_Go_level_2/shortener/internal/app/redirectBL"
	"DZ_Backend_dev_Go_level_2/shortener/internal/app/repository/followingBL"
	"DZ_Backend_dev_Go_level_2/shortener/internal/app/repository/shortenerBL"
	"context"
	"fmt"
	"time"
)

type Handlers struct {
	redirectBL *redirectBL.Redirect
}

func NewHandlers(redirectBL *redirectBL.Redirect) *Handlers {
	h := &Handlers{
		redirectBL: redirectBL,
	}
	return h
}

type Shortener struct {
	ShortLink  string    `json:"short_link"`
	FullLink   string    `json:"full_link"`
	StatLink   string    `json:"stat_link"`
	TotalCount int       `json:"total_count"`
	CreatedAt  time.Time `json:"created_at"`
}

func (h *Handlers) CreateShortener(ctx context.Context, short Shortener) (Shortener, error) {
	shortenerBL := shortenerBL.Shortener{
		FullLink: short.FullLink,
	}

	newShort, err := h.redirectBL.CreateShortLink(ctx, shortenerBL)
	if err != nil {
		return Shortener{}, fmt.Errorf("error when creating: %w", err)
	}

	return Shortener{
		ShortLink:  newShort.ShortLink,
		FullLink:   newShort.FullLink,
		StatLink:   newShort.StatLink,
		TotalCount: newShort.TotalCount,
		CreatedAt:  newShort.CreatedAt,
	}, nil
}

type Redirect struct {
	ShortLink string `json:"short_link"`
	FullLink string `json:"full_link"`
}

func (h *Handlers) Redirect(ctx context.Context, short Redirect) (Redirect, error) {
	shortenerBL := shortenerBL.Shortener{
		ShortLink: short.ShortLink,
	}

	getFullink, err := h.redirectBL.GetFullLink(ctx, shortenerBL)
	if err != nil {
		return short, fmt.Errorf("error when get URL: %w", err)
	}

	return Redirect{
		ShortLink:  short.ShortLink,
		FullLink:   getFullink.FullLink,
	}, nil
}

type Statistic struct {
	ShortLink  string                  `json:"short_link"`
	TotalCount int                     `json:"total_count"`
	CreatedAt  time.Time               `json:"created_at"`
	FollowList []followingBL.Following `json:"follow_list"`
}

func (h *Handlers) GetStatisticList(ctx context.Context, statisticLink string) (Statistic, error) {
	statistic, err := h.redirectBL.GetStatisticList(ctx, statisticLink)
	if err != nil {
		return Statistic{}, fmt.Errorf("error get statistic: %w", err)
	}

	return Statistic{
		ShortLink:  statistic.ShortLink,
		TotalCount: statistic.TotalCount,
		CreatedAt:  statistic.CreatedAt,
		FollowList: statistic.FollowList,
	}, nil
}
