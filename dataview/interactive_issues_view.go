package dataview

import (
	"fmt"
	"issues/v2/db"

	dg "github.com/bwmarrin/discordgo"
)

// requires a message to have been sent BEFORE adding the buttons,
// as it depends on state.MessageID
func MakeInteractiveIssuesView(issues []db.Issue, state *db.ProjectViewState, options IssuesViewOptions, dummy bool) []dg.MessageComponent {
	// define the buttons, with the message ids
	msgID := state.MessageID
	if dummy {
		msgID = "DUMMY"
	}
	queryButtons := dg.ActionsRow{
		Components: []dg.MessageComponent{
			// TODO: check for page position
			dg.Button{Label: "Show closed", Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-show-closed:%s", msgID), Disabled: dummy},
			dg.Button{Label: "Sort by code", Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-sort-by:%s:code", msgID), Disabled: dummy}, // TODO:
			dg.Button{Label: "Order asc", Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-order:%s:asc", msgID), Disabled: dummy},       // TODO:
			dg.Button{Label: "Filters...", Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-filters:%s", msgID), Disabled: dummy},
			dg.Button{Label: "My issues", Style: dg.SuccessButton, CustomID: fmt.Sprintf("issues-show-mine:%s", msgID), Disabled: dummy},
		},
	}
	arrowButtons := dg.ActionsRow{
		Components: []dg.MessageComponent{
			dg.Button{Emoji: &dg.ComponentEmoji{Name: "⏮️"}, Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-bigleft:%s", msgID), Disabled: dummy},
			dg.Button{Emoji: &dg.ComponentEmoji{Name: "⬅️"}, Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-left:%s", msgID), Disabled: dummy},
			dg.Button{Emoji: &dg.ComponentEmoji{Name: "➡️"}, Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-right:%s", msgID), Disabled: dummy},
			dg.Button{Emoji: &dg.ComponentEmoji{Name: "⏭️"}, Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issues-bigright:%s", msgID), Disabled: dummy},
		},
	}

	// generate the view
	view := MakeIssuesView(issues, state, options)

	return []dg.MessageComponent{queryButtons, view, arrowButtons}
}
