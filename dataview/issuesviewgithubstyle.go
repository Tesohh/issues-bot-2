package dataview

import (
	"fmt"
	"issues/v2/db"

	dg "github.com/bwmarrin/discordgo"
)

type IssuesViewGithubStyleOptions struct {
	TitleOverride string // if set, will replace the default "Issues (filter)" title
}

const MaxIssuesPerPage = 15

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

	content := ""
	for _, issue := range issues {
		content += fmt.Sprintf("\n - %s %s %s", issue.PrettyLink(), issue.CutTitle(3), issue.PrettyTags(3, 3))
		//  TODO: trim text to make sure it is acertain len
		//  TODO: add all info (tags, priority and category emoji)
	}
	components = append(components, dg.TextDisplay{Content: content})

	return dg.Container{Components: components}
}
