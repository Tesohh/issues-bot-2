package data

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/slash"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

type RoleDef struct {
	Emoji string
	Key   string
	Color int
}

func (role RoleDef) DisplayName() string {
	return fmt.Sprintf("%s %s", strings.ToUpper(role.Key), role.Emoji)
}

func (role RoleDef) AsChoice() *dg.ApplicationCommandOptionChoice {
	return &dg.ApplicationCommandOptionChoice{
		Name:  role.DisplayName(),
		Value: role.Key,
	}
}

func (role RoleDef) AsSelectChoice() dg.SelectMenuOption {
	return dg.SelectMenuOption{
		Label: strings.ToUpper(role.Key),
		Emoji: &dg.ComponentEmoji{Name: role.Emoji},
		Value: role.Key,
	}
}

func (role RoleDef) ToDiscordRoleParams() *dg.RoleParams {
	return &dg.RoleParams{
		Name:        role.DisplayName(),
		Color:       &role.Color,
		Mentionable: slash.Ptr(true),
	}
}

var GenericCategoryRole = RoleDef{Emoji: "🧻", Key: "generic", Color: (0xfffffc)}
var FeatureCategoryRole = RoleDef{Emoji: "💎", Key: "feat", Color: (0x00afb9)}
var FixCategoryRole = RoleDef{Emoji: "🐞", Key: "fix", Color: (0xD63830)}
var ChoreCategoryRole = RoleDef{Emoji: "🧹", Key: "chore", Color: (0xFF7F50)}

var CategoryRoles = []RoleDef{
	GenericCategoryRole, FeatureCategoryRole, FixCategoryRole, ChoreCategoryRole,
}

var LowPriorityRole = RoleDef{Emoji: "⏬", Key: "low", Color: (0x0077b6)}
var NormalPriorityRole = RoleDef{Emoji: "📗", Key: "normal", Color: (0x81B29A)}
var ImportantPriorityRole = RoleDef{Emoji: "⚠️", Key: "important", Color: (0xffba08)}
var CriticalPriorityRole = RoleDef{Emoji: "🛑", Key: "critical", Color: (0xd00000)}

var PriorityRoles = []RoleDef{
	LowPriorityRole, NormalPriorityRole, ImportantPriorityRole, CriticalPriorityRole,
}

var DiscussionRole = RoleDef{Emoji: "💬", Key: "discussion", Color: (0xCC4BC2)}
var NobodyRole = RoleDef{Emoji: "❔", Key: "nobody", Color: (0xdcdcdc)} // gainsboro

var CategoryOptionChoices = []*dg.ApplicationCommandOptionChoice{
	GenericCategoryRole.AsChoice(),
	FeatureCategoryRole.AsChoice(),
	FixCategoryRole.AsChoice(),
	ChoreCategoryRole.AsChoice(),
}

var PriorityOptionChoices = []*dg.ApplicationCommandOptionChoice{
	CriticalPriorityRole.AsChoice(),
	ImportantPriorityRole.AsChoice(),
	NormalPriorityRole.AsChoice(),
	LowPriorityRole.AsChoice(),
}

var CategoryOptionSelectChoices = []dg.SelectMenuOption{
	GenericCategoryRole.AsSelectChoice(),
	FeatureCategoryRole.AsSelectChoice(),
	FixCategoryRole.AsSelectChoice(),
	ChoreCategoryRole.AsSelectChoice(),
}

var PriorityOptionSelectChoices = []dg.SelectMenuOption{
	CriticalPriorityRole.AsSelectChoice(),
	ImportantPriorityRole.AsSelectChoice(),
	NormalPriorityRole.AsSelectChoice(),
	LowPriorityRole.AsSelectChoice(),
}

var StatusOptionSelectChoices = []dg.SelectMenuOption{
	{
		Label: db.IssueStatusNames[0],
		Emoji: &dg.ComponentEmoji{Name: db.IssueStatusIcons[0]},
		Value: "0",
	},
	{
		Label: db.IssueStatusNames[1],
		Emoji: &dg.ComponentEmoji{Name: db.IssueStatusIcons[1]},
		Value: "1",
	},
	{
		Label: db.IssueStatusNames[2],
		Emoji: &dg.ComponentEmoji{Name: db.IssueStatusIcons[2]},
		Value: "2",
	},
	{
		Label: db.IssueStatusNames[3],
		Emoji: &dg.ComponentEmoji{Name: db.IssueStatusIcons[3]},
		Value: "3",
	},
}
