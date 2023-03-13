// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: product_crud.sql

package repository

import (
	"context"
	"database/sql"
)

const createProduct = `-- name: CreateProduct :one

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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand
`

type CreateProductParams struct {
	Name        string
	Description string
	Price       int64
	Thumbnail   string
	Inventory   int32
	SupplierID  int64
	CategoryID  int64
	Brand       sql.NullString
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProduct,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.Thumbnail,
		arg.Inventory,
		arg.SupplierID,
		arg.CategoryID,
		arg.Brand,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Thumbnail,
		&i.Inventory,
		&i.SupplierID,
		&i.CategoryID,
		&i.CreatedAt,
		&i.Brand,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM product WHERE id = $1 and supplier_id = $2
`

type DeleteProductParams struct {
	ID         int64
	SupplierID int64
}

func (q *Queries) DeleteProduct(ctx context.Context, arg DeleteProductParams) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, arg.ID, arg.SupplierID)
	return err
}

const deleteProductByID = `-- name: DeleteProductByID :exec
DELETE FROM product WHERE id = $1
`

func (q *Queries) DeleteProductByID(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProductByID, id)
	return err
}

const descInventory = `-- name: DescInventory :exec
UPDATE product
SET
    inventory = inventory - $1
WHERE id = $2 and inventory >= $1
`

type DescInventoryParams struct {
	Inventory int32
	ID        int64
}

func (q *Queries) DescInventory(ctx context.Context, arg DescInventoryParams) error {
	_, err := q.db.ExecContext(ctx, descInventory, arg.Inventory, arg.ID)
	return err
}

const getAllProduct = `-- name: GetAllProduct :many

SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product
`

func (q *Queries) GetAllProduct(ctx context.Context) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getAllProduct)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Thumbnail,
			&i.Inventory,
			&i.SupplierID,
			&i.CategoryID,
			&i.CreatedAt,
			&i.Brand,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getListProductByIDs = `-- name: GetListProductByIDs :many

SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product WHERE id IN ($1)
`

func (q *Queries) GetListProductByIDs(ctx context.Context, id string) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Thumbnail,
			&i.Inventory,
			&i.SupplierID,
			&i.CategoryID,
			&i.CreatedAt,
			&i.Brand,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductByCategory = `-- name: GetProductByCategory :many

SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product WHERE category_id = $1 LIMIT $2 OFFSET $3
`

type GetProductByCategoryParams struct {
	CategoryID int64
	Limit      int32
	Offset     int32
}

func (q *Queries) GetProductByCategory(ctx context.Context, arg GetProductByCategoryParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductByCategory, arg.CategoryID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Thumbnail,
			&i.Inventory,
			&i.SupplierID,
			&i.CategoryID,
			&i.CreatedAt,
			&i.Brand,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductByID = `-- name: GetProductByID :one

SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product WHERE id = $1
`

func (q *Queries) GetProductByID(ctx context.Context, id int64) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProductByID, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Thumbnail,
		&i.Inventory,
		&i.SupplierID,
		&i.CategoryID,
		&i.CreatedAt,
		&i.Brand,
	)
	return i, err
}

const getProductBySupplier = `-- name: GetProductBySupplier :many

SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product WHERE supplier_id = $1 LIMIT $2 OFFSET $3
`

type GetProductBySupplierParams struct {
	SupplierID int64
	Limit      int32
	Offset     int32
}

func (q *Queries) GetProductBySupplier(ctx context.Context, arg GetProductBySupplierParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductBySupplier, arg.SupplierID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Thumbnail,
			&i.Inventory,
			&i.SupplierID,
			&i.CategoryID,
			&i.CreatedAt,
			&i.Brand,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductInventory = `-- name: GetProductInventory :one
SELECT inventory FROM product
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProductInventory(ctx context.Context, id int64) (int32, error) {
	row := q.db.QueryRowContext(ctx, getProductInventory, id)
	var inventory int32
	err := row.Scan(&inventory)
	return inventory, err
}

const getRecommendProduct = `-- name: GetRecommendProduct :many

SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product LIMIT $1 OFFSET $2
`

type GetRecommendProductParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetRecommendProduct(ctx context.Context, arg GetRecommendProductParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getRecommendProduct, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Thumbnail,
			&i.Inventory,
			&i.SupplierID,
			&i.CategoryID,
			&i.CreatedAt,
			&i.Brand,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const incInventory = `-- name: IncInventory :exec
UPDATE product
SET
    inventory = inventory + $1
WHERE id = $2
`

type IncInventoryParams struct {
	Inventory int32
	ID        int64
}

func (q *Queries) IncInventory(ctx context.Context, arg IncInventoryParams) error {
	_, err := q.db.ExecContext(ctx, incInventory, arg.Inventory, arg.ID)
	return err
}

const updateProduct = `-- name: UpdateProduct :exec

UPDATE product
SET
    name = $2,
    price = $3,
    inventory = $4,
    brand = $5
WHERE id = $1 and supplier_id = $6
`

type UpdateProductParams struct {
	ID         int64
	Name       string
	Price      int64
	Inventory  int32
	Brand      sql.NullString
	SupplierID int64
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.ExecContext(ctx, updateProduct,
		arg.ID,
		arg.Name,
		arg.Price,
		arg.Inventory,
		arg.Brand,
		arg.SupplierID,
	)
	return err
}
