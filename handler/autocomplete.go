package handler

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/man"
	"log/slog"
	"strings"
	"time"

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
		case "dependson":
			start := time.Now()
			ch, err := s.Channel(i.ChannelID)
			if err != nil {
				slog.Error("error while fetching channel during issue completions", "err", err)
				return
			}
			fmt.Printf("(channel) time elapsed ms: %d\n", time.Since(start).Milliseconds())

			parent, err := s.Channel(ch.ParentID)
			if err != nil {
				slog.Error("error while fetching parent channel during issue completions", "err", err)
				return
			}
			fmt.Printf("(parent) time elapsed ms: %d\n", time.Since(start).Milliseconds())
			theFamily := []string{parent.ID}

			grandParent, err := s.Channel(parent.ParentID)
			if err == nil {
				theFamily = append(theFamily, grandParent.ID)
				return
			}
			fmt.Printf("(grandparent) time elapsed ms: %d\n", time.Since(start).Milliseconds())

			project, err := db.Projects.
				Select("id, prefix").
				Where("discord_category_channel_id IN ?", theFamily).
				First(db.Ctx)

			fmt.Printf("(project) time elapsed ms: %d\n", time.Since(start).Milliseconds())

			search := subcommand.Options[0].StringValue()
			issues, err := db.Issues.
				Select("id, status, code, title").
				Where("project_id = ?", project.ID).
				Where("title LIKE ?", "%"+search+"%").
				Limit(5).
				Find(db.Ctx)
			if err != nil {
				slog.Error("error while fetching issue completions", "err", err)
				return
			}
			fmt.Printf("(issues) time elapsed ms: %d\n", time.Since(start).Milliseconds())

			for i := range issues {
				issues[i].Project = project

				choices = append(choices, &dg.ApplicationCommandOptionChoice{
					Name:  issues[i].ChannelName(),
					Value: fmt.Sprint(issues[i].ID),
				})
			}

			fmt.Printf("(done) time elapsed ms: %d\n", time.Since(start).Milliseconds())
		}
	case "man":
		search := command.Options[0].StringValue()
		for _, page := range man.Pages {
			if strings.Contains(strings.ToLower(page.Title), strings.ToLower(search)) {
				choices = append(choices, &dg.ApplicationCommandOptionChoice{
					Name:  page.Title,
					Value: page.ID,
				})
			}
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

// func issueAutocompletion(search )
