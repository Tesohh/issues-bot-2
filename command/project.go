package command

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/logic"
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

		var err error
		switch subcommand.Name {
		case "new":
			name := options["name"].StringValue()
			repourl := ""
			if repourlRaw, ok := options["repourl"]; ok {
				repourl = repourlRaw.StringValue()
			}

			err = ProjectNew(s, i, prefix, name, repourl)
		case "rename":
			name := options["name"].StringValue()
			err = ProjectRename(s, i, prefix, name)
		case "delete":
			confirm := options["confirm"].BoolValue()
			err = ProjectDelete(s, i, prefix, confirm)
		}

		return err
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
	generalChannel, err := s.GuildChannelCreateComplex(i.GuildID, dg.GuildChannelCreateData{
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
		GeneralChannelID:         generalChannel.ID,
		AutoListMessageID:        "", // TODO:
		GuildID:                  i.GuildID,
	}

	err = db.Projects.Create(db.Ctx, &project)
	if err != nil {
		return err
	}

	for i := range project.Issues {
		project.Issues[i].Project.GuildID = project.GuildID
	}

	state := db.ProjectViewState{
		ProjectID:   project.ID,
		Project:     project,
		Filter:      db.DefaultFilter(),
		Sorter:      db.DefaultSorter(),
		Permanent:   true,
		ListNameFmt: "# AutoList™️ for %s `[%s]`",
	}

	err = logic.InitIssueViewDetached(s, i, inputChannel.ID, &state)
	if err != nil {
		return err
	}

	embed := dg.MessageEmbed{
		Title:       fmt.Sprintf("Created project %s [`%s`]", name, strings.ToUpper(prefix)),
		Description: fmt.Sprintf("Check out <#%s>", inputChannel.ID),
	}
	return slash.ReplyWithEmbed(s, i, embed, false)
}

func ProjectRename(s *dg.Session, i *dg.Interaction, prefix string, name string) error {
	query := db.Projects.Where("prefix = ? AND guild_id = ?", prefix, i.GuildID)

	rowsAffected, err := query.Update(db.Ctx, "name", name)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrProjectNotFound
	}

	project, err := query.Select("discord_category_channel_id").First(db.Ctx)
	if err != nil {
		return err
	}

	_, err = s.ChannelEdit(project.DiscordCategoryChannelID, &dg.ChannelEdit{Name: name})
	if err != nil {
		return err
	}

	// TODO: update autolist

	embed := dg.MessageEmbed{
		Title: fmt.Sprintf("Successfully renamed `%s` to %s", prefix, name),
	}
	return slash.ReplyWithEmbed(s, i, embed, false)
}

func ProjectDelete(s *dg.Session, i *dg.Interaction, prefix string, confirmation bool) error {
	if !confirmation {
		return slash.ReplyWithEmbed(s, i, dg.MessageEmbed{
			Title: "alright, no actions taken",
		}, true)
	}

	// get project true id
	project, err := db.Projects.
		Select("id, discord_category_channel_id, general_channel_id, issues_input_channel_id").
		Where("prefix = ? AND guild_id = ?", prefix, i.GuildID).
		First(db.Ctx)
	if err != nil {
		return err
	}

	// delete all issues under this project
	_, err = db.Issues.Where("project_id = ?", project.ID).Delete(db.Ctx)
	if err != nil {
		return err
	}

	// delete the project in the DB
	_, err = db.Projects.Where("id = ?", project.ID).Delete(db.Ctx)
	if err != nil {
		return err
	}

	// delete the projectviewstates
	_, err = db.ProjectViewStates.Where("project_id = ?", project.ID).Delete(db.Ctx)
	if err != nil {
		return err
	}

	// delete the discord channels
	for _, id := range []string{project.DiscordCategoryChannelID, project.GeneralChannelID, project.IssuesInputChannelID} {
		_, err := s.ChannelDelete(id)
		if err != nil {
			return err
		}
	}

	embed := dg.MessageEmbed{
		Title: fmt.Sprintf("Deleted project `%s`", strings.ToUpper(prefix)),
	}

	return slash.ReplyWithEmbed(s, i, embed, false)
}
