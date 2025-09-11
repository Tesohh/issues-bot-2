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

var New = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "new",
		Description: "Adds a new issue if you're in a #xxx-issues channel, or adds a task if you're in a issue thread",
		Options: []*dg.ApplicationCommandOption{
			{
				Name:        "title",
				Description: "the title of the issue",
				Required:    true,
				Type:        dg.ApplicationCommandOptionString,
			},
			{
				Name:        "category",
				Description: "what's the category of this issue (ie. FEATURE, FIX, OPTIMIZATION, REFACTOR, CHORE)",
				Type:        dg.ApplicationCommandOptionString,
				Choices:     data.CategoryOptionChoices,
			},
			{
				Name:        "priority",
				Description: "what's the priority of this issue (ie. LOW, NEXT VERSION, RELEASE, IMPORTANT, CRITICAL)",
				Type:        dg.ApplicationCommandOptionString,
				Choices:     data.PriorityOptionChoices,
			},
			{
				Name:        "tags",
				Description: "comma separated tags (eg. `tag1, tag2`)",
				Type:        dg.ApplicationCommandOptionString,
			},
			{
				Name:        "assign",
				Description: "select user to assign this to. if you want to assign this to nobody, use the `nobody` flag",
				Type:        dg.ApplicationCommandOptionUser,
			},
			{
				Name:        "discussion",
				Description: "is this a discussion?",
				Type:        dg.ApplicationCommandOptionBoolean,
			},
			{
				Name:        "nobody",
				Description: "set this to true to assign the issue to nobody",
				Type:        dg.ApplicationCommandOptionBoolean,
			},
		},
	},
	Disabled: false,
	Func: func(s *dg.Session, i *dg.Interaction) error {
		opts := slash.GetOptionMap(i)

		project, err := db.Projects.Where("issues_input_channel_id = ?", i.ChannelID).First(db.Ctx)
		if err == gorm.ErrRecordNotFound {
			return ErrNotInIssueInputChannel
		} else if err != nil {
			return err
		}

		guild, err := db.Guilds.Where("id = ?", i.GuildID).First(db.Ctx)
		if err != nil {
			return err
		}

		title := opts["title"].StringValue()
		tags := ""
		if tagsOpt, ok := opts["tags"]; ok {
			tagsSplit := strings.Split(tagsOpt.StringValue(), ",")

			// remove duplicate tags
			slices.Sort(tagsSplit)
			tagsSplit = slices.Compact(tagsSplit)

			for i := range tagsSplit {
				tagsSplit[i] = strings.Trim(tagsSplit[i], " +")
			}
			tags = strings.Join(tagsSplit, ",")
		}

		kind := db.IssueKindNormal
		if discussionOpt, ok := opts["discussion"]; ok {
			if discussionOpt.BoolValue() {
				kind = db.IssueKindDiscussion
			}
		}

		categoryRoleID := guild.GenericCategoryRoleID
		if categoryRoleOpt, ok := opts["category"]; ok {
			key := categoryRoleOpt.StringValue()
			role, err := db.Roles.
				Select("id").
				Where("key = ?", key).
				Where("guild_id = ?", i.GuildID).
				Where("kind = ?", db.RoleKindCategory).
				First(db.Ctx)
			if err == nil {
				categoryRoleID = role.ID
			}
		}

		priorityRoleID := guild.NormalPriorityRoleID
		if priorityRoleOpt, ok := opts["priority"]; ok {
			key := priorityRoleOpt.StringValue()
			role, err := db.Roles.
				Select("id").
				Where("key = ?", key).
				Where("guild_id = ?", i.GuildID).
				Where("kind = ?", db.RoleKindPriority).
				First(db.Ctx)
			if err == nil {
				priorityRoleID = role.ID
			}
		}

		assignees := []db.User{}
		if nobodyOpt, ok := opts["nobody"]; !ok || (ok && !nobodyOpt.BoolValue()) { // if "nobody" is not defined or false
			assigneeID := i.Member.User.ID
			if assigneeOpt, ok := opts["assign"]; ok {
				assigneeID = assigneeOpt.Value.(string)
			}
			assignees = append(assignees, db.User{ID: assigneeID})
		}

		issue := &db.Issue{
			Title:           title,
			Tags:            tags,
			Kind:            kind,
			ProjectID:       project.ID,
			RecruiterUserID: i.Member.User.ID,
			AssigneeUsers:   assignees,
			CategoryRoleID:  categoryRoleID,
			PriorityRoleID:  priorityRoleID,
		}

		code, err := logic.GetIssueCode(issue)
		if err != nil {
			return fmt.Errorf("error in issue db insertion: %w", err)
		}
		issue.Code = &code

		slash.ReplyWithEmbed(s, i, dg.MessageEmbed{
			Title: title,
		}, true)
		err = s.InteractionResponseDelete(i)
		if err != nil {
			slog.Warn("couldn't delete ack message. no big deal", "err", err)
		}

		issue.Project = project
		thread, err := logic.CreateThreadFromIssue(issue, s)
		if err != nil {
			return fmt.Errorf("error in thread creation: %w", err)
		}

		err = logic.InitIssueThread(issue, &guild, thread, s)
		if err != nil {
			return err
		}

		err = db.Issues.Create(db.Ctx, issue)
		if err != nil {
			return err
		}

		return logic.UpdateAllInteractiveIssuesViews(s, project.ID)
	},
}
