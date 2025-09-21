package logic

import (
	"issues/v2/db"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

func UpdateEverythingAboutSingleIssue(s *dg.Session, i *dg.Interaction, issue *db.Issue) error {
	guild, err := db.Guilds.Select("nobody_role_id").Where("id = ?", i.GuildID).First(db.Ctx)
	if err != nil {
		return err
	}

	relationships, err := GetIssueRelationshipsOfKind(issue, db.RelationshipKindDependency)
	if err != nil {
		return err
	}

	err = UpdateIssueThreadDetail(s, issue, relationships, guild.NobodyRoleID)
	if err != nil {
		return err
	}

	err = UpdateAllInteractiveIssuesViews(s, issue.ProjectID)
	if err != nil {
		return err
	}

	go func() {
		err = UpdateDependencyDetails(s, i, issue)
		if err != nil {
			slog.Error("error while updating dependency details after UpdateEverythingAboutSingleIssue", "issue.ID", issue.ID, "err", err)
			return
		}
	}()

	return nil
}
