package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upLesson5Migrate, downLesson5Migrate)
}

func upLesson5Migrate(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`CREATE TABLE "users"
(
    "user_id" INT,
    "name"    VARCHAR,
    "age"     INT,
    "spouse"  INT
);`)

	_, err = tx.Exec(`CREATE UNIQUE INDEX "users_user_id" ON "users" ("user_id");`)

	_, err = tx.Exec(`CREATE TABLE "activities"
(
    "user_id" INT,
    "date"    TIMESTAMP,
    "name"    VARCHAR
);`)

	_, err = tx.Exec(`CREATE INDEX "activities_user_id_date" ON "activities" ("user_id", "date");`)

	return err
}

func downLesson5Migrate(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DROP TABLE "users";`)
	_, err = tx.Exec(`DROP TABLE "activities";`)

	return err
}
