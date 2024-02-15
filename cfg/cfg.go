package cfg

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	cfg     *configuration
	textMap = make(textLangMap)
	GetText = textMap.GetText
)

func Config() *configuration {
	return cfg
}

func Init(cfgPath string, langPath string) error {
	var err error
	f, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("Error: open [%s] failed, err: %v", cfgPath, err)
	}

	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return fmt.Errorf("Error: Unmarshal %s failed, err: %v", cfgPath, err)
	}

	f, err = os.ReadFile(langPath)
	if err != nil {
		return fmt.Errorf("Error: open [%s] failed, err: %v", cfgPath, err)
	}

	err = yaml.Unmarshal(f, &textMap)
	if err != nil {
		return fmt.Errorf("Error: Unmarshal %s failed, err: %v", cfgPath, err)
	}

	return nil
}
