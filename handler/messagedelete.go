package handler

import (
	"issues/v2/db"

	"github.com/bwmarrin/discordgo"
)

func MessageDelete(s *discordgo.Session, msg *discordgo.MessageDelete) {
	db.ProjectViewStates.Where("message_id = ?", msg.ID).Delete(db.Ctx)
}
