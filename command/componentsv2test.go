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
		_ = container

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
					CategoryRole:   db.Role{Emoji: "💎"},
					PriorityRoleID: "1404946111998984263",
					PriorityRole:   db.Role{Emoji: "⚠️"},
					Project:        db.Project{GuildID: "1404937966853427390"},
					ThreadID:       "1409194601516105808",
				})
		}

		view := dataview.MakeIssuesViewGithubStyle(sampleIssues, dataview.ProjectViewState{}, dataview.IssuesViewGithubStyleOptions{
			TitleOverride: "AutoList™️ for LOREM",
		})
		arrowbuttons := dg.ActionsRow{
			Components: []dg.MessageComponent{
				dg.Button{Emoji: &dg.ComponentEmoji{Name: "⏮️"}, Style: dg.SecondaryButton, Disabled: true, CustomID: "bigleft"},
				dg.Button{Emoji: &dg.ComponentEmoji{Name: "⬅️"}, Style: dg.SecondaryButton, CustomID: "left"},
				dg.Button{Emoji: &dg.ComponentEmoji{Name: "➡️"}, Style: dg.SecondaryButton, CustomID: "right"},
				dg.Button{Emoji: &dg.ComponentEmoji{Name: "⏭️"}, Style: dg.SecondaryButton, CustomID: "bigright"},
				dg.Button{Label: "My issues", Style: dg.SuccessButton, CustomID: "show-my-issues"},
			},
		}
		queryButtons := dg.ActionsRow{
			Components: []dg.MessageComponent{
				dg.Button{Label: "Show closed", Style: dg.SecondaryButton, CustomID: "showclosed"},
				dg.Button{Label: "Sort by code", Style: dg.SecondaryButton, CustomID: "sort-by:code"},
				dg.Button{Label: "Order asc", Style: dg.SecondaryButton, CustomID: "order:asc"},
				dg.Button{Label: "Filters...", Style: dg.SecondaryButton, CustomID: "filters"},
				dg.Button{Label: "Query...", Style: dg.SecondaryButton, CustomID: "query"},
			},
		}
		// url.ParseQuery("issuelistgoto?message=123123&page=2")
		return slash.ReplyWithComponents(s, i, false, queryButtons, view, arrowbuttons)
	},
}
