package command

import (
	"fmt"
	"issues/v2/data"
	"issues/v2/db"
	"issues/v2/logic"
	"issues/v2/slash"
	"log/slog"
	"slices"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var codeOpt = dg.ApplicationCommandOption{
	Type:        dg.ApplicationCommandOptionInteger,
	Name:        "code",
	Description: "the code of the issue to edit. is inferred if you're in an issue thread",
	Required:    false,
	MinValue:    slash.Ptr(float64(0)),
	MaxValue:    0,
}

var Issue = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "issue",
		Description: "various issue editing commands",
		Options: []*dg.ApplicationCommandOption{
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "assign",
				Description: "assigns the issue to a user and removes them if they are already assigned",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionUser,
						Name:        "assignee",
						Description: "user to assign/remove",
						Required:    true,
					},
					&codeOpt,
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "category",
				Description: "changes the category of the issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "role",
						Description: "role to set as category",
						Choices:     data.CategoryOptionChoices,
						Required:    true,
					},
					&codeOpt,
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "priority",
				Description: "changes the priority of the issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "role",
						Description: "role to set as priority",
						Choices:     data.PriorityOptionChoices,
						Required:    true,
					},
					&codeOpt,
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "rename",
				Description: "changes the title of the issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "title",
						Description: "the new title",
						MinLength:   slash.Ptr(1),
						Required:    true,
					},
					&codeOpt,
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommandGroup,
				Name:        "mark",
				Description: "marks the issue as ...",
				Options: []*dg.ApplicationCommandOption{
					{Type: dg.ApplicationCommandOptionSubCommand, Name: "todo", Description: "🟩 todo", Options: []*dg.ApplicationCommandOption{&codeOpt}},
					{Type: dg.ApplicationCommandOptionSubCommand, Name: "doing", Description: "🟦 doing", Options: []*dg.ApplicationCommandOption{&codeOpt}},
					{Type: dg.ApplicationCommandOptionSubCommand, Name: "done", Description: "🟪 done", Options: []*dg.ApplicationCommandOption{&codeOpt}},
					{Type: dg.ApplicationCommandOptionSubCommand, Name: "cancelled", Description: "🟥 cancelled", Options: []*dg.ApplicationCommandOption{&codeOpt}},
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "tag",
				Description: "toggles a tag on the issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "tag",
						Description: "the tag to toggle",
						Required:    true,
					},
					&codeOpt,
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "tags",
				Description: "replaces tags with the list of tags provided",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "tags",
						Description: "the comma separated tags to replace. Input a single , to delete all tags",
						Required:    true,
					},
					&codeOpt,
				},
			},
		},
	},
	Func: func(s *dg.Session, i *dg.Interaction) error {
		subcommand := i.ApplicationCommandData().Options[0]
		options := slash.GetOptionMapRaw(subcommand.Options)

		query := db.Issues.
			Preload("Tags", nil).
			Preload("AssigneeUsers", nil).
			Preload("Project", func(db gorm.PreloadBuilder) error {
				db.Select("ID", "Prefix")
				return nil
			})

		var codeOpt *dg.ApplicationCommandInteractionDataOption

		if opt, ok := options["code"]; ok {
			codeOpt = opt
		} else if opt := subcommand.Options[0].GetOption("code"); opt != nil {
			codeOpt = opt
		}

		if codeOpt != nil {
			code := codeOpt.IntValue()
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return err
			}
			project, err := db.Projects.
				Select("id").
				Where("discord_category_channel_id = ?", channel.ParentID).
				Or("issues_input_channel_id = ?", channel.ParentID).
				First(db.Ctx)
			if err != nil {
				return err
			}

			query = query.
				Where("code = ?", code).
				Where("project_id = ?", project.ID)
		} else {
			query = query.Where("thread_id = ?", i.ChannelID)
		}

		issue, err := query.First(db.Ctx)
		if err == gorm.ErrRecordNotFound {
			return ErrNotInIssueThread
		}

		switch subcommand.Name {
		case "assign":
			assignee := options["assignee"].UserValue(nil)
			err = IssueAssign(s, i, &issue, assignee)
		case "category", "priority":
			key := options["role"].StringValue()
			role, err := db.Roles.
				Select("id, kind").
				Where("key = ?", key).
				Where("guild_id = ?", i.GuildID).
				First(db.Ctx)
			if err != nil {
				return err
			}
			err = IssueCategoryOrPriority(s, i, &issue, &role, subcommand.Name)
		case "rename":
			title := options["title"].StringValue()
			err = IssueRename(s, i, &issue, title)
		case "mark":
			arg := subcommand.Options[0].Name
			err = IssueMark(s, i, &issue, arg)
		case "tag":
			tag := options["tag"].StringValue()
			err = IssueTag(s, i, &issue, tag)
		case "tags":
			tags := options["tags"].StringValue()
			err = IssueTags(s, i, &issue, tags)
		}

		if err != nil {
			return err
		}

		guild, err := db.Guilds.Select("nobody_role_id").Where("id = ?", i.GuildID).First(db.Ctx)
		if err != nil {
			return err
		}
		err = logic.UpdateIssueThreadDetail(s, &issue, guild.NobodyRoleID)
		if err != nil {
			return err
		}

		err = logic.UpdateAllInteractiveIssuesViews(s, issue.ProjectID)
		if err != nil {
			return err
		}

		return nil
	},
}

func IssueAssign(s *dg.Session, i *dg.Interaction, issue *db.Issue, assignee *dg.User) error {
	index := slices.IndexFunc(issue.AssigneeUsers, func(user db.User) bool {
		return user.ID == assignee.ID
	})

	msgFmt := ""
	if index == -1 {
		issue.AssigneeUsers = append(issue.AssigneeUsers, db.User{ID: assignee.ID})
		err := db.Conn.Table("issue_assignees").
			Create(map[string]any{
				"issue_id": issue.ID,
				"user_id":  assignee.ID,
			}).Error
		if err != nil {
			return err
		}

		msgFmt = "<@%s> added <@%s> to assignees"
	} else {
		issue.AssigneeUsers = slices.Delete(issue.AssigneeUsers, index, index+1)
		err := db.Conn.Table("issue_assignees").
			Where("issue_id = ?", issue.ID).
			Where("user_id = ?", assignee.ID).
			Delete(map[string]any{}).Error
		if err != nil {
			return err
		}
		msgFmt = "<@%s> removed <@%s> from assignees"
	}

	msg := fmt.Sprintf(msgFmt, i.Member.User.ID, assignee.ID)
	return slash.ReplyWithText(s, i, msg, false)
}

func IssueCategoryOrPriority(s *dg.Session, i *dg.Interaction, issue *db.Issue, role *db.Role, subcommand string) error {
	switch subcommand {
	case "priority":
		if role.Kind != db.RoleKindPriority {
			return fmt.Errorf("%w (expected priority, got %s)", ErrWrongRole, role.Kind)
		}
		if issue.PriorityRoleID == role.ID {
			msg := fmt.Sprintf("Priority was already <@&%s>, no actions taken.", role.ID)
			return slash.ReplyWithText(s, i, msg, true)
		}
		issue.PriorityRoleID = role.ID
		_, err := db.Issues.Where("id = ?", issue.ID).Update(db.Ctx, "priority_role_id", role.ID)
		if err != nil {
			return err
		}
	case "category":
		if role.Kind != db.RoleKindCategory {
			return fmt.Errorf("%w (expected category, got %s)", ErrWrongRole, role.Kind)
		}
		if issue.CategoryRoleID == role.ID {
			msg := fmt.Sprintf("Category was already <@&%s>, no actions taken.", role.ID)
			return slash.ReplyWithText(s, i, msg, true)
		}
		issue.CategoryRoleID = role.ID
		_, err := db.Issues.Where("id = ?", issue.ID).Update(db.Ctx, "category_role_id", role.ID)
		if err != nil {
			return err
		}
	}

	msg := fmt.Sprintf("<@%s> updated %s to <@&%s>", i.Member.User.ID, subcommand, role.ID)
	return slash.ReplyWithText(s, i, msg, false)
}

func IssueRename(s *dg.Session, i *dg.Interaction, issue *db.Issue, title string) error {
	if issue.Title == title {
		msg := "Title was already that, no actions taken."
		return slash.ReplyWithText(s, i, msg, true)
	}

	issue.Title = title
	_, err := db.Issues.Where("id = ?", issue.ID).Update(db.Ctx, "title", title)
	if err != nil {
		return err
	}

	_, err = s.ChannelEdit(issue.ThreadID, &dg.ChannelEdit{Name: issue.ChannelName()})
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("<@%s> updated the title to \"%s\"", i.Member.User.ID, title)
	return slash.ReplyWithText(s, i, msg, false)
}

var marksPerIssue = map[uint]int{}

func IssueMark(s *dg.Session, i *dg.Interaction, issue *db.Issue, subcommand string) error {
	var issueStatus db.IssueStatus
	var archive = false
	var lock = false
	var autoArchiveDuration = 10080
	switch subcommand {
	case "todo":
		issueStatus = db.IssueStatusTodo
	case "doing":
		issueStatus = db.IssueStatusDoing
	case "done":
		issueStatus = db.IssueStatusDone
		archive = true
		autoArchiveDuration = 60
	case "cancelled":
		issueStatus = db.IssueStatusCancelled
		archive = true
		lock = true
		autoArchiveDuration = 60
	}

	if issue.Status == issueStatus {
		msg := fmt.Sprintf("Status was already %s, no actions taken.", subcommand)
		return slash.ReplyWithText(s, i, msg, true)
	}

	issue.Status = issueStatus
	_, err := db.Issues.Where("id = ?", issue.ID).Update(db.Ctx, "status", issueStatus)
	if err != nil {
		return err
	}

	go func() {
		_, err = s.ChannelEdit(issue.ThreadID, &dg.ChannelEdit{
			Name:                issue.ChannelName(),
			AutoArchiveDuration: autoArchiveDuration,
			Locked:              &lock,
		})
		if err != nil {
			slog.Error("error while trying to edit thread.", "err", err)
			return
		}
	}()

	// and finally send the embed
	var alsoWillArchiveString string
	if archive {
		alsoWillArchiveString = ", thread will be archived in 1 hour if inactive"
	}

	if _, ok := marksPerIssue[issue.ID]; !ok {
		marksPerIssue[issue.ID] = 0
	}
	marksPerIssue[issue.ID] += 1

	var warnString string
	if marksPerIssue[issue.ID] >= 3 {
		warnString = "-# ⚠️ due to discord ratelimits the channel name might not update for a long time"
	}

	msg := fmt.Sprintf("<@%s> marked the issue as %s %s%s\n%s", i.Member.User.ID, db.IssueStatusIcons[issueStatus], db.IssueStatusNames[issueStatus], alsoWillArchiveString, warnString)
	return slash.ReplyWithText(s, i, msg, false)
}

func IssueTag(s *dg.Session, i *dg.Interaction, issue *db.Issue, name string) error {
	name = strings.Trim(name, "+ ")
	name = strings.ToLower(name)

	index := slices.IndexFunc(issue.Tags, func(ltag db.Tag) bool {
		return ltag.Name == name
	})

	msgFmt := ""
	if index == -1 { // doesn't exist, create it
		tag := db.Tag{
			Name:      name,
			ProjectID: issue.ProjectID,
		}

		err := db.Conn.Model(issue).Association("Tags").Append(&tag)
		if err != nil {
			return err
		}

		msgFmt = "<@%s> added tag `+%s`"
	} else {
		err := db.Conn.
			Model(issue).
			Association("Tags").
			Delete(&db.Tag{Name: name, ProjectID: issue.ProjectID})
		if err != nil {
			return err
		}

		msgFmt = "<@%s> removed tag `+%s`"
	}

	msg := fmt.Sprintf(msgFmt, i.Member.User.ID, name)
	return slash.ReplyWithText(s, i, msg, false)
}

func IssueTags(s *dg.Session, i *dg.Interaction, issue *db.Issue, tagsRaw string) error {
	// parse and remove duplicates
	tagNames := db.ParseTags(tagsRaw)
	for i := range tagNames {
		tagNames[i] = strings.ToLower(tagNames[i])
	}
	slices.Sort(tagNames)
	tagNames = slices.Compact(tagNames)

	tags := []db.Tag{}
	for _, tagName := range tagNames {
		tags = append(tags, db.Tag{Name: tagName, ProjectID: issue.ProjectID})
	}

	err := db.Conn.
		Model(issue).
		Association("Tags").
		Replace(&tags)
	if err != nil {
		return err
	}

	prettyTags := ""
	if len(tags) == 0 {
		prettyTags = "[no tags]"
	} else {
		prettyTags = issue.PrettyTags(999, 999)
	}
	msg := fmt.Sprintf("<@%s> replaced tags with %s", i.Member.User.ID, prettyTags)
	return slash.ReplyWithText(s, i, msg, false)
}
