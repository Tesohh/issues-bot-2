package logic

import (
	"issues/v2/dataview"
	"issues/v2/db"
	"issues/v2/slash"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

func InitIssueThread(issue *db.Issue, guild *db.Guild, thread *dg.Channel, s *dg.Session, i *dg.Interaction) error {
	mentionees := []string{issue.RecruiterUserID}
	for _, user := range issue.AssigneeUsers {
		mentionees = append(mentionees, user.ID)
	}

	deleteMe, _ := s.ChannelMessageSendComplex(thread.ID, &dg.MessageSend{
		Content: slash.MentionMany(mentionees, "@", ""),
		Flags:   dg.MessageFlagsSuppressNotifications,
	})
	_ = s.ChannelMessageDelete(thread.ID, deleteMe.ID)

	components := dataview.MakeIssueThreadDetail(issue, guild.NobodyRoleID)
	msg, err := s.ChannelMessageSendComplex(thread.ID, &dg.MessageSend{
		Components: components,
		Flags:      dg.MessageFlagsIsComponentsV2,
	})
	if err != nil {
		slog.Warn("error while initializing thread", "thread", thread.ID, "issue.ID", issue.ID, "err", err)
		return nil
	}

	_, err = db.Issues.Where("id = ?", issue.ID).Update(db.Ctx, "message_id", msg.ID)
	if err != nil {
		return err
	}

	return nil
}
