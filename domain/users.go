package domain

import (
	"time"
)

type User struct {
	id         uint
	name       string
	email      string
	password   string
	created_at time.Time
	updated_at time.Time
	wallet_id  string
}

func CreateUsers(reqBody []byte) error {
}
