package botbot

import (
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// Context of a message
type Context struct {
	s *discordgo.Session
	m *discordgo.MessageCreate

	channelName string
	args        []string
	app         map[string]interface{}
}

// Channel name that is known to users
func (c *Context) Channel() string {
	return c.channelName
}

// Arg in the message
func (c *Context) Arg(i int) string {
	if len(c.args) < i+1 {
		return ""
	}
	return c.args[i]
}

// Attachment reader of the message
func (c *Context) Attachment(i int) (io.ReadCloser, error) {
	if len(c.m.Attachments) < i+1 {
		return nil, errors.Errorf("no attachment at %d", i)
	}
	url := c.m.Attachments[i].URL

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Send a response to the same channel
func (c *Context) Send(msg string) error {
	_, err := c.s.ChannelMessageSend(c.m.ChannelID, msg)
	return err
}

// Value of a global context var
func (c *Context) Value(name string) interface{} {
	return c.app[name]
}

func (c *Context) sub() *Context {
	return &Context{
		s:           c.s,
		m:           c.m,
		channelName: c.channelName,
		args:        c.args[1:],
		app:         c.app,
	}
}
