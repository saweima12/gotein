package handler

import (
	"fmt"
	"gotein/pkg/tgbot"

	"github.com/mymmrac/telego"
)

type BotHandler struct {
	bot *tgbot.TeleBot
}

func (bo *BotHandler) SetBot(b *tgbot.TeleBot) {
	bo.bot = b
}

func (bo *BotHandler) Handle(u *telego.Update) {
	if u.Message != nil {
		mh := tgbot.GetMsgHelper(u.Message)
		bo.handleMessage(mh)
	}
	if u.InlineQuery != nil {
		qh := tgbot.GetInlineHelper(u.InlineQuery)
		bo.handleInline(qh)
	}
}

func (bo *BotHandler) handleMessage(m *tgbot.MsgHelper) {
	fmt.Printf("%+v", m.From.ID)
}

func (bo *BotHandler) handleInline(m *tgbot.InlineHelper) {

}

func (bo *BotHandler) API() *telego.Bot {
	return bo.bot.API()
}
