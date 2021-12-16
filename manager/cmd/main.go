package main

import (
	"database/sql"
	"fmt"
	"lesson05/manager/internal/db/postgresDB"
	"lesson05/manager/internal/db/postgresDB/shardPG"
	_ "lesson05/migrate"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose"
)

func migrateMasterDB() error {
	pgMasterDSN := "postgres://test:test@localhost:8081/test?sslmode=disable"

	pg, err := sql.Open("pgx", pgMasterDSN)
	if err != nil {
		return err
	}
	defer pg.Close()

	err = pg.Ping()
	if err != nil {
		pg.Close()
		return err
	}

	//// Так подключаются SQL миграции
	//err = goose.Up(pg, "migrate")
	//if err != nil {
	//	return err
	//}

	// Так подключаются миграции на языке go
	err = goose.Up(pg, "/var")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	m := shardPG.NewManager(10)

	m.Add(&shardPG.Shard{"port=8100 user=test password=test dbname=test sslmode=disable", 0})
	m.Add(&shardPG.Shard{"port=8110 user=test password=test dbname=test sslmode=disable", 1})
	m.Add(&shardPG.Shard{"port=8120 user=test password=test dbname=test sslmode=disable", 2})

	uu := []*postgresDB.User{
		{1, "Joe Biden", 78, 10},
		{10, "Jill Biden", 69, 1},
		{13, "Donald Trump", 74, 25},
		{25, "Melania Trump", 78, 13},
	}

	for _, u := range uu {
		err := u.Create()
		if err != nil {
			fmt.Println(fmt.Errorf("error on create user %v: %w", u, err))
		}
	}

	//err := migrateMasterDB()
	//if err != nil {
	//	log.Println("func migrateMasterDB error ", err)
	//}
}
