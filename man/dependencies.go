package man

import dg "github.com/bwmarrin/discordgo"

var Dependencies = Page{
	ID:    "dependencies",
	Title: "Dependencies",
	Func: func(s *dg.Session, i *dg.Interaction) ([]dg.MessageComponent, error) {
		return []dg.MessageComponent{
			text(`
This is a new feature of |YIELD 2|. 
Dependencies ensure that one or more issues must be completed before another can be closed. 
This also makes it very easy to see the relationships between issues.`),
			media("https://i.ibb.co/spzv21Ws/immagine.png"),

			text(`
This is the same system that also powers [[Tasks]].

⚠️ If an issue has at least one open dependency (|🟩 todo| or |🟦 doing|), it cannot be closed until those dependencies are resolved.`),

			text(`
## Adding dependencies
There are three ways to add dependencies:
1. **When creating an issue with |/new|:**
  Use the |dependson| flag in |/new| to link a **single** pre-existing issue.
2. **When creating an issue with [[shorthand syntax]]:**
  Add one or more lines like:
  |depends on <issue thread mention>|
  -# see |/man page:Shorthand Syntax| for more
3. **On an existing issue:**
  Use |/issue dependson| to add or remove (toggle) dependencies`),
		}, nil
	},
}
