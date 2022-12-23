-- name: CreateCategory :exec
INSERT INTO category ("id", "name") 
VALUES ($1, $2);
