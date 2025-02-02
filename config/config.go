package config

import "time"

const perm = 0666

type Config struct {
	path          string
	Database      databaseConfig      `json:"database"`
	Telegram      telegramConfig      `json:"telegram"`
	Shikimori     shikimoriConfig     `json:"shikimori"`
	Notifications notificationsConfig `json:"notification_config"`
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

type notificationsConfig struct {
	CheckDelay time.Duration `json:"check_delay"`
	MailDelay  time.Duration `json:"mail_delay"`
}
