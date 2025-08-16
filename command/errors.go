package command

import (
	"errors"
)

var (
	ErrDuplicateProject = errors.New("project is duplicate for this guild")
	ErrProjectNotFound  = errors.New("project was not found")
)
