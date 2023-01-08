-- name: CreateCategory :exec
INSERT INTO category ("id", "name", "thumbnail") 
VALUES ($1, $2, $3);

-- name: GetAllCategory :many
SELECT * FROM category;