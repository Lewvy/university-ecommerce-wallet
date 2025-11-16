package config

import (
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strconv"
	"time"
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

	ESDSN string `env:"ES_DSN"`
}

func NewConfig() (cfg Config, err error) {

	cfg.Port = os.Getenv("PORT")
	cfg.CacheDSN = os.Getenv("CACHE_DSN")
	cfg.CloudinaryURL = os.Getenv("CLOUDINARY_URL")
	cfg.DBString = os.Getenv("GOOSE_DBSTRING")

	cfg.MailerHost = os.Getenv("MAILER_HOST")
	mailerPortStr := os.Getenv("MAILER_PORT")

	cfg.MailerSender = os.Getenv("MAILER_SENDER")
	cfg.MailerUsername = os.Getenv("MAILER_USERNAME")
	cfg.MailerPassword = os.Getenv("MAILER_PASSWORD")

	if mailerPortStr == "" {
		return Config{}, fmt.Errorf("config error: MAILER_PORT environment variable is missing")
	}

	cfg.MailerPort, err = strconv.Atoi(mailerPortStr)
	if err != nil {
		return Config{}, fmt.Errorf("config error: invalid MAILER_PORT value '%s': %w", mailerPortStr, err)
	}

	cfg.ESDSN = os.Getenv("ES_DSN")
	if cfg.ESDSN == "" {
		return Config{}, fmt.Errorf("config error: ES_DSN environment variable is missing")
	}

	cfg.StartTime = time.Now()
	flag.StringVar(&cfg.Env, "env", "dev", "set development environment")
	flag.Parse()

	return cfg, nil
}
