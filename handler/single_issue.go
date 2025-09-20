package handler

import (
	"issues/v2/command"
	"issues/v2/db"
	"issues/v2/logic"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

func issueSetStatus(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	issue, err := db.IssueQueryWithDependencies().Where("id = ?", args[1]).First(db.Ctx)
	if err != nil {
		return err
	}

	subcommand := ""
	switch args[2] {
	case "0":
		subcommand = "todo"
	case "1":
		subcommand = "doing"
	case "2":
		subcommand = "done"
	case "3":
		subcommand = "cancelled"
	}
	err = command.IssueMark(s, i.Interaction, &issue, subcommand, false)
	if err != nil {
		return err
	}

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

	err = logic.UpdateAllInteractiveIssuesViews(s, issue.ProjectID)
	if err != nil {
		return err
	}

	go func() {
		err = logic.UpdateDependencyDetails(s, i.Interaction, &issue)
		if err != nil {
			slog.Error("error while updating dependency details after running Set Status button", "issue.ID", issue.ID, "err", err)
			return
		}
	}()

	return nil
}

func issueToggleAuthorAssignee(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	issue, err := db.IssueQueryWithDependencies().Where("id = ?", args[1]).First(db.Ctx)
	if err != nil {
		return err
	}

	err = command.IssueAssign(s, i.Interaction, &issue, &dg.User{ID: i.Member.User.ID}, false)
	if err != nil {
		return err
	}

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

	err = logic.UpdateAllInteractiveIssuesViews(s, issue.ProjectID)
	if err != nil {
		return err
	}

	go func() {
		err = logic.UpdateDependencyDetails(s, i.Interaction, &issue)
		if err != nil {
			slog.Error("error while updating dependency details after running Assign Me button", "issue.ID", issue.ID, "err", err)
			return
		}
	}()

	return nil
}
