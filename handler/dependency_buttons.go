package handler

import (
	"issues/v2/db"
	"issues/v2/logic"

	dg "github.com/bwmarrin/discordgo"
)

func issueDepsGoto(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	_, err := db.Issues.Where("id = ?", args[1]).Update(db.Ctx, "UIDepsCurrentPage", args[2])
	if err != nil {
		return err
	}

	issue, err := db.IssueQueryWithDependencies().Where("id = ?", args[1]).First(db.Ctx)

	guild, err := db.Guilds.Select("nobody_role_id").Where("id = ?", i.GuildID).First(db.Ctx)
	if err != nil {
		return err
	}

	relationships, err := logic.GetIssueRelationshipsOfKind(&issue, db.RelationshipKindDependency)
	if err != nil {
		return err
	}

	err = logic.UpdateIssueThreadDetail(s, &issue, relationships, guild.NobodyRoleID)
	if err != nil {
		return err
	}

	return nil
}

func issueDepSetStatus(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	_, err := db.Issues.Where("id = ?", args[2]).Update(db.Ctx, "status", args[3])
	if err != nil {
		return err
	}

	issue, err := db.IssueQueryWithDependencies().Where("id = ?", args[1]).First(db.Ctx)

	guild, err := db.Guilds.Select("nobody_role_id").Where("id = ?", i.GuildID).First(db.Ctx)
	if err != nil {
		return err
	}

	relationships, err := logic.GetIssueRelationshipsOfKind(&issue, db.RelationshipKindDependency)
	if err != nil {
		return err
	}

	err = logic.UpdateIssueThreadDetail(s, &issue, relationships, guild.NobodyRoleID)
	if err != nil {
		return err
	}

	return nil
}
