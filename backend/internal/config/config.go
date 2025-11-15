package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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
	MailerPort     int
	MailerUsername string
	MailerPassword string
	CloudinaryURL  string
}

func NewConfig() (cfg Config, err error) {
	cfg.Port = os.Getenv("PORT")
	cfg.CacheDSN = os.Getenv("CACHE_DSN")
	cfg.CloudinaryURL = os.Getenv("CLOUDINARY_URL")

	cfg.DBString = os.Getenv("GOOSE_DBSTRING")

	cfg.MailerHost = os.Getenv("MAILER_HOST")
	mailerPortStr := os.Getenv("MAILER_PORT")
	cfg.MailerPort, err = strconv.Atoi(mailerPortStr)
	if err != nil {
		return Config{}, fmt.Errorf("config error: invalid MAILER_PORT value '%s': %w", mailerPortStr, err)
	}
	cfg.MailerSender = os.Getenv("MAILER_SENDER")
	cfg.MailerUsername = os.Getenv("MAILER_USERNAME")
	cfg.MailerPassword = os.Getenv("MAILER_PASSWORD")

	cfg.StartTime = time.Now()
	flag.StringVar(&cfg.Env, "env", "dev", "set development environment")

	flag.Parse()
	return cfg, nil
}
