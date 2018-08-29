package botbot

import (
	"bytes"
	"fmt"
)

func helpCommand(name string, cmds []*Command) *Command {
	return &Command{
		Name:        "help",
		Subcommands: cmds,
		Action: func(ctx *Context) error {
			buf := bytes.NewBuffer(nil)
			buf.WriteString("```\n")
			fmt.Fprintf(buf, "Help for %s\n", name)
			fmt.Fprint(buf, "Commands:\n\n")
			for _, c := range cmds {
				fmt.Fprintf(buf, "%s: %s\n", c.Name, c.Description)
			}
			buf.WriteString("```")
			return ctx.Send(buf.String())
		},
	}
}
