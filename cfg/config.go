package cfg

type configuration struct {
	IsDev       bool  `yaml:"is_dev"`
	OwnerID     int64 `yaml:"owner_id"`
	BotConfig   `yaml:"bot_config"`
	MeiliConfig `yaml:"meili_config"`
}

type BotConfig struct {
	BotToken  string `yaml:"bot_token"`
	HookURL   string `yaml:"hook_url"`
	IsWebhook bool   `yaml:"is_webhook"`
}

type MeiliConfig struct {
	Host string `yaml:"host"`
	Key  string `yaml:"key"`
}
