package data

import (
	"context"
	db "ecommerce/internal/data/gen"

	"github.com/jackc/pgx/v5"
)

type OrderStore interface {
	CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.Order, error)
	CreateOrderItem(ctx context.Context, arg db.CreateOrderItemParams) error
	GetOrderByID(ctx context.Context, id int64) (db.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]db.OrderItem, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]db.Order, error)
	WithTx(tx pgx.Tx) OrderStore
}

type sqlOrderStore struct {
	q *db.Queries
}

func NewOrderStore(queries *db.Queries) OrderStore {
	return &sqlOrderStore{
		q: queries,
	}
}

func (s *sqlOrderStore) WithTx(tx pgx.Tx) OrderStore {
	return &sqlOrderStore{
		q: db.New(tx),
	}
}

func (s *sqlOrderStore) CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
	return s.q.CreateOrder(ctx, arg)
}

func (s *sqlOrderStore) CreateOrderItem(ctx context.Context, arg db.CreateOrderItemParams) error {
	return s.q.CreateOrderItem(ctx, arg)
}

func (s *sqlOrderStore) GetOrderByID(ctx context.Context, id int64) (db.Order, error) {
	return s.q.GetOrderByID(ctx, id)
}

func (s *sqlOrderStore) GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]db.OrderItem, error) {
	return s.q.GetOrderItemsByOrderID(ctx, orderID)
}

func (s *sqlOrderStore) GetOrdersByUserID(ctx context.Context, userID int64) ([]db.Order, error) {
	return s.q.GetOrdersByUserID(ctx, userID)
}
