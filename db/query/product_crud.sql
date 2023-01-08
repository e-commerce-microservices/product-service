-- name: CreateProduct :one
INSERT INTO product (name, description, price, thumbnail, inventory, supplier_id, category_id)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING *;


-- name: GetProductByID :one
SELECT * FROM product WHERE id = $1;

-- name: GetAllProduct :many
SELECT * FROM product;

-- name: GetProductByCategory :many
SELECT * FROM product WHERE category_id = $1;