package data

import (
	"context"
	db "ecommerce/internal/data/gen"
)

type AuthStore interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
}

type sqlAuthStore struct {
	q *db.Queries
}

func (s *sqlAuthStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return s.q.CreateUser(ctx, arg)
}

func NewAuthStore(queries *db.Queries) AuthStore {
	return &sqlAuthStore{
		q: queries,
	}
}
