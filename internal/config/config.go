package config

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	Port      string
	Env       string
	StartTime time.Time
}

func NewConfig() (cfg Config) {
	cfg.Port = os.Getenv("PORT")
	cfg.StartTime = time.Now()
	flag.StringVar(&cfg.Env, "env", "dev", "set development environment")
	flag.Parse()
	return cfg
}
