package dataview

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/helper"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

type IssuesViewGithubStyleOptions struct {
	TitleOverride string // if set, will replace the default "Issues (filter)" title
}

const MaxIssuesPerPage = 20

func MakeIssuesViewGithubStyle(issues []db.Issue, state ProjectViewState, options IssuesViewGithubStyleOptions) dg.Container {
	title := "# Issues"
	if len(options.TitleOverride) > 0 {
		title = "# " + options.TitleOverride
	}

	subtitle := fmt.Sprintf("\n-# (%s, %s)", state.Filter, state.Sorter)

	issues = state.Filter.Apply(issues)
	issues = state.Sorter.Apply(issues)

	components := []dg.MessageComponent{
		dg.TextDisplay{Content: title + subtitle},
	}

	longestCode := 0
	for _, issue := range issues {
		length := helper.DigitsLen(int(*issue.Code))
		if length > longestCode {
			longestCode = length
		}
	}

	content := ""
	for _, issue := range issues {
		line := fmt.Sprintf("\n - %s %s %s %s", issue.PrettyLink(longestCode), issue.RoleEmojis(), issue.CutTitle(30), issue.PrettyTags(3, 7))
		content += line
		slog.Debug("length of line", "len", len(line))
		//  TODO: trim text to make sure it is acertain len
		//  TODO: add all info (tags, priority and category emoji)
	}
	components = append(components, dg.TextDisplay{Content: content})

	return dg.Container{Components: components}
}
