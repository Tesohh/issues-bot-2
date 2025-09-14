package man

import dg "github.com/bwmarrin/discordgo"

var Shorthand = Page{
	ID:    "shorthand",
	Title: "Shorthand Syntax",
	Content: []dg.MessageComponent{
		text(`
The shorthand syntax is the *most efficient way* to create new **Issues**, **Discussions (WIP)** and **Tasks (WIP)**.
It is recommended to** always use this** instead of |/new| for a much *faster *and *frictionless *experience.
Users of |v1| will be already familiar with it, as it remains mostly unchanged.

The principle is the same as |/new|: you need to be in the |#xxx-issues| channel for your project so the bot can infer your project.`),
		media("https://i.ibb.co/99RBTWr2/anatomy.png"),

		text(`Notes and tips:
- Start with a |- | like you would in a todo list
- You can place components (category, priority, etc.) in any order, but the sequence above makes the most sense
  - Technically you can even put them in between your title, but of course that will lead to broken text.
- Place dependencies (|depends on ...|) on a new line (note: dependencies are WIP)
- Unlike |/new|, you can assign the issue to multiple people

the |/new| equivalent would be: |/new title:add NUKE 3D raytraced graphics category:FEAT priority:CRITICAL assign:@tesohh dependson:#CVV-1 add nuke|`),
		text(`## Other
			### Nobody`),

		media("https://i.ibb.co/Q7DqGn6t/nobody.png"),
		text(`This is used to assign the issue to nobody, not even yourself as the recruiter.
|/new| equivalent: |/new title:make the frontend nicer category:FEAT  priority:IMPORTANT  nobody:True|`),

		text("### Discussions"),
		media("https://i.ibb.co/QvHG2mg4/Group-1.png"),
		text(`This is used to create a new discussion. (see |/man page:Discussions| for more)
|/new| equivalent: |/new title:What beer to get for the barbecue? discussion:True|`),

		text("### Tags"),
		media("https://i.ibb.co/W41bMLpL/tags.png"),
		text(`Lastly, this is how you add tags; by prefixing them with a |+| sign.
|/new| equivalent: |/new title:Write docker-compose script category:CHORE tags:devops|`),
	},
}
