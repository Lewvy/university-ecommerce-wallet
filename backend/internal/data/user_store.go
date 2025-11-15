package data

import (
	"context"
	db "ecommerce/internal/data/gen"

	"github.com/jackc/pgx/v5"
)

type UserStore interface {
	GetUserAuthByEmail(ctx context.Context, email string) (db.GetUserAuthByEmailRow, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByID(ctx context.Context, id int) (db.GetUserByIDRow, error)
	VerifyUserEmail(ctx context.Context, id int) error
	UpdateUserEmail(ctx context.Context, id int, updated_email string) error
	WithTx(tx pgx.Tx) UserStore
}

type sqlUserStore struct {
	q *db.Queries
}

func (s *sqlUserStore) GetUserAuthByEmail(ctx context.Context, email string) (db.GetUserAuthByEmailRow, error) {
	return s.q.GetUserAuthByEmail(ctx, email)
}

func (s *sqlUserStore) VerifyUserEmail(ctx context.Context, id int) error {
	return s.q.VerifyUserEmail(ctx, int32(id))
}
func (s *sqlUserStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (_ db.User, _ error) {
	return s.q.CreateUser(ctx, arg)
}
func (s *sqlUserStore) GetUserByID(ctx context.Context, id int) (_ db.GetUserByIDRow, _ error) {
	return s.q.GetUserByID(ctx, int32(id))
}

func (s *sqlUserStore) UpdateUserEmail(ctx context.Context, id int, updated_email string) error {
	params := db.UpdateUserEmailParams{
		Email: updated_email,
		ID:    int32(id),
	}

	return s.q.UpdateUserEmail(ctx, params)

}

func (s *sqlUserStore) WithTx(tx pgx.Tx) UserStore {
	return &sqlUserStore{
		q: db.New(tx),
	}
}

func NewUserStore(queries *db.Queries) UserStore {
	return &sqlUserStore{
		q: queries,
	}
}
