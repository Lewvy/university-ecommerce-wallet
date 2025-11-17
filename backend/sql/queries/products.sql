-- name: CreateProduct :one
INSERT INTO products (
    seller_id, 
	seller_name,
	seller_phone,
    name, 
    description, 
	condition,
	category,
    price, 
    stock, 
    image_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetProductsByPriceRange :many
SELECT * FROM products
WHERE is_active = TRUE
  AND price BETWEEN $1 AND $2 
ORDER BY 
  CASE WHEN $3 = 'asc' THEN price END ASC,
  CASE WHEN $3 = 'desc' THEN price END DESC,
  created_at DESC;

-- name: GetProductsByCategory :many
SELECT * FROM products
WHERE category = $1 AND is_active = TRUE
ORDER BY created_at DESC;

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

-- name: GetProductsByIDs :many
SELECT * FROM products
WHERE id = ANY($1::bigint[]);
