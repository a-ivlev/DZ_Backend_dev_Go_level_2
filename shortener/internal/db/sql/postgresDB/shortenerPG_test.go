// +build integration

package postgresDB_test

import "testing"

//func connect() *sql.DB {
//	db, err := sql.Open("pgx", "postgres://postgres:1110@127.0.0.1/shortener-db?sslmode=disable")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = db.Ping()
//	if err != nil {
//		db.Close()
//		log.Fatal(err)
//	}
//	return db
//}

func TestPostgresDB_SearchShortLink(t *testing.T) {
	ctx := context.Background()
	dbpool := connect(ctx)
	defer dbpool.Close()

	tests := []struct {
		name string
		store *postgresDB.PostgresDB
		ctx	context.Context
		client models.Client
		check func(*testing.T, models.ClientID, error)
	}{
		{
			name: "success",
			store: postgresDB.NewPostgresDB(dbpool),
			ctx: context.Background(),
			client: models.Client{
				FirstName: "Егор",
				MiddleName: "Алексевич",
				LastName: "Саблин",
				Phone: "+79225555353",
			},
			check: func(t *testing.T, hint models.ClientID, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, hint)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			hint, err := tt.store.InsertClient(tt.ctx, tt.client)
			tt.check(t, hint, err)
		})
	}
}
