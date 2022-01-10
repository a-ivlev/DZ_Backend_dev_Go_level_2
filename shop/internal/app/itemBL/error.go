package itemBL

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrItemNotFound = errors.New("item not found")
	ErrNoRows       = errors.New("no rows in result set")
	ErrDeleted      = errors.New("the element has already been removed")
)
