package dataview

import (
	"fmt"
	"issues/v2/db"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

type GroupIssuesBy uint

const (
	GroupIssuesByStatus   GroupIssuesBy = 0
	GroupIssuesByCategory GroupIssuesBy = 1
	GroupIssuesByHybrid   GroupIssuesBy = 2 // groups by category, but separates issues in Working status
)

type IssuesViewOptions struct {
	TitleOverride         string // if set, will replace the default "Issues (filter)" title
	DefaultPriorityRoleID string // if set, issues with this priority will not show the role
	GroupIssuesBy         GroupIssuesBy
}

type issuesGroup struct {
	title  string
	issues []*db.Issue
}

func MakeIssuesView(unfilteredIssues []db.Issue, filter IssueFilter, options IssuesViewOptions) dg.Container {
	title := fmt.Sprintf("# Issues %s", filter.String())
	if len(options.TitleOverride) > 0 {
		title = options.TitleOverride
	}

	issues := filter.Apply(unfilteredIssues)

	groups := map[any]*issuesGroup{}
	var groupAccessor func(issue db.Issue) any
	var groupTitleFormatter func(issue db.Issue) string
	var groupsOrderMaker func(groups map[any]*issuesGroup) []any

	switch options.GroupIssuesBy {
	case GroupIssuesByStatus:
		groupAccessor = func(issue db.Issue) any { return issue.Status }
		groupTitleFormatter = func(issue db.Issue) string { return db.IssueStatusNames[issue.Status] }
		groupsOrderMaker = func(groups map[any]*issuesGroup) []any {
			return []any{db.IssueStatusTodo, db.IssueStatusWorking, db.IssueStatusDone, db.IssueStatusKilled}
		}
	case GroupIssuesByCategory:
		groupAccessor = func(issue db.Issue) any { return issue.CategoryRoleID }
		groupTitleFormatter = func(issue db.Issue) string { return fmt.Sprintf("<@&%s>", issue.CategoryRoleID) }
		groupsOrderMaker = func(groups map[any]*issuesGroup) []any {
			order := []any{options.DefaultPriorityRoleID}
			for k := range groups {
				if k != options.DefaultPriorityRoleID {
					order = append(order, k)
				}
			}
			return order
		}

	case GroupIssuesByHybrid:
		groupAccessor = func(issue db.Issue) any {
			if issue.Status == db.IssueStatusWorking {
				return "Working"
			} else {
				return issue.CategoryRoleID
			}
		}
		groupTitleFormatter = func(issue db.Issue) string {
			if issue.Status == db.IssueStatusWorking {
				return "Working"
			} else {
				return fmt.Sprintf("<@&%s>", issue.CategoryRoleID)
			}
		}
		groupsOrderMaker = func(groups map[any]*issuesGroup) []any {
			order := []any{"Working"}
			for k := range groups {
				if k != "Working" {
					order = append(order, k)
				}
			}
			return order
		}

	default:
		panic("trying to group by something that isnt Status, Category or Hybrid (should be unreachable)")
	}

	for _, issue := range issues {
		group, ok := groups[groupAccessor(issue)]
		if !ok {
			groups[groupAccessor(issue)] = &issuesGroup{title: groupTitleFormatter(issue), issues: []*db.Issue{}}
			group = groups[groupAccessor(issue)]
		}
		group.issues = append(group.issues, &issue)
	}

	components := []dg.MessageComponent{
		dg.TextDisplay{Content: title},
	}

	groupsOrder := groupsOrderMaker(groups)
	for _, key := range groupsOrder {
		group, ok := groups[key]
		if !ok {
			continue
		}

		content := fmt.Sprintf("### %s", group.title)
		for _, issue := range group.issues {
			priority := ""
			if issue.PriorityRoleID != options.DefaultPriorityRoleID {
				priority = fmt.Sprintf("<@&%s>", issue.PriorityRoleID)
			}
			content += fmt.Sprintf("\n - %s %s", issue.PrettyLink(), priority)
		}
		slog.Debug("dataview content length", "len(content)", len(content), "len(group.Issues)", len(group.issues))
		components = append(components, dg.TextDisplay{Content: content})
		components = append(components, dg.Separator{})
	}

	// remove last separator
	if components[len(components)-1].Type() == dg.SeparatorComponent {
		components = components[:len(components)-1]
	}

	return dg.Container{Components: components}
}
