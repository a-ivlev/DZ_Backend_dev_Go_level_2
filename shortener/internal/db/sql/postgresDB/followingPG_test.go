// +build integration

package postgresDB_test

import (
	"CourseProjectBackendDevGoLevel-1/shortener/internal/app/repository/followingBL"
	"CourseProjectBackendDevGoLevel-1/shortener/internal/db/sql/postgresDB"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPostgresDB_CreateFollow(t *testing.T) {

	tests := []struct {
		name string
		store *postgresDB.PostgresDB
		ctx	context.Context
		follow *followingBL.Following
		check func(*testing.T, followingBL.Following, error)
	}{
		{
			name: "success",
			store: postgresDB.NewPostgresDB(),
			ctx: context.Background(),
			follow: followingBL.Following{
				ID: uuid.Parse("9ba87644-294d-4f3e-b5df-51c9d8c330b3"),
				ShortenerID: uuid.Parse("9ba87644-294d-4f3e-b5df-51c9d8c330b3"),
				StatLink: "f408fe9bad0a4699a83ef219cfb0f6a0",
				IPaddress: "127.0.0.1",
				Count: 0,
				FollowLinkAt: "2021-11-13 15:43:11.621 +0300",
			},
			check: func(t *testing.T, follow followingBL.Following, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, hint)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			hint, err := tt.store.CreateFollow(tt.ctx, tt.follow)
			tt.check(t, hint, err)
		})
	}
}

