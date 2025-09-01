package dataview

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/helper"
	"issues/v2/slash"

	dg "github.com/bwmarrin/discordgo"
)

type IssuesViewOptions struct {
	TitleOverride string // if set, will replace the default "Issues (filter)" title
}

const MaxIssuesPerPage = 20
const MaxTitleLength = 70
const MaxTagsCount = 3
const MaxTagLength = 8

func MakeIssuesView(issues []db.Issue, state *db.ProjectViewState, options IssuesViewOptions) dg.Container {
	title := fmt.Sprintf("# Issues in %s", state.Project.Name)
	if len(options.TitleOverride) > 0 {
		title = "# " + options.TitleOverride
	}

	subtitle := fmt.Sprintf("\n-# (%s, %s)", state.Filter, state.Sorter)

	issues = state.Filter.Apply(issues)
	issues = state.Sorter.Apply(issues)
	issues = helper.Paginate(issues, MaxIssuesPerPage, state.CurrentPage)

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
		tags := issue.PrettyTags(MaxTagsCount, MaxTagLength)
		line := fmt.Sprintf("\n - %s %s %s %s",
			issue.PrettyLink(longestCode),
			issue.RoleEmojis(),
			issue.CutTitle(MaxTitleLength-len(tags)),
			tags,
		)

		content += line
		//  TODO: trim text to make sure it is acertain len
		//  TODO: add all info (tags, priority and category emoji)
	}

	pageText := fmt.Sprintf("\n-# page %d/%d", state.CurrentPage+1, helper.Pages(issues, MaxIssuesPerPage))
	components = append(components, dg.TextDisplay{Content: content}, dg.TextDisplay{Content: pageText})

	return slash.StandardizeContainer(dg.Container{Components: components})
}
