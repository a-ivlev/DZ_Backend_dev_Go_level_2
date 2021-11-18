package postgresDB

import (
	"CourseProjectBackendDevGoLevel-1/shortener/internal/app/repository/followingBL"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var _ followingBL.FollowingStore = &PostgresDB{}

type FollowingPG struct {
	ID           uuid.UUID `db:"id"`
	ShortenerID  uuid.UUID `db:"shortener_id"`
	StatLink     string    `db:"stat_link"`
	IPaddress    string    `db:"ip_address"`
	Count        int       `db:"count"`
	FollowLinkAt time.Time `db:"follow_link_at"`
}

func (pg *PostgresDB) CreateFollow(ctx context.Context, follow followingBL.Following) (*followingBL.Following, error) {
	followDB := &FollowingPG{
		ID:           follow.ID,
		ShortenerID:  follow.ShortenerID,
		StatLink:     follow.StatLink,
		IPaddress:    follow.IPaddress,
		Count:        follow.Count,
		FollowLinkAt: follow.FollowLinkAt,
	}

	_, err := pg.db.ExecContext(ctx, `INSERT INTO following
    (id, shortener_id, stat_link, ip_address, count, follow_link_at)
    values ($1, $2, $3, $4, $5, $6);`,
		followDB.ID,
		followDB.ShortenerID,
		followDB.StatLink,
		followDB.IPaddress,
		followDB.Count,
		followDB.FollowLinkAt,
	)
	if err != nil {
		return nil, err
	}

	return &follow, nil
}

func (pg *PostgresDB) ReadFollow(ctx context.Context, uid uuid.UUID) (*followingBL.Following, error) {
	followDB := &FollowingPG{}
	rows, err := pg.db.QueryContext(ctx, `SELECT id, shortener_id, stat_link, ip_address, count, follow_link_at
	FROM following WHERE id=$1;`, uid)
	if err != nil {
		return nil, err
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&followDB.ID,
			&followDB.ShortenerID,
			&followDB.StatLink,
			&followDB.IPaddress,
			&followDB.Count,
			&followDB.FollowLinkAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	return &followingBL.Following{
		ID:           followDB.ID,
		ShortenerID:  followDB.ShortenerID,
		StatLink:     followDB.StatLink,
		IPaddress:    followDB.IPaddress,
		Count:        followDB.Count,
		FollowLinkAt: followDB.FollowLinkAt,
	}, nil
}

func (pg *PostgresDB) UpdateFollow(ctx context.Context, following followingBL.Following) (*followingBL.Following, error) {
	followDB := &FollowingPG{
		ID:           following.ID,
		ShortenerID:  following.ShortenerID,
		StatLink:     following.StatLink,
		IPaddress:    following.IPaddress,
		Count:        following.Count,
		FollowLinkAt: following.FollowLinkAt,
	}

	_, err := pg.db.ExecContext(ctx, `
	UPDATE following SET ip_address=$2, count=$3, follow_link_at=$4 WHERE id=$1;`,
		followDB.ID,
		followDB.IPaddress,
		followDB.Count,
		followDB.FollowLinkAt,
	)
	if err != nil {
		return nil, err
	}

	return &following, nil
}

func (pg *PostgresDB) SearchElement(ctx context.Context, statisticLink string, ipAddress string) (*followingBL.Following, error) {
	followDB := &FollowingPG{}
	rows, err := pg.db.QueryContext(ctx, `SELECT id, shortener_id, stat_link, ip_address, count, follow_link_at
	FROM following WHERE NOT EXISTS (stat_link LIKE $1 AND ip_address=$2);`, "%"+statisticLink, ipAddress)
	if err != nil {
		return nil, err
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&followDB.ID,
			&followDB.ShortenerID,
			&followDB.StatLink,
			&followDB.IPaddress,
			&followDB.Count,
			&followDB.FollowLinkAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	return &followingBL.Following{
		ID:           followDB.ID,
		ShortenerID:  followDB.ShortenerID,
		StatLink:     followDB.StatLink,
		IPaddress:    followDB.IPaddress,
		Count:        followDB.Count,
		FollowLinkAt: followDB.FollowLinkAt,
	}, nil
}

func (pg *PostgresDB) GetFollowingList(ctx context.Context, statisticLink string) ([]followingBL.Following, error) {
	// В слайс followList будут собраны все строки, полученные из базы
	var followList []followingBL.Following
	var followDB = &FollowingPG{}

	rows, err := pg.db.QueryContext(ctx, `
	SELECT id, shortener_id, stat_link, ip_address, count, follow_link_at
	FROM following WHERE stat_link LIKE $1;`, "%"+statisticLink)
	if err != nil {
		return nil, sql.ErrNoRows
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	for rows.Next() {
		// Scan записывает значения столбцов в свойства структуры following
		if err := rows.Scan(
			&followDB.ID,
			&followDB.ShortenerID,
			&followDB.StatLink,
			&followDB.IPaddress,
			&followDB.Count,
			&followDB.FollowLinkAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		followBL := followingBL.Following{
			ID:           followDB.ID,
			ShortenerID:  followDB.ShortenerID,
			StatLink:     followDB.StatLink,
			IPaddress:    followDB.IPaddress,
			Count:        followDB.Count,
			FollowLinkAt: followDB.FollowLinkAt,
		}
		followList = append(followList, followBL)
	}
	// Проверка, что во время выборки данных не происходило ошибок
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read response: %w", rows.Err())
	}

	return followList, nil
}
