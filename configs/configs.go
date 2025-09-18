package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	App struct {
		Name         string
		Port         int
		Version      string
		WriteTimeout time.Duration
		ReadTimeout  time.Duration
		IdleTimeout  time.Duration
		Environment  string
		LogPath      string
	}
	Database struct {
		URI             string
		MaxConns        int32
		MinConns        int32
		MaxConnLifetime time.Duration
		MaxConnIdleTime time.Duration
	}
	StartTime time.Time
}

var Cfg = &Config{}

func NewConfigLoader() error {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading configs file, %w", err)
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %w", err)
	}

	Cfg.StartTime = time.Now()

	return nil
}
