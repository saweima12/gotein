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

	// initialize bot proxy.
	cfg := cfg.Config()

	h := handler.New(&cfg.MeiliConfig)
	bot, err := tgbot.NewBot(&cfg.BotConfig, h)
	if err != nil {
		panic(err)
	}
	// Initialize logger.
	logger.InitLogger(cfg.IsDev)

	// Handle signal.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	if cfg.IsWebhook {
		err = bot.ViaWebhook(c)
	} else {
		err = bot.ViaPolling(c)
	}

	if err != nil {
		panic(err)
	}

}
