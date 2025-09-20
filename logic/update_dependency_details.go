package logic

import (
	"issues/v2/db"

	dg "github.com/bwmarrin/discordgo"
)

func UpdateDependencyDetails(s *dg.Session, i *dg.Interaction, issue *db.Issue) error {
	relationships, err := GetIssueRelationshipsOfKind(issue, db.RelationshipKindDependency)
	if err != nil {
		return err
	}

	guild, err := db.Guilds.Select("nobody_role_id").Where("id = ?", i.GuildID).First(db.Ctx)
	if err != nil {
		return err
	}

	// dependants
	for _, r := range relationships.Inbound {
		target, err := db.IssueQueryWithDependencies().Where("id = ?", r.FromIssueID).First(db.Ctx)
		if err != nil {
			return err
		}

		targetRelationships, err := GetIssueRelationshipsOfKind(&target, db.RelationshipKindDependency)
		if err != nil {
			return err
		}

		err = UpdateIssueThreadDetail(s, &target, targetRelationships, guild.NobodyRoleID)
		if err != nil {
			return err
		}
	}

	// dependencies
	for _, r := range relationships.Outbound {
		target, err := db.IssueQueryWithDependencies().Where("id = ?", r.ToIssueID).First(db.Ctx)
		if err != nil {
			return err
		}

		targetRelationships, err := GetIssueRelationshipsOfKind(&target, db.RelationshipKindDependency)
		if err != nil {
			return err
		}

		err = UpdateIssueThreadDetail(s, &target, targetRelationships, guild.NobodyRoleID)
		if err != nil {
			return err
		}
	}

	return nil

}
