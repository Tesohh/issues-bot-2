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
	if err == gorm.ErrRecordNotFound {
		return nil // safely ignore the ErrRecordNotFound
	} else if err != nil {
		return err
	}
	_ = project

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

	// remove duplicate tags
	slices.Sort(captures.Tags)
	captures.Tags = slices.Compact(captures.Tags)
	tags := strings.Join(captures.Tags, ",")
	_ = tags // DELETEME:

	// define the issue
	issue := db.Issue{
		Title: title,
		// Tags:            tags,
		Status:          db.IssueStatusTodo,
		ProjectID:       project.ID,
		Project:         project,
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

	err = logic.InitIssueThread(&issue, &guild, thread, s)
	if err != nil {
		return err
	}

	// add the issue to db
	err = db.Issues.Create(db.Ctx, &issue)
	if err != nil {
		return err
	}

	// update all lists
	return logic.UpdateAllInteractiveIssuesViews(s, project.ID)
}
