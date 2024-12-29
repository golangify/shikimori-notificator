package config

import (
	"encoding/json"
	"os"
)

const perm = 0666

type Config struct {
	path        string
	Database    databaseConfig    `json:"database"`
	Telegram    telegramConfig    `json:"telegram"`
	Shikimori   shikimoriConfig   `json:"shikimori"`
	Notificator notificatorConfig `json:"notificator"`
}

type databaseConfig struct {
	DatabaseString string `json:"database_string"`
}

type telegramConfig struct {
	BotApiToken string `json:"bot_api_token"`
}

type shikimoriConfig struct {
	Cookie    string `json:"cookie"`
	XsrfToken string `json:"xsrf_token"`
	UserAgent string `json:"user_agent"`
}

type notificatorConfig struct {
}

func LoadFromJsonFile(path string) (*Config, error) {
	cfg := &Config{path: path}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = cfg.Save()
			if err != nil {
				return nil, err
			}
			return cfg, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	cfg.Save()

	return cfg, nil
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, perm)
}
