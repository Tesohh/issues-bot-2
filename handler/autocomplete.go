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
			var err error
			choices, err = issueAutocomplete(s, i, subcommand.Options[0])
			if err != nil {
				slog.Error("error while executing issueAutocomplete", "err", err)
				return
			}
		}

	case "new":
		var err error
		choices, err = issueAutocomplete(s, i, command.GetOption("dependson"))
		if err != nil {
			slog.Error("error while executing issueAutocomplete", "err", err)
			return
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

func issueAutocomplete(s *dg.Session, i *dg.InteractionCreate, searchOpt *dg.ApplicationCommandInteractionDataOption) ([]*dg.ApplicationCommandOptionChoice, error) {
	parentID, err := getChannelParentID(s, i.ChannelID)
	if err != nil {
		return nil, err
	}
	theFamily := []string{parentID}

	grandParentID, err := getChannelParentID(s, parentID)
	if err == nil {
		theFamily = append(theFamily, grandParentID)
	}

	project, err := db.Projects.
		Select("id, prefix").
		Where("discord_category_channel_id IN ?", theFamily).
		First(db.Ctx)

	// search := subcommand.Options[0].StringValue()
	search := searchOpt.StringValue()
	issues, err := db.Issues.
		Select("id, status, code, title").
		Where("project_id = ?", project.ID).
		Where("title LIKE ?", "%"+search+"%").
		Limit(5).
		Find(db.Ctx)
	if err != nil {
		return nil, err
	}

	choices := []*dg.ApplicationCommandOptionChoice{}
	for i := range issues {
		issues[i].Project = project

		choices = append(choices, &dg.ApplicationCommandOptionChoice{
			Name:  issues[i].ChannelName(),
			Value: fmt.Sprint(issues[i].ID),
		})
	}

	return choices, nil
}

type channelParentCacheEntry struct {
	lastUpdate time.Time
	parentID   string
}

var channelParentCache = map[string]channelParentCacheEntry{}

func getChannelParentID(s *dg.Session, channelID string) (string, error) {
	if entry, ok := channelParentCache[channelID]; ok {
		if time.Since(entry.lastUpdate) > 1*time.Hour {
			go func() {
				ch, err := s.Channel(channelID)
				if err != nil {
					slog.Error("error while background fetching", "channelID", channelID, "err", err)
					return
				}
				channelParentCache[channelID] = channelParentCacheEntry{lastUpdate: time.Now(), parentID: ch.ParentID}
			}()
		}
		return entry.parentID, nil
	} else {
		ch, err := s.Channel(channelID)
		if err != nil {
			return "", err
		}
		channelParentCache[channelID] = channelParentCacheEntry{lastUpdate: time.Now(), parentID: ch.ParentID}
		return ch.ParentID, nil
	}
}
