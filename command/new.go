package command

import (
	"issues/v2/slash"

	dg "github.com/bwmarrin/discordgo"
)

var New = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "new",
		Description: "Adds a new issue if you're in a #xxx-issues channel, or adds a task if you're in a issue thread",
		Options: []*dg.ApplicationCommandOption{
			{
				Name:        "title",
				Description: "the title of the issue",
				Required:    true,
				Type:        dg.ApplicationCommandOptionString,
			},
			{
				Name:        "category",
				Description: "what's the category of this issue (ie. FEATURE, FIX, OPTIMIZATION, REFACTOR, CHORE)",
				Type:        dg.ApplicationCommandOptionRole,
			},
			{
				Name:        "priority",
				Description: "what's the priority of this issue (ie. LOW, NEXT VERSION, RELEASE, IMPORTANT, CRITICAL)",
				Type:        dg.ApplicationCommandOptionRole,
			},
			{
				Name:        "discussion",
				Description: "is this a discussion?",
				Type:        dg.ApplicationCommandOptionBoolean,
			},
			{
				Name:        "tags",
				Description: "comma separated tags (eg. `tag1, tag2`)",
				Type:        dg.ApplicationCommandOptionString,
			},
			{
				Name:        "assign",
				Description: "select user to assign this to. if you want to assign this to nobody, use the `nobody` flag",
				Type:        dg.ApplicationCommandOptionUser,
			},
			{
				Name:        "nobody",
				Description: "set this to true to assign the issue to nobody",
				Type:        dg.ApplicationCommandOptionBoolean,
			},
		},
	},
	Disabled: false,
	Func: func(s *dg.Session, i *dg.Interaction) error {
		panic("TODO")
	},
}
