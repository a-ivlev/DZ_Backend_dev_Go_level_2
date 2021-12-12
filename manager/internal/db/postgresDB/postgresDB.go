package postgresDB

import (
	"database/sql"
)

type User struct {
	UserId int
	Name   string
	Age    int
	Spouse int
}

func (u *User) connection() (*sql.DB, error) {
	s, err := u.ShardById(u.UserId)
	if err != nil {
		return nil, err
	}
	return u.Connection(s.Address)
}

func (u *User) Create() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`INSERT INTO "users" VALUES ($1, $2, $3, $4)`, u.UserId, u.Name, u.Age, u.Spouse)
	return err
}

func (u *User) Read() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	r := c.QueryRow(`SELECT "name", "age", "spouse" FROM "users" WHERE "user_id" = $1`, u.UserId)
	return r.Scan(
		&u.Name,
		&u.Age,
		&u.Spouse,
	)
}

func (u *User) Update() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`UPDATE "users" SET "name" = $2, "age" = $3, "spouse" = $4 WHERE "user_id" = $1`,
		u.UserId, u.Name, u.Age, u.Spouse)
	return err
}

func (u *User) Delete() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`DELETE FROM "users" WHERE "user_id" = $1`, u.UserId)
	return err
}
