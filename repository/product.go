package repository

import (
	"database/sql"
	"log"
	"test/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type productRepository struct {
	db *sqlx.DB
}

// ProductRepository ...
type ProductRepository interface {
	GetAllProduct(*sqlx.DB) ([]*model.Product, error)
}

// NewProductRepository ...
func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{
		db,
	}
}

// GetAllProduct ...
func (productRepository *productRepository) GetAllProduct(db *sqlx.DB) ([]*model.Product, error) {
	var products []*model.Product
	rows, sqlError := db.Query(`
			SELECT product_id, product_name, basic_price, created_date
			FROM public.products;
		`,
	)

	if sqlError != nil {
		log.Print("SQL error on GetAllProduct => ", sqlError)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ID int
			var name string
			var basicPrice float64
			var createdDate time.Time

			sqlError := rows.Scan(
				&ID,
				&name,
				&basicPrice,
				&createdDate,
			)

			if sqlError != nil {
				log.Print("SQL error on GetAllProduct => ", sqlError)
			} else {
				products = append(
					products,
					&model.Product{
						ProductID:          ID,
						ProductName:        name,
						BasicPrice: basicPrice,
						CreatedDate: createdDate,
					},
				)
			}
		}
	}

	if sqlError == sql.ErrNoRows {
		sqlError = nil
	}

	return products, sqlError
}

