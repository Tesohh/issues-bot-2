package handler

import (
	"issues/v2/slash"
	"log/slog"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

type messageComponentHandlerFunc func(s *dg.Session, i *dg.InteractionCreate, args []string) error

type messageComponentHandler struct {
	argsCount int
	ack       bool
	handler   messageComponentHandlerFunc
}

var messageComponentHandlers = map[string]messageComponentHandler{
	"ping": {1, true, func(s *dg.Session, i *dg.InteractionCreate, args []string) error {
		_, err := s.ChannelMessageSend(i.ChannelID, args[1])
		return err
	}},

	"issues-goto":          {2, true, issuesGoto},
	"issues-set-statuses":  {2, true, issuesSetStatuses},
	"issues-sort-by":       {2, true, issuesSortBy},
	"issues-order":         {2, true, issuesOrder},
	"issues-filter-people": {1, true, issuesFilterPeople},
	"issues-filter-data":   {1, true, issuesFilterData},

	"issue-set-status":             {2, true, issueSetStatus},
	"issue-toggle-author-assignee": {1, true, issueToggleAuthorAssignee},

	"issue-deps-goto":      {2, true, issueDepsGoto},
	"issue-dep-set-status": {3, true, issueDepSetStatus},
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

	if handler.argsCount+1 > len(args) {
		slog.Error("message component called with not enough arguments", "id", args[0], "args", args)
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
		slog.Error("message component handler error", "err", err, "args", args)
	}
	if handler.ack {
		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: dg.InteractionResponseUpdateMessage,
		})
	}
}
