package catalogBL

import (
	"fmt"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/repository/groupUserBL"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/repository/usersBL"
	"github.com/google/uuid"
)

type Catalog struct {
	userStr  *usersBL.UsersStore
	groupStr *groupUserBL.GroupUsersStore
}

func NewCatalog(userStr *usersBL.UsersStore, groupStr *groupUserBL.GroupUsersStore) *Catalog {
	return &Catalog{
		userStr:  userStr,
		groupStr: groupStr,
	}
}

func (c *Catalog) CreateUser(user usersBL.UserBL) (*usersBL.UserBL, error) {
	createUser, err := c.userStr.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}

	return createUser, nil
}

func (c *Catalog) CreateGroupUser(groupUser groupUserBL.GroupUserBL) (*groupUserBL.GroupUserBL, error) {
	createUser, err := c.groupStr.CreateGroupUser(groupUser)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}

	return createUser, nil
}

func (c *Catalog) AddUserToGroup(idGroup uuid.UUID, user usersBL.UserBL) (*groupUserBL.GroupUserBL, error) {
	userGroup, err := c.groupStr.ReadGroupUser(idGroup)
	if err != nil {
		return nil, fmt.Errorf("read group ID %s error: %w", idGroup, err)
	}

	userGroup.UsersList[user.ID] = user
	user.GroupList[userGroup.ID] = *userGroup

	return userGroup, nil
}

func (c *Catalog) DeleteUserFromGroup(idGroup uuid.UUID, user usersBL.UserBL) (*groupUserBL.GroupUserBL, error) {
	userGroup, err := c.groupStr.ReadGroupUser(idGroup)
	if err != nil {
		return nil, fmt.Errorf("read group ID %s error: %w", idGroup, err)
	}

	delete(userGroup.UsersList, user.ID)

	return userGroup, nil
}

func (c *Catalog) SearchUserByName(userName string) (*usersBL.UserBL, error) {
	user, err := c.userStr.SearchUserByName(userName)
	if err != nil {
		return nil, fmt.Errorf("search user name %s error: %w", userName, err)
	}

	return user, nil
}

func (c *Catalog) SearchUsersByGroupName(groupName string) ([]usersBL.UserBL, error) {
	var users []usersBL.UserBL
	group, err := c.groupStr.SearchGroupByName(groupName)
	if err != nil {
		return nil, fmt.Errorf("search group name %s error: %w", groupName, err)
	}

	for _, user := range group.UsersList {
		users = append(users, user)
	}

	return users, nil
}

func (c *Catalog) SearchGroupByName(groupName string) (*groupUserBL.GroupUserBL, error) {
	group, err := c.groupStr.SearchGroupByName(groupName)
	if err != nil {
		return nil, fmt.Errorf("search group name %s error: %w", groupName, err)
	}

	return group, nil
}

func (c *Catalog) SearchGroupByUserName(userName string) ([]groupUserBL.GroupUserBL, error) {
	var groups []groupUserBL.GroupUserBL
	user, err := c.userStr.SearchUserByName(userName)
	if err != nil {
		return nil, fmt.Errorf("search group name %s error: %w", userName, err)
	}

	for _, group := range user.GroupList {
		groups = append(groups, group)
	}

	return groups, nil
}
