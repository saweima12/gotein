package tgbot

import (
	"fmt"
	"gotein/cfg"
	"gotein/logger"
	"os"

	"github.com/mymmrac/telego"
	"github.com/panjf2000/ants"
)

type TeleBot struct {
	id      int64
	cfg     *cfg.Configuration
	api     *telego.Bot
	pool    *ants.PoolWithFunc
	handler UpdateHandler
}

func NewWorker(api *telego.Bot, h UpdateHandler, cfg *cfg.Configuration) (*TeleBot, error) {
	result := &TeleBot{
		api:     api,
		handler: h,
		cfg:     cfg,
	}

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
	err := te.api.SetWebhook(&telego.SetWebhookParams{
		URL: te.cfg.HookURL + "/bot" + te.api.Token(),
	})
	if err != nil {
		return err
	}

	// Receive information about webhook
	info, _ := te.api.GetWebhookInfo()
	portStr := fmt.Sprintf(":%v", te.cfg.Port)
	logger.Infof("Webhook Info: %+v\n", info)

	logger.Infof(`Listen port -> "%v"`, portStr)

	go func() {
		te.api.StartWebhook(portStr)
	}()

	defer func() {
		te.api.StopWebhook()
	}()

	updates, err := te.api.UpdatesViaWebhook("/bot" + te.api.Token())
	if err != nil {
		return err
	}

	logger.Info("Start to recive webhhok updates.")
	te.reciveUpdate(updates, stopCh)
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
