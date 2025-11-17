package service

import (
	"context"
	"ecommerce/internal/cache"
	"ecommerce/internal/data"
	"errors"

	// db "ecommerce/internal/data/gen"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type CartItem struct {
	ProductID int64  `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Name      string `json:"name"`
	Price     int32  `json:"price"`
	ImageUrl  string `json:"image_url"`
	Stock     int32  `json:"stock"`
	SellerID  int64  `json:"seller_id"`
}

type CartService struct {
	Store  data.ProductStore
	Cache  *cache.ValkeyCache
	Logger *slog.Logger
}

func NewCartService(store data.ProductStore, cache *cache.ValkeyCache, logger *slog.Logger) *CartService {
	return &CartService{
		Store:  store,
		Cache:  cache,
		Logger: logger,
	}
}

func (s *CartService) AddToCart(ctx context.Context, userID int64, productID int64, quantity int) error {
	logger := s.Logger.With("user_id", userID, "product_id", productID, "quantity", quantity)

	product, err := s.Store.GetProductByID(ctx, productID)

	if product.SellerID == userID {
		logger.Warn("User attempted to add their own item to cart", "product_id", productID, "seller_id", product.SellerID)
		return ErrCannotBuyOwnItem
	}

	if err != nil {
		logger.Warn("Failed to add to cart, product not found")
		return data.ErrRecordNotFound
	}
	if product.Stock < int32(quantity) {
		logger.Warn("Failed to add to cart, insufficient stock")
		return errors.New("insufficient stock")
	}

	logger.Info("Adding item to cart")
	return s.Cache.AddToCart(ctx, userID, productID, quantity)
}

func (s *CartService) GetCart(ctx context.Context, userID int64) ([]CartItem, error) {
	logger := s.Logger.With("user_id", userID)
	logger.Info("Fetching user cart")

	cart, err := s.Cache.GetCart(ctx, userID)
	if err != nil {
		logger.Error("Failed to get simple cart from cache", "error", err)
		return nil, err
	}

	if len(cart) == 0 {
		return []CartItem{}, nil
	}

	productIDs := make([]int64, 0, len(cart))
	quantityMap := make(map[int64]int, len(cart))

	for idStr, qtyStr := range cart {
		productID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.Warn("Invalid productID in cart, skipping", "product_id", idStr)
			continue
		}
		quantity, err := strconv.Atoi(qtyStr)
		if err != nil {
			logger.Warn("Invalid quantity in cart, skipping", "product_id", idStr, "quantity", qtyStr)
			continue
		}

		productIDs = append(productIDs, productID)
		quantityMap[productID] = quantity
	}

	if len(productIDs) == 0 {
		return []CartItem{}, nil
	}

	pgProductIDs := pgtype.Array[int64]{
		Elements: productIDs,
		Valid:    true,
		Dims:     []pgtype.ArrayDimension{{Length: int32(len(productIDs)), LowerBound: 1}},
	}

	products, err := s.Store.GetProductsByIDs(ctx, pgProductIDs.Elements)
	if err != nil {
		logger.Error("Failed to get product details for cart", "error", err)
		return nil, err
	}

	fullCart := make([]CartItem, 0, len(products))
	for _, p := range products {
		fullCart = append(fullCart, CartItem{
			ProductID: p.ID,
			Quantity:  quantityMap[p.ID],
			Name:      p.Name,
			Price:     p.Price,
			ImageUrl:  p.ImageUrl.String,
			Stock:     p.Stock,
			SellerID:  p.SellerID,
		})
	}

	return fullCart, nil
}

func (s *CartService) UpdateCartItemQuantity(ctx context.Context, userID, productID int64, quantity int) error {
	s.Logger.Info("Updating cart item quantity", "user_id", userID, "product_id", productID, "quantity", quantity)
	return s.Cache.UpdateCartItemQuantity(ctx, userID, productID, quantity)
}

func (s *CartService) DeleteCartItem(ctx context.Context, userID, productID int64) error {
	s.Logger.Info("Deleting item from cart", "user_id", userID, "product_id", productID)
	return s.Cache.DeleteCartItem(ctx, userID, productID)
}

func (s *CartService) ClearCart(ctx context.Context, userID int64) error {
	s.Logger.Info("Clearing user cart", "user_id", userID)
	return s.Cache.ClearCart(ctx, userID)
}
