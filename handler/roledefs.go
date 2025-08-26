package handler

import (
	"fmt"
	"issues/v2/slash"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type roleDef struct {
	Emoji string
	Key   string
	Color int
}

func (role roleDef) DisplayName() string {
	return fmt.Sprintf("%s %s", strings.ToUpper(role.Key), role.Emoji)
}

func (role roleDef) ToDiscordRoleParams() *discordgo.RoleParams {
	return &discordgo.RoleParams{
		Name:        role.DisplayName(),
		Color:       &role.Color,
		Mentionable: slash.Ptr(true),
	}
}

var categoryRoles = []roleDef{
	{
		Emoji: "🧻",
		Key:   "generic",
		Color: (0xfffffc),
	},
	{
		Emoji: "💎",
		Key:   "feat",
		Color: (0x00afb9),
	},
	{
		Emoji: "🐞",
		Key:   "fix",
		Color: (0xD63830),
	},
	{
		Emoji: "🧹",
		Key:   "chore",
		Color: (0xFF7F50),
	},
}
var priorityRoles = []roleDef{
	{
		Emoji: "⏬",
		Key:   "low",
		Color: (0x0077b6),
	},
	{
		Emoji: "📗",
		Key:   "normal",
		Color: (0x81B29A),
	},
	{
		Emoji: "⚠️",
		Key:   "important",
		Color: (0xffba08),
	},
	{
		Emoji: "🛑",
		Key:   "critical",
		Color: (0xd00000),
	},
}

var discussionRole = roleDef{
	Emoji: "💬",
	Key:   "discussion",
	Color: (0xCC4BC2),
}
var nobodyRole = roleDef{
	Emoji: "❔",
	Key:   "nobody",
	Color: (0xdcdcdc), // gainsboro
}
