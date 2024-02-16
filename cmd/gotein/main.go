package main

import (
	"flag"
	"gotein/cfg"
	"gotein/logger"
	"gotein/pkg/handler"
	"gotein/pkg/tgbot"
	"os"
	"os/signal"
	"syscall"

	"github.com/mymmrac/telego"
)

func main() {
	var err error
	// Get configPath parameter.
	cfgPath := flag.String("config", "config.yml", "configuration file path.")
	langPath := flag.String("lang", "lang.yml", "language file path.")
	flag.Parse()

	// Initialize configuration.
	err = cfg.Init(*cfgPath, *langPath)
	if err != nil {
		panic(err)
	}

	cfg := cfg.Config()
	// Initialize bot
	api, err := telego.NewBot(cfg.BotToken)
	if err != nil {
		panic(err)
	}

	// Initialize handler & worker
	h := handler.New(api, cfg)
	proxy, err := tgbot.NewWorker(api, h)
	if err != nil {
		panic(err)
	}
	// Initialize logger.
	logger.InitLogger(cfg.IsDev)

	// Handle signal.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	if cfg.IsWebhook {
		err = proxy.ViaWebhook(c)
	} else {
		err = proxy.ViaPolling(c)
	}

	if err != nil {
		panic(err)
	}

}
