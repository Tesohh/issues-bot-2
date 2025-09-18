package man

import (
	"issues/v2/db"

	dg "github.com/bwmarrin/discordgo"
)

var PrioritiesAndCategories = Page{
	ID:    "priorities-and-categories",
	Title: "Priorities and Categories",
	Func: func(s *dg.Session, i *dg.Interaction) ([]dg.MessageComponent, error) {
		guild, err := db.Guilds.Where("id = ?", i.GuildID).First(db.Ctx)
		if err != nil {
			return nil, err
		}

		return []dg.MessageComponent{
			text(` ## Priorities
<@&%s> - worth considering if you have extra time
<@&%s> - default priority
<@&%s> - requires more attention that usual
<@&%s> - requires immediate attention
				`, guild.LowPriorityRoleID,
				guild.NormalPriorityRoleID,
				guild.ImportantPriorityRoleID,
				guild.CriticalPriorityRoleID),
			dg.Separator{},
			text(`## Categories
<@&%s> - default category. ideally, you should never use this
<@&%s> - *enhancement to the project*: new features, refactors, ideas
<@&%s> - unexpected bugs or behaviours that pop up
<@&%s> - anything unrelated to code: documentation, presentations, submissions etc.`,
				guild.GenericCategoryRoleID,
				guild.FeatCategoryRoleID,
				guild.FixCategoryRoleID,
				guild.ChoreCategoryRoleID),
			dg.Separator{},
			text(`## Other roles
These roles are not categories or priorities but you can use them with the [[Shorthand Syntax]] to change the properties of the issues.
See |/man page:Shorthand Syntax| for more.

<@&%s> - Changes the issue type to a [[Discussion]].
<@&%s> - Makes the issue not have any assignees, not even the author / recruiter.`,
				guild.DiscussionRoleID, guild.NobodyRoleID),
		}, nil
	},
}
