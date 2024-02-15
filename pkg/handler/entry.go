package handler

import (
	"fmt"
	"gotein/cfg"
	"gotein/pkg/repositories"
	"gotein/pkg/services"
	"gotein/pkg/tgbot"
	"strconv"

	"github.com/meilisearch/meilisearch-go"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

type BotHandler struct {
	bot       *tgbot.TeleBot
	meiliServ services.MeiliServ
}

func New(meiliCfg *cfg.MeiliConfig) *BotHandler {
	meiliCli := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   meiliCfg.Host,
		APIKey: meiliCfg.Key,
	})

	meiliRepo := repositories.NewMeiliRepo(meiliCli)
	meiliServ := services.NewMeiliServ(meiliRepo)

	return &BotHandler{
		meiliServ: meiliServ,
	}
}

func (bo *BotHandler) SetBot(b *tgbot.TeleBot) {
	bo.bot = b
}

func (bo *BotHandler) API() *telego.Bot {
	return bo.bot.API()
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
	fmt.Printf("%+v\n", m.From.ID)

	idStr := strconv.Itoa(int(m.From.ID))

	fmt.Println(bo.meiliServ.IsAllowUser(idStr))
	fmt.Println(bo.meiliServ.IsListenUser(idStr))
}

func (bo *BotHandler) handleInline(m *tgbot.InlineHelper) {
	queryStr := m.Query
	bo.meiliServ.SearchMedia(queryStr)
	// generate response.
	params := telegoutil.InlineQuery(queryStr)
	// Send answer
	bo.API().AnswerInlineQuery(params)
	fmt.Println(queryStr)
}
