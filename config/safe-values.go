package config

import "time"

// защита во избежание блокировки от ddos
func (c *Config) setSafeValues() {
	const (
		minCheckDelay = time.Second
		minMailDelay  = time.Second / 3
	)

	if c.Notifications.CheckDelay < minCheckDelay {
		c.Notifications.CheckDelay = minCheckDelay
	}
	if c.Notifications.MailDelay < minMailDelay {
		c.Notifications.MailDelay = minMailDelay
	}
}
