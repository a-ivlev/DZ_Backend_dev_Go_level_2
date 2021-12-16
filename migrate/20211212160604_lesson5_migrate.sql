-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "users"
(
    "user_id" INT,
    "name"    VARCHAR,
    "age"     INT,
    "spouse"  INT
);

CREATE UNIQUE INDEX "users_user_id" ON "users" ("user_id");

CREATE TABLE "activities"
(
    "user_id" INT,
    "date"    TIMESTAMP,
    "name"    VARCHAR
);

CREATE INDEX "activities_user_id_date" ON "activities" ("user_id", "date");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE "users";

DROP TABLE "activities";
-- +goose StatementEnd
