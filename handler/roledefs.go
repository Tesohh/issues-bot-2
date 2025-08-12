package handler

import (
	"issues/v2/slash"

	"github.com/bwmarrin/discordgo"
)

var categoryRoles = []*discordgo.RoleParams{
	{
		Name:  "GENERIC",
		Color: slash.Ptr(0xfffffc),
	},
	{
		Name:  "FEATURE",
		Color: slash.Ptr(0x00afb9),
	},
	{
		Name:  "FIX",
		Color: slash.Ptr(0xff8800),
	},
	{
		Name:  "OPTIMIZATION",
		Color: slash.Ptr(0x468C98),
	},
	{
		Name:  "REFACTOR",
		Color: slash.Ptr(0x7F2982),
	},
	{
		Name:  "CHORE",
		Color: slash.Ptr(0xda627d),
	},
}
var priorityRoles = []*discordgo.RoleParams{
	{
		Name:  "LOW",
		Color: slash.Ptr(0x0077b6),
	},
	{
		Name:  "NEXT VERSION",
		Color: slash.Ptr(0x9448BC),
	},
	{
		Name:  "RELEASE",
		Color: slash.Ptr(0x80DAEB),
	},
	{
		Name:  "NORMAL",
		Color: slash.Ptr(0x81B29A),
	},
	{
		Name:  "IMPORTANT",
		Color: slash.Ptr(0xffba08),
	},
	{
		Name:  "CRITICAL",
		Color: slash.Ptr(0xd00000),
	},
}

var discussionRole = &discordgo.RoleParams{
	Name:  "DISCUSSION",
	Color: slash.Ptr(0xCC4BC2),
}
var nobodyRole = &discordgo.RoleParams{
	Name:  "NOBODY",
	Color: slash.Ptr(0xdcdcdc), // gainsboro
}
