// Package config provides configuration management for the application.
// It defines structs for storing configuration settings and functions
// for loading settings from environment variables.
package config

import (
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Store StoreConfig
		App   AppConfig
		Fiber fiber.Config
	}

	AppConfig struct {
		LogLevel string
		Host     string
		Port     string
	}

	StoreConfig struct {
		DBConnectionURL string
		LogLevel        string
	}
)

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg := Config{
		Store: StoreConfig{
			DBConnectionURL: v.GetString("store.db.connection_url"),
			LogLevel:        v.GetString("store.log.level"),
		},
		Fiber: fiber.Config{
			Prefork:     v.GetBool("fiber.prefork_enabled"),
			ReadTimeout: time.Duration(v.GetInt("fiber.read_timeout")) * time.Second,
			JSONDecoder: sonic.Unmarshal,
			JSONEncoder: sonic.Marshal,
		},
		App: AppConfig{
			LogLevel: v.GetString("app.log.level"),
			Host:     v.GetString("app.host"),
			Port:     v.GetString("app.port"),
		},
	}

	return &cfg, nil
}
