package tgbot

import (
	"github.com/mymmrac/telego"
)

type UpdateHandler interface {
	Handle(u *telego.Update)
}
