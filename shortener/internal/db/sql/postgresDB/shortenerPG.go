package postgresDB

import (
	"CourseProjectBackendDevGoLevel-1/shortener/internal/app/repository/shortenerBL"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

var _ shortenerBL.ShortenerStore = &PostgresDB{}

type ShortenerPG struct {
	ID         uuid.UUID `db:"id"`
	ShortLink  string    `db:"short_link"`
	FullLink   string    `db:"full_link"`
	StatLink   string    `db:"stat_link"`
	TotalCount int       `db:"total_count"`
	CreatedAt  time.Time `db:"created_at"`
}

func (pg *PostgresDB) CreateShort(ctx context.Context, short shortenerBL.Shortener) (*shortenerBL.Shortener, error) {
	shortDB := &ShortenerPG{
		ID:         short.ID,
		ShortLink:  short.ShortLink,
		FullLink:   short.FullLink,
		StatLink:   short.StatLink,
		TotalCount: short.TotalCount,
		CreatedAt:  short.CreatedAt,
	}

	result, err := pg.db.ExecContext(ctx, `INSERT INTO shortener
    (id, short_link, full_link, stat_link, total_count, created_at)
    values ($1, $2, $3, $4, $5, $6);`,
		shortDB.ID,
		shortDB.ShortLink,
		shortDB.FullLink,
		shortDB.StatLink,
		shortDB.TotalCount,
		shortDB.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", rows)
	}

	return &short, nil
}

func (pg *PostgresDB) UpdateShort(ctx context.Context, short shortenerBL.Shortener) (*shortenerBL.Shortener, error) {
	shortDB := &ShortenerPG{
		ID:         short.ID,
		ShortLink:  short.ShortLink,
		FullLink:   short.FullLink,
		StatLink:   short.StatLink,
		TotalCount: short.TotalCount,
		CreatedAt:  short.CreatedAt,
	}

	_, err := pg.db.ExecContext(ctx, `UPDATE shortener SET total_count=$2 WHERE id=$1;`,
		shortDB.ID,
		shortDB.TotalCount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update shortener total_count and created_at db: %w", err)
	}

	return &short, nil
}

func (pg *PostgresDB) SearchShortLink(ctx context.Context, shortLink string) (*shortenerBL.Shortener, error) {
	shortDB := &ShortenerPG{}

	const sql = `select id, short_link, full_link, stat_link, total_count, created_at
	from shortener where short_link like $1;`
	rows, err := pg.db.QueryContext(ctx, sql, "%"+shortLink)
	if err != nil {
		return nil, fmt.Errorf("failed to select request db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&shortDB.ID,
			&shortDB.ShortLink,
			&shortDB.FullLink,
			&shortDB.StatLink,
			&shortDB.TotalCount,
			&shortDB.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	return &shortenerBL.Shortener{
		ID:         shortDB.ID,
		ShortLink:  shortDB.ShortLink,
		FullLink:   shortDB.FullLink,
		StatLink:   shortDB.StatLink,
		TotalCount: shortDB.TotalCount,
		CreatedAt:  shortDB.CreatedAt,
	}, nil
}

func (pg *PostgresDB) SearchStatLink(ctx context.Context, statisticLink string) (*shortenerBL.Shortener, error) {
	shortDB := &ShortenerPG{}
	rows, err := pg.db.QueryContext(ctx, `SELECT id, short_link, full_link, stat_link, total_count, created_at
	FROM shortener WHERE stat_link LIKE $1;`, "%"+statisticLink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&shortDB.ID,
			&shortDB.ShortLink,
			&shortDB.FullLink,
			&shortDB.StatLink,
			&shortDB.TotalCount,
			&shortDB.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	return &shortenerBL.Shortener{
		ID:         shortDB.ID,
		ShortLink:  shortDB.ShortLink,
		FullLink:   shortDB.FullLink,
		StatLink:   shortDB.StatLink,
		TotalCount: shortDB.TotalCount,
		CreatedAt:  shortDB.CreatedAt,
	}, nil
}
