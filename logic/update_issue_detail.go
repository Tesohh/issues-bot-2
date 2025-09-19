package logic

import (
	"issues/v2/dataview"
	"issues/v2/db"

	dg "github.com/bwmarrin/discordgo"
)

func UpdateIssueThreadDetail(s *dg.Session, issue *db.Issue, relationships db.RelationshipsByDirection, nobodyRoleID string) error {
	detail := dataview.MakeIssueThreadDetail(issue, relationships, nobodyRoleID)
	_, err := s.ChannelMessageEditComplex(&dg.MessageEdit{
		Components: &detail,
		Channel:    issue.ThreadID,
		ID:         issue.MessageID,
	})

	return err
}
