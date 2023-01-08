// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: product_crud.sql

package repository

import (
	"context"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO product (name, description, price, thumbnail, inventory, supplier_id, category_id)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at
`

type CreateProductParams struct {
	Name        string
	Description string
	Price       int64
	Thumbnail   string
	Inventory   int32
	SupplierID  int64
	CategoryID  int64
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
	)
	return i, err
}

const getAllProduct = `-- name: GetAllProduct :many
SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at FROM product
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
SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at FROM product WHERE category_id = $1
`

func (q *Queries) GetProductByCategory(ctx context.Context, categoryID int64) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductByCategory, categoryID)
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
SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at FROM product WHERE id = $1
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
	)
	return i, err
}
