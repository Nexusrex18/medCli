package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Cache   CacheConfig   `mapstructure:"cache"`
	Display DisplayConfig `mapstructure:"display"`
	CSV     CSVConfig     `mapstructure:"csv"`
}

type CSVConfig struct {
	FilePath string `mapstructure:"file_path"`
}

type CacheConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	TTL      string `mapstructure:"ttl"`
	MaxItems int    `mapstructure:"max_items"`
}

type DisplayConfig struct {
	Theme       string `mapstructure:"theme"`
	Animations  bool   `mapstructure:"animations"`
	PageSize    int    `mapstructure:"page_size"`
	AutoRefresh bool   `mapstructure:"auto_refresh"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME/.medCli")
	v.AddConfigPath(".")

	// Set defaults
	v.SetDefault("display.theme", "dark")
	v.SetDefault("display.animations", true)
	v.SetDefault("display.page_size", 10)
	v.SetDefault("display.auto_refresh", true)
	v.SetDefault("cache.enabled", true)
	v.SetDefault("cache.ttl", "1h")
	v.SetDefault("cache.max_items", 1000)
	v.SetDefault("csv.file_path", "medicine_data.csv")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; use defaults
			log.Println("No config file found, using defaults")
		} else {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}