package cfg

type Configuration struct {
	Port        int    `yaml:"port"`
	IsDev       bool   `yaml:"is_dev"`
	OwnerID     int64  `yaml:"owner_id"`
	DomainURL   string `yaml:"domain_url"`
	CacheDir    string `yaml:"cache_dir"`
	RegisterKey string `yaml:"register_key"`
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
