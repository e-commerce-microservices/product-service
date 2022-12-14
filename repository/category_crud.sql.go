// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: category_crud.sql

package repository

import (
	"context"
	"database/sql"
)

const createCategory = `-- name: CreateCategory :exec
INSERT INTO category ("id", "name", "thumbnail") 
VALUES ($1, $2, $3)
`

type CreateCategoryParams struct {
	ID        int64
	Name      string
	Thumbnail sql.NullString
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) error {
	_, err := q.db.ExecContext(ctx, createCategory, arg.ID, arg.Name, arg.Thumbnail)
	return err
}

const getAllCategory = `-- name: GetAllCategory :many
SELECT id, name, created_at, thumbnail FROM category
`

func (q *Queries) GetAllCategory(ctx context.Context) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, getAllCategory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.Thumbnail,
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
