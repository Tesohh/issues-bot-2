package dataview

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/helper"
	"issues/v2/slash"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

const MaxIssuesPerPage = 20
const MaxTitleLength = 70
const MaxTagsCount = 3
const MaxTagLength = 8

func MakeIssuesView(issues []db.Issue, totalIssueCount int, state *db.ProjectViewState) dg.Container {
	title := strings.Replace("# "+state.ListNameFmt, "$n", state.Project.Name, 1)
	title = strings.Replace(title, "$p", strings.ToUpper(state.Project.Prefix), 1)

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

	permanentStr := ""
	if state.Permanent {
		permanentStr = "♾️ "
	}
	pageText := fmt.Sprintf("\n-# %spage %d/%d", permanentStr, state.CurrentPage+1, (totalIssueCount/MaxIssuesPerPage)+1)
	if state.DeletedAt.Valid {
		pageText += "\n-# ⚠️WARNING\n-# this list has been purged and cannot be interacted with\n-# kindly delete this message!\n-# if you wish to make lists permanent, use the `permanent` flag in `/list issues` next time"
	}
	components = append(components, dg.TextDisplay{Content: content}, dg.TextDisplay{Content: pageText})

	return slash.StandardizeContainer(dg.Container{Components: components})
}
