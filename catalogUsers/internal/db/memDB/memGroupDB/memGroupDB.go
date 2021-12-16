package memGroupDB

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/repository/groupUserBL"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/repository/usersBL"
	"github.com/google/uuid"
)

type group struct {
	ID               uuid.UUID
	NameGroup        string
	DescriptionGroup string
	UsersList        map[uuid.UUID]usersBL.UserBL
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeleteAt         time.Time
}

type memGroupDB struct {
	sync.Mutex
	db map[uuid.UUID]group
}

func NewMemGroup() *memGroupDB {
	return &memGroupDB{
		db: make(map[uuid.UUID]group),
	}
}

var _ groupUserBL.GroupUserStore = &memGroupDB{}

func (mg *memGroupDB) CreateGroupUser(groupUser groupUserBL.GroupUserBL) (*groupUserBL.GroupUserBL, error) {
	mg.Lock()
	defer mg.Unlock()

	g := &group{
		ID:               groupUser.ID,
		NameGroup:        groupUser.NameGroup,
		DescriptionGroup: groupUser.DescriptionGroup,
		UsersList:        groupUser.UsersList,
		CreatedAt:        groupUser.CreatedAt,
		UpdatedAt:        groupUser.UpdatedAt,
		DeleteAt:         groupUser.DeleteAt,
	}

	mg.db[g.ID] = *g

	if _, ok := mg.db[g.ID]; !ok {
		return nil, errors.New("error adding group to the database")
	}

	return &groupUser, nil
}

func (mg *memGroupDB) ReadGroupUser(id uuid.UUID) (*groupUserBL.GroupUserBL, error) {
	mg.Lock()
	defer mg.Unlock()

	var g group
	var ok bool
	if g, ok = mg.db[id]; !ok {
		return nil, errors.New("error read group to the database")
	}
	return &groupUserBL.GroupUserBL{
		ID:               g.ID,
		NameGroup:        g.NameGroup,
		DescriptionGroup: g.DescriptionGroup,
		UsersList:        g.UsersList,
		CreatedAt:        g.CreatedAt,
		UpdatedAt:        g.UpdatedAt,
		DeleteAt:         g.DeleteAt,
	}, nil
}

func (mg *memGroupDB) UpdateGroupUser(groupUser groupUserBL.GroupUserBL) (*groupUserBL.GroupUserBL, error) {
	mg.Lock()
	defer mg.Unlock()

	g := &group{
		ID:               groupUser.ID,
		NameGroup:        groupUser.NameGroup,
		DescriptionGroup: groupUser.DescriptionGroup,
		UsersList:        groupUser.UsersList,
		CreatedAt:        groupUser.CreatedAt,
		UpdatedAt:        groupUser.UpdatedAt,
		DeleteAt:         groupUser.DeleteAt,
	}

	if _, ok := mg.db[g.ID]; !ok {
		return nil, errors.New("error read group to the database")
	}
	mg.db[g.ID] = *g

	return &groupUserBL.GroupUserBL{
		ID:               g.ID,
		NameGroup:        g.NameGroup,
		DescriptionGroup: g.DescriptionGroup,
		UsersList:        g.UsersList,
		CreatedAt:        g.CreatedAt,
		UpdatedAt:        g.UpdatedAt,
		DeleteAt:         g.DeleteAt,
	}, nil
}

func (mg *memGroupDB) DeleteGroupUser(id uuid.UUID) error {
	mg.Lock()
	defer mg.Unlock()

	if _, ok := mg.db[id]; !ok {
		return errors.New("error read group to the database")
	}

	delete(mg.db, id)

	return nil
}

func (mg *memGroupDB) SearchGroupByName(nameGroup string) (*groupUserBL.GroupUserBL, error) {
	for _, g := range mg.db {
		if strings.Contains(g.NameGroup, nameGroup) {
			return &groupUserBL.GroupUserBL{
				ID:               g.ID,
				NameGroup:        g.NameGroup,
				DescriptionGroup: g.DescriptionGroup,
				UsersList:        g.UsersList,
				CreatedAt:        g.CreatedAt,
				UpdatedAt:        g.UpdatedAt,
				DeleteAt:         g.DeleteAt,
			}, nil
		}
	}

	return nil, errors.New("error, not found group to the database")
}
