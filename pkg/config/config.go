package config

import (
	"flag"
	"time"
)

// Config contains the application configuration
type Config struct {
	ServerAddress string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
}

// Load loads the configuration from flags and environment variables
func Load() *Config {
	addr := flag.String("addr", "localhost:8080", "http service address")
	flag.Parse()

	return &Config{
		ServerAddress: *addr,
		ReadTimeout:   15 * time.Second,
		WriteTimeout:  15 * time.Second,
		IdleTimeout:   60 * time.Second,
	}
}
