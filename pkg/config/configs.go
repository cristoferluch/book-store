package config

import "time"

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
