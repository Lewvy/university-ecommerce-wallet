-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    total_amount,
    status
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: CreateOrderItem :exec
INSERT INTO order_items (
    order_id,
    product_id,
    seller_id,
    quantity,
    price_at_purchase
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrderItemsByOrderID :many
SELECT * FROM order_items
WHERE order_id = $1;

-- name: GetOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC;
