package command

import (
	"fmt"
	"issues/v2/slash"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

type node struct {
	name        string
	description string
	options     []*dg.ApplicationCommandOption
	children    []*node
	parent      *node
}

func (n *node) fullName() string {
	if n.parent == nil {
		return "/" + n.name
	}
	return n.parent.fullName() + " " + n.name
}

func (n *node) prettyTreeString(depth int) string {
	padding := strings.Repeat("\t", depth)
	str := fmt.Sprintf("%s`%s` - %s", padding, n.fullName(), n.description)
	for _, c := range n.children {
		str += "\n" + c.prettyTreeString(depth+1)
	}
	return str
}

func collect(parent *node, opt *dg.ApplicationCommandOption) {
	node := &node{
		name:        opt.Name,
		description: opt.Description,
		parent:      parent,
	}

	parent.children = append(parent.children, node)
	for _, opt := range opt.Options {
		if opt.Type == dg.ApplicationCommandOptionSubCommand || opt.Type == dg.ApplicationCommandOptionSubCommandGroup {
			collect(node, opt)
		} else {
			node.options = append(node.options, opt)
		}
	}
}

var Help = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "help",
		Description: "displays list of all commands",
	},
	Disabled: false,
	Func: func(s *dg.Session, i *dg.Interaction) error {
		nodes := []*node{}

		for _, command := range Commands {
			node := &node{
				name:        command.Name,
				description: command.Description,
			}

			nodes = append(nodes, node)
			for _, opt := range command.Options {
				if opt.Type == dg.ApplicationCommandOptionSubCommand || opt.Type == dg.ApplicationCommandOptionSubCommandGroup {
					collect(node, opt)
				} else {
					node.options = append(node.options, opt)
				}
			}
		}

		str := ""
		for _, node := range nodes {
			str += node.prettyTreeString(0) + "\n" + "\n"
		}

		container := dg.Container{
			Components: []dg.MessageComponent{
				dg.TextDisplay{Content: "# Help"},
				dg.Separator{},
				dg.TextDisplay{Content: str},
			},
		}
		return slash.ReplyWithComponents(s, i, true, container)
	},
}
