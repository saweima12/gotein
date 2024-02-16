package handler

import (
	"gotein/cfg"
	"gotein/logger"
	"gotein/pkg/repositories"
	"gotein/pkg/services"
	"gotein/pkg/tgbot"

	"github.com/meilisearch/meilisearch-go"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type BotHandler struct {
	api       *telego.Bot
	cfg       *cfg.Configuration
	mediaServ services.MediaServ
	userServ  services.UserServ
	cacheServ services.CacheServ
}

func New(bot *telego.Bot, cfg *cfg.Configuration) *BotHandler {
	meiliCfg := cfg.MeiliConfig
	meiliCli := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   meiliCfg.Host,
		APIKey: meiliCfg.Key,
	})

	meiliRepo := repositories.NewMeiliRepo(meiliCli)
	cacheServ := services.NewCacheServ(bot)
	userServ := services.NewUserServ(meiliRepo)
	mediaServ := services.NewMediaServ(meiliRepo, cacheServ, cfg)

	return &BotHandler{
		api:       bot,
		mediaServ: mediaServ,
		userServ:  userServ,
		cacheServ: cacheServ,
	}
}

func (bo *BotHandler) Handle(u *telego.Update) {
	if u.InlineQuery != nil {
		qh := tgbot.GetInlineHelper(u.InlineQuery)
		bo.handleInline(qh)
	}

	if u.Message != nil {
		mh := tgbot.GetMsgHelper(u.Message)
		bo.handleMessage(mh)
	}
}

func (bo *BotHandler) handleInline(m *tgbot.InlineHelper) {
	queryStr := m.Query
	result := bo.mediaServ.SearchMedia(queryStr)
	// generate response.
	params := tu.InlineQuery(m.ID, result...)
	// Send answer
	err := bo.API().AnswerInlineQuery(params)
	if err != nil {
		logger.Errorf("Answer inlineQuery failed, err: %v", err)
	}
}

func (bo *BotHandler) handleMessage(m *tgbot.MsgHelper) {
	if m.Chat.Type == "supergroup" {
		bo.handleGroupMessage(m)
	} else {
		bo.handleChatMessage(m)
	}
}

func (bo *BotHandler) API() *telego.Bot {
	return bo.api
}
