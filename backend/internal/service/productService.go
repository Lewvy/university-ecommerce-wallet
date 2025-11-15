package service

import (
	"context"
	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"errors"
	"fmt"
	"log/slog"
)

type ProductService struct {
	Store  data.ProductStore
	Logger *slog.Logger
}

type CreateProductParams struct {
	SellerID     int64
	Name         string
	Description  string
	Price        int32
	Stock        int32
	ThumbnailURL string
	ImageURLs    []string
}

type ProductDetails struct {
	db.Product
	Images []db.ProductImage `json:"images"`
}

func NewProductService(store data.ProductStore, logger *slog.Logger) *ProductService {
	return &ProductService{
		Store:  store,
		Logger: logger,
	}
}

func (s *ProductService) CreateProductWithImages(
	ctx context.Context,
	txStore data.ProductStore,
	params CreateProductParams,
) (db.Product, error) {

	s.Logger.Info("Creating new product in transaction", "seller_id", params.SellerID, "name", params.Name)

	productParams := db.CreateProductParams{
		SellerID:    params.SellerID,
		Name:        params.Name,
		Description: data.NewPGText(params.Description),
		Price:       params.Price,
		Stock:       params.Stock,
		ImageUrl:    data.NewPGText(params.ThumbnailURL),
	}

	newProduct, err := txStore.CreateProduct(ctx, productParams)
	if err != nil {
		s.Logger.Error("Failed to create product record in transaction", "error", err)
		return db.Product{}, err
	}

	s.Logger.Info("Creating product images", "product_id", newProduct.ID, "count", len(params.ImageURLs))
	for i, imgURL := range params.ImageURLs {
		imgParams := db.CreateProductImageParams{
			ProductID:    newProduct.ID,
			ImageUrl:     imgURL,
			DisplayOrder: int32(i),
		}

		err = txStore.CreateProductImage(ctx, imgParams)
		if err != nil {
			s.Logger.Error("Failed to create product image record in transaction", "error", err, "image_url", imgURL)

			return db.Product{}, fmt.Errorf("failed to insert image %s: %w", imgURL, err)
		}
	}

	s.Logger.Info("Product and images created successfully", "product_id", newProduct.ID)
	return newProduct, nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]db.Product, error) {
	s.Logger.Info("Fetching all products")
	products, err := s.Store.GetAllProducts(ctx)
	if err != nil {
		s.Logger.Error("Failed to retrieve all products from store", "error", err)
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProductDetails(ctx context.Context, productID int64) (ProductDetails, error) {
	s.Logger.Info("Fetching product details", "product_id", productID)

	var details ProductDetails

	product, err := s.Store.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			s.Logger.Warn("Product not found", "product_id", productID)
		} else {
			s.Logger.Error("Failed to get product by ID", "product_id", productID, "error", err)
		}
		return details, err
	}

	images, err := s.Store.GetProductImages(ctx, productID)
	if err != nil {

		if !errors.Is(err, data.ErrRecordNotFound) {
			s.Logger.Warn("Failed to get product images", "product_id", productID, "error", err)
		}
	}

	details.Product = product
	details.Images = images

	return details, nil
}
