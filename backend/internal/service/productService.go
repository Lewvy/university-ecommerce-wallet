package service

import (
	"context"
	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultUploadWorkerCap = 5

func (s *ProductService) UploadImagesConcurrent(
	ctx context.Context,
	files []*multipart.FileHeader,
	maxConcurrency int,
) ([]string, error) {
	if len(files) == 0 {
		return nil, nil
	}
	if maxConcurrency <= 0 {
		maxConcurrency = defaultUploadWorkerCap
	}
	if maxConcurrency > len(files) {
		maxConcurrency = len(files)
	}

	type job struct {
		idx  int
		file *multipart.FileHeader
	}
	type result struct {
		idx int
		url string
		err error
	}

	jobs := make(chan job)
	results := make(chan result, len(files))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	workerCount := maxConcurrency
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {

				f, err := j.file.Open()
				if err != nil {
					results <- result{idx: j.idx, err: err}

					cancel()
					continue
				}

				func() {
					defer f.Close()
					url, err := s.CloudSvc.UploadImage(ctx, f, j.file.Filename)
					if err != nil {
						results <- result{idx: j.idx, err: err}
						cancel()
						return
					}
					results <- result{idx: j.idx, url: url, err: nil}
				}()
			}
		}()
	}

	go func() {
		for i, fh := range files {
			select {
			case <-ctx.Done():
				return
			default:
				jobs <- job{idx: i, file: fh}
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	imageURLs := make([]string, len(files))
	var firstErr error
	received := 0
	for r := range results {
		received++
		if r.err != nil {
			if firstErr == nil {
				firstErr = r.err
			}
			continue
		}
		imageURLs[r.idx] = r.url
	}

	if firstErr != nil {
		return nil, firstErr
	}

	for i, u := range imageURLs {
		if u == "" {
			return nil, errors.New("missing upload result for image index " + strconv.Itoa(i))
		}
	}

	return imageURLs, nil
}

func (s *ProductService) CreateProductWithFiles(
	ctx context.Context,
	params CreateProductParams,
	files []*multipart.FileHeader,
	maxUploadConcurrency int,
) (db.Product, error) {

	s.Logger.Info("Uploading product images concurrently", "count", len(files))
	imageURLs, err := s.UploadImagesConcurrent(ctx, files, maxUploadConcurrency)
	if err != nil {
		s.Logger.Error("Image uploads failed", "error", err)
		return db.Product{}, err
	}
	if len(imageURLs) == 0 {
		return db.Product{}, errors.New("no images uploaded")
	}

	params.ThumbnailURL = imageURLs[0]
	params.ImageURLs = imageURLs

	return s.CreateProductWithTransaction(ctx, params)
}

type ProductService struct {
	Store    data.ProductStore
	Pool     *pgxpool.Pool
	Logger   *slog.Logger
	CloudSvc CloudService
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

func NewProductService(store data.ProductStore, cloud CloudService, pool *pgxpool.Pool, logger *slog.Logger) *ProductService {
	return &ProductService{
		Store:    store,
		Pool:     pool,
		Logger:   logger,
		CloudSvc: cloud,
	}
}

func (s *ProductService) CreateProductWithTransaction(
	ctx context.Context,
	params CreateProductParams,
) (db.Product, error) {
	s.Logger.Info("Starting transaction for product creation", "seller_id", params.SellerID, "name", params.Name)

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		s.Logger.Error("Failed to begin transaction", "error", err)
		return db.Product{}, errors.New("database connection error")
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.Logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	txStore := s.Store.WithTx(tx)

	newProduct, err := s.CreateProductWithImages(ctx, txStore, params)
	if err != nil {

		return db.Product{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		s.Logger.Error("Failed to commit transaction", "error", err)
		return db.Product{}, errors.New("database commit failed")
	}

	s.Logger.Info("Transaction committed successfully", "product_id", newProduct.ID)
	return newProduct, nil
}

func (s *ProductService) createProduct(
	ctx context.Context,
	txStore data.ProductStore,
	params CreateProductParams,
) (db.Product, error) {
	s.Logger.Info("Creating main product record", "seller_id", params.SellerID, "name", params.Name)

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
		s.Logger.Error("Failed to create product record", "error", err)
		return db.Product{}, err
	}
	return newProduct, nil
}

func (s *ProductService) createProductImages(
	ctx context.Context,
	txStore data.ProductStore,
	productID int64,
	imageURLs []string,
) error {
	s.Logger.Info("Creating product images records", "product_id", productID, "count", len(imageURLs))

	for i, imgURL := range imageURLs {
		imgParams := db.CreateProductImageParams{
			ProductID:    productID,
			ImageUrl:     imgURL,
			DisplayOrder: int32(i),
		}

		err := txStore.CreateProductImage(ctx, imgParams)
		if err != nil {
			s.Logger.Error("Failed to create product image record", "error", err, "image_url", imgURL)
			return fmt.Errorf("failed to insert image %s: %w", imgURL, err)
		}
	}
	return nil
}

func (s *ProductService) CreateProductWithImages(
	ctx context.Context,
	txStore data.ProductStore,
	params CreateProductParams,
) (db.Product, error) {
	s.Logger.Info("Starting combined creation of product and images in store", "seller_id", params.SellerID, "name", params.Name)

	newProduct, err := s.createProduct(ctx, txStore, params)
	if err != nil {
		return db.Product{}, err
	}

	err = s.createProductImages(ctx, txStore, newProduct.ID, params.ImageURLs)
	if err != nil {
		return db.Product{}, err
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

func (s *ProductService) GetProductsBySeller(ctx context.Context, sellerID int64) ([]db.Product, error) {
	s.Logger.Info("Fetching products by seller", "seller_id", sellerID)

	products, err := s.Store.GetProductsBySeller(ctx, sellerID)
	if err != nil {
		s.Logger.Error("Failed to get products by seller", "error", err)
		return nil, err
	}

	return products, nil
}
