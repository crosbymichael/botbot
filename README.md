# botbot

A package inspired by urfave/cli for writting discord bots in the same manner as you write CLI applicaitons.


## Example

```go
bot, err := botbot.New("test", clix.GlobalString("token"))
if err != nil {
	return err
}
bot.Commands = []*botbot.Command{
	timeCommand,
}
if err := bot.Start(); err != nil {
	return err
}
<-s
return bot.Close()
```

```go
var timeCommand = &botbot.Command{
	Name:        "time",
	Description: "returns the current time",
	Action: func(ctx *botbot.Context) error {
		now := time.Now()
		return ctx.Send(now.Format(time.RFC3339))
	},
}
```
