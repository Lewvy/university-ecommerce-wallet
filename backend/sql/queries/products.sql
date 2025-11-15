-- name: CreateProduct :one
INSERT INTO products (
    seller_id, 
    name, 
    description, 
    price, 
    stock, 
    image_url
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: CreateProductImage :exec
INSERT INTO product_images (
    product_id, 
    image_url, 
    display_order
) VALUES (
    $1, $2, $3
);

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 AND is_active = TRUE;

-- name: GetProductImages :many
SELECT * FROM product_images
WHERE product_id = $1
ORDER BY display_order ASC;

-- name: GetAllProducts :many
SELECT * FROM products
WHERE is_active = TRUE
ORDER BY created_at DESC;


-- name: GetProductsBySeller :many
SELECT * FROM products
WHERE seller_id = $1
ORDER BY created_at DESC;
