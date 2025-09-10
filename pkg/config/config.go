package config

import (
	"flag"
	"time"
)

// Config contiene la configuración de la aplicación
type Config struct {
	ServerAddress string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
}

// Load carga la configuración desde flags y variables de entorno
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
