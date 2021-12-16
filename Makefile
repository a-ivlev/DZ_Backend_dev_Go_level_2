MIGRATION_DIR := "migrate"

# Команда для запуска name=lesson5_migrate make new-migrate
new-migrate:
	goose -dir $(MIGRATION_DIR) create $(name) sql

# Команда для запуска PG_DSN=postgres://test:test@localhost:8081/test?sslmode=disable make migrate-up
migrate-up:
	goose -dir $(MIGRATION_DIR) postgres "$(PG_DSN)" up

# Команда для запуска PG_DSN=postgres://test:test@localhost:8081/test?sslmode=disable make migrate-down
migrate-down:
	goose -dir $(MIGRATION_DIR) postgres "$(PG_DSN)" down



# Команда для запуска name=lesson5_migrate make new-go-migrate
new-go-migrate:
	goose -dir $(MIGRATION_DIR) create $(name) go
