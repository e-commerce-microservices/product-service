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

-- name: GetProductByCategoryAndTime :many
SELECT * FROM product WHERE category_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3 ;

-- name: GetProductByCategoryAndPriceInc :many
SELECT * FROM product WHERE category_id = $1 ORDER BY price LIMIT $2 OFFSET $3;

-- name: GetProductByCategoryAndPriceDesc :many
SELECT * FROM product WHERE category_id = $1 ORDER BY price DESC LIMIT $2 OFFSET $3 ;

-- name: GetRecommendProduct :many

SELECT * FROM product LIMIT $1 OFFSET $2;

-- name: GetProductBySupplier :many

SELECT * FROM product WHERE supplier_id = $1 LIMIT $2 OFFSET $3;

-- name: GetProductBySupplierAndCategory :many

SELECT * FROM product WHERE supplier_id = $1 AND category_id = $4 LIMIT $2 OFFSET $3;

-- name: GetProductBySupplierAndTime :many

SELECT * FROM product WHERE supplier_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetProductBySupplierAndTimeAndCategory :many

SELECT * FROM product WHERE supplier_id = $1 AND category_id = $4 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetProductBySupplierAndPriceDesc :many

SELECT * FROM product WHERE supplier_id = $1 ORDER BY price DESC LIMIT $2 OFFSET $3;

-- name: GetProductBySupplierAndPriceDescAndCategory :many

SELECT * FROM product WHERE supplier_id = $1 AND category_id = $4 ORDER BY price DESC LIMIT $2 OFFSET $3;

-- name: GetProductBySupplierAndPriceInc :many

SELECT * FROM product WHERE supplier_id = $1 ORDER BY price LIMIT $2 OFFSET $3;

-- name: GetProductBySupplierAndPriceIncAndCategory :many

SELECT * FROM product WHERE supplier_id = $1 AND category_id=$4 ORDER BY price LIMIT $2 OFFSET $3;


-- name: UpdateProduct :exec

UPDATE product
SET
    name = $2,
    price = $3,
    inventory = $4,
    brand = $5
WHERE id = $1 and supplier_id = $6;

-- name: GetListProductByIDs :many

SELECT * FROM product WHERE id IN ($1);

-- name: DescInventory :exec
UPDATE product
SET
    inventory = inventory - $1
WHERE id = $2 and inventory >= $1;

-- name: GetProductInventory :one
SELECT inventory FROM product
WHERE id = $1 LIMIT 1;

-- name: IncInventory :exec
UPDATE product
SET
    inventory = inventory + $1
WHERE id = $2;

-- name: DeleteProduct :exec
DELETE FROM product WHERE id = $1 and supplier_id = $2;

-- name: DeleteProductByID :exec
DELETE FROM product WHERE id = $1;