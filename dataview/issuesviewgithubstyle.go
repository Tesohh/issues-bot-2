package dataview

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/helper"
	"issues/v2/slash"

	dg "github.com/bwmarrin/discordgo"
)

const MaxIssuesPerPage = 20
const MaxTitleLength = 70
const MaxTagsCount = 3
const MaxTagLength = 8

func MakeIssuesView(issues []db.Issue, totalIssueCount int, state *db.ProjectViewState) dg.Container {
	titleFmt := "# Issues in %s `[%s]`"
	if len(state.ListNameFmt) > 0 {
		titleFmt += state.ListNameFmt
	}

	title := fmt.Sprintf(titleFmt, state.Project.Name, state.Project.Prefix)
	subtitle := fmt.Sprintf("\n-# (%s, %s)", state.Filter, state.Sorter)

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
	}
	if len(content) == 0 {
		content = "There are no issues here. **Get to work!**"
	}

	pageText := fmt.Sprintf("\n-# page %d/%d", state.CurrentPage+1, (totalIssueCount/MaxIssuesPerPage)+1)
	components = append(components, dg.TextDisplay{Content: content}, dg.TextDisplay{Content: pageText})

	return slash.StandardizeContainer(dg.Container{Components: components})
}
