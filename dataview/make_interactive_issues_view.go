package dataview

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/helper"
	"slices"

	dg "github.com/bwmarrin/discordgo"
)

func makeStatusButton(state *db.ProjectViewState, dummy bool) dg.Button {
	label := ""
	style := dg.SecondaryButton
	statuses := ""
	if slices.Equal(state.Filter.Statuses, []db.IssueStatus{db.IssueStatusTodo, db.IssueStatusWorking}) {
		label = "Show closed"
		statuses = "done,killed"
	} else if slices.Equal(state.Filter.Statuses, []db.IssueStatus{db.IssueStatusDone, db.IssueStatusKilled}) {
		label = "Show open"
		statuses = "todo,working"
	} else {
		label = "Reset statuses"
		statuses = "todo,working"
	}
	return dg.Button{
		Label:    label,
		Style:    style,
		CustomID: fmt.Sprintf("issues-set-statuses:%s:%s", state.MessageID, statuses),
		Disabled: dummy,
	}
}

func makeSortByButton(state *db.ProjectViewState, dummy bool) dg.Button {
	label := ""
	style := dg.SecondaryButton
	sortby := ""

	switch state.Sorter.SortBy {
	case db.IssueSortByCode:
		label = "Sort by date"
		sortby = "date"
	case db.IssueSortByDate:
		label = "Sort by code"
		sortby = "code"
	default:
		label = "Reset sort by"
		sortby = "code"
	}

	return dg.Button{
		Label:    label,
		Style:    style,
		CustomID: fmt.Sprintf("issues-sort-by:%s:%s", state.MessageID, sortby),
		Disabled: dummy,
	}
}

func makeSortOrderButton(state *db.ProjectViewState, dummy bool) dg.Button {
	label := ""
	style := dg.SecondaryButton
	order := ""

	switch state.Sorter.SortOrder {
	case db.SortOrderAscending:
		label = "Order desc"
		order = "desc"
	case db.SortOrderDescending:
		label = "Order asc"
		order = "asc"
	default:
		label = "Reset order"
		order = "asc"
	}

	return dg.Button{
		Label:    label,
		Style:    style,
		CustomID: fmt.Sprintf("issues-order:%s:%s", state.MessageID, order),
		Disabled: dummy,
	}
}

// requires a message to have been sent BEFORE adding the buttons,
// as it depends on state.MessageID
func MakeInteractiveIssuesView(issues []db.Issue, state *db.ProjectViewState, dummy bool) []dg.MessageComponent {
	// define the buttons, with the message ids
	msgID := state.MessageID
	if dummy {
		msgID = "DUMMY"
	}

	prepared := state.Filter.Apply(issues)
	prepared = state.Sorter.Apply(prepared)
	pages := helper.Pages(issues, MaxIssuesPerPage)

	leftDisable := dummy || state.CurrentPage <= 0
	rightDisable := dummy || state.CurrentPage >= pages-1

	queryButtons := dg.ActionsRow{
		Components: []dg.MessageComponent{
			// TODO: check for page position
			makeStatusButton(state, dummy),
			makeSortByButton(state, dummy),
			makeSortOrderButton(state, dummy),
			dg.Button{Label: "Filters...", Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-filters:%s", msgID), Disabled: dummy},
		},
	}
	arrowButtons := dg.ActionsRow{
		Components: []dg.MessageComponent{
			dg.Button{Emoji: &dg.ComponentEmoji{Name: "⏮️"}, Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-goto:%s:%d:bigleft", msgID, 0), Disabled: leftDisable},
			dg.Button{Emoji: &dg.ComponentEmoji{Name: "⬅️"}, Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-goto:%s:%d:left", msgID, state.CurrentPage-1), Disabled: leftDisable},
			dg.Button{Emoji: &dg.ComponentEmoji{Name: "➡️"}, Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-goto:%s:%d:right", msgID, state.CurrentPage+1), Disabled: rightDisable},
			dg.Button{Emoji: &dg.ComponentEmoji{Name: "⏭️"}, Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-goto:%s:%d:bigright", msgID, pages-1), Disabled: rightDisable},
		},
	}

	prepared = helper.Paginate(prepared, MaxIssuesPerPage, state.CurrentPage)

	// generate the view
	view := MakeIssuesView(prepared, len(issues), state)

	return []dg.MessageComponent{queryButtons, view, arrowButtons}
}
