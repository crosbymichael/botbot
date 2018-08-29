package botbot

import "fmt"

// Command for a message
type Command struct {
	// Name of the command
	Name string
	// Description of the command shown in help output
	Description string
	// Channel that this command is bound to
	// if no channel is specified, this command can be called
	// from any channel
	Channel string
	// Action of the command
	Action func(*Context) error
	// Subcommands
	Subcommands []*Command
}

func (c *Command) hasSubCommand(name string) bool {
	for _, s := range c.Subcommands {
		if s.Name == name {
			return true
		}
	}
	return false
}

func (c *Command) run(ctx *Context) error {
	if c.Channel != "" && ctx.Channel() != c.Channel {
		return ctx.Send(fmt.Sprintf("%s not allowed in this channel", c.Name))
	}
	if len(ctx.args) == 0 || !c.hasSubCommand(ctx.args[0]) {
		if c.Action != nil {
			return c.Action(ctx)
		}
		return helpCommand(c.Name, c.Subcommands).run(ctx)
	}
	for _, cmd := range c.Subcommands {
		if cmd.Name == ctx.args[0] {
			return cmd.run(ctx.sub())
		}
	}
	return helpCommand(c.Name, c.Subcommands).run(ctx)
}
