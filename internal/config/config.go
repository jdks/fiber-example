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
		Fiber fiber.Config
	}

	StoreConfig struct {
		DBConnectionURL string
	}
)

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.AddConfigPath("config/")
	v.SetConfigName("local")
	v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg := Config{
		Store: StoreConfig{
			DBConnectionURL: v.GetString("db"),
		},
		Fiber: fiber.Config{
			Prefork:     v.GetBool("fiber.prefork_enabled"),
			ReadTimeout: time.Duration(v.GetInt("fiber.read_timeout")) * time.Second,
			JSONDecoder: sonic.Unmarshal,
			JSONEncoder: sonic.Marshal,
		},
	}

	return &cfg, nil
}
