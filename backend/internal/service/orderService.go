package service

import (
	"context"
	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCartEmpty        = errors.New("cart is empty")
	ErrCannotBuyOwnItem = errors.New("user cannot buy their own item")
)

type OrderService struct {
	OrderStore    data.OrderStore
	ProductStore  data.ProductStore
	WalletService *WalletService
	CartService   *CartService
	Pool          *pgxpool.Pool
	Logger        *slog.Logger
}

func NewOrderService(
	os data.OrderStore,
	ps data.ProductStore,
	ws *WalletService,
	cs *CartService,
	pool *pgxpool.Pool,
	logger *slog.Logger,
) *OrderService {
	return &OrderService{
		OrderStore:    os,
		ProductStore:  ps,
		WalletService: ws,
		CartService:   cs,
		Pool:          pool,
		Logger:        logger,
	}
}

func (s *OrderService) CreateOrderFromCart(ctx context.Context, buyerID int64) (db.Order, error) {
	logger := s.Logger.With("buyer_id", buyerID)
	logger.Info("Attempting to create order from cart")

	cartItems, err := s.CartService.GetCart(ctx, buyerID)
	if err != nil {
		logger.Error("Failed to get cart", "error", err)
		return db.Order{}, err
	}
	if len(cartItems) == 0 {
		return db.Order{}, ErrCartEmpty
	}

	sellerPayments := make(map[int64]int64)
	var grandTotal int64 = 0

	for _, item := range cartItems {
		if item.Stock < int32(item.Quantity) {
			logger.Warn("Insufficient stock for product", "product_id", item.ProductID, "stock", item.Stock, "wanted", item.Quantity)
			return db.Order{}, errors.New("insufficient stock for " + item.Name)
		}
		itemTotal := int64(item.Price) * int64(item.Quantity)
		sellerPayments[item.SellerID] += itemTotal
		grandTotal += itemTotal
	}
	grandTotal *= 100

	logger.Info("Calculated payment totals", "grand_total", grandTotal, "seller_count", len(sellerPayments))

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		logger.Error("Failed to begin transaction", "error", err)
		return db.Order{}, err
	}
	defer tx.Rollback(ctx)

	txWalletStore := data.NewWalletStore(db.New(tx))
	txOrderStore := data.NewOrderStore(db.New(tx))

	buyerID32 := int32(buyerID)
	_, err = s.WalletService.Debit(ctx, buyerID32, grandTotal)
	if err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			logger.Warn("Buyer has insufficient funds")
			return db.Order{}, ErrInsufficientFunds
		}
		logger.Error("Failed to debit buyer", "error", err)
		return db.Order{}, err
	}
	logger.Info("Buyer debited successfully")

	for sellerID, amount := range sellerPayments {
		sellerID32 := int32(sellerID)
		_, err = s.WalletService.creditWalletInternal(ctx, txWalletStore, sellerID32, amount, "order_payout", "completed", &buyerID32)
		if err != nil {
			logger.Error("Failed to credit seller", "seller_id", sellerID, "error", err)
			return db.Order{}, errors.New("failed to pay seller")
		}
	}
	logger.Info("Sellers credited successfully")

	orderParams := db.CreateOrderParams{
		UserID:      buyerID,
		TotalAmount: grandTotal,
		Status:      "completed",
	}
	order, err := txOrderStore.CreateOrder(ctx, orderParams)
	if err != nil {
		logger.Error("Failed to create order record", "error", err)
		return db.Order{}, err
	}
	logger.Info("Order record created", "order_id", order.ID)

	for _, item := range cartItems {
		itemParams := db.CreateOrderItemParams{
			OrderID:         order.ID,
			ProductID:       item.ProductID,
			SellerID:        item.SellerID,
			Quantity:        int32(item.Quantity),
			PriceAtPurchase: item.Price,
		}
		if err := txOrderStore.CreateOrderItem(ctx, itemParams); err != nil {
			logger.Error("Failed to create order item record", "product_id", item.ProductID, "error", err)
			return db.Order{}, err
		}

	}
	logger.Info("Order items created successfully")

	if err := tx.Commit(ctx); err != nil {
		logger.Error("Failed to commit transaction", "error", err)
		return db.Order{}, err
	}

	if err := s.CartService.ClearCart(ctx, buyerID); err != nil {

		logger.Error("Failed to clear cart after successful order", "error", err)
	}

	logger.Info("Order processing complete", "order_id", order.ID)
	return order, nil
}
