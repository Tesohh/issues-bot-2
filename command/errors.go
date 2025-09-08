package command

import (
	"errors"
)

var (
	ErrDuplicateProject       = errors.New("project is duplicate for this guild")
	ErrProjectNotFound        = errors.New("project was not found")
	ErrNotInIssueInputChannel = errors.New("tried to create issue in a channel that isn't an issue input one (eg. `#xxx-issues`)")
	ErrPrefixNotSpecified     = errors.New("project not found, as command wasn't executed in a project context and `prefix` was not provided")
	ErrNotInIssueThread       = errors.New("issue not found, as command wasn't executed in an issue's thread and `code` was not specified")
)
