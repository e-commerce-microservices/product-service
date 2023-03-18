-- name: CreateCategory :exec
INSERT INTO category ("id", "name", "thumbnail") 
VALUES ($1, $2, $3);

-- name: GetAllCategory :many
SELECT * FROM category;

-- name: GetCategoryBySupplier :many
SELECT DISTINCT "category_id", "category"."name" FROM "product" JOIN "category" ON "product"."category_id" = "category"."id"
WHERE "supplier_id" = $1;