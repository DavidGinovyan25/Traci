package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	JWTSecret       string        `mapstructure:"JWT_SECRET"`
	TokenTTL        time.Duration `mapstructure:"TOKEN_TTL"`
	ConnectTimeout  time.Duration `mapstructure:"CONNECT_TIMEOUT"`
	PostgresURL     string        `mapstructure:"POSTGRES_URL"`
	HTTPAddr        string        `mapstructure:"HTTP_ADDR"`
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
}

func NewConfig() *Config { return &Config{} }

func (c *Config) Load() error {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetDefault("HTTP_ADDR", ":8080")
	v.SetDefault("TOKEN_TTL", "1h")
	v.SetDefault("CONNECT_TIMEOUT", "5s")
	v.SetDefault("SHUTDOWN_TIMEOUT", "10s")
	for _, key := range []string{"POSTGRES_URL", "HTTP_ADDR", "SHUTDOWN_TIMEOUT", "JWT_SECRET", "TOKEN_TTL", "CONNECT_TIMEOUT"} {
		if err := v.BindEnv(key); err != nil {
			return fmt.Errorf("bind configuration: %w", err)
		}
	}
	if err := v.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("не удалось прочитать параметры: %w", err)
	}
	var loaded Config
	if err := v.Unmarshal(&loaded); err != nil {
		return fmt.Errorf("не удалось распарсить параметры: %w", err)
	}
	if loaded.PostgresURL == "" {
		return fmt.Errorf("POSTGRES_URL is required")
	}
	if loaded.HTTPAddr == "" {
		return fmt.Errorf("HTTP_ADDR is required")
	}
	if loaded.ShutdownTimeout <= 0 {
		return fmt.Errorf("SHUTDOWN_TIMEOUT must be positive")
	}
	if len(loaded.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must contain at least 32 bytes")
	}
	if loaded.TokenTTL < time.Second {
		return fmt.Errorf("TOKEN_TTL must be at least one second")
	}
	if loaded.ConnectTimeout <= 0 {
		return fmt.Errorf("CONNECT_TIMEOUT must be positive")
	}
	*c = loaded
	return nil
}
