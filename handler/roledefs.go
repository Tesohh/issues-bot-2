package handler

import (
	"issues/v2/slash"

	"github.com/bwmarrin/discordgo"
)

var categoryRoles = []*discordgo.RoleParams{
	{
		Name:  "GENERIC 🧻",
		Color: slash.Ptr(0xfffffc),
	},
	{
		Name:  "FEAT 💎",
		Color: slash.Ptr(0x00afb9),
	},
	{
		Name:  "FIX 🐞",
		Color: slash.Ptr(0xD63830),
	},
	{
		Name:  "CHORE 🧹",
		Color: slash.Ptr(0xFF7F50),
	},
}
var priorityRoles = []*discordgo.RoleParams{
	{
		Name:  "LOW ⏬",
		Color: slash.Ptr(0x0077b6),
	},
	{
		Name:  "NORMAL 📗",
		Color: slash.Ptr(0x81B29A),
	},
	{
		Name:  "IMPORTANT ⚠️",
		Color: slash.Ptr(0xffba08),
	},
	{
		Name:  "CRITICAL 🛑",
		Color: slash.Ptr(0xd00000),
	},
}

var discussionRole = &discordgo.RoleParams{
	Name:  "DISCUSSION 💬",
	Color: slash.Ptr(0xCC4BC2),
}
var nobodyRole = &discordgo.RoleParams{
	Name:  "NOBODY ❔",
	Color: slash.Ptr(0xdcdcdc), // gainsboro
}
