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

	chatCmdMap  tgbot.CommandMap
	groupCmdMap tgbot.CommandMap
}

func New(bot *telego.Bot, cfg *cfg.Configuration) *BotHandler {
	meiliCfg := cfg.MeiliConfig
	meiliCli := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   meiliCfg.Host,
		APIKey: meiliCfg.Key,
	})

	meiliRepo := repositories.NewMeiliRepo(meiliCli)
	cacheServ := services.NewCacheServ(bot, cfg)
	userServ := services.NewUserServ(meiliRepo)
	mediaServ := services.NewMediaServ(meiliRepo, cacheServ, cfg)

	h := &BotHandler{
		api:       bot,
		cfg:       cfg,
		mediaServ: mediaServ,
		userServ:  userServ,
		cacheServ: cacheServ,

		chatCmdMap:  make(tgbot.CommandMap),
		groupCmdMap: make(tgbot.CommandMap),
	}
	h.Init()

	return h
}

func (bo *BotHandler) Init() {
	bo.chatCmdMap["ak"] = bo.akCommand
	bo.chatCmdMap["sk"] = bo.skCommand
	bo.chatCmdMap["rm"] = bo.rmCommand
	bo.chatCmdMap["register"] = bo.registerCommand

	bo.groupCmdMap["listen"] = bo.listenCommand
	bo.groupCmdMap["addkeyword"] = bo.addkeywordCommand
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
		logger.GetLogger().Errorf("Answer inlineQuery failed, err: %v", err)
	}
}

func (bo *BotHandler) handleMessage(m *tgbot.MsgHelper) {
	logger.Debug(m.Chat.Type)
	if m.Chat.Type == "supergroup" || m.Chat.Type == "channel" {
		bo.handleGroupMessage(m)
	} else {
		bo.handleChatMessage(m)
	}
}

func (bo *BotHandler) API() *telego.Bot {
	return bo.api
}
