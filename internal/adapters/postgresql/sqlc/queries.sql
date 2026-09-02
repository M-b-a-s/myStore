-- name: ListProducts :many
SELECT * FROM products;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (name, price_in_cents, quantity)
VALUES ($1, $2, $3)
RETURNING *;