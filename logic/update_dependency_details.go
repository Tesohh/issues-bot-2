package logic

import (
	"fmt"
	"issues/v2/db"

	dg "github.com/bwmarrin/discordgo"
)

func UpdateDependencyDetails(s *dg.Session, guildID string, issue *db.Issue) error {
	relationships, err := GetIssueRelationshipsOfKind(issue, db.RelationshipKindDependency)
	if err != nil {
		return fmt.Errorf("UpdateDependencyDetails -> GetIssueRelationshipsOfKind -> %w", err)
	}

	guild, err := db.Guilds.Select("nobody_role_id").Where("id = ?", guildID).First(db.Ctx)
	if err != nil {
		return fmt.Errorf("UpdateDependencyDetails -> fetch guild -> %w", err)
	}

	// dependants
	for _, r := range relationships.Inbound {
		target, err := db.IssueQueryWithDependencies().Where("id = ?", r.FromIssueID).First(db.Ctx)
		if err != nil {
			return fmt.Errorf("UpdateDependencyDetails -> get dependant target -> %w", err)
		}

		targetRelationships, err := GetIssueRelationshipsOfKind(&target, db.RelationshipKindDependency)
		if err != nil {
			return fmt.Errorf("UpdateDependencyDetails -> get dependant target relationships -> %w", err)
		}

		err = UpdateIssueThreadDetail(s, &target, targetRelationships, guild.NobodyRoleID)
		if err != nil {
			return fmt.Errorf("UpdateDependencyDetails -> update dependant's detail -> %w", err)
		}
	}

	// dependencies
	for _, r := range relationships.Outbound {
		target, err := db.IssueQueryWithDependencies().Where("id = ?", r.ToIssueID).First(db.Ctx)
		if err != nil {
			return fmt.Errorf("UpdateDependencyDetails -> get dependency target -> %w", err)
		}

		if target.Kind == db.IssueKindTask {
			continue
		}

		targetRelationships, err := GetIssueRelationshipsOfKind(&target, db.RelationshipKindDependency)
		if err != nil {
			return fmt.Errorf("UpdateDependencyDetails -> get dependency relationships -> %w", err)
		}

		err = UpdateIssueThreadDetail(s, &target, targetRelationships, guild.NobodyRoleID)
		if err != nil {
			return fmt.Errorf("UpdateDependencyDetails -> update dependency detail -> %w", err)
		}
	}

	return nil

}
