package logic

import (
	"issues/v2/db"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

// expects issue.Project to be set
//
// also updates the DB entry for that issue with the new ThreadID
func CreateThreadFromIssue(issue *db.Issue, s *dg.Session, i *dg.Interaction) (*dg.Channel, error) {
	thread, err := s.ThreadStart(issue.Project.IssuesInputChannelID, issue.ChannelName(), dg.ChannelTypeGuildPublicThread, 10080)
	if err != nil {
		return nil, err
	}
	issue.ThreadID = thread.ID

	err = s.ChannelMessageDelete(issue.Project.IssuesInputChannelID, thread.ID)
	if err != nil {
		slog.Warn("couldn't delete thread start message. no big deal", "err", err)
	}

	return thread, nil
}
