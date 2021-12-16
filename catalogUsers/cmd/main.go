package main

import (
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/catalogBL"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/repository/groupUserBL"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/app/repository/usersBL"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/db/memDB/memGroupDB"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/catalogUsers/internal/db/memDB/memUserDB"
)

func main() {
	memUserDB := memUserDB.NewMemUserDB()
	memGroupDB := memGroupDB.NewMemGroup()

	userBL := usersBL.NewUsersStore(memUserDB)
	groupBL := groupUserBL.NewGroupUsersStore(memGroupDB)

	catalog := catalogBL.NewCatalog(userBL, groupBL)

}
