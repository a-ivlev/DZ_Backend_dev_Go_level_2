package postgresDB

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib" // PostgresSQL driver
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB() *PostgresDB {
	dsn := os.Getenv("POSTGRES_DSN")

	if dsn == "" {
		log.Fatal("POSTGRES_DSN is not set in the environment variable")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

const	(DDL = `
CREATE TABLE IF NOT EXISTS public.shortener (
	id uuid NOT NULL,
	short_link varchar(9) NOT NULL,
	full_link varchar NOT NULL,
	stat_link varchar(33) NOT NULL,
	total_count int4 NULL,
	created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS public."following" (
	id uuid NOT NULL,
	shortener_id uuid NOT NULL,
	stat_link varchar(33) NOT NULL,
	ip_address varchar(16) NULL,
	count int4 NULL,
	follow_link_at timestamptz NULL
);`)
	_, err = db.Exec(DDL)
	if err != nil {
		log.Fatal("failed to create shortener and following to Postgres db: ", err)
	}
	pg := &PostgresDB{
		db: db,
	}

	return pg
}

func (pg *PostgresDB) Close() {
	pg.db.Close()
}
