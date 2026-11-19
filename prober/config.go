package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	TargetURL string        `mapstructure:"TARGET_URL"`
	Workers   int           `mapstructure:"WORKERS"`
	Timeout   time.Duration `mapstructure:"TIMEOUT"`
	AuthURL   string        `mapstructure:"AUTH_URL"`
	Email     string        `mapstructure:"EMAIL"`
	Password  string        `mapstructure:"PASSWORD"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() error {
	v := viper.New()
	v.SetConfigFile(".env")

	keys := []string{
		"TARGET_URL",
		"WORKERS",
		"TIMEOUT",
		"AUTH_URL",
		"EMAIL",
		"PASSWORD",
	}
	for _, key := range keys {
		v.SetDefault(key, "")
		if err := v.BindEnv(key); err != nil {
			return fmt.Errorf("bind configuration %s: %w", key, err)
		}
	}

	v.SetDefault("WORKERS", 1)
	v.SetDefault("TIMEOUT", 10*time.Second)

	if err := v.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read prober configuration: %w", err)
	}

	var loaded Config
	if err := v.Unmarshal(&loaded); err != nil {
		return fmt.Errorf("parse prober configuration: %w", err)
	}
	if loaded.TargetURL == "" {
		return fmt.Errorf("TARGET_URL is required")
	}
	targetURL, err := url.Parse(loaded.TargetURL)
	if err != nil || (targetURL.Scheme != "http" && targetURL.Scheme != "https") || targetURL.Host == "" {
		return fmt.Errorf("TARGET_URL must be an absolute HTTP or HTTPS URL")
	}
	if loaded.AuthURL == "" {
		return fmt.Errorf("AUTH_URL is required")
	}
	authURL, err := url.Parse(loaded.AuthURL)
	if err != nil || (authURL.Scheme != "http" && authURL.Scheme != "https") || authURL.Host == "" {
		return fmt.Errorf("AUTH_URL must be an absolute HTTP or HTTPS URL")
	}
	if loaded.Email == "" {
		return fmt.Errorf("EMAIL is required")
	}
	if loaded.Password == "" {
		return fmt.Errorf("PASSWORD is required")
	}
	if loaded.Workers < 1 {
		return fmt.Errorf("WORKERS must be greater than zero")
	}
	if loaded.Timeout <= 0 {
		return fmt.Errorf("TIMEOUT must be positive")
	}

	*c = loaded
	return nil
}
