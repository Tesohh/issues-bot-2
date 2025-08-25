package command

import (
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
					dg.Separator{
						Divider: slash.Ptr(true),
						Spacing: slash.Ptr(dg.SeparatorSpacingSizeSmall),
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
		return slash.ReplyWithComponents(s, i, false, container)
	},
}
