package man

import dg "github.com/bwmarrin/discordgo"

var Tasks = Page{
	ID:    "tasks",
	Title: "Tasks",
	Func: func(_ *dg.Session, _ *dg.Interaction) ([]dg.MessageComponent, error) {
		return []dg.MessageComponent{
			text(`
Tasks are a very convenient way to divide a larger issue into smaller actionable pieces.
They function similarly to issues but only have a title and a status, which can only be |todo| or |done|.
They don't increment the codes, don't pollute your lists and don't create an individual thread.
Under the hood they use the same system as [[Dependencies]].
## Creating tasks
There are two ways to create a task. 
For both you need to be inside the thread of the dependant (parent) issue.

1. **Shorthand syntax**
just start a message with |- | like usual. The only parameter is the title. Don't try to use @roles! They are not accepted by the parser and will just show up as an ugly string.

2. **|/task new|**
## |/task| command
Other than |/task new|, with the |/task| command you can:

- |/task rename|
Similar to |/issue rename|

- |/task toggle|
Toggles between the |todo| and |done| status of a task.
You can also do that by just clicking on the button on the right of the task.

- |/task promote|
Promotes a task to a full issue. In case things get outta hand.
The issue will get it's own code, thread and everything.
It will inherit [[Priority]] and [[Category]] from it's parent.`),
		}, nil
	},
}
