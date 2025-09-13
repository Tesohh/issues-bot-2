package command

import (
	"issues/v2/slash"
)

var Commands = map[string]*slash.Command{
	"project": &Project,
	"issue":   &Issue,
	"new":     &New,
	"list":    &List,
	"reset":   &Reset,
}

func init() {
	Commands["help"] = &Help
}
