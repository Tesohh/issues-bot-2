package command

import (
	"fmt"
	"issues/v2/dataview"
	"issues/v2/db"
	"issues/v2/slash"
	"math/rand/v2"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
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
		for i := range dataview.MaxIssuesPerPage {
			sampleIssues = append(sampleIssues,
				db.Issue{
					ID:     uint(i + 1),
					Code:   slash.Ptr(uint(i + 1)),
					Status: db.IssueStatus(rand.Int32N(4)),
					Tags: strings.ToLower(
						gofakeit.RandomString([]string{
							fmt.Sprintf("%s, %s, %s, %s, %s", gofakeit.HackerNoun(), gofakeit.HackerNoun(), gofakeit.HackerNoun(), gofakeit.HackerNoun(), gofakeit.HackerNoun()),
							fmt.Sprintf("%s, %s, %s, %s", gofakeit.HackerNoun(), gofakeit.HackerNoun(), gofakeit.HackerNoun(), gofakeit.HackerNoun()),
							fmt.Sprintf("%s, %s, %s", gofakeit.HackerNoun(), gofakeit.HackerNoun(), gofakeit.HackerNoun()),
							fmt.Sprintf("%s, %s", gofakeit.AppVersion(), gofakeit.HackerNoun()),
							fmt.Sprintf("%s, %s", gofakeit.HackerNoun(), gofakeit.HackerNoun()),
							fmt.Sprintf("%s", gofakeit.HackerNoun()),
							fmt.Sprintf("%s", gofakeit.HackerAbbreviation()),
							fmt.Sprintf("%s", gofakeit.AppVersion()),
							"",
						}),
					),
					Title:          fmt.Sprintf("%s the %s %s and the %s", gofakeit.HackerVerb(), gofakeit.HackerAdjective(), gofakeit.HackerNoun(), gofakeit.HackerNoun()),
					CategoryRoleID: "1404946100275777556",
					CategoryRole:   db.Role{Emoji: gofakeit.RandomString([]string{"🧻", "💎", "🐞", "🧹"})},
					PriorityRoleID: "1404946111998984263",
					PriorityRole:   db.Role{Emoji: gofakeit.RandomString([]string{"⏬", "📗", "⚠️", "🛑"})},
					Project:        db.Project{GuildID: "1404937966853427390"},
					ThreadID:       "1409194601516105808",
				})
		}

		db.ProjectViewStates.Create(db.Ctx, &db.ProjectViewState{
			MessageID:   gofakeit.AdverbDegree() + gofakeit.Email() + gofakeit.IPv4Address(),
			ProjectID:   1,
			CurrentPage: 0,
			Filter: db.IssueFilter{
				Statuses:        []db.IssueStatus{db.IssueStatusDone},
				Title:           "the",
				Tags:            []string{"gut"},
				PriorityRoleIDs: []string{"234234"},
				CategoryRoleIDs: []string{"923482398", "23498728937"},
				AssigneeIDs:     []string{"23489237894"},
			},
			Sorter: db.IssueSorter{
				SortBy:    db.IssueSortByDate,
				SortOrder: db.SortOrderAscending,
			},
		})

		view := dataview.MakeIssuesView(sampleIssues, len(sampleIssues), &db.ProjectViewState{})
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
				dg.Button{Label: "Show closed", Style: dg.SecondaryButton, CustomID: "ping:harris"},
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
