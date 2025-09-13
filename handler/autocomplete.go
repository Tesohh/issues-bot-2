package handler

import (
	"issues/v2/db"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

func Autocomplete(s *dg.Session, i *dg.InteractionCreate) {
	command := i.ApplicationCommandData()
	choices := []*dg.ApplicationCommandOptionChoice{}
	respond := true

	switch command.Name {
	case "issue":
		subcommand := command.Options[0]
		switch subcommand.Name {
		case "tag":
			value := ""
			if len(subcommand.Options) > 0 {
				value = subcommand.Options[0].StringValue()
			}

			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				slog.Error("error while getting channel for issue tags autocomplete", "err", err)
				return
			}

			project, err := db.Projects.
				Select("id").
				Where("discord_category_channel_id = ?", channel.ParentID).
				Or("issues_input_channel_id = ?", channel.ParentID).
				First(db.Ctx)
			if err != nil {
				slog.Error("error while getting project for issue tags autocompelte", "err", err)
				return
			}

			// get all tags from this project
			tags, err := db.Tags.Where("project_id = ?", project.ID).Find(db.Ctx)
			if err != nil {
				slog.Error("error while getting all tags for issue tags autocompelte", "err", err)
				return
			}
			_ = tags

			// TODO: get the issue, check it's tags and then do completion
		}
	}

	if respond {
		err := s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: dg.InteractionApplicationCommandAutocompleteResult,
			Data: &dg.InteractionResponseData{Choices: choices},
		})
		if err != nil {
			slog.Error("error while responding to autocomplete", "err", err)
			return
		}
	}
}
