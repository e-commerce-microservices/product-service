package service

import (
	"database/sql"
	"fmt"

	"github.com/e-commerce-microservices/product-service/repository"
)

// ProductFilter ...
type ProductFilter struct {
	ByTime      bool
	ByPriceInc  bool
	ByPriceDesc bool
}

func (p ProductFilter) GenerateListProduct(db *sql.DB, categoryID int64, limit, offset int32) ([]repository.Product, error) {
	var getProductByCategory = `SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product WHERE category_id = $1 LIMIT $2 OFFSET $3`
	if p.ByTime {
		getProductByCategory = fmt.Sprintf("SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product WHERE category_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3")
	} else if p.ByPriceInc {
		getProductByCategory = fmt.Sprintf("SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product WHERE category_id = $1 ORDER BY price LIMIT $2 OFFSET $3")
	} else if p.ByPriceDesc {
		getProductByCategory = fmt.Sprintf("SELECT id, name, description, price, thumbnail, inventory, supplier_id, category_id, created_at, brand FROM product WHERE category_id = $1 ORDER BY price DESC LIMIT $2 OFFSET $3")
	}

	rows, err := db.Query(getProductByCategory, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []repository.Product
	for rows.Next() {
		var i repository.Product
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
