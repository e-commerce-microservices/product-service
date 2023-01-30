-- name: CreateProduct :one

INSERT INTO
    product (
        name,
        description,
        price,
        thumbnail,
        inventory,
        supplier_id,
        category_id,
        brand
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetProductByID :one

SELECT * FROM product WHERE id = $1;

-- name: GetAllProduct :many

SELECT * FROM product;

-- name: GetProductByCategory :many

SELECT * FROM product WHERE category_id = $1 LIMIT $2 OFFSET $3;

-- name: GetRecommendProduct :many

SELECT * FROM product LIMIT $1 OFFSET $2;

-- name: GetProductBySupplier :many

SELECT * FROM product WHERE supplier_id = $1 LIMIT $2 OFFSET $3;

-- name: UpdateProduct :exec

UPDATE product
SET
    name = $2,
    price = $3,
    inventory = $4,
    brand = $5
WHERE id = $1 and supplier_id = $6;

-- name: GetListProductByIDs :many

SELECT * FROM product WHERE id IN (pg.Array(?));