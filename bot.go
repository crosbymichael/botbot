package botbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/shlex"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Handler of messages
type Handler interface {
	Run(*Context) error
}

// New returns a new bot
func New(handle, token string) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	b := &Bot{
		Handle:     handle,
		AppContext: make(map[string]interface{}),
		d:          dg,
		channels:   make(map[string]string),
	}
	dg.AddHandler(b.handler)
	return b, nil
}

// Bot for handling discord messages
type Bot struct {
	Commands   []*Command
	Handle     string
	AppContext map[string]interface{}

	d        *discordgo.Session
	channels map[string]string
}

// Close the server
func (b *Bot) Close() error {
	return b.d.Close()
}

// Start listening and processing messages
func (b *Bot) Start() error {
	b.Commands = append(b.Commands, helpCommand(b.Handle, b.Commands))
	return b.d.Open()
}

func (b *Bot) handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// don't handle messages from the bot
	if m.Author.ID == s.State.User.ID {
		return
	}
	if err := b.handleMessage(s, m); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"channel": m.ChannelID,
		}).Error("handle message")
	}
}

func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	parts, err := shlex.Split(m.Content)
	if err != nil {
		logrus.WithError(err).Warn("unable to parse message")
		return nil
	}
	if len(parts) < 1 {
		return nil
	}
	if parts[0] != b.Handle {
		return nil
	}
	parts = parts[1:]
	cname, ok := b.channels[m.ChannelID]
	if !ok {
		ch, err := s.Channel(m.ChannelID)
		if err != nil {
			return errors.Wrap(err, "get channel name")
		}
		cname = ch.Name
		b.channels[m.ChannelID] = cname
	}
	ctx := &Context{
		s:           s,
		m:           m,
		channelName: cname,
		app:         b.AppContext,
	}
	if len(parts) < 1 {
		return b.command("help").run(ctx)
	}
	var cmdErr error
	ctx.args = parts[1:]
	switch name := parts[0]; name {
	case "", "help":
		cmdErr = b.command("help").run(ctx)
	default:
		cmd := b.command(name)
		if cmd == nil {
			return b.command("help").run(ctx)
		}
		cmdErr = cmd.run(ctx)
	}
	if cmdErr != nil {
		s.ChannelMessageSend(m.ChannelID, cmdErr.Error())
	}
	return nil
}

func (b *Bot) command(name string) *Command {
	for _, c := range b.Commands {
		if name == c.Name {
			return c
		}
	}
	return nil
}
