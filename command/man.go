package command

import (
	"issues/v2/man"
	"issues/v2/slash"

	dg "github.com/bwmarrin/discordgo"
)

var Man = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "man",
		Description: "read manual pages",
		Options: []*dg.ApplicationCommandOption{
			{
				Type:         dg.ApplicationCommandOptionString,
				Name:         "page",
				Description:  "what page do you want to read?",
				Required:     true,
				Autocomplete: true,
			},
		},
	},
	Func: func(s *dg.Session, i *dg.Interaction) error {
		options := slash.GetOptionMap(i)
		id := options["page"].StringValue()
		page, ok := man.Pages[id]
		if !ok {
			return ErrManPageDoesNotExist
		}

		components := []dg.MessageComponent{
			dg.TextDisplay{Content: "# " + page.Title},
		}
		components = append(components, page.Content...)

		return slash.ReplyWithComponents(s, i, true, dg.Container{
			Components: components,
		})
	},
}
