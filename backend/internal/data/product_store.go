package data

import (
	"context"
	db "ecommerce/internal/data/gen"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func pgErr(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError
	return pgErr, errors.As(err, &pgErr)
}

type ProductStore interface {
	CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error)
	GetProductsByIDs(ctx context.Context, ids []int64) ([]db.Product, error)
	CreateProductImage(ctx context.Context, arg db.CreateProductImageParams) error
	GetProductByID(ctx context.Context, id int64) (db.Product, error)
	GetProductImages(ctx context.Context, productID int64) ([]db.ProductImage, error)
	GetProductsBySeller(ctx context.Context, sellerID int64) ([]db.Product, error)

	GetAllProducts(ctx context.Context) ([]db.Product, error)
	WithTx(tx pgx.Tx) ProductStore
}

type sqlProductStore struct {
	q *db.Queries
}

func NewProductStore(queries *db.Queries) ProductStore {
	return &sqlProductStore{
		q: queries,
	}
}

func (s *sqlProductStore) GetProductsByIDs(ctx context.Context, ids []int64) ([]db.Product, error) {
	return s.q.GetProductsByIDs(ctx, ids)
}

func (s *sqlProductStore) GetProductsBySeller(ctx context.Context, sellerID int64) ([]db.Product, error) {
	products, err := s.q.GetProductsBySeller(ctx, sellerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []db.Product{}, nil
		}
		return nil, err
	}
	return products, nil
}

func (s *sqlProductStore) WithTx(tx pgx.Tx) ProductStore {
	return &sqlProductStore{
		q: db.New(tx),
	}
}

func (s *sqlProductStore) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
	return s.q.CreateProduct(ctx, arg)
}

func (s *sqlProductStore) CreateProductImage(ctx context.Context, arg db.CreateProductImageParams) error {
	return s.q.CreateProductImage(ctx, arg)
}

func (s *sqlProductStore) GetAllProducts(ctx context.Context) ([]db.Product, error) {
	return s.q.GetAllProducts(ctx)
}

func (s *sqlProductStore) GetProductByID(ctx context.Context, id int64) (db.Product, error) {
	product, err := s.q.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Product{}, ErrRecordNotFound
		}
		return db.Product{}, err
	}
	return product, nil
}

func (s *sqlProductStore) GetProductImages(ctx context.Context, productID int64) ([]db.ProductImage, error) {
	images, err := s.q.GetProductImages(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []db.ProductImage{}, nil
		}
		return nil, err
	}
	return images, nil
}
