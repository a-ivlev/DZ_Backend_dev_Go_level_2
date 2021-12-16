package usersBL

import (
	"fmt"
	"time"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/repository/groupUserBL"

	"github.com/google/uuid"
)

type UserBL struct {
	ID         uuid.UUID
	FirstName  string
	MiddleName string
	LastName   string
	Phone      string
	Email      string
	GroupList  map[uuid.UUID]groupUserBL.GroupUserBL
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeleteAt   time.Time
}

type UserStore interface {
	CreateUser(user UserBL) (*UserBL, error)
	ReadUser(id uuid.UUID) (*UserBL, error)
	UpdateUser(user UserBL) (*UserBL, error)
	DeleteUser(id uuid.UUID) error
	SearchUserByName(userName string) (*UserBL, error)
}

type UsersStore struct {
	userStr UserStore
}

func NewUsersStore(userStr UserStore) *UsersStore {
	return &UsersStore{
		userStr: userStr,
	}
}

func (us *UsersStore) CreateUser(user UserBL) (*UserBL, error) {
	user.ID = uuid.New()
	newUser, err := us.userStr.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}

	return newUser, nil
}

func (us *UsersStore) SearchUserByName(userName string) (*UserBL, error) {
	user, err := us.userStr.SearchUserByName(userName)
	if err != nil {
		return nil, fmt.Errorf("search user by name %s error: %w", userName, err)
	}

	return user, nil
}
