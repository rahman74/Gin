package repository

import (
	"database/sql"
	"log"
	"test/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type paymentMethodsRepository struct {
	db *sqlx.DB
}

// PaymentMethodsRepository ...
type PaymentMethodsRepository interface {
	GetAllPaymentMethods(*sqlx.DB) ([]*model.PaymentMethods, error)
}

// NewPaymentMethodsRepository ...
func NewPaymentMethodsRepository(db *sqlx.DB) PaymentMethodsRepository {
	return &paymentMethodsRepository{
		db,
	}
}

// GetAllOrders ...
func (paymentMethodsRepository *paymentMethodsRepository) GetAllPaymentMethods(db *sqlx.DB) ([]*model.PaymentMethods, error) {
	var paymentMethods []*model.PaymentMethods
	rows, sqlError := db.Query(`
			SELECT payment_method_id, method_name, code, created_date
			FROM public.payment_methods;	
		`,
	)

	if sqlError != nil {
		log.Print("SQL error on GetAllPaymentMethods => ", sqlError)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ID int
			var createdDate time.Time
			var name, code string

			sqlError := rows.Scan(
				&ID,
				&name,
				&code,
				&createdDate,
			)

			if sqlError != nil {
				log.Print("SQL error on GetAllPaymentMethods => ", sqlError)
			} else {
				paymentMethods = append(
					paymentMethods,
					&model.PaymentMethods{
						MethodName:          name,
						PaymentMethodID:        ID,
						Code: code,
						CreatedDate: createdDate,
					},
				)
			}
		}
	}

	if sqlError == sql.ErrNoRows {
		sqlError = nil
	}

	return paymentMethods, sqlError
}

