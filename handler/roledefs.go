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

var GenericCategoryRole = roleDef{
	Emoji: "🧻",
	Key:   "generic",
	Color: (0xfffffc),
}

var FeatureCategoryRole = roleDef{
	Emoji: "💎",
	Key:   "feat",
	Color: (0x00afb9),
}

var FixCategoryRole = roleDef{
	Emoji: "🐞",
	Key:   "fix",
	Color: (0xD63830),
}

var ChoreCategoryRole = roleDef{
	Emoji: "🧹",
	Key:   "chore",
	Color: (0xFF7F50),
}

var CategoryRoles = []roleDef{
	GenericCategoryRole, FeatureCategoryRole, FixCategoryRole, ChoreCategoryRole,
}

var LowPriorityRole = roleDef{
	Emoji: "⏬",
	Key:   "low",
	Color: (0x0077b6),
}
var NormalPriorityRole = roleDef{
	Emoji: "📗",
	Key:   "normal",
	Color: (0x81B29A),
}
var ImportantPriorityRole = roleDef{
	Emoji: "⚠️",
	Key:   "important",
	Color: (0xffba08),
}
var CriticalPriorityRole = roleDef{
	Emoji: "🛑", //‼️
	Key:   "critical",
	Color: (0xd00000),
}

var PriorityRoles = []roleDef{
	LowPriorityRole, NormalPriorityRole, ImportantPriorityRole, CriticalPriorityRole,
}

var DiscussionRole = roleDef{
	Emoji: "💬",
	Key:   "discussion",
	Color: (0xCC4BC2),
}
var NobodyRole = roleDef{
	Emoji: "❔",
	Key:   "nobody",
	Color: (0xdcdcdc), // gainsboro
}
