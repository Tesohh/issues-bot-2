package handler

import (
	"log/slog"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

type MessageComponentHandlerFunc func(s *dg.Session, i *dg.InteractionCreate, args []string) error

var messageComponentHandlers = map[string]MessageComponentHandlerFunc{
	"ping": func(s *dg.Session, i *dg.InteractionCreate, args []string) error {
		_, err := s.ChannelMessageSend(i.ChannelID, "pong")
		return err
	},
}

// component custom ids need to be in this format: action:arg0:arg1
func MessageComponent(s *dg.Session, i *dg.InteractionCreate) {
	id := i.MessageComponentData().CustomID
	args := strings.Split(id, ":")

	if len(args) == 0 {
		slog.Warn("just got a message component interaction, with 0 args")
		return
	}

	handler, ok := messageComponentHandlers[args[0]]
	if !ok {
		slog.Error("message component handler for this customid not found", "id", args[0])
		return
	}

	err := handler(s, i, args)
	if err != nil {
		slog.Error("message component handler error", "err", err, "args", args)
	}
}
