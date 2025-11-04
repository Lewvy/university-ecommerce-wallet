package data

import (
	"context"
	"database/sql"
	"ecommerce/domain"
	db "ecommerce/internal/data/gen"
	"errors"
)

type UserStore interface {
	GetUserAuthByEmail(ctx context.Context, email string) (db.GetUserAuthByEmailRow, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByID(ctx context.Context, id int) (db.GetUserByIDRow, error)
}

type sqlUserStore struct {
	q *db.Queries
}

func (s *sqlUserStore) GetUserAuthByEmail(ctx context.Context, email string) (db.GetUserAuthByEmailRow, error) {
	return s.q.GetUserAuthByEmail(ctx, email)
}
func (s *sqlUserStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (_ db.User, _ error) {
	return s.q.CreateUser(ctx, arg)
}
func (s *sqlUserStore) GetUserByID(ctx context.Context, id int) (_ db.GetUserByIDRow, _ error) {
	return s.q.GetUserByID(ctx, int32(id))
}

func NewUserStore(queries *db.Queries) UserStore {
	return &sqlUserStore{
		q: queries,
	}
}
