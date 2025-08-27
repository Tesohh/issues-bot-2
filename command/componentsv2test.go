package command

import (
	"issues/v2/dataview"
	"issues/v2/db"
	"issues/v2/slash"

	dg "github.com/bwmarrin/discordgo"
)

var ComponentsV2Test = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "components_v2_test",
		Description: "components_v2_test",
	},
	Func: func(s *dg.Session, i *dg.Interaction) error {
		container :=
			dg.Container{
				Components: []dg.MessageComponent{
					// Text
					dg.TextDisplay{Content: "👋 Hello from **Components V2**!"},

					dg.ActionsRow{
						Components: []dg.MessageComponent{
							dg.Button{
								CustomID: "btn_ping",
								Label:    "Ping",
								Style:    dg.PrimaryButton,
							},
							dg.Button{
								CustomID: "btn_danger",
								Label:    "Danger",
								Style:    dg.DangerButton,
							},
						},
					},

					dg.Section{
						Components: []dg.MessageComponent{
							dg.TextDisplay{Content: "⚪️ Example task #1"},
						},
						Accessory: dg.Button{
							Label:    "Select",
							Style:    dg.SecondaryButton,
							CustomID: "select1",
						},
					},
					dg.Section{
						Components: []dg.MessageComponent{
							dg.TextDisplay{Content: "⚪️ Example task #2"},
						},
						Accessory: dg.Button{
							Label:    "Select",
							Style:    dg.SecondaryButton,
							CustomID: "select2",
						},
					},
					dg.Section{
						Components: []dg.MessageComponent{
							dg.TextDisplay{Content: "⚪️ Example task #3"},
						},
						Accessory: dg.Button{
							Label:    "Select",
							Style:    dg.SecondaryButton,
							CustomID: "select3",
						},
					},
					dg.Section{
						Components: []dg.MessageComponent{
							dg.TextDisplay{Content: "⚪️ Example task #4"},
						},
						Accessory: dg.Button{
							Label:    "Select",
							Style:    dg.SecondaryButton,
							CustomID: "select4",
						},
					},

					// Separator
					dg.Separator{},

					// Media gallery
					dg.MediaGallery{
						Items: []dg.MediaGalleryItem{
							{
								Media: dg.UnfurledMediaItem{
									URL: "https://upload.wikimedia.org/wikipedia/commons/thumb/6/6a/Gatto_Siberiano_cuccioli_%28cropped%29.JPG/500px-Gatto_Siberiano_cuccioli_%28cropped%29.JPG",
								}},
						},
					},
				},
			}

		sampleIssues := []db.Issue{}
		for range 15 {
			sampleIssues = append(sampleIssues,
				db.Issue{
					ID:             3,
					Code:           slash.Ptr(uint(3)),
					Status:         db.IssueStatusTodo,
					Tags:           "gut, besser, gosser",
					Title:          "lorem ipsum dolor sit amet",
					CategoryRoleID: "1404946100275777556",
					PriorityRoleID: "1404946111998984263",
					Project:        db.Project{GuildID: "1404937966853427390"},
					ThreadID:       "1409194601516105808",
				})
		}

		view := dataview.MakeIssuesViewGithubStyle(sampleIssues, dataview.ProjectViewState{}, dataview.IssuesViewGithubStyleOptions{
			TitleOverride: "TESTLIST...",
		})
		buttons := dg.ActionsRow{
			Components: []dg.MessageComponent{
				// dg.Button{Label: "", Style: dg.SecondaryButton, CustomID: "show-my-issues"},
				dg.Button{Label: "Group by Hybrid", Style: dg.PrimaryButton, Emoji: &dg.ComponentEmoji{Name: "🔁"}, CustomID: "cycle-group-by"},
				dg.Button{Label: "Show completed", Style: dg.PrimaryButton, CustomID: "show-completed"},
				dg.Button{Label: "My issues", Style: dg.SuccessButton, CustomID: "show-my-issues"},
			},
		}
		return slash.ReplyWithComponents(s, i, false, container, view, buttons)
	},
}
