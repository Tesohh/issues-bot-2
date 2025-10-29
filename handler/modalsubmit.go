package handler

import (
	"issues/v2/slash"
	"log/slog"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

type modalSubmitHandlerFunc func(s *dg.Session, i *dg.InteractionCreate, args []string) error

type modalSubmitHandler struct {
	argsCount int
	ack       bool
	handler   modalSubmitHandlerFunc
}

var modalSubmitHandlers = map[string]modalSubmitHandler{
	"ping": {1, true, func(s *dg.Session, i *dg.InteractionCreate, args []string) error {
		_, err := s.ChannelMessageSend(i.ChannelID, args[1])
		return err
	}},
	"issues_filter_people_submit": {1, true, issuesFilterPeopleSubmit},
	"issues_filter_data_submit":   {1, true, issuesFilterDataSubmit},
}

// component custom ids need to be in this format: action:arg0:arg1
func ModalSubmit(s *dg.Session, i *dg.InteractionCreate) {
	id := i.ModalSubmitData().CustomID
	args := strings.Split(id, ":")

	if len(args) == 0 {
		slog.Warn("just got a modal submit interaction, with 0 args")
		return
	}

	handler, ok := modalSubmitHandlers[args[0]]
	if !ok {
		slog.Error("modal submit handler for this customid not found", "id", args[0])
		return
	}

	if handler.argsCount+1 > len(args) {
		slog.Error("modal submit called with not enough arguments", "id", args[0], "args", args)
		return
	}

	err := handler.handler(s, i, args)
	if err != nil {
		embed := dg.MessageEmbed{
			Title:       "Error",
			Description: err.Error(),
			Color:       0xFF0000,
		}
		slash.ReplyWithEmbed(s, i.Interaction, embed, true)
		slog.Error("modal submit handler error", "err", err, "args", args)
	}
	if handler.ack {
		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: dg.InteractionResponseUpdateMessage,
		})
	}
}
