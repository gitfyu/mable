package mable

import (
	"github.com/gitfyu/mable/chat"
)

func handlePlay(c *conn) error {
	return c.Disconnect(&chat.Msg{
		Text:  "TODO",
		Color: chat.ColorYellow,
	})
}
