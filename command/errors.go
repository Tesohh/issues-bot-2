package command

import (
	"errors"
)

var (
	ErrDuplicateProject = errors.New("project is duplicate for this guild")
)
