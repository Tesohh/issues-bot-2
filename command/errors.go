package command

import (
	"errors"
)

var (
	ErrDuplicateProject       = errors.New("project is duplicate for this guild")
	ErrProjectNotFound        = errors.New("project was not found")
	ErrNotInIssueInputChannel = errors.New("tried to create issue in a channel that isn't an issue input one (eg. `#xxx-issues`)")
)
