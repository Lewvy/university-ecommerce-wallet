package config

import (
	"flag"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Port      string
	Env       string
	StartTime time.Time
	DBString  string
	DB        *pgxpool.Pool
}

func NewConfig() (cfg Config) {
	cfg.Port = os.Getenv("PORT")
	cfg.DBString = os.Getenv("GOOSE_DBSTRING")
	cfg.StartTime = time.Now()
	flag.StringVar(&cfg.Env, "env", "dev", "set development environment")
	flag.Parse()
	return cfg
}
