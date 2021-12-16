package groupUserBL

import (
	"fmt"
	"time"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/repository/usersBL"

	"github.com/google/uuid"
)

type GroupUserBL struct {
	ID               uuid.UUID
	NameGroup        string
	DescriptionGroup string
	UsersList        map[uuid.UUID]usersBL.UserBL
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeleteAt         time.Time
}

type GroupUserStore interface {
	CreateGroupUser(groupUser GroupUserBL) (*GroupUserBL, error)
	ReadGroupUser(id uuid.UUID) (*GroupUserBL, error)
	UpdateGroupUser(groupUser GroupUserBL) (*GroupUserBL, error)
	DeleteGroupUser(id uuid.UUID) error
	SearchGroupByName(nameGroup string) (*GroupUserBL, error)
}

type GroupUsersStore struct {
	groupUserStore GroupUserStore
}

func NewGroupUsersStore(groupUsersStr GroupUserStore) *GroupUsersStore {
	return &GroupUsersStore{
		groupUserStore: groupUsersStr,
	}
}

func (gs *GroupUsersStore) CreateGroupUser(groupUser GroupUserBL) (*GroupUserBL, error) {
	groupUser.ID = uuid.New()
	newGroupUser, err := gs.groupUserStore.CreateGroupUser(groupUser)
	if err != nil {
		return nil, fmt.Errorf("create group user error: %w", err)
	}

	return newGroupUser, nil
}

func (gs *GroupUsersStore) ReadGroupUser(id uuid.UUID) (*GroupUserBL, error) {
	groupUser, err := gs.groupUserStore.ReadGroupUser(id)
	if err != nil {
		return nil, fmt.Errorf("read group user id %s error: %w", id, err)
	}

	return groupUser, nil
}

func (gs *GroupUsersStore) SearchGroupByName(nameGroup string) (*GroupUserBL, error) {
	group, err := gs.groupUserStore.SearchGroupByName(nameGroup)
	if err != nil {
		return nil, fmt.Errorf("search group by name %s error: %w", nameGroup, err)
	}

	return group, nil
}
