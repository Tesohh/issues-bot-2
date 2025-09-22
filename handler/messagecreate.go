package handler

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/logic"
	"log/slog"
	"regexp"
	"slices"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func MessageCreate(s *dg.Session, m *dg.MessageCreate) {
	err := messageCreate(s, m)
	if err != nil {
		slog.Error(err.Error())
		embed := dg.MessageEmbed{
			Title:       "Error",
			Description: err.Error(),
			Color:       0xFF0000,
		}
		s.ChannelMessageSendEmbedReply(m.ChannelID, &embed, m.Reference())
	}
}

var parser = regexp.MustCompile(`(?mi)<@&(?P<RoleID>\d+)>|<@(?P<UserID>\d+)>|\+(?P<Tag>\w+)|depends on\s*<#(?P<DependsOnChannelID>\d+)>`)

type parserCaptures struct {
	Raws                []string
	RoleIDs             []string
	UserIDs             []string
	Tags                []string
	DependsOnChannelIDs []string
}

type parserMentions struct {
	priorityRoleID string
	categoryRoleID string
	nobody         bool
	discussion     bool
}

func messageCreate(s *dg.Session, m *dg.MessageCreate) error {
	if !strings.HasPrefix(m.Content, "-") {
		return nil
	}

	m.Content = strings.TrimLeft(m.Content, "- ")

	// check if we're in a project
	project, err := db.Projects.Where("issues_input_channel_id = ?", m.ChannelID).First(db.Ctx)
	if err == nil {
		return addIssueFromShorthand(s, m, &project)
	} else if err != gorm.ErrRecordNotFound {
		// bubble up every error except for RecordNotFound
		return err
	}

	issue, err := db.IssueQueryWithDependencies().Where("thread_id = ?", m.ChannelID).First(db.Ctx)
	if err == nil {
		return addTaskFromShorthand(s, m, &issue)
	} else if err != gorm.ErrRecordNotFound {
		// bubble up every error except for RecordNotFound
		return nil
	}

	// the user is neither trying to add an issue, neither a task.
	// they probably are just trying to input a list or something
	return nil
}

func addTaskFromShorthand(s *dg.Session, m *dg.MessageCreate, issue *db.Issue) error {
	title := strings.Join(strings.Fields(m.Content), " ")

	// create the task in the DB
	task := db.Issue{
		Code:            nil,
		Title:           title,
		Status:          db.IssueStatusTodo,
		Kind:            db.IssueKindTask,
		ProjectID:       issue.ProjectID,
		RecruiterUserID: issue.RecruiterUserID,
		AssigneeUsers:   issue.AssigneeUsers,
		CategoryRoleID:  issue.CategoryRoleID,
		PriorityRoleID:  issue.PriorityRoleID,
	}

	err := db.Issues.Create(db.Ctx, &task)
	if err != nil {
		return err
	}

	// create the relationship in the DB
	relationship := db.Relationship{
		FromIssueID: issue.ID,
		ToIssueID:   task.ID,
		Kind:        db.RelationshipKindDependency,
	}

	err = db.Relationships.Create(db.Ctx, &relationship)
	if err != nil {
		return err
	}

	s.ChannelMessageDelete(m.ChannelID, m.ID)

	msg := fmt.Sprintf("<@%s> added task `%s`", m.Author.ID, task.CutTitle(25))
	_, err = s.ChannelMessageSend(issue.ThreadID, msg)
	if err != nil {
		return err
	}

	return logic.UpdateEverythingAboutSingleIssue(s, m.GuildID, issue)
}

func addIssueFromShorthand(s *dg.Session, m *dg.MessageCreate, project *db.Project) error {
	// get the guild
	guild, err := db.Guilds.
		Select("id, generic_category_role_id, normal_priority_role_id, nobody_role_id").
		Where("id = ?", m.GuildID).
		First(db.Ctx)
	if err != nil {
		return err
	}

	// delete the message
	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)

	// parse the message
	captures := parserCaptures{}
	ptrs := []*[]string{&captures.Raws, &captures.RoleIDs, &captures.UserIDs, &captures.Tags, &captures.DependsOnChannelIDs}

	matches := parser.FindAllStringSubmatch(m.Content, -1)
	for _, match := range matches {
		for i, subexp := range match {
			if subexp != "" {
				*ptrs[i] = append(*ptrs[i], subexp)
			}
		}
	}

	// remove captures from the message
	for _, raw := range captures.Raws {
		m.Content = strings.ReplaceAll(m.Content, raw, "")
	}
	title := strings.Join(strings.Fields(m.Content), " ")

	// figure out roles
	roles, err := db.Roles.Select("id, kind").Where("id in ?", captures.RoleIDs).Find(db.Ctx)

	mentions := parserMentions{
		priorityRoleID: guild.NormalPriorityRoleID,
		categoryRoleID: guild.GenericCategoryRoleID,
	}
	for _, role := range roles {
		switch role.Kind {
		case db.RoleKindPriority:
			mentions.priorityRoleID = role.ID
		case db.RoleKindCategory:
			mentions.categoryRoleID = role.ID
		case db.RoleKindNobody:
			mentions.nobody = true
		case db.RoleKindDiscussion:
			mentions.discussion = true // TODO:
		}
	}

	// do assignees
	assignees := []db.User{}
	if !mentions.nobody && len(captures.UserIDs) == 0 {
		assignees = append(assignees, db.User{ID: m.Author.ID})
	} else {
		for _, userID := range captures.UserIDs {
			assignees = append(assignees, db.User{ID: userID})
		}
	}

	// remove duplicate tags and normalize
	for i := range captures.Tags {
		captures.Tags[i] = strings.ToLower(captures.Tags[i])
	}
	slices.Sort(captures.Tags)
	captures.Tags = slices.Compact(captures.Tags)

	tags := []db.Tag{}
	for _, tagName := range captures.Tags {
		tags = append(tags, db.Tag{Name: tagName, ProjectID: project.ID})
	}

	// define the issue
	issue := db.Issue{
		Title:           title,
		Tags:            tags,
		Status:          db.IssueStatusTodo,
		ProjectID:       project.ID,
		Project:         *project,
		RecruiterUserID: m.Author.ID,
		AssigneeUsers:   assignees,
		CategoryRoleID:  mentions.categoryRoleID,
		PriorityRoleID:  mentions.priorityRoleID,
	}

	code, err := logic.GetIssueCode(&issue)
	if err != nil {
		return fmt.Errorf("error in issue db insertion: %w", err)
	}
	issue.Code = &code

	// create and initialize the thread
	thread, err := logic.CreateThreadFromIssue(&issue, s)
	if err != nil {
		return fmt.Errorf("error in thread creation: %w", err)
	}

	err = logic.InitIssueThread(&issue, db.RelationshipsByDirection{}, &guild, thread, s)
	if err != nil {
		return err
	}

	// add the issue to db
	err = db.Issues.Create(db.Ctx, &issue)
	if err != nil {
		return err
	}

	relationships := []db.Relationship{}
	for _, relCapture := range captures.DependsOnChannelIDs {
		// get the target issue from the capture channelID
		target, err := db.Issues.
			Preload("Tags", nil).
			Preload("AssigneeUsers", nil).
			Preload("Project", func(db gorm.PreloadBuilder) error {
				db.Select("ID", "Prefix", "GuildID")
				return nil
			}).
			Where("thread_id = ?", relCapture).
			First(db.Ctx)
		if err == gorm.ErrRecordNotFound {
			// if it's not an issue, ignore it and move on
			continue
		} else if err != nil {
			return err
		}

		// but if it is,
		// add the relationship
		relationship := db.Relationship{
			FromIssueID: issue.ID,
			ToIssueID:   target.ID,
			Kind:        db.RelationshipKindDependency,
		}
		err = db.Relationships.Create(db.Ctx, &relationship)
		if err != nil {
			return err
		}
		relationship.FromIssue = issue
		relationship.ToIssue = target
		relationships = append(relationships, relationship)

		// refresh target's view (in a goroutine)
		go func() {
			relationships, err := logic.GetIssueRelationshipsOfKind(&target, db.RelationshipKindDependency)
			if err != nil {
				slog.Error("error while getting relationship view after shorthand", "issue", issue.ID, "target", target.ID)
				return
			}

			err = logic.UpdateIssueThreadDetail(s, &target, relationships, guild.NobodyRoleID)
			if err != nil {
				slog.Error("error while updating target view after shorthand", "issue", issue.ID, "target", target.ID)
				return
			}
		}()
	}

	// if there are any relationships, update THIS issue's view
	if len(relationships) > 0 {
		err = logic.UpdateIssueThreadDetail(s, &issue, db.RelationshipsByDirection{Outbound: relationships}, guild.NobodyRoleID)
		if err != nil {
			return err
		}
	}

	// update all lists
	return logic.UpdateAllInteractiveIssuesViews(s, project.ID)
}
