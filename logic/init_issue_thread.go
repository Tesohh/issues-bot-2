package logic

import (
	"issues/v2/dataview"
	"issues/v2/db"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

func InitIssueThread(issue *db.Issue, guild *db.Guild, thread *dg.Channel, s *dg.Session, i *dg.Interaction) error {
	// we don't need to send any temporary messges. you are automatically mentioned as expected with cv2

	components := dataview.MakeIssueThreadDetail(issue, guild.NobodyRoleID)
	msg, err := s.ChannelMessageSendComplex(thread.ID, &dg.MessageSend{
		Components: components,
		Flags:      dg.MessageFlagsIsComponentsV2,
	})
	if err != nil {
		slog.Warn("error while initializing thread", "thread", thread.ID, "issue.ID", issue.ID, "err", err)
		return nil
	}
	issue.MessageID = msg.ID

	return nil
}
