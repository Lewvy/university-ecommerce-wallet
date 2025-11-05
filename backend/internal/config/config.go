package config

import (
	"flag"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Port           string
	Env            string
	StartTime      time.Time
	DBString       string
	DB             *pgxpool.Pool
	CacheDSN       string
	MailerSender   string
	MailerHost     string
	MailerPort     string
	MailerUsername string
	MailerPassword string
}

func NewConfig() (cfg Config) {
	cfg.Port = os.Getenv("PORT")
	cfg.CacheDSN = os.Getenv("CACHE_DSN")

	cfg.DBString = os.Getenv("GOOSE_DBSTRING")

	cfg.MailerHost = os.Getenv("MAILER_HOST")
	cfg.MailerPort = os.Getenv("MAILER_PORT")
	cfg.MailerSender = os.Getenv("MAILER_SENDER")
	cfg.MailerUsername = os.Getenv("MAILER_USERNAME")
	cfg.MailerPassword = os.Getenv("MAILER_PASSWORD")

	cfg.StartTime = time.Now()
	flag.StringVar(&cfg.Env, "env", "dev", "set development environment")
	flag.Parse()
	return cfg
}
