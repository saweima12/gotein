package tgbot

import (
	"gotein/cfg"
	"gotein/logger"
	"os"

	"github.com/mymmrac/telego"
	"github.com/panjf2000/ants"
)

type UpdateHandler interface {
	SetBot(b *TeleBot)
	Handle(u *telego.Update)
}

type TeleBot struct {
	id      int64
	api     *telego.Bot
	pool    *ants.PoolWithFunc
	handler UpdateHandler
}

func NewBot(botCfg *cfg.BotConfig, h UpdateHandler) (*TeleBot, error) {
	api, err := telego.NewBot(botCfg.BotToken)
	if err != nil {
		return nil, err
	}

	result := &TeleBot{
		api:     api,
		handler: h,
	}

	// Get the infomation of bot.
	user, err := api.GetMe()
	if err != nil {
		return nil, err
	}
	result.id = user.ID

	h.SetBot(result)
	// Create handler pool.
	np, err := ants.NewPoolWithFunc(20, func(i interface{}) {
		u := i.(telego.Update)
		result.handler.Handle(&u)
	})
	if err != nil {
		return nil, err
	}

	// Insert pool
	result.pool = np

	return result, nil
}

func (te *TeleBot) ViaWebhook(stopCh chan os.Signal) error {
	err := te.api.SetWebhook(&telego.SetWebhookParams{})
	if err != nil {
		return err
	}

	logger.Info("Start to recive webhhok updates.")
	return nil
}

func (te *TeleBot) ViaPolling(stopCh chan os.Signal) error {
	updates, err := te.api.UpdatesViaLongPolling(nil)
	if err != nil {
		return err
	}
	logger.Info("Start to polling updates.")

	te.reciveUpdate(updates, stopCh)
	return nil
}

func (te *TeleBot) reciveUpdate(updates <-chan telego.Update, stopCh chan os.Signal) {

	for {
		select {
		case update := <-updates:
			te.pool.Invoke(update)
		case <-stopCh:
			return
		}
	}
}

func (te *TeleBot) API() *telego.Bot {
	return te.api
}

func (te *TeleBot) ID() int64 {
	return te.id
}
