package command

import (
	"issues/v2/slash"
)

var Commands = map[string]*slash.Command{
	"project": &Project,
	"issue":   &Issue,
	"task":    &Task,
	"new":     &New,
	"list":    &List,
	"reset":   &Reset,
	"man":     &Man,
}

func init() {
	Commands["help"] = &Help
}
