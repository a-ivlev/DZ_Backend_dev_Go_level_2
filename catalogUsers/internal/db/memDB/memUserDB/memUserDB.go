package memUserDB

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/repository/groupUserBL"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/repository/usersBL"
	"github.com/google/uuid"
)

type user struct {
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

type memUserDB struct {
	sync.Mutex
	db map[uuid.UUID]user
}

func NewMemUserDB() *memUserDB {
	return &memUserDB{db: make(map[uuid.UUID]user)}
}

var _ usersBL.UserStore = &memUserDB{}

func (mu *memUserDB) CreateUser(usBL usersBL.UserBL) (*usersBL.UserBL, error) {
	mu.Lock()
	defer mu.Unlock()

	u := &user{
		ID:         usBL.ID,
		FirstName:  usBL.FirstName,
		MiddleName: usBL.MiddleName,
		LastName:   usBL.LastName,
		Phone:      usBL.Phone,
		Email:      usBL.Email,
		GroupList:  usBL.GroupList,
		CreatedAt:  usBL.CreatedAt,
		UpdatedAt:  usBL.UpdatedAt,
		DeleteAt:   usBL.DeleteAt,
	}

	mu.db[u.ID] = *u

	if _, ok := mu.db[u.ID]; !ok {
		return nil, errors.New("error adding user to the database")
	}

	return &usBL, nil
}

func (mu *memUserDB) ReadUser(id uuid.UUID) (*usersBL.UserBL, error) {
	mu.Lock()
	defer mu.Unlock()

	var u user
	var ok bool
	if u, ok = mu.db[id]; !ok {
		return nil, errors.New("error read user to the database")
	}
	return &usersBL.UserBL{
		ID:         u.ID,
		FirstName:  u.FirstName,
		MiddleName: u.MiddleName,
		LastName:   u.LastName,
		Phone:      u.Phone,
		Email:      u.Email,
		GroupList:  u.GroupList,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		DeleteAt:   u.DeleteAt,
	}, nil
}

func (mu *memUserDB) UpdateUser(usBL usersBL.UserBL) (*usersBL.UserBL, error) {
	mu.Lock()
	defer mu.Unlock()

	u := &user{
		ID:         usBL.ID,
		FirstName:  usBL.FirstName,
		MiddleName: usBL.MiddleName,
		LastName:   usBL.LastName,
		Phone:      usBL.Phone,
		Email:      usBL.Email,
		GroupList:  usBL.GroupList,
		CreatedAt:  usBL.CreatedAt,
		UpdatedAt:  usBL.UpdatedAt,
		DeleteAt:   usBL.DeleteAt,
	}

	if _, ok := mu.db[u.ID]; !ok {
		return nil, errors.New("error read user to the database")
	}

	mu.db[u.ID] = *u

	return &usersBL.UserBL{
		ID:         u.ID,
		FirstName:  u.FirstName,
		MiddleName: u.MiddleName,
		LastName:   u.LastName,
		Phone:      u.Phone,
		Email:      u.Email,
		GroupList:  u.GroupList,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		DeleteAt:   u.DeleteAt,
	}, nil
}

func (mu *memUserDB) DeleteUser(id uuid.UUID) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := mu.db[id]; !ok {
		return errors.New("error read user to the database")
	}

	delete(mu.db, id)

	return nil
}

func (mu *memUserDB) SearchUserByName(userName string) (*usersBL.UserBL, error) {
	for _, u := range mu.db {
		if strings.Contains(u.LastName, userName) {
			return &usersBL.UserBL{
				ID:         u.ID,
				FirstName:  u.FirstName,
				MiddleName: u.MiddleName,
				LastName:   u.LastName,
				Phone:      u.Phone,
				Email:      u.Email,
				GroupList:  u.GroupList,
				CreatedAt:  u.CreatedAt,
				UpdatedAt:  u.UpdatedAt,
				DeleteAt:   u.DeleteAt,
			}, nil
		}
	}

	return nil, errors.New("error, not found user to the database")
}
