package texter

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type Texter interface {
	TextFor(user tb.User, code string) string
}
