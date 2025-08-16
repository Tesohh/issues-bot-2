package command

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/slash"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

var Project = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "project",
		Description: "project management command",
		Options: []*dg.ApplicationCommandOption{
			{
				Name:        "new",
				Description: "creates a new project",
				Options: []*dg.ApplicationCommandOption{
					{
						Name:        "prefix",
						Description: "3 letter prefix for your project (eg. KER, PYC, CVV)",
						Required:    true,
						Type:        dg.ApplicationCommandOptionString,
						MinLength:   slash.Ptr(3),
						MaxLength:   3,
					},
					{
						Name:        "name",
						Description: "full name for your project (eg. Open Classeviva)",
						Required:    true,
						Type:        dg.ApplicationCommandOptionString,
					},
					{
						Name:        "repourl",
						Description: "(optional) URL to your repo",
						Type:        dg.ApplicationCommandOptionString,
					},
				},
				Type: dg.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "rename",
				Description: "renames an existing project",
				Options: []*dg.ApplicationCommandOption{
					{
						Name:        "prefix",
						Description: "3 letter prefix for the existing project (eg. KER, PYC, CVV)",
						Required:    true,
						Type:        dg.ApplicationCommandOptionString,
						MinLength:   slash.Ptr(3),
						MaxLength:   3,
					},
					{
						Name:        "name",
						Description: "new full name for your project (eg. Open Classeviva)",
						Required:    true,
						Type:        dg.ApplicationCommandOptionString,
					},
				},
				Type: dg.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "delete",
				Description: "deletes an existing project",
				Options: []*dg.ApplicationCommandOption{
					{
						Name:        "prefix",
						Description: "3 letter prefix for the existing project (eg. KER, PYC, CVV)",
						Required:    true,
						Type:        dg.ApplicationCommandOptionString,
						MinLength:   slash.Ptr(3),
						MaxLength:   3,
					},
					{
						Type:        dg.ApplicationCommandOptionBoolean,
						Name:        "confirm",
						Description: "are you really sure?",
						Required:    true,
					},
				},
				Type: dg.ApplicationCommandOptionSubCommand,
			},
		},
	},
	Disabled: false,
	Func: func(s *dg.Session, i *dg.Interaction) error {
		subcommand := i.ApplicationCommandData().Options[0]
		options := slash.GetOptionMapRaw(subcommand.Options)
		prefix := strings.ToLower(options["prefix"].StringValue())

		switch subcommand.Name {
		case "new":
			name := options["name"].StringValue()
			repourl := ""
			if repourlRaw, ok := options["repourl"]; ok {
				repourl = repourlRaw.StringValue()
			}

			err := ProjectNew(s, i, prefix, name, repourl)
			if err != nil {
				return err
			}
		case "rename":
		case "delete":
		}
		return nil
	},
}

func ProjectNew(s *dg.Session, i *dg.Interaction, prefix string, name string, repourl string) error {
	count, err := db.Projects.Where("prefix = ? AND guild_id = ?", prefix, i.GuildID).Count(db.Ctx, "*")
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrDuplicateProject
	}

	// at this point, the project does not already exist
	// so now create the channels

	category, err := s.GuildChannelCreate(i.GuildID, name, dg.ChannelTypeGuildCategory)
	if err != nil {
		return err
	}
	prefix = strings.ToLower(prefix)

	generalChannelName := fmt.Sprintf("%s-general", prefix)
	_, err = s.GuildChannelCreateComplex(i.GuildID, dg.GuildChannelCreateData{
		Name:     generalChannelName,
		Type:     dg.ChannelTypeGuildText,
		Topic:    repourl,
		ParentID: category.ID,
	})
	if err != nil {
		return err
	}

	inputChannelName := fmt.Sprintf("%s-issues", prefix)
	inputChannel, err := s.GuildChannelCreateComplex(i.GuildID, dg.GuildChannelCreateData{
		Name:     inputChannelName,
		Type:     dg.ChannelTypeGuildText,
		Topic:    repourl,
		ParentID: category.ID,
	})
	if err != nil {
		return err
	}

	project := db.Project{
		Name:                     name,
		Prefix:                   prefix,
		RepoURL:                  repourl,
		DiscordCategoryChannelID: category.ID,
		IssuesInputChannelID:     inputChannel.ID,
		AutoListMessageID:        "", // TODO:
		GuildID:                  i.GuildID,
	}

	err = db.Projects.Create(db.Ctx, &project)
	if err != nil {
		return err
	}

	embed := dg.MessageEmbed{
		Title:       fmt.Sprintf("Created project %s [%s]", name, strings.ToUpper(prefix)),
		Description: fmt.Sprintf("Check out <#%s>", inputChannel.ID),
	}

	err = slash.ReplyWithEmbed(s, i, embed, false)
	if err != nil {
		return err
	}

	return nil
}
